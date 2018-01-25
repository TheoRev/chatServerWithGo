package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	cssHandle := http.FileServer(http.Dir("./front/css"))
	jsHandle := http.FileServer(http.Dir("./front/js"))

	mux := mux.NewRouter()
	mux.HandleFunc("/hola", Hola).Methods("GET")
	mux.HandleFunc("/hola_json", HolaJSON).Methods("GET")
	mux.HandleFunc("/static", loadStatic).Methods("GET")
	mux.HandleFunc("/validate", validate).Methods("POST")
	mux.HandleFunc("/chat/{user_name}", webSocket).Methods("GET")

	http.Handle("/", mux)
	http.Handle("/css/", http.StripPrefix("/css/", cssHandle))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandle))
	log.Println("El server esta corriendo en el puerto: 3030")
	log.Fatal(http.ListenAndServe(":3030", nil))
}

func createUser(username string, ws *websocket.Conn) User {
	return User{username, ws}
}

// AddUser agrega usuarios a la variable de persistencia
func AddUser(user User) {
	users.Lock()
	defer users.Unlock()

	users.m[user.User_Name] = user
}

func removeUser(user_name string) {
	users.Lock()
	defer users.Unlock()
	delete(users.m, user_name)
}

func webSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	un := vars["user_name"]
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		return
	}

	current_user := createUser(un, ws)
	AddUser(current_user)
	log.Println("Nuevo usuario agregado")

	for {
		type_message, message, err := ws.ReadMessage()
		if err != nil {
			removeUser(un)
			return
		}
		final_message := concatMessage(un, message)
		sendMessage(type_message, toArrayByte(final_message))
	}
}

func concatMessage(username string, arreglo []byte) string {
	return username + " : " + string(arreglo[:])
}

func sendMessage(type_message int, message []byte) {
	users.RLock()
	defer users.RUnlock()

	for _, user := range users.m {
		err := user.WebSocket.WriteMessage(type_message, message)
		if err != nil {
			return
		}
	}
}

func toArrayByte(value string) []byte {
	return []byte(value)
}

func validate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user_name := r.FormValue("user_name")
	response := Response{}
	if userExist(user_name) {
		response.IsValid = false
	} else {
		response.IsValid = true
	}
	json.NewEncoder(w).Encode(response)
}

var users = struct {
	m map[string]User
	sync.RWMutex
}{m: make(map[string]User)}

func userExist(user string) bool {
	users.RLock()
	defer users.RUnlock()
	if _, ok := users.m[user]; ok {
		return true
	}
	return false
}

type User struct {
	User_Name string
	WebSocket *websocket.Conn
}

func loadStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./front/index.html")
}

// Response structura
type Response struct {
	Mensaje string `json:"message"`
	Status  int    `json:"status"`
	IsValid bool   `json:"isvalid"`
}

// Hola funcion de saludo
func Hola(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola mundo desde Go"))
}

// HolaJSON devuelve una estructura en formato json
func HolaJSON(w http.ResponseWriter, r *http.Request) {
	response := createRespose("Esto essta en formato Json", 200, true)
	json.NewEncoder(w).Encode(response)
}

func createRespose(message string, status int, valid bool) Response {
	return Response{message, status, valid}
}
