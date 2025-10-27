package data

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kyleochata/conservetp/users-backend/src/types"
)

type UsersData struct {
	db *sql.DB
}

func NewUsersData(db *sql.DB) *UsersData {
	return &UsersData{db: db}
}

func (ud *UsersData) GetAllUsers() ([]types.User, error) {
	rows, err := ud.db.Query("SELECT id, name, email, created_at, last_login FROM users WHERE is_active = true")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	users := []types.User{}
	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.LastLogin); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}

func (ud *UsersData) CreateUser(user types.CreateUserData) (*types.UserResponse, error) {
	var newUser types.User
	err := ud.db.QueryRow(
		"INSERT INTO users (email, name, pwd) VALUES ($1, $2, $3) RETURNING id, email, name, created_at",
		user.Email, user.Name, user.Pwd,
	).Scan(&newUser.ID, &newUser.Email, &newUser.Name, &newUser.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert new user into table: %w", err)
	}
	return &types.UserResponse{
		User: &newUser,
	}, nil
}

func (ud *UsersData) DeleteUser(id string) error {
	result, err := ud.db.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("fialed to remove user with id: %v: %w", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %v not found", id)
	}
	return nil
}

func (ud UsersData) GetUserById(id string) (*types.UserResponse, error) {
	fmt.Println("Start getUserbyid data")
	if id == "" {
		return nil, fmt.Errorf("Error getting userbyid: empty id")
	}
	fmt.Println(id)
	var user types.User
	err := ud.db.QueryRow(
		`SELECT id, name, email, created_at, updated_at, last_login, is_active 
         FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.IsActive,
	)

	fmt.Println(user.ID)
	fmt.Printf("time: %s", user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("Error getting userbyid: %s: %w", id, err)
	}
	return &types.UserResponse{
		User: &user,
	}, nil
}

// func (ud UsersData) GetUserByEmail(email string) (User, error) {
// 	if email == "" {
// 		return User{}, fmt.Errorf("Error getting userbyemail: empty email")
// 	}
// 	var user User
// 	err := ud.db.QueryRow(
// 		"SELECT id, name, email WHERE email = $1",
// 		email,
// 	).Scan(&user.ID, &user.Name, &user.Email)
// 	if err != nil {
// 		return User{}, fmt.Errorf("Error getting userbyemail: %s: %w", email, err)
// 	}
// 	return user, nil
// }

func (ud *UsersData) UpdateUserInfo(id string, updUser types.CreateUserData) (*types.UserResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("Error updating user: empty id")
	}
	var qBuilder strings.Builder
	params := []interface{}{}
	paramCount := 1

	qBuilder.WriteString("UPDATE users SET")

	if updUser.Name != "" {
		fmt.Fprintf(&qBuilder, " name = $%d,", paramCount)
		params = append(params, updUser.Name)
		paramCount++
	}
	if updUser.Email != "" {
		fmt.Fprintf(&qBuilder, " email = $%d,", paramCount)
		params = append(params, updUser.Email)
		paramCount++
	}
	if updUser.Pwd != "" {
		fmt.Fprintf(&qBuilder, " pwd = $%d,", paramCount)
		params = append(params, updUser.Pwd)
		paramCount++
	}
	if paramCount == 1 {
		return ud.GetUserById(id)
	}
	query := strings.TrimSuffix(qBuilder.String(), ",")
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, email", paramCount)
	params = append(params, id)

	var user types.User
	err := ud.db.QueryRow(query, params...).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("Error updating user: %w", err)
	}
	return &types.UserResponse{
		User: &user,
	}, nil
}

// func (ud *UsersData) ChangeUserActivity(id string, isActive bool) (User, error) {
// 	if id == "" {
// 		return User{}, fmt.Errorf("Error changing active status: id empty")
// 	}
// 	query := fmt.Sprintf("UPDATE users SET is_active = %d", !isActive)
// 	query += "WHERE id = $1 RETURNING id, name, email, is_active"
// 	var user User
// 	err := ud.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.IsActive)
// 	if err != nil {
// 		return user, fmt.Errorf("Error changing userStatus: queryRow/Scan: %w", err)
// 	}
// 	return user, nil
// }
