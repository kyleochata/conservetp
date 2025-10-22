package services

import (
	"fmt"

	"github.com/kyleochata/conservetp/users-backend/src/data"
)

type UsersService struct {
	usersData     *data.UsersData
	addressesData *data.AddressesData
}

func NewUsersService(usersData *data.UsersData, addressesData *data.AddressesData) *UsersService {
	return &UsersService{usersData: usersData, addressesData: addressesData}
}

func (us *UsersService) GetAllUsers() ([]data.User, error) {
	return us.usersData.GetAllUsers()
}

func (us *UsersService) CreateUser(name, email, pwd string) (*data.User, error) {
	if name == "" || email == "" || pwd == "" {
		return nil, fmt.Errorf("name, email, password are required")
	}

	//TODO: hash pwd

	return us.usersData.CreateUser(name, email, pwd)
}

func (us *UsersService) DeleteUser(id string) error {
	if id == "" {
		return fmt.Errorf("Must supply an id to delete a user")
	}
	return us.usersData.DeleteUser(id)
}
