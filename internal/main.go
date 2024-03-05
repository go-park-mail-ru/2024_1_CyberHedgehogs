package main

import (
	"github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/handler"
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	userTable := &repo.UserTable{}
	sessionTable := &repo.SessionTable{}
	api := handler.NewAuthHandler(userTable, sessionTable)
	r := mux.NewRouter()
	r.HandleFunc("/register", api.Registration).Methods("POST")
	r.HandleFunc("/login", api.Login).Methods("POST")
	r.HandleFunc("/logout", api.Logout).Methods("POST")
	http.ListenAndServe(":8080", r)
}
