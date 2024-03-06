package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/handler"
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistration(t *testing.T) {
	userTable := &repo.UserTable{Users: make(map[string]*repo.User)}
	sessionTable := &repo.SessionTable{Sessions: make(map[string]*repo.Session)}

	api := handler.NewAuthHandler(userTable, sessionTable)

	tests := []struct {
		name           string
		user           repo.User
		expectedStatus int
		expectedBody   string
		preTest        func()
	}{
		{
			name: "Successful registration",
			user: repo.User{
				Login:    "newUser",
				Username: "New User",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name: "Registration with existing user",
			user: repo.User{
				Login:    "existingUser",
				Username: "Existing User",
				Email:    "existing@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "Error adding user"}`,
			preTest: func() {
				userTable.Users["existingUser"] = &repo.User{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preTest != nil {
				tt.preTest()
			}

			userJSON, _ := json.Marshal(tt.user)
			request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
			response := httptest.NewRecorder()

			r := mux.NewRouter()
			r.HandleFunc("/register", api.Registration).Methods("POST")
			r.ServeHTTP(response, request)

			assert.Equal(t, tt.expectedStatus, response.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, response.Body.String(), tt.expectedBody)
			}

			if tt.preTest != nil {
				userTable.Users = make(map[string]*repo.User)
			}
		})
	}
}
