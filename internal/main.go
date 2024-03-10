package main

import (
	"log"
	"net/http"

	_ "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/docs"
	"github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/handler"
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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

// @title			Saudade API
// @version			1.0
// @description 	API Server for Patreon like Application

// @host 			localhost:3031
// @BasePath 		/
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

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3031/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":3031", r))

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
