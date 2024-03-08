package main

import (
	"github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/handler"
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	userTable := &repo.UserTable{}
	userTable.Users = make(map[string]*repo.User)
	sessionTable := &repo.SessionTable{}
	sessionTable.Sessions = make(map[string]*repo.Session)
	api := handler.NewAuthHandler(userTable, sessionTable)
	r := mux.NewRouter()
	r.HandleFunc("/register", api.Registration).Methods("POST")
	r.HandleFunc("/login", api.Login).Methods("POST")
	r.HandleFunc("/logout", api.Logout).Methods("POST")
	r.HandleFunc("/profile", api.GetUserProfile).Methods("GET")
	r.HandleFunc("/posts", api.GetUserPosts).Methods("GET")
	http.ListenAndServe(":8081", r)
}
