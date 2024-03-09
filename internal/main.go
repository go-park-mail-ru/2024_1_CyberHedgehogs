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
	r.HandleFunc("/register", api.Registration).Methods("POST")
	r.HandleFunc("/login", api.Login).Methods("POST")
	r.HandleFunc("/logout", api.Logout).Methods("POST")
	r.HandleFunc("/profile", api.GetUserProfile).Methods("GET")
	r.HandleFunc("/posts", api.GetUserPosts).Methods("GET")
	/*

		"Access-Control-Allow-Origin": "http://localhost:3030",
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Credentials": true,
	*/

	http.ListenAndServe(":3031",
		handlers.CORS(

			//handlers.ContentTypeHandler(, "application/json"),
			handlers.AllowedOrigins([]string{"http://localhost:3030"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			//	handlers.
			//	handlers."Set-Cookie": `session_id=${crypto.randomUUID()}; max-age=60`
			handlers.AllowCredentials(),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedHeaders([]string{"Content-Type"}),
		)(r))
}
