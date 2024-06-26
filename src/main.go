package main

import (
	"fmt"
	"net/http"

	"example.com/m/src/attacks"
	"example.com/m/src/dashboard"
	"example.com/m/src/db"
	"example.com/m/src/emails"
	"example.com/m/src/users"
	"example.com/m/src/util"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var router *mux.Router = mux.NewRouter()

func init() {
	// User routes
	router.HandleFunc("/users", users.ListUsers).Methods(http.MethodGet)
	router.HandleFunc("/user", users.AddUser).Methods(http.MethodPut)

	// Email routes
	router.HandleFunc("/emails", emails.ListEmails).Methods(http.MethodGet)
	router.HandleFunc("/emails/{"+util.URLParameterEmailId+"}", emails.GetEmail).Methods(http.MethodGet)
	router.HandleFunc("/emails", emails.CreateNewEmail).Methods(http.MethodPut)

	// Attack routes
	router.HandleFunc("/attacks/triggerPending", attacks.TriggerPendingAttacks).Methods(http.MethodPost)
	router.HandleFunc("/attacks/history", attacks.ListPreviousAttacks).Methods(http.MethodGet)
	router.HandleFunc("/attacks/now", attacks.TriggerAttackNow).Methods(http.MethodPost)
	router.HandleFunc("/attacks/future", attacks.ScheduleFutureAttack).Methods(http.MethodPut)
	router.HandleFunc("/attacks/clicked/{"+util.URLParameterAttackId+"}/{"+util.URLParameterUserEmail+"}", attacks.RecordAttackResults).Methods(http.MethodGet)

	// Report routes
	router.HandleFunc("/dashboard/data", dashboard.GetGaugeData).Methods(http.MethodGet)

	// webpage routes
	// fs := http.FileServer(http.Dir("./static"))
	// router.Handle("/vedikacorp", fs)
}

func main() {
	// Initialize the Mongo connection and client
	db.InitClient()
	defer db.DisconnectClient()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
	})
	handler := c.Handler(router)

	fmt.Println("[main] Listening on port 80...")

	http.ListenAndServe(":80", handler)
}
