package repository

import (
	"errors"
	"math/rand"
	"regexp"
	"sync"
	"time"
)

var nextUserID uint = 1

type User struct {
	ID       uint   `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

type UserTable struct {
	Users map[string]*User
	mu    sync.Mutex
}

type Info struct {
	Message string `json:"message"`
}

func (table *UserTable) AddUser(user User) error {
	table.mu.Lock()
	defer table.mu.Unlock()

	err := ValidateNewUser(user, table)
	if err != nil {
		return err
	}
	user.ID = nextUserID
	nextUserID++
	table.Users[user.Login] = &user
	return nil
}

func ValidateNewUser(user User, table *UserTable) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
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

	allowedRoles := []string{"user", "creator"}
	roleIsValid := false
	for _, r := range allowedRoles {
		if user.Role == r {
			roleIsValid = true
			break
		}
	}
	if !roleIsValid {
		return errors.New("недопустимая роль")
	}

	return nil
}

type Session struct {
	ExpirationDate time.Time
	User           *User  `json:"user"`
	Id             string `json:"session_id"`
}
type SessionTable struct {
	Sessions map[string]Session
	Users    map[string]*User
	mu       sync.Mutex
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

func (table *SessionTable) sessionsCleanup() {
	table.mu.Lock()
	defer table.mu.Unlock()
	for sessionID, session := range table.Sessions {
		if session.ExpirationDate.Before(time.Now()) {
			delete(table.Sessions, sessionID)
		}
	}
	return
}

func (table *SessionTable) DeleteSession(sessionID string) error {
	table.mu.Lock()
	defer table.mu.Unlock()
	if _, ok := table.Sessions[sessionID]; ok {
		delete(table.Sessions, sessionID)
		return nil
	}
	return errors.New("no such session")
}

// todo при уадлении пользователя надо удалить и его сессию!
func (table *SessionTable) CheckSession(sessionID string) (user *User, err error) {
	table.mu.Lock()
	defer table.mu.Unlock()
	if session, ok := table.Sessions[sessionID]; ok {
		if session.ExpirationDate.Before(time.Now()) {
			delete(table.Sessions, sessionID)
			return nil, errors.New("session expired")
		}
		return session.User, nil
	}
	return nil, errors.New("session not found")
}

// todo вызов удаления сессии
func (table *SessionTable) AddSession(user User) (sessionID string, err error) {
	table.mu.Lock()
	defer table.mu.Unlock()

	if tableUser, ok := table.Users[user.Login]; ok {
		sessionID = RandStringRunes(28)
		// todo проверка что мне "повезло" невероятным образом два ID выбить одинаковых
		ses := Session{User: tableUser, ExpirationDate: time.Now().Add(10 * time.Minute), Id: sessionID} //todo длиннее сессию
		table.Sessions[sessionID] = ses
		return sessionID, nil
	}
	return "", errors.New("user not found")
}

type SessionsController struct {
	users    *UserTable
	sessions *SessionTable
	mu       sync.Mutex
}
