package data

import (
	"database/sql"
	"fmt"
)

type AddressesData struct {
	db *sql.DB
}

func NewAddressesData(db *sql.DB) *AddressesData {
	return &AddressesData{db: db}
}

type Address struct {
	ID        string `json:"id"`
	Street    string `json:"street"`
	AptNumber string `json:"apt_num"`
	Zipcode   string `json:"zipccode"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	IsPrimary bool   `json:"is_primary"`
}

func (ad AddressesData) GetAllAddresses() ([]Address, error) {
	rows, err := ad.db.Query("SELECT id, user_id, street, apt_num, zipcode, city, state, country, is_primary FROM addresses")
	if err != nil {
		return nil, fmt.Errorf("Failed to get all addresses: %w", err)
	}
	defer rows.Close()
	var addresses []Address
	for rows.Next() {
		var address Address
		if err := rows.Scan(
			&address.ID, &address.Street, &address.AptNumber, &address.Zipcode,
			&address.City, &address.State, &address.Country, &address.IsPrimary,
		); err != nil {
			return nil, fmt.Errorf("Failed to get all addresses: scan failure: %w", err)
		}
		addresses = append(addresses, address)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating addresses: %w", err)
	}
	return addresses, nil
}

func (ad *AddressesData) CreateNewAddress(userId string, address Address) (Address, error) {
	if userId == "" {
		return Address{}, fmt.Errorf("Failed to Create New Address: empty userId")
	}
	if userId != address.ID {
		address.ID = userId
	}
	err := ad.db.QueryRow(
		"INSERT INTO addresses (user_id, street, apt_num, zipcode, city, state, country, is_primary) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id, street, apt_num, zipcode, city, state, country, is_primary",
		address.Street, userId, address.AptNumber, address.Zipcode, address.City, address.State, address.Country, address.IsPrimary,
	).Scan(
		&address.ID, &address.Street, &address.AptNumber, &address.Zipcode,
		&address.City, &address.State, &address.Country, &address.IsPrimary,
	)
	if err != nil {
		return Address{}, fmt.Errorf("Failed to Insert New address: %w", err)
	}
	return address, nil
}

func (ad AddressesData) GetNumberOfAddressBy(filterTag, filterValue string) (int, error) {
	if filterTag == "" || filterValue == "" {
		return -1, fmt.Errorf("Error getting number of addreses by: empty filter string(s)")
	}
	query := fmt.Sprintf(
		"SELECT id, user_id, street, apt_num, zipcode, city, state, country, is_primary FROM addresses WHERE %s = %s",
		filterTag, filterValue,
	)
	rows, err := ad.db.Query(query)
	if err != nil {
		return -1, fmt.Errorf("Failed to get number of addresses by %s=%s: %w", filterTag, filterValue, err)
	}
	var addresses []Address
	for rows.Next() {
		var address Address
		if err := rows.Scan(
			&address.ID, &address.Street, &address.AptNumber, &address.Zipcode,
			&address.City, &address.State, &address.Country, &address.IsPrimary,
		); err != nil {
			return -1, fmt.Errorf("Failed to scan addresses: %w", err)
		}
		addresses = append(addresses, address)
	}
	return len(addresses), nil
}
