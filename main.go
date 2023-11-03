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
	listAddr string
}

var ugrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	pathToCred = "./.env"
	store      *sessions.CookieStore
)

type apiFn func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFn(fn apiFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Println(err)
			// TODO: Return pretty JSON errors with message and status
			fmt.Fprintf(w, "This was not suppose to happen")
		}
	}
}

func main() {
	env, _ := readEnv()
	store = sessions.NewCookieStore([]byte(env.SecretKey))
	router := mux.NewRouter()

	router.HandleFunc("/", makeHTTPHandleFn(handleIndex))
	router.HandleFunc("/api/auth", makeHTTPHandleFn(handleAuth))
	router.HandleFunc("/api/ws", makeHTTPHandleFn(handNewWsConn))

	err := http.ListenAndServe(":3001", router)

	fmt.Println(err)
}

func handleIndex(w http.ResponseWriter, r *http.Request) error {
	sess, _ := store.Get(r, "credentials")
	// TODO: Check if the user is authenticated
}

func handleAuth(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid HTTP method: %s", r.Method)
	}
	// Retrieve credentials from a file
	env, err := readEnv()
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

	if _, err := setCookie(r, w); err != nil {
		return err
	}

	return nil
}

func setCookie(r *http.Request, w http.ResponseWriter) (uuid.UUID, error) {
	u := uuid.New()

	session, _ := store.Get(r, "credentials")
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

func readEnv() (Env, error) {
	b, err := os.ReadFile(pathToCred)
	if err != nil {
		return Env{}, fmt.Errorf("cannot access %s file or the the %s file is empty", pathToCred, pathToCred)
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
