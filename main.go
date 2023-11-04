package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

type Server struct {
	connectedUsers Users
	listAddr       string
	pathToCred     string
	store          *sessions.CookieStore
}

type Users struct {
	users []User
}
type User struct {
	uuid uuid.UUID
}

func (u Users) isUserValid(otherUser string) bool {
	for _, user := range u.users {
		if user.uuid.String() == otherUser {
			return true
		}
	}
	return false
}

var ugrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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
			fmt.Fprint(w, Message{Status: false, Message: err.Error()}.toJSON())
		}
	}
}

func main() {
	s := &Server{connectedUsers: Users{}, listAddr: ":3000", pathToCred: "./.env"}
	s.Run()
}

func (s *Server) Run() {
	env, _ := s.readEnv()
	s.store = sessions.NewCookieStore([]byte(env.SecretKey))
	router := mux.NewRouter()

	router.HandleFunc("/api/auth", makeHTTPHandleFn(s.handleAuth))
	router.HandleFunc("/api/ws", makeHTTPHandleFn(handNewWsConn))

	err := http.ListenAndServe(s.listAddr, router)

	fmt.Println(err)
}

func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) error {
	sess, err := s.store.Get(r, "credentials")
	if err != nil {
		return err
	}

	uuid, ok := sess.Values["uuid"]

	isAuth := sess.Values["auth"] == true && ok && s.connectedUsers.isUserValid(fmt.Sprint(uuid))
	isBadMethod := r.Method != http.MethodPost

	switch {
	case isAuth:
		fmt.Fprint(w, Message{Status: true, Message: "You are authorized."}.toJSON())
	case isBadMethod:
		return fmt.Errorf("invalid HTTP method: %s", r.Method)
	default:
		// Retrieve credentials from a file
		env, err := s.readEnv()
		if err != nil {
			return err
		}

		// Check if user is authorized to proceed
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
	}
	return nil
}

func (s *Server) setCookie(r *http.Request, w http.ResponseWriter) (uuid.UUID, error) {
	u := uuid.New()

	session, _ := s.store.Get(r, "credentials")
	session.Values["authorized"] = true
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

func handNewWsConn(w http.ResponseWriter, r *http.Request) error {
	conn, err := ugrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	reader(conn)
	return nil
}

func reader(conn *websocket.Conn) {
	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Message: ", string(p))

		conn.WriteMessage(msgType, []byte("Hello from Server! "))
	}
}
