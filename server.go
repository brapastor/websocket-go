package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"sync"
	"github.com/gorilla/websocket"
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
type Response struct {
	Messagge string `json:"message"`
	Status int `json:"status"`
	IsValid bool `json:"isvalid"`
}

var Users = struct {
	m map[string] User
	sync.RWMutex
}{m: make(map[string] User)}

type User struct {
	User_Name string
	WebSocket *websocket.Conn
}

func UserExist(user_name string) bool  {
	Users.RLock() // bloqueamos la estructura
	defer Users.RUnlock()
	if _, ok := Users.m[user_name]; ok {
		return true
	}
	return false
}
func CreateUser(user_name string, ws *websocket.Conn) User {
	return User{user_name, ws}
}
func SendMessage(type_message int, message []byte)  {
	Users.RLock();
	defer Users.RUnlock()
	for _, user := range Users.m{
		if err := user.WebSocket.WriteMessage(type_message, message); err != nil{
			return
		}
	}
}
func AddUser(user User)  {
	Users.Lock()
	defer Users.Unlock()
	Users.m[user.User_Name] = user
}
func RemoveUser(user_name string)  {
	log.Println("El usuario se fue")
	Users.Lock()
	defer Users.Unlock()
	delete(Users.m, user_name)
}

func ToArryByte(value string) []byte {
	return []byte(value)
}

func ConcatMessage(user_name string, array []byte) string  {
	return user_name + " : " + string(array[:])
}

func WebSocket(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	user_name := vars["user_name"] //obtenemos el username de la url

	//Creamos un websocker
	ws, err := upgrader.Upgrade(w,r,nil)
	if err != nil{
		log.Println(err)
		return
	}
	//Creamos un usuario
	current_user := CreateUser(user_name, ws)
	//Almacenamos el usuario en el map
	AddUser(current_user)
	log.Println("Usuario agregado")

	//Queda a escucha si el cliente nos manda un mensaje
	for{
		type_message, message, err := ws.ReadMessage()
		if err != nil{
			RemoveUser(user_name)
			return
		}
		final_message := ConcatMessage(user_name, message)
		SendMessage(type_message, ToArryByte(final_message))
	}
}
func main()  {
	cssHandle := http.FileServer(http.Dir("./Front/css/"))
	jsHandle := http.FileServer(http.Dir("./Front/js/"))

	mux := mux.NewRouter()
	mux.HandleFunc("/hola", HolaMundo).Methods("GET")
	mux.HandleFunc("/holajson", HolaMundoJson).Methods("GET")
	mux.HandleFunc("/", LoadStatic).Methods("GET")
	mux.HandleFunc("/validate", Validate).Methods("POST")
	mux.HandleFunc("/chat/{user_name}", WebSocket).Methods("GET")


	//Utilizar url de mux
	http.Handle("/", mux)
	http.Handle("/css/", http.StripPrefix("/css/", cssHandle))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandle))
	log.Println("EL SERVIDOR SE ENCUENTRA EN EL PUERTO 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func Validate(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	user_name := r.FormValue("user_name")
	response := Response{}
	if UserExist(user_name){
		response.IsValid = false
	}else {
		response.IsValid = true
	}
	json.NewEncoder(w).Encode(response)
}

func LoadStatic(w http.ResponseWriter, r *http.Request)  {
	http.ServeFile(w,r,"./Front/index.html")
}

func HolaMundoJson(w http.ResponseWriter, r *http.Request)  {
	response := CreateResponse("esto esta en formato json", 200, true)
	json.NewEncoder(w).Encode(response)
}

func HolaMundo(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("Hola mundo desde go"))
}


func CreateResponse(message string, status int, valid bool)  Response {
	return Response{message, status, valid}
}