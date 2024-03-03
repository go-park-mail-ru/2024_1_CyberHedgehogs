package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddUser(t *testing.T) {
	userTable := UserTable{Users: make(map[string]*User)}
	user := User{
		Login:    "testUser",
		Username: "Test Username",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}

	err := userTable.AddUser(user)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(userTable.Users))
	assert.Equal(t, "testUser", userTable.Users["testUser"].Login)
}
func TestValidateNewUser(t *testing.T) {
	userTable := UserTable{Users: make(map[string]*User)}
	validUser := User{
		Login:    "validUser",
		Username: "Valid Username",
		Email:    "valid@example.com",
		Password: "validPass123",
		Role:     "user",
	}

	invalidUser := User{
		Login:    "iu",
		Username: "IU",
		Email:    "invalidEmail",
		Password: "123",
		Role:     "invalidRole",
	}

	err := ValidateNewUser(validUser, &userTable)
	assert.Nil(t, err)

	err = ValidateNewUser(invalidUser, &userTable)
	assert.NotNil(t, err)
}
func TestSessionManagement(t *testing.T) {
	sessionTable := SessionTable{Sessions: make(map[string]Session), Users: make(map[string]*User)}

	user := User{
		Login: "testSessionUser",
	}
	sessionTable.Users[user.Login] = &user

	sessionID, err := sessionTable.AddSession(user)
	assert.Nil(t, err)
	assert.NotEmpty(t, sessionID)

	checkedUser, err := sessionTable.CheckSession(sessionID)
	assert.Nil(t, err)
	assert.Equal(t, user.Login, checkedUser.Login)

	err = sessionTable.DeleteSession(sessionID)
	assert.Nil(t, err)

	_, err = sessionTable.CheckSession(sessionID)
	assert.NotNil(t, err)
}

func TestSessionExpiration(t *testing.T) {
	sessionTable := SessionTable{Sessions: make(map[string]Session), Users: make(map[string]*User)}

	user := User{
		Login: "expiringUser",
	}
	sessionTable.Users[user.Login] = &user

	sessionID, _ := sessionTable.AddSession(user)
	sessionTable.Sessions[sessionID] = Session{
		ExpirationDate: time.Now().Add(-1 * time.Minute),
		User:           &user,
		Id:             sessionID,
	}

	_, err := sessionTable.CheckSession(sessionID)
	assert.NotNil(t, err)
	assert.Equal(t, "session expired", err.Error())
}
