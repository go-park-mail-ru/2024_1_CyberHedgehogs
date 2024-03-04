package repository

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var nextUserID uint = 1

const SessionLiveTime = 10 * time.Minute

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type User struct {
	ID       uint   `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserTable struct {
	Users map[string]*User
	mu    sync.Mutex
}

type Info struct {
	Message string `json:"message"`
}

func (table *UserTable) AddUser(user User) error {
	user.ID = nextUserID
	nextUserID++
	table.mu.Lock()
	table.Users[user.Login] = &user
	table.mu.Unlock()
	return nil
}

func (table *UserTable) ValidateNewUser(user User) error {
	if !emailRegex.MatchString(user.Email) {
		return errors.New("неверный формат email")
	}

	if len(user.Login) < 3 {
		return errors.New("login должен быть не менее 3 символов")
	}
	_, exists := table.Users[user.Login]
	if exists {
		return errors.New("пользователь с таким логином уже существует")
	}
	if len(user.Username) < 3 {
		return errors.New("username должен быть не менее 3 символов")
	}

	if len(user.Password) < 6 {
		return errors.New("пароль должен быть не менее 6 символов")
	}

	return nil
}

type UserSessionInfo struct {
	ID    uint   `json:"id,omitempty"`
	Login string `json:"login"`
}
type Session struct {
	ExpirationDate time.Time
	User           *UserSessionInfo `json:"user_info"`
	ID             string           `json:"session_id"`
}

type SessionTable struct {
	Sessions map[string]Session
	mu       sync.Mutex
}

func (table *SessionTable) IsActiveSession(sessionID string) (user *UserSessionInfo) {
	nowTime := time.Now()
	table.mu.Lock()
	if session, ok := table.Sessions[sessionID]; ok {
		if session.ExpirationDate.Before(nowTime) {
			delete(table.Sessions, sessionID)
			table.mu.Unlock()
			return nil
		}
		table.mu.Unlock()
		return session.User
	}
	table.mu.Unlock()
	return nil
}

type SessionsGoController struct {
	SessionsTabl *SessionTable
	UsersTabl    *UserTable
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (table *SessionsGoController) sessionsCleanup() {
	table.SessionsTabl.mu.Lock()
	for sessionID, sessn := range table.SessionsTabl.Sessions {
		if sessn.ExpirationDate.Before(time.Now()) {
			delete(table.SessionsTabl.Sessions, sessionID)
		}
	}
	table.UsersTabl.mu.Unlock()
	return
}

func (table *SessionsGoController) DeleteSession(sessionID string) error {
	table.SessionsTabl.mu.Lock()
	if _, ok := table.SessionsTabl.Sessions[sessionID]; ok {
		delete(table.SessionsTabl.Sessions, sessionID)
		table.SessionsTabl.mu.Unlock()
		return nil
	}
	table.UsersTabl.mu.Unlock()
	return nil
}

func (table *SessionsGoController) CheckSession(sessionID string) (user *UserSessionInfo) {
	if foundUser := table.SessionsTabl.IsActiveSession(sessionID); foundUser != nil {
		return foundUser
	}
	return nil
}

func (table *SessionsGoController) AddSession(user User) (sessionID string, err error) {
	if user.Login == "" {
		return "", errors.New("unvalid login or password")
	}
	table.UsersTabl.mu.Lock()
	tableUser, ok := table.UsersTabl.Users[user.Login]
	if !ok {
		table.UsersTabl.mu.Unlock()
		return "", errors.New("unvalid login or password")

	}

	table.UsersTabl.mu.Unlock()
	if tableUser.Password != user.Password {
		return "", errors.New("unvalid login or password")
	}
	sessionID = RandStringRunes(28)
	newUser := UserSessionInfo{ID: tableUser.ID, Login: tableUser.Login}
	ses := Session{User: &newUser, ExpirationDate: time.Now().Add(SessionLiveTime), ID: sessionID}
	table.SessionsTabl.mu.Lock()
	table.SessionsTabl.Sessions[sessionID] = ses
	table.SessionsTabl.mu.Unlock()
	return sessionID, nil
}

type Renderer struct{}

func (r *Renderer) DecodeJSON(body io.ReadCloser, data any) error {
	defer body.Close()
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return err
	}
	return nil
}

func (r *Renderer) EncodeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
