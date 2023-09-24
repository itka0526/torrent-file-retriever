package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
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

var pathToCred = "./.secret"

type apiFn func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFn(fn apiFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Println(err)
			fmt.Fprintf(w, "This was not suppose to happen. ")
		}
	}
}

func main() {
	router := mux.NewRouter()

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Println("hello")
		fmt.Fprintln(w, "Hello")
	})

	router.HandleFunc("/api/auth", makeHTTPHandleFn(handleAuth))
	router.HandleFunc("/api/ws", makeHTTPHandleFn(handNewWsConn))

	err := http.ListenAndServe(":3001", nil)

	fmt.Println(err)
}

func handleAuth(_ http.ResponseWriter, r *http.Request) error {
	// This function should set a cookie on the client
	// This cookie will be used for ensuring authentication for subsequent websocket requests
	b, err := os.ReadFile(pathToCred)
	if err != nil {
		return err
	}

	creds := strings.Split(string(b), "\n")
	username := strings.Split(creds[0], ":")[1]
	password := strings.Split(creds[1], ":")[1]

	log.Println(r.Header, username, password)

	return nil
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
