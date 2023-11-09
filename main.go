package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/itka0526/gtorrent/src"
)

const Unauthorized = "please login"

type Server struct {
	listAddr   string
	pathToCred string
	store      *sessions.CookieStore
	wsHub      *src.Hub
}

type apiFn func(http.ResponseWriter, *http.Request) error

type Message struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func (t Message) toJSON() string {
	b, err := json.Marshal(t)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func makeHTTPHandleFn(fn apiFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, Message{Status: false, Message: err.Error()}.toJSON())
		}
	}
}

func main() {
	hub := src.NewHub()
	s := &Server{listAddr: ":3000", pathToCred: "./.env", wsHub: hub}
	s.Run()
}

func (s *Server) Run() {
	env, _ := s.readEnv()
	s.store = sessions.NewCookieStore([]byte(env.SecretKey))
	s.store.Options.MaxAge = int(time.Hour) * 24

	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("static/assets"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	r.HandleFunc("/api/auth", makeHTTPHandleFn(s.handleAuth)).Methods("POST")
	r.HandleFunc("/api/upload", makeHTTPHandleFn(s.handleUpload)).Methods("POST")
	r.HandleFunc("/api/delete", makeHTTPHandleFn(s.handleDelete)).Methods("POST")
	r.HandleFunc("/api/download", makeHTTPHandleFn(s.handleDownload)).Methods("POST")
	r.HandleFunc("/api/magnet", makeHTTPHandleFn(s.handleMagnet)).Methods("POST")
	r.HandleFunc("/api/ws", s.handleWs)

	err := http.ListenAndServe(s.listAddr, r)
	fmt.Println(err)
}

func (s *Server) handleWs(w http.ResponseWriter, r *http.Request) {
	sess, err := s.store.Get(r, "credentials")
	if err != nil {
		fmt.Fprint(w, Message{Status: false, Message: "please re-login again"}.toJSON())
	}
	// TODO: make some kind of error handling here
	if _, ok := sess.Values["uuid"]; ok {
		src.NewClient(s.wsHub, sess.Values["uuid"].(string), w, r)
	} else {
		fmt.Fprint(w, Message{Status: false, Message: "failed to create a websocket connection"}.toJSON())
	}
}

func (s *Server) handleMagnet(w http.ResponseWriter, r *http.Request) error {
	if !s.isSessionActive(r) {
		return fmt.Errorf(Unauthorized)
	}
	rb, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var JSONReqBody struct {
		Url string `json:"url"`
	}

	if err := json.Unmarshal(rb, &JSONReqBody); err != nil {
		return err
	}

	if err := src.DownloadMagnet(JSONReqBody.Url); err != nil {
		return err
	}

	data := src.GetFiles()
	wsMsg, _ := json.Marshal(src.WSMessage{ResType: "get_files_res", Data: string(data)})
	s.wsHub.Broadcast <- wsMsg
	fmt.Fprint(w, Message{Status: true, Message: "Torrent file added."}.toJSON())
	return nil
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) error {
	if !s.isSessionActive(r) {
		return fmt.Errorf(Unauthorized)
	}

	rb, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var JSONReqBody src.MyFileInfo

	if err := json.Unmarshal(rb, &JSONReqBody); err != nil {
		return err
	}

	file, err := src.GetFile(JSONReqBody)
	if err != nil {
		return err
	}

	if JSONReqBody.IsDir {
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", JSONReqBody.Name))
	}

	w.Write(file)
	return nil
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) error {
	if !s.isSessionActive(r) {
		return fmt.Errorf(Unauthorized)
	}

	rb, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var JSONReqBody src.MyFileInfo

	if err := json.Unmarshal(rb, &JSONReqBody); err != nil {
		return err
	}

	if err := src.DeleteFile(JSONReqBody); err != nil {
		return err
	}

	data := src.GetFiles()
	wsMsg, _ := json.Marshal(src.WSMessage{ResType: "get_files_res", Data: string(data)})
	s.wsHub.Broadcast <- wsMsg

	fmt.Fprint(w, Message{Status: true, Message: fmt.Sprintf(`"%s" was removed.`, JSONReqBody.Name)}.toJSON())
	return nil
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) error {
	if !s.isSessionActive(r) {
		return fmt.Errorf(Unauthorized)
	}

	var fileNames []string
	json.Unmarshal([]byte(r.FormValue("file_names")), &fileNames)

	for _, fn := range fileNames {
		f, fh, err := r.FormFile(fn)
		if err != nil {
			return fmt.Errorf("cannot read %s", fn)
		}

		defer f.Close()

		if err := src.SaveFile(&f, fh); err != nil {
			return err
		}
	}

	data := src.GetFiles()
	wsMsg, _ := json.Marshal(src.WSMessage{ResType: "get_files_res", Data: string(data)})
	s.wsHub.Broadcast <- wsMsg

	fmt.Fprint(w, Message{Status: true, Message: "Files were successfully uploaded."}.toJSON())
	return nil
}

func (s *Server) isSessionActive(r *http.Request) bool {
	sess, err := s.store.Get(r, "credentials")
	if err != nil {
		return false
	}
	_, ok := sess.Values["uuid"]

	isAuth := sess.Values["auth"] == true && ok

	return isAuth
}

func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) error {
	if s.isSessionActive(r) {
		fmt.Fprint(w, Message{Status: true, Message: "You are authorized."}.toJSON())
		return nil
	}
	// Retrieve credentials from a file
	env, err := s.readEnv()
	if err != nil {
		return err
	}

	// Check if user is authorized
	var sentCreds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sentCreds); err != nil {
		return fmt.Errorf("credentials structure is wrong")
	}

	if sentCreds.Username != env.Username || sentCreds.Password != env.Password {
		return fmt.Errorf("credentials are wrong")
	}

	if _, err := s.setCookie(r, w); err != nil {
		return err
	}

	fmt.Fprint(w, Message{Status: true, Message: "You are authorized."}.toJSON())
	return nil
}

func (s *Server) setCookie(r *http.Request, w http.ResponseWriter) (uuid.UUID, error) {
	u := uuid.New()

	session, _ := s.store.Get(r, "credentials")
	session.Values["auth"] = true
	session.Values["uuid"] = u.String()
	sessErr := session.Save(r, w)

	if sessErr != nil {
		return u, fmt.Errorf("cannot set a cookie")
	}

	return u, nil
}

type Env struct {
	Username  string
	Password  string
	SecretKey string
}

func (s *Server) readEnv() (Env, error) {
	b, err := os.ReadFile(s.pathToCred)
	if err != nil {
		return Env{}, fmt.Errorf("cannot access %s file or the the %s file is empty", s.pathToCred, s.pathToCred)
	}
	creds := strings.Split(string(b), "\n")
	env := Env{
		Username:  strings.Split(creds[0], ":")[1],
		Password:  strings.Split(creds[1], ":")[1],
		SecretKey: strings.Split(creds[2], ":")[1],
	}
	return env, nil
}
