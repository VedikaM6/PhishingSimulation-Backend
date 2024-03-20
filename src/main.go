package main

import (
	"fmt"
	"net/http"

	"example.com/m/src/db"
	"example.com/m/src/emails"
	"example.com/m/src/users"
	"example.com/m/src/util"
	"github.com/gorilla/mux"
)

var router *mux.Router = mux.NewRouter()

func init() {
	// User routes
	router.HandleFunc("/users", users.ListUsers).Methods(http.MethodGet)
	router.HandleFunc("/user", users.AddUser).Methods(http.MethodPut)

	// Email routes
	router.HandleFunc("/emails", emails.ListEmails).Methods(http.MethodGet)
	router.HandleFunc("/emails/{"+util.URLParameterEmailId+"}", emails.GetAttackEmail).Methods(http.MethodGet)

	// Attack routes
	router.HandleFunc("/triggerAttack", attacks.TriggerAttacks).Methods(http.MethodPOST)
}

func main() {
	// Initialize the Mongo connection and client
	db.InitClient()
	defer db.DisconnectClient()

	fmt.Println("[main] Listening on port 80...")

	http.ListenAndServe(":80", router)
}
