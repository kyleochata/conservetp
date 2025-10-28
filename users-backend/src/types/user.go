package types

import (
	"time"
)

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Pwd       string     `json:"pwd,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	IsActive  bool       `json:"is_active"`
}

type CreateUserRequest struct {
	Name    string                `json:"name"`
	Email   string                `json:"email"`
	Pwd     string                `json:"pwd"`
	Address *CreateAddressRequest `json:"address,omitempty"`
}

func (cu CreateUserRequest) GetPwd() string {
	return cu.Pwd
}
func (cu *CreateUserRequest) SetPwd(pwd string) {
	if pwd == "" {
		return
	}
	cu.Pwd = pwd
	return
}

type UpdateUserRequest struct {
	Name    string                `json: "name,omitempty"`
	Email   string                `json:"email,omitempty"`
	Pwd     string                `json: "pwd,omitempty"`
	Address *UpdateAddressRequest `json:"address,omitempty"`
}

func (uu UpdateUserRequest) GetPwd() string {
	return uu.Pwd
}
func (uu *UpdateUserRequest) SetPwd(pwd string) {
	if pwd == "" {
		return
	}
	uu.Pwd = pwd
	return
}

type CreateUserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

func (cud CreateUserData) GetPwd() string {
	return cud.Pwd
}
func (cud *CreateUserData) SetPwd(pwd string) {
	if pwd == "" {
		return
	}
	cud.Pwd = pwd
	return
}

type UserResponse struct {
	User    *User            `json:"user"`
	Address *AddressResponse `json:"address,omitempty"`
}
