package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kyleochata/conservetp/users-backend/src/services"
	"github.com/kyleochata/conservetp/users-backend/src/types"
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
		uh.getAllUsers(w, r)
	case http.MethodPost:
		uh.createSingleUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func (uh *UsersHandler) HandleUserId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Printf("id = %s", id)
	if id == "" {
		http.Error(w, "User id path extraction failed", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		uh.getUserById(w, r)
	case http.MethodPut:
		uh.updateUserById(w, r, id)
	case http.MethodDelete:
		uh.deleteUserById(w, r, id)
	default:
		http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
	}
}

// Service calls

func (uh UsersHandler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.usersService.GetAllUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed  to get users: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (uh UsersHandler) createSingleUser(w http.ResponseWriter, r *http.Request) {
	var req types.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Email == "" || req.Pwd == "" {
		http.Error(w, "invalid request body: Missing user fields (name, email, password)", http.StatusBadRequest)
		return
	}

	if req.Address != nil {
		user, err := uh.usersService.CreateUserWAddress(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)
	} else {
		user, err := uh.usersService.CreateUserWOutAddress(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(*user)
	}
}

// Single User calls

func (uh UsersHandler) getUserById(w http.ResponseWriter, r *http.Request) {
	uId := r.URL.Path
	fmt.Printf("uId: %s\n", uId)
	user, err := uh.usersService.GetSingleUserById(uId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(*user); err != nil {
		http.Error(w, "Failed to encod response", http.StatusInternalServerError)
		return
	}
}

func (uh UsersHandler) updateUserById(w http.ResponseWriter, r *http.Request, userId string) {
	var req types.CreateUserRequest
	//if there's an address provided
	if req.Address != nil {
		user, err := uh.usersService.UpdateUserInfoWAddr(userId, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(user)
	} else {
		// no address provided
		user, err := uh.usersService.UpdateUserInfo(userId, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
	return
}

func (uh UsersHandler) deleteUserById(w http.ResponseWriter, r *http.Request, userId string) {
	err := uh.usersService.DeleteUserById(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "applicaton/json")
	json.NewEncoder(w).Encode(fmt.Sprintf("User with deleted successfully.\nID: %s", userId))
}
