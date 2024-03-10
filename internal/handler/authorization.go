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
		Name:     SessionCookieName,
		Value:    sessionID,
		Expires:  time.Now().Add(SessionCookieExpires),
		Secure:   true,
		HttpOnly: true,
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
		answ := repo.Response{false, "no cookie", http.StatusBadRequest}
		api.render.EncodeJSON(w, http.StatusBadRequest, answ)
		return
	}

	err = api.sesManager.DeleteSession(sessionCookie.Value)
	if err != nil {
		answ := repo.Response{false, "error deleting session", http.StatusInternalServerError}
		api.render.EncodeJSON(w, http.StatusInternalServerError, answ)

		return
	}
	answ := repo.Response{true, "", http.StatusNoContent}
	api.render.EncodeJSON(w, http.StatusNoContent, answ)
}

func (api *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user repo.User
	if err := api.render.DecodeJSON(r.Body, &user); err != nil {
		answ := repo.Response{false, "error: Invalid JSON format", http.StatusBadRequest}
		api.render.EncodeJSON(w, http.StatusBadRequest, answ)
		return
	}
	tableUser := api.userRepo.ValidateUserCredentials(&user)
	if tableUser == nil {
		answ := repo.Response{false, "error: wrong credentials", http.StatusUnauthorized}
		api.render.EncodeJSON(w, http.StatusUnauthorized, answ)
		return
	}
	userInfo := repo.UserSessionInfo{UserID: tableUser.ID, Login: tableUser.Login}
	sessionID, err := api.sesManager.AddSession(&userInfo)
	if err != nil {
		answ := repo.Response{false, "wrong login or password", http.StatusUnauthorized}
		api.render.EncodeJSON(w, http.StatusUnauthorized, answ)
		return
	}
	setSessionCookie(w, sessionID)
	answ := repo.Response{true, "", http.StatusOK}
	api.render.EncodeJSON(w, http.StatusOK, answ)
}

func (api *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var user repo.User
	err := api.render.DecodeJSON(r.Body, &user)
	if err != nil {
		answ := repo.Response{false, "error: Invalid JSON format", http.StatusBadRequest}
		api.render.EncodeJSON(w, http.StatusBadRequest, answ)
		return
	}
	tableUser, err := api.userRepo.AddUser(&user)
	if err != nil {
		answ := repo.Response{false, "error: Error adding user", http.StatusInternalServerError}
		api.render.EncodeJSON(w, http.StatusInternalServerError, answ)

		return
	}
	userInfo := repo.UserSessionInfo{UserID: tableUser.ID, Login: tableUser.Login}

	sessionID, err := api.sesManager.AddSession(&userInfo)
	if err != nil {
		answ := repo.Response{false, "failed to create session", http.StatusInternalServerError}
		api.render.EncodeJSON(w, http.StatusInternalServerError, answ)
		return
	}

	setSessionCookie(w, sessionID)

	answ := repo.Response{true, "", http.StatusOK}
	api.render.EncodeJSON(w, http.StatusOK, answ)

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
		answ := repo.Response{false, "no cookie", http.StatusBadRequest}
		api.render.EncodeJSON(w, http.StatusBadRequest, answ)
		return
	}

	userInfo := api.sesManager.CheckSession(sessionCookie.Value)
	if userInfo == nil {
		answ := repo.Response{false, "session not found", http.StatusUnauthorized}
		api.render.EncodeJSON(w, http.StatusUnauthorized, answ)
		return
	}

	userProfile := &repo.User{
		ID:       userInfo.UserID,
		Login:    userInfo.Login,
		Username: "mock_username",
		Email:    "mock@example.com",
	}

	answ := repo.Response{true, "", userProfile}
	api.render.EncodeJSON(w, http.StatusOK, answ)
}

func (api *AuthHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		answ := repo.Response{false, "no cookie", http.StatusBadRequest}
		api.render.EncodeJSON(w, http.StatusBadRequest, answ)
		return
	}

	userInfo := api.sesManager.CheckSession(sessionCookie.Value)
	if userInfo == nil {
		answ := repo.Response{false, "session not found", http.StatusUnauthorized}
		api.render.EncodeJSON(w, http.StatusUnauthorized, answ)
		return
	}

	userPosts := map[int]*Post{
		1: {ID: 1, Title: "Mock Post 1", Description: "Description of Mock Post 1", Author: userInfo.Login},
		2: {ID: 2, Title: "Mock Post 2", Description: "Description of Mock Post 2", Author: userInfo.Login},
	}

	answ := repo.Response{true, "", userPosts}
	api.render.EncodeJSON(w, http.StatusOK, answ)
}
