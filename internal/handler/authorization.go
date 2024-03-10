package handler

import (
	"net/http"
	"time"

	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
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

// Logout godoc
// @Summary Logout
// @Description Logout user and delete session
// @Tags auth
// @Produce json
// @Param session_id header string true "Session ID"
// @Success 204 {object} repo.Response "Session deleted successfully"
// @Failure 400 {object} repo.Response "No cookie provided"
// @Failure 500 {object} repo.Response "Error deleting session"
// @Router /logout [post]
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

// Login handles the user login request.
// @Summary Login
// @Description Login with user credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param request body repo.User true "User credentials"
// @Success 200 {object} repo.Response "Successfully logged in"
// @Failure 400 {object} repo.Response "Invalid JSON format"
// @Failure 401 {object} repo.Response "Wrong credentials"
// @Failure 500 {object} repo.Response "Internal server error"
// @Router /login [post]
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

// Registration handles the user registration request.
// @Summary Registration
// @Tags auth
// @Description Create a new user account
// @ID create-account
// @Accept json
// @Produce json
// @Param user body repo.User true "User information"
// @Success 200 {object} repo.Response "Account successfully created"
// @Failure 400 {object} repo.Response "Bad Request"
// @Failure 500 {object} repo.Response "Adding user error"
// @Router /register [post]
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

// GetUserProfile
// @Summary Get user profile
// @Description Get user profile information based on session ID
// @Tags user
// @Produce json
// @Param session_id header string true "Session ID"
// @Success 200 {object} repo.Response "User profile retrieved successfully"
// @Failure 400 {object} repo.Response "No cookie provided"
// @Failure 401 {object} repo.Response "Session not found"
// @Router /profile [get]
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

// GetUserPosts
// @Summary Get user posts
// @Description Get posts created by the authenticated user
// @Tags user
// @Produce json
// @Param session_id header string true "Session ID"
// @Success 200 {object} repo.Response "User posts retrieved successfully"
// @Failure 400 {object} repo.Response "No cookie provided"
// @Failure 401 {object} repo.Response "Session not found"
// @Router /posts [get]
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
