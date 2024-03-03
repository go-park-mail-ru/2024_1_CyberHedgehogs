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

func TestLogin(t *testing.T) {
	userTable := &rep.UserTable{Users: make(map[string]*rep.User)}
	sessionHandler := NewSessionHandler(userTable)

	user := rep.User{Login: "testLogin", Password: "testPass"}
	userTable.Users[user.Login] = &user

	loginData := map[string]string{"login": "testLogin", "password": "testPass"}
	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	sessionHandler.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	cookie := w.Result().Cookies()
	assert.NotEmpty(t, cookie, "Cookie 'session_id' should not be empty")
}
