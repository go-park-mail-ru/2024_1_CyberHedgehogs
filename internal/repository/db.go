package repository

import (
	"encoding/json"
	"errors"
	"github.com/oklog/ulid/v2"
	"io"
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
	mu    sync.RWMutex
}

type Info struct {
	Message string `json:"message"`
}

func (table *UserTable) ValidateUserCredentials(user *User) *User {
	table.mu.Lock()
	defer table.mu.Unlock()
	if us, ok := table.Users[user.Login]; ok && us.Password == user.Password { // todo  хэши паролей
		return us
	}
	return nil
}

func (table *UserTable) AddUser(user *User) (*User, error) {
	err := table.ValidateNewUser(user)
	if err != nil {
		return user, err
	}
	user.ID = nextUserID
	nextUserID++
	table.mu.Lock()
	table.Users[user.Login] = user
	table.mu.Unlock()
	return user, nil
}

func (table *UserTable) ValidateNewUser(user *User) error {
	if !emailRegex.MatchString(user.Email) {
		return errors.New("неверный формат email")
	}

	if len(user.Login) < 3 {
		return errors.New("login должен быть не менее 3 символов")
	}
	table.mu.Lock()
	_, exists := table.Users[user.Login]
	table.mu.Unlock()
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
	UserID uint   `json:"user_id,omitempty"`
	Login  string `json:"login"`
}

func (userInfo *UserSessionInfo) GetUserIDLogin(user *User) {
	userInfo.UserID = user.ID
	userInfo.Login = user.Login
}

type Session struct {
	ExpirationDate time.Time
	UserInfo       *UserSessionInfo `json:"user_info"`
	ID             string           `json:"session_id"`
}

type SessionTable struct {
	Sessions map[string]*Session
	mu       sync.RWMutex
}

func (table *SessionTable) CheckSession(sessionID string) (user *UserSessionInfo) {
	nowTime := time.Now()
	table.mu.Lock()
	defer table.mu.Unlock()

	if session, ok := table.Sessions[sessionID]; ok {
		if session.ExpirationDate.Before(nowTime) {
			delete(table.Sessions, sessionID)
			return nil
		}
		return session.UserInfo
	}
	return nil
}

func (table *SessionTable) SessionsCleanup() {
	table.mu.Lock()
	for sessionID, sessn := range table.Sessions {
		if sessn.ExpirationDate.Before(time.Now()) {
			delete(table.Sessions, sessionID)
		}
	}
	table.mu.Unlock()
}

func (table *SessionTable) DeleteSession(sessionID string) error {
	table.mu.Lock()
	defer table.mu.Unlock()
	if _, ok := table.Sessions[sessionID]; ok {
		delete(table.Sessions, sessionID)
		return nil
	}
	return nil
}

func (table *SessionTable) AddSession(user *UserSessionInfo) (sessionID string, err error) {
	sessionID = ulid.Make().String()
	ses := Session{ExpirationDate: time.Now().Add(SessionLiveTime), UserInfo: user, ID: sessionID}
	table.mu.Lock()
	table.Sessions[sessionID] = &ses
	table.mu.Unlock()
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
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
