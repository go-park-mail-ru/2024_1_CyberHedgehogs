package main

import (
	"github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/handler"
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

/*
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do stuff before the handlers
		h.ServeHTTP(w, r)
		// do stuff after the hadlers

	})
}
func Middleware2(s string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// do stuff
			fmt.Println(s)
			h.ServeHTTP(w, r)
		})
	}
}

*/

func main() {

	userTable := &repo.UserTable{}
	userTable.Users = make(map[string]*repo.User)
	sessionTable := &repo.SessionTable{}
	sessionTable.Sessions = make(map[string]*repo.Session)
	api := handler.NewAuthHandler(userTable, sessionTable)

	r := mux.NewRouter()
	/*
		subRouter := r.PathPrefix("/cors/").Subrouter()
		subRouter.Use(Middleware2(""))

	*/
	r.HandleFunc("/register", api.Registration).Methods("POST")
	r.HandleFunc("/login", api.Login).Methods("POST")
	r.HandleFunc("/logout", api.Logout).Methods("POST")
	r.HandleFunc("/profile", api.GetUserProfile).Methods("GET")
	r.HandleFunc("/posts", api.GetUserPosts).Methods("GET")
	http.ListenAndServe(":8081",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
		)(r))
}
