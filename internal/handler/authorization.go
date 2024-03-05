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

type UserRepository interface {
	AddUser(user *repo.User) (*repo.User, error)
	ValidateUserCredentials(user *repo.User) *repo.User
}

type SessionManager interface {
	AddSession(user *repo.UserSessionInfo) (sessionID string, err error)
	DeleteSession(sessionID string) error
	CheckSession(sessionID string) (user *repo.UserSessionInfo)
	SessionsCleanup()
}

type authHandler struct {
	sesManager SessionManager
	userRepo   UserRepository
	render     repo.Renderer
}

func NewAuthHandler(ur UserRepository, sm SessionManager) authHandler {
	handler := authHandler{sesManager: sm, userRepo: ur}
	return handler
}

func (api *authHandler) Logout(w http.ResponseWriter, r *http.Request) {

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, `no cookie`, http.StatusBadRequest)
		return
	}

	err = api.sesManager.DeleteSession(sessionCookie.Value)
	if err != nil {
		http.Error(w, `error deleting session`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (api *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user repo.User
	if err := api.render.DecodeJSON(r.Body, &user); err != nil {
		http.Error(w, `"error": "Invalid JSON format"`, http.StatusBadRequest)
		return
	}
	tableUser := api.userRepo.ValidateUserCredentials(&user)
	if tableUser == nil {
		http.Error(w, `"error": wrong credentials"`, http.StatusUnauthorized)
		return
	}
	userInfo := repo.UserSessionInfo{UserID: tableUser.ID, Login: tableUser.Login}
	sessionID, err := api.sesManager.AddSession(&userInfo)
	if err != nil {
		http.Error(w, "wrong login or password", http.StatusUnauthorized)
		return
	}
	setSessionCookie(w, sessionID)
	w.WriteHeader(http.StatusOK)
}

func (api *authHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var user repo.User
	err := api.render.DecodeJSON(r.Body, &user)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}
	_, err = api.userRepo.AddUser(&user) // todo надо авторизовать пользователя
	if err != nil {
		http.Error(w, `{"error": "Error adding user"}`, http.StatusBadRequest)
		return
	}

	mes := repo.Info{
		Message: "success",
	}

	api.render.EncodeJSON(w, http.StatusCreated, mes)
}
