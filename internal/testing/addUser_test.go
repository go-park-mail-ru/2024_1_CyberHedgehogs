package repository_test

import (
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddUser(t *testing.T) {
	tests := []struct {
		name          string
		userToAdd     *repo.User
		expectedID    uint
		expectError   bool
		expectedError string
		setupFunc     func(*repo.UserTable)
	}{
		{
			name: "Add valid user",
			userToAdd: &repo.User{
				Login:    "newUser",
				Username: "New User",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			expectedID:  1,
			expectError: false,
		},
		{
			name: "Add user with existing login",
			userToAdd: &repo.User{
				Login:    "existingUser",
				Username: "Another User",
				Email:    "another@example.com",
				Password: "password123",
			},
			expectError:   true,
			expectedError: "пользователь с таким логином уже существует",
			setupFunc: func(table *repo.UserTable) {
				existingUser := &repo.User{
					Login:    "existingUser",
					Username: "Existing User",
					Email:    "existing@example.com",
					Password: "password123",
				}
				table.Users[existingUser.Login] = existingUser
			},
		},
		{
			name: "Add user with invalid email",
			userToAdd: &repo.User{
				Login:    "uniqueUser",
				Username: "Unique User",
				Email:    "notanemail",
				Password: "password123",
			},
			expectError:   true,
			expectedError: "неверный формат email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userTable := &repo.UserTable{Users: make(map[string]*repo.User)}

			if tt.setupFunc != nil {
				tt.setupFunc(userTable)
			}

			addedUser, err := userTable.AddUser(tt.userToAdd)

			if tt.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, addedUser.ID)
				assert.Equal(t, tt.userToAdd.Login, addedUser.Login)
				assert.Equal(t, tt.userToAdd.Email, addedUser.Email)
				assert.NotNil(t, userTable.Users[tt.userToAdd.Login])
			}
		})
	}
}
