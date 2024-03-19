package main

import (
	"fmt"
	"net/http"

	"example.com/m/src/db"
	"example.com/m/src/users"
	"github.com/gorilla/mux"
)

var router *mux.Router = mux.NewRouter()

func init() {
	// User routes
	router.HandleFunc("/users", users.ListUsers).Methods(http.MethodGet)
	router.HandleFunc("/user", users.AddUser).Methods(http.MethodPut)

	// Attack routes
}

func main() {
	// Initialize the Mongo connection and client
	db.InitClient()
	defer db.DisconnectClient()

	fmt.Println("[main] Listening on port 80...")

	http.ListenAndServe(":80", router)
}
