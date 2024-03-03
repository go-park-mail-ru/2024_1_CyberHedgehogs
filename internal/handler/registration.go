package handler

import (
	"encoding/json"
	rep "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"io"
	"net/http"
)

type UserHandler struct {
	users *rep.UserTable
}

func NewUserHandler(table *rep.UserTable) *UserHandler {
	return &UserHandler{table}
}

func (api *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "Error reading request body"}`, http.StatusBadRequest)
		return
	}

	var user rep.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	err = rep.ValidateNewUser(user, api.users)
	if err != nil {
		mes := rep.Info{
			Message: err.Error(),
		}
		jsonResponse, err2 := json.Marshal(mes)
		if err2 != nil {
			http.Error(w, `{"error": "Error marshaling response"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	err = api.users.AddUser(user)
	if err != nil {
		http.Error(w, `{"error": "Error adding user"}`, http.StatusInternalServerError)
		return
	}
	mes := rep.Info{
		Message: "success",
	}

	jsonResponse, err2 := json.Marshal(mes)
	if err2 != nil {
		http.Error(w, `{"error": "Error marshaling response"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonResponse)
}
