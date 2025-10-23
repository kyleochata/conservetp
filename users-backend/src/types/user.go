package types

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Pwd       string    `json:"pwd,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login,omitempty"`
	IsActive  bool      `json:"is_active"`
}

type CreateUserRequest struct {
	Name    string                `json:"name"`
	Email   string                `json:"email"`
	Pwd     string                `json:"pwd"`
	Address *CreateAddressRequest `json:"address,omitempty"`
}

type CreateUserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

type UserResponse struct {
	User    *User            `json:"user"`
	Address *AddressResponse `json:"address,omitempty"`
}
