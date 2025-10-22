package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kyleochata/conservetp/users-backend/src/services"
)

type UsersHandler struct {
	usersService *services.UsersService
}

func NewUsersHandler(usersService *services.UsersService) *UsersHandler {
	return &UsersHandler{usersService: usersService}
}

func (uh *UsersHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uh.getUsers(w, r)
	case http.MethodPost:
		uh.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (uh *UsersHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.usersService.GetAllUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get users: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (uh *UsersHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Pwd   string `json:"pwd"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	user, err := uh.usersService.CreateUser(request.Name, request.Email, request.Pwd)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
