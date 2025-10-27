package services

import (
	"fmt"

	"github.com/kyleochata/conservetp/users-backend/src/data"
	"github.com/kyleochata/conservetp/users-backend/src/types"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	usersData     *data.UsersData
	addressesData *data.AddressesData
}

func NewUsersService(usersData *data.UsersData, addressesData *data.AddressesData) *UsersService {
	return &UsersService{usersData: usersData, addressesData: addressesData}
}

func (us *UsersService) GetAllUsers() ([]types.User, error) {
	return us.usersData.GetAllUsers()
}

func (us *UsersService) GetSingleUserById(id string) (*types.UserResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("Error getting user by id: empty user - (s)")
	}
	fmt.Println("Start getuserbyid service")
	fmt.Println(id)
	getUserByIdResponse, err := us.usersData.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return getUserByIdResponse, nil
}

// change createUser to either call createUserNoAddress vs createUserwAddress
func (us *UsersService) CreateUserWAddress(user types.CreateUserRequest) (types.UserResponse, error) {
	// create the user entry first to get userid for address foreign key
	userCreationResponse, err := us.CreateUserWOutAddress(user)
	if err != nil {
		return types.UserResponse{}, err
	}
	//use userid to create new address
	addrData := user.Address
	addr, err := us.addressesData.CreateNewAddress(addrData, userCreationResponse.User.ID)
	if err != nil {
		return types.UserResponse{}, err
	}
	fmt.Println(addr)
	userCreationResponse.Address = addr
	return *userCreationResponse, nil
}

func (us *UsersService) CreateUserWOutAddress(user types.CreateUserRequest) (*types.UserResponse, error) {
	if user.Pwd == "" {
		return nil, fmt.Errorf("Error creating user w/ address (s): pwd missing")
	}

	userData := types.CreateUserData{
		Name:  user.Name,
		Email: user.Email,
		Pwd:   user.Pwd,
	}
	if err := us.hashPassword(&userData); err != nil {
		return nil, err
	}

	newUserResponse, err := us.usersData.CreateUser(userData)
	if err != nil {
		return nil, err
	}
	return newUserResponse, nil
}

func (us *UsersService) DeleteUserById(id string) error {
	if id == "" {
		return fmt.Errorf("Must supply an id to delete a user")
	}
	return us.usersData.DeleteUser(id)
}

func (us *UsersService) UpdateUserInfoWAddr(id string, user types.CreateUserRequest) (types.UserResponse, error) {
	if id == "" {
		return types.UserResponse{}, fmt.Errorf("Error Updating User infor: empty id")
	}
	updateUserResponse, err := us.UpdateUserInfo(id, user)
	if err != nil {
		return types.UserResponse{}, err
	}
	oldAddrId := updateUserResponse.Address.Address.ID
	updateAddr := user.Address

	updateAddrResponse, err := us.updateAddress(oldAddrId, updateAddr, updateUserResponse.User.ID)
	if err != nil {
		return types.UserResponse{}, err
	}
	updateUserResponse.Address = updateAddrResponse
	return updateUserResponse, nil
}

func (us UsersService) updateAddress(addrId string, newAddr *types.CreateAddressRequest, userId string) (*types.AddressResponse, error) {
	//check addr and id exists
	if userId == "" {
		return nil, fmt.Errorf("Error updating address. Empty userId")
	}
	if addrId == "" {
		return nil, fmt.Errorf("Error updating address. Empty addrId")
	}
	if newAddr == nil {
		return nil, fmt.Errorf("Error updating address. Empty userId")
	}

	_, err := us.usersData.GetUserById(userId) // change it to have a addressincluded flag that'll signal including addresses with the query. default = false
	if err != nil {
		return nil, err
	}

	prevAddr, err := us.addressesData.GetAddressById(addrId)
	if err != nil {
		return nil, err
	}

	newAddr = us.parseAddressChanges(prevAddr, newAddr)

	//return call data layer
	updateAddrResponse, err := us.addressesData.UpdateAddress(addrId, userId, newAddr)
	if err != nil {
		return nil, err
	}

	return updateAddrResponse, nil
}

func (us *UsersService) UpdateUserInfo(id string, user types.CreateUserRequest) (types.UserResponse, error) {
	if id == "" {
		return types.UserResponse{}, fmt.Errorf("Error Updating User infor: empty id")
	}

	rUserData := types.CreateUserData{
		Name:  user.Name,
		Email: user.Email,
		Pwd:   user.Pwd,
	}

	if err := us.hashPassword(&rUserData); err != nil {
		return types.UserResponse{}, err
	}

	updateUserResponse, err := us.usersData.UpdateUserInfo(id, rUserData)
	if err != nil {
		return types.UserResponse{}, err
	}

	return *updateUserResponse, nil
}

func (us UsersService) hashPassword(userData *types.CreateUserData) error {
	if userData.Pwd == "" {
		return fmt.Errorf("Password for user empty")
	}
	hPwd, err := bcrypt.GenerateFromPassword([]byte(userData.Pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userData.Pwd = string(hPwd)
	return nil
}

func (us UsersService) parseAddressChanges(prev *types.AddressResponse, new *types.CreateAddressRequest) *types.CreateAddressRequest {
	if prev == nil || new == nil {
		return new
	}
	updated := &types.CreateAddressRequest{}
	updated.Street = us.updateStringField(prev.Address.Street, new.Street)
	updated.AptNum = us.updateStringField(*prev.Address.AptNum, new.AptNum)
	updated.Zipcode = us.updateStringField(prev.Address.Zipcode, new.Zipcode)
	updated.City = us.updateStringField(prev.Address.City, new.City)
	updated.State = us.updateStringField(prev.Address.State, new.State)
	updated.Country = us.updateStringField(prev.Address.Country, new.Country)
	updated.IsPrimary = new.IsPrimary

	return updated
}

func (us UsersService) updateStringField(old, new string) string {
	if new != "" && old != new {
		return new
	}
	return old
}
