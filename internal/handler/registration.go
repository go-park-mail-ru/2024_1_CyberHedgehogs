package handler

import (
	"encoding/json"
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"net/http"
)

type UserHandler struct {
	users *repo.UserTable
}

type RegistrationManager interface {
	AddUser(user repo.User) error
	ValidateNewUser(user repo.User) error
}

type RegistrationHandler struct {
	table RegistrationManager
}

func NewRegistrationHandler(table *repo.UserTable) RegistrationHandler {
	handler := RegistrationHandler{table: table}
	return handler
}

func (api *RegistrationHandler) Registration(w http.ResponseWriter, r *http.Request) {

	var user repo.User
	render := repo.Renderer{}
	err := render.DecodeJSON(r.Body, &user)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	err = api.table.ValidateNewUser(user)
	if err != nil {
		mes := repo.Info{
			Message: err.Error(),
		}
		jsonResponse, _ := json.Marshal(mes)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	err = api.table.AddUser(user)
	if err != nil {
		http.Error(w, `{"error": "Error adding user"}`, http.StatusInternalServerError)
		return
	}
	mes := repo.Info{
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
