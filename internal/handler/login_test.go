package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/handler"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *mux.Router {
	userTable := &repository.UserTable{Users: make(map[string]*repository.User)}
	sessionTable := &repository.SessionTable{Sessions: make(map[string]*repository.Session)}

	testUser := &repository.User{
		Login:    "testUser",
		Username: "Test User",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userTable.Users[testUser.Login] = testUser

	authHandler := handler.NewAuthHandler(userTable, sessionTable)

	r := mux.NewRouter()
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	return r
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
	}{
		{
			name: "Successful login",
			payload: map[string]string{
				"login":    "testUser",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid login credentials",
			payload: map[string]string{
				"login":    "testUser",
				"password": "wrongPassword",
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	r := setupRouter()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payloadBytes))
			req.Header.Add("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}
