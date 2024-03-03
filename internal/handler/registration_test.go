package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	rep "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestRegistration(t *testing.T) {
	userTable := rep.UserTable{Users: make(map[string]*rep.User)}
	api := NewUserHandler(&userTable)

	user := rep.User{
		Login:    "testLogin",
		Username: "testUsername",
		Email:    "test@example.com",
		Password: "testPass",
		Role:     "user",
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	api.Registration(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response rep.Info
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "success", response.Message)
}
