package handler

import (
	"fmt"
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

type AuthHandler struct {
	sesManager SessionManager
	userRepo   UserRepository
	render     repo.Renderer
}

func NewAuthHandler(ur UserRepository, sm SessionManager) AuthHandler {
	handler := AuthHandler{sesManager: sm, userRepo: ur}
	return handler
}

func (api *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

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

func (api *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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

func (api *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var user repo.User
	err := api.render.DecodeJSON(r.Body, &user)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}
	tableUser, err := api.userRepo.AddUser(&user)
	if err != nil {
		http.Error(w, `{"error": "Error adding user"}`, http.StatusInternalServerError)
		return
	}
	userInfo := repo.UserSessionInfo{UserID: tableUser.ID, Login: tableUser.Login}

	sessionID, err := api.sesManager.AddSession(&userInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create session: %v", err), http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, sessionID)
	w.WriteHeader(http.StatusOK)
}

type Post struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	AuthorID    int    `json:"author_id,omitempty"`
}

func (api *AuthHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "no cookie", http.StatusBadRequest)
		return
	}

	userInfo := api.sesManager.CheckSession(sessionCookie.Value)
	if userInfo == nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	userProfile := &repo.User{
		ID:       userInfo.UserID,
		Login:    userInfo.Login,
		Username: "mock_username",
		Email:    "mock@example.com",
	}

	api.render.EncodeJSON(w, http.StatusOK, userProfile)
}

func (api *AuthHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "no cookie", http.StatusBadRequest)
		return
	}

	userInfo := api.sesManager.CheckSession(sessionCookie.Value)
	if userInfo == nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	userPosts := map[int]*Post{
		1: {ID: 1, Title: "Mock Post 1", Description: "Description of Mock Post 1", Author: userInfo.Login},
		2: {ID: 2, Title: "Mock Post 2", Description: "Description of Mock Post 2", Author: userInfo.Login},
	}

	api.render.EncodeJSON(w, http.StatusOK, userPosts)
}
