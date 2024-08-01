package functions

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var (
	db         *sql.DB
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func init() {
	// Initialize database connection
	var err error
	db, err = sql.Open("postgres", "user=username dbname=mwalimuconnect sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Create users table if not exists
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            email TEXT PRIMARY KEY,
            password_hash TEXT NOT NULL,
            ethereum_address TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		panic(fmt.Sprintf("Failed to create users table: %v", err))
	}
}

// func storeUserCredentials(email, password, ethereumAddress string) error {
// 	// Input validation
// 	if err := validateInputs(email, password, ethereumAddress); err != nil {
// 		return err
// 	}

// 	// Hash the password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return fmt.Errorf("failed to hash password: %v", err)
// 	}

// 	// Prepare SQL statement
// 	stmt, err := db.Prepare("INSERT INTO users (email, password_hash, ethereum_address) VALUES ($1, $2, $3)")
// 	if err != nil {
// 		return fmt.Errorf("failed to prepare SQL statement: %v", err)
// 	}
// 	defer stmt.Close()

// 	// Execute the prepared statement
// 	_, err = stmt.Exec(email, string(hashedPassword), ethereumAddress)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "unique constraint") {
// 			return errors.New("email already exists")
// 		}
// 		return fmt.Errorf("failed to store user credentials: %v", err)
// 	}

// 	return nil
// }

func validateInputs(email, password, ethereumAddress string) error {
	// Validate email
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	// Validate password
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Validate Ethereum address
	if !isValidEthereumAddress(ethereumAddress) {
		return errors.New("invalid Ethereum address")
	}

	return nil
}

func isValidEthereumAddress(address string) bool {
	// Basic Ethereum address validation
	// This checks if the address is a 40-character hex string prefixed with "0x"
	match, _ := regexp.MatchString("^0x[0-9a-fA-F]{40}$", address)
	return match
}
