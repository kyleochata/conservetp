package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kyleochata/conservetp/users-backend/src/data"
	"github.com/kyleochata/conservetp/users-backend/src/handlers"
	"github.com/kyleochata/conservetp/users-backend/src/services"
	_ "github.com/lib/pq"
)

func main() {
	// Wait for database to be ready
	if err := waitForDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db, err := connectToDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	fmt.Println("✅ Successfully connected to database!")

	// Initialize layers
	usersData := data.NewUsersData(db)
	usersService := services.NewUsersService(usersData)
	usersHandler := handlers.NewUsersHandler(usersService)

	// Setup routes
	http.HandleFunc("/users", usersHandler.HandleUsers)
	//http.HandleFunc("/users/", usersHandler.HandleUserByID)
	//http.HandleFunc("/health", healthCheck)

	fmt.Println("🚀 Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func waitForDB() error {
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		db, err := connectToDB()
		if err != nil {
			log.Printf("Attempt %d/%d: Database not ready yet: %v", i+1, maxAttempts, err)
			time.Sleep(2 * time.Second)
			continue
		}
		db.Close()
		return nil
	}
	return fmt.Errorf("database not ready after %d attempts", maxAttempts)
}

func connectToDB() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "users-db"),
		getEnv("DB_PORT", "5432"),
		getEnv("POSTGRES_USER", "postgres"),
		getEnv("POSTGRES_PASSWORD", ""),
		getEnv("POSTGRES_DB", "users_db"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func testUsersTable(db *sql.DB) error {
	// Try to select from users table
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("users table query failed: %w", err)
	}

	fmt.Printf("✅ Users table exists with %d records\n", count)
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
