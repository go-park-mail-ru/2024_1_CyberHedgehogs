package handler

import (
	"encoding/json"
	rep "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"io"
	"net/http"
	"time"
)

type SessionHandler struct {
	sessionTable *rep.SessionTable
}

func NewSessionHandler(users *rep.UserTable) *SessionHandler { // users по указателю?

	sessionTable := &rep.SessionTable{
		Sessions: make(map[string]rep.Session),
		Users:    users.Users,
	}

	return &SessionHandler{
		sessionTable: sessionTable,
	}
}

func (api *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, `{"error": "Session cookie not found"}`, http.StatusBadRequest)
		} else {
			http.Error(w, `{"error": "Error retrieving session cookie"}`, http.StatusBadRequest)
		}
		return
	}

	err = api.sessionTable.DeleteSession(sessionCookie.Value)
	if err != nil {
		http.Error(w, `error deleting session`, http.StatusBadRequest)
		return
	}

}

func (api *SessionHandler) Login(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "Error reading request body"}`, http.StatusBadRequest)
		return
	}

	var user rep.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	if user.Login == "" || user.Password == "" {
		http.Error(w, `{"error": "login and password fields cannot be empty"}`, http.StatusBadRequest)
		return
	}

	sessionID, err := api.sessionTable.AddSession(user)
	if err != nil {
		http.Error(w, "wrong login or password", http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(10 * time.Minute),
	}
	http.SetCookie(w, cookie)

}
