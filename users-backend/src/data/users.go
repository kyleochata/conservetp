package data

import (
	"database/sql"
	"fmt"
)

type UsersData struct {
	db *sql.DB
}

func NewUsersData(db *sql.DB) *UsersData {
	return &UsersData{db: db}
}

type User struct {
	ID    string `json: "id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (ud *UsersData) GetAllUsers() ([]User, error) {
	rows, err := ud.db.Query("SELECT id, name, email FROM users WHERE is_active = true")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}

func (ud *UsersData) CreateUser(name, email, pwd string) (*User, error) {
	var user User
	err := ud.db.QueryRow(
		"INSERT INTO users (name, email, pwd) VALUES ($1, $2, $3) RETURNING id, name, email",
		name, email, pwd,
	).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
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
