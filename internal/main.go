package main

import (
	"github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/handler"
	rep "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	userTable := rep.UserTable{}
	r := mux.NewRouter()
	api := handler.NewUserHandler(&userTable)
	api2 := handler.NewSessionHandler(&userTable)
	r.HandleFunc("POST /register", api.Registration)
	r.HandleFunc("POST /login", api2.Login)
	r.HandleFunc("POST /login", api2.Logout)
	http.ListenAndServe(":8080", r)
}
