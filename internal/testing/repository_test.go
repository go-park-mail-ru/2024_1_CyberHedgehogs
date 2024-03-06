package repository_test

import (
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddUserExtended(t *testing.T) {
	tests := []struct {
		name       string
		setupUsers func(table *repo.UserTable)
		user       *repo.User
		want       *repo.User
		wantErr    bool
		errMessage string
	}{
		{
			name: "duplicate login",
			setupUsers: func(table *repo.UserTable) {
				existingUser := &repo.User{
					Login:    "existingUser",
					Username: "Existing User",
					Email:    "existing@example.com",
					Password: "password123",
				}
				table.Users[existingUser.Login] = existingUser
			},
			user: &repo.User{
				Login:    "existingUser",
				Username: "New User",
				Email:    "new@example.com",
				Password: "password123",
			},
			wantErr:    true,
			errMessage: "пользователь с таким логином уже существует",
		},
		{
			name: "invalid email format",
			user: &repo.User{
				Login:    "newUser",
				Username: "New User",
				Email:    "invalidemail",
				Password: "password123",
			},
			wantErr:    true,
			errMessage: "неверный формат email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userTable := &repo.UserTable{Users: make(map[string]*repo.User)}
			if tt.setupUsers != nil {
				tt.setupUsers(userTable)
			}

			got, err := userTable.AddUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
