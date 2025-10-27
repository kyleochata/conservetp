package data

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kyleochata/conservetp/users-backend/src/types"
)

type AddressesData struct {
	db *sql.DB
}

func NewAddressesData(db *sql.DB) *AddressesData {
	return &AddressesData{db: db}
}

func (ad AddressesData) GetAllAddresses() ([]types.Address, error) {
	rows, err := ad.db.Query("SELECT id, user_id, street, apt_num, zipcode, city, state, country, is_primary FROM addresses")
	if err != nil {
		return nil, fmt.Errorf("Failed to get all addresses: %w", err)
	}
	defer rows.Close()
	var addresses []types.Address
	for rows.Next() {
		var address types.Address
		if err := rows.Scan(
			&address.ID, &address.UserID, &address.Street, &address.AptNum, &address.Zipcode,
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

func (ad *AddressesData) CreateNewAddress(address *types.CreateAddressRequest, userId string) (*types.AddressResponse, error) {
	if userId == "" {
		return nil, fmt.Errorf("Failed to Create New Address: empty userId")
	}
	if address == nil {
		return nil, fmt.Errorf("Failed to create new address: empty address")
	}

	var newAddr types.Address

	err := ad.db.QueryRow(
		`INSERT INTO addresses (user_id, street, apt_num, zipcode, city, state, country, is_primary) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, street, apt_num, zipcode, city, state, country, is_primary`,
		userId, address.Street, address.AptNum, address.Zipcode, address.City, address.State, address.Country, address.IsPrimary,
	).Scan(
		&newAddr.ID, &newAddr.Street, &newAddr.AptNum, &newAddr.Zipcode,
		&newAddr.City, &newAddr.State, &newAddr.Country, &newAddr.IsPrimary,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to Insert New address: %w", err)
	}
	return &types.AddressResponse{
		Address: &newAddr,
	}, nil
}

// func (ad AddressesData) GetNumberOfAddressBy(filterTag, filterValue string) (int, error) {
// 	if filterTag == "" || filterValue == "" {
// 		return -1, fmt.Errorf("Error getting number of addreses by: empty filter string(s)")
// 	}
// 	query := fmt.Sprintf(
// 		"SELECT id, user_id, street, apt_num, zipcode, city, state, country, is_primary FROM addresses WHERE %s = %s",
// 		filterTag, filterValue,
// 	)
// 	rows, err := ad.db.Query(query)
// 	if err != nil {
// 		return -1, fmt.Errorf("Failed to get number of addresses by %s=%s: %w", filterTag, filterValue, err)
// 	}
// 	var addresses []Address
// 	for rows.Next() {
// 		var address Address
// 		if err := rows.Scan(
// 			&address.ID, &address.Street, &address.AptNumber, &address.Zipcode,
// 			&address.City, &address.State, &address.Country, &address.IsPrimary,
// 		); err != nil {
// 			return -1, fmt.Errorf("Failed to scan addresses: %w", err)
// 		}
// 		addresses = append(addresses, address)
// 	}
// 	return len(addresses), nil
// }

func (ad *AddressesData) GetAddressById(id string) (*types.AddressResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("Error getting address by id: empty id data")
	}
	var res types.Address
	err := ad.db.QueryRow(
		"SELECT id, user_id, street, apt_num, zipcode, city, state, country, is_primary FROM addresses WHERE id = $1",
		id,
	).Scan(
		&res.ID, &res.UserID, &res.Street, &res.AptNum,
		&res.Zipcode, &res.City, &res.State, &res.Country, &res.IsPrimary,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Address not found")
		}
		return nil, fmt.Errorf("Failed to get addres by id: %s: %w", id, err)
	}
	return &types.AddressResponse{Address: &res}, nil

}

func (ad *AddressesData) UpdateAddress(addrId, userId string, newAddr *types.CreateAddressRequest) (*types.AddressResponse, error) {
	if addrId == "" || userId == "" {
		return nil, fmt.Errorf("Error updating address: empty id (addr & user)")
	}
	if newAddr == nil {
		return nil, fmt.Errorf("Error updating address: empty newAddr change")
	}

	exists, err := ad.DoesAddrExistGeneric(StrFilter{Field: "id", Value: addrId}, StrFilter{Field: "user_id", Value: userId})
	if err != nil {
		return nil, fmt.Errorf("Error checking if addr exits: %w", err)
	}
	if !exists {
		return ad.CreateNewAddress(newAddr, userId)
	}

	var updatedAddr types.Address
	err = ad.db.QueryRow(
		`UPDATE addresses 
		 SET street = $1, apt_num = $2, zipcode = $3, city = $4, state = $5, country = $6, is_primary = $7
		 WHERE id = $8 AND user_id = $9
		 RETURNING id, user_id, street, apt_num, zipcode, city, state, country, is_primary`,
		newAddr.Street, newAddr.AptNum, newAddr.Zipcode, newAddr.City, newAddr.State, newAddr.Country, newAddr.IsPrimary,
		addrId, userId,
	).Scan(
		&updatedAddr.ID, &updatedAddr.UserID, &updatedAddr.Street, &updatedAddr.AptNum,
		&updatedAddr.Zipcode, &updatedAddr.City, &updatedAddr.State, &updatedAddr.Country, &updatedAddr.IsPrimary,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed at updateAddrRes: %w", err)
	}

	return &types.AddressResponse{
		Address: &updatedAddr,
	}, nil
}

type Filter interface {
	GetField() string
	GetValue() interface{}
}
type StrFilter struct {
	Field string
	Value string
}

func (s StrFilter) GetField() string      { return s.Field }
func (s StrFilter) GetValue() interface{} { return s.Value }

type BoolFilter struct {
	Field string
	Value bool
}

func (b BoolFilter) GetField() string      { return b.Field }
func (b BoolFilter) GetValue() interface{} { return b.Value }

func (ad AddressesData) DoesAddrExistGeneric(filters ...Filter) (bool, error) {
	if len(filters) == 0 {
		return false, fmt.Errorf("Error checking addr exists: empty filters map")
	}

	fValues := make([]interface{}, 0, len(filters))
	var qBuilder strings.Builder
	qBuilder.WriteString("SELECT EXISTS(SELECT 1 FROM addresses WHERE")

	for i, filter := range filters {
		s := fmt.Sprintf(" %s = $%d AND", filter.GetField(), i+1)
		qBuilder.WriteString(s)
		fValues = append(fValues, filter.GetValue())
	}

	query := strings.TrimSuffix(qBuilder.String(), "AND")
	query += ")"

	var exists bool
	err := ad.db.QueryRow(query, fValues...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("No address found with: %v", filters)
		}
		return false, fmt.Errorf("Error querying row: \nqStr: %s\nerr: %w", query, err)
	}

	return exists, nil
}
