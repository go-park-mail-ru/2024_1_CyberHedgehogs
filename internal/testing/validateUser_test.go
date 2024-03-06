package repository_test

import (
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateNewUser(t *testing.T) {
	tests := []struct {
		name       string
		setupUsers func(table *repo.UserTable)
		user       *repo.User
		wantErr    bool
		errMessage string
	}{
		{
			name: "invalid email format",
			user: &repo.User{
				Email: "bademail",
			},
			wantErr:    true,
			errMessage: "неверный формат email",
		},
		{
			name: "login too short",
			user: &repo.User{
				Login: "ab",
				Email: "valid@example.com",
			},
			wantErr:    true,
			errMessage: "login должен быть не менее 3 символов",
		},
		{
			name: "username too short",
			user: &repo.User{
				Login:    "validLogin",
				Username: "ab",
				Email:    "valid@example.com",
			},
			wantErr:    true,
			errMessage: "username должен быть не менее 3 символов",
		},
		{
			name: "password too short",
			user: &repo.User{
				Login:    "validLogin",
				Username: "validUsername",
				Email:    "valid@example.com",
				Password: "123",
			},
			wantErr:    true,
			errMessage: "пароль должен быть не менее 6 символов",
		},
		{
			name: "duplicate login",
			setupUsers: func(table *repo.UserTable) {
				existingUser := &repo.User{
					Login:    "existingLogin",
					Username: "Existing User",
					Email:    "existing@example.com",
					Password: "password123",
				}
				table.Users[existingUser.Login] = existingUser
			},
			user: &repo.User{
				Login:    "existingLogin",
				Username: "New User",
				Email:    "new@example.com",
				Password: "newpassword",
			},
			wantErr:    true,
			errMessage: "пользователь с таким логином уже существует",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userTable := &repo.UserTable{Users: make(map[string]*repo.User)}
			if tt.setupUsers != nil {
				tt.setupUsers(userTable)
			}

			err := userTable.ValidateNewUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
