package types

import (
	"time"
)

type Address struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Street    string    `json:"street"`
	AptNum    *string   `json:"apt_num,omitempty"`
	Zipcode   string    `json:"zipcode"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Country   string    `json:"country"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAddressRequest struct {
	Street    string `json:"street"`
	AptNum    string `json:"apt_num,omitempty"`
	Zipcode   string `json:"zipcode"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country,omitempty"`
	IsPrimary bool   `json:"is_primary,omitempty"`
}

type AddressResponse struct {
	Address *Address `json:"address"`
}
