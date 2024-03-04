package handler

import (
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"net/http"
	"time"
)

const (
	SessionCookieName    = "session_id"
	SessionCookieExpires = repo.SessionLiveTime
)

func setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:    SessionCookieName,
		Value:   sessionID,
		Expires: time.Now().Add(SessionCookieExpires),
	}
	http.SetCookie(w, cookie)
}

type SessionManager interface {
	AddSession(user repo.User) (sessionID string, err error)
	DeleteSession(sessionID string) error
	CheckSession(sessionID string) (user *repo.UserSessionInfo)
}

type SessionHandler struct {
	manager SessionManager
}

func NewSessionsGoController(sessions *repo.SessionTable, users *repo.UserTable) SessionManager {
	controller := &repo.SessionsGoController{SessionsTabl: sessions, UsersTabl: users}
	return controller
}

func NewSessionHandler(sesManager SessionManager) SessionHandler {

	sessionHandler := SessionHandler{
		manager: sesManager,
	}
	return sessionHandler
}

func (api *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `no cookie`, http.StatusOK)
		return
	}

	err = api.manager.DeleteSession(sessionCookie.Value)
	if err != nil {
		http.Error(w, `error deleting session`, http.StatusBadRequest)
		return
	}

}

func (api *SessionHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user repo.User
	renderer := repo.Renderer{}

	if err := renderer.DecodeJSON(r.Body, &user); err != nil {
		http.Error(w, `"error": "Invalid JSON format"`, http.StatusBadRequest)
		return
	}

	sessionID, err := api.manager.AddSession(user)
	if err != nil {
		http.Error(w, "wrong login or password", http.StatusBadRequest)
		return
	}
	setSessionCookie(w, sessionID)
}
