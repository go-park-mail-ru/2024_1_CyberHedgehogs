package main

import (
	"github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/handler"
	rep "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	userTabl := rep.UserTable{}
	sessionTabl := rep.SessionTable{}
	sesContorller := handler.NewSessionsGoController(&sessionTabl, &userTabl)
	r := mux.NewRouter()

	api := handler.NewRegistrationHandler(&userTabl)
	api2 := handler.NewSessionHandler(sesContorller)

	r.HandleFunc("/register", api.Registration).Methods("POST")
	r.HandleFunc("/login", api2.Login).Methods("POST")
	r.HandleFunc("/logout", api2.Logout).Methods("POST")

	http.ListenAndServe(":8080", r)
}
