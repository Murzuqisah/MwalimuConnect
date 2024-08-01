package functions

import (
	"crypto/ecdsa"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/lib/pq" // PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Initialize database connection
	var err error
	db, err = sql.Open("postgres", "user=username dbname=mwalimuconnect sslmode=disable")
	if err != nil {
		panic(err)
	}

	// Create users table if not exists
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            email TEXT PRIMARY KEY,
            password_hash TEXT NOT NULL,
            ethereum_address TEXT NOT NULL
        )
    `)
	if err != nil {
		panic(err)
	}
}

type EthereumAccount struct {
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}

func createEthereumAccount() (*EthereumAccount, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ethereum key: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &EthereumAccount{
		PrivateKey: privateKey,
		Address:    address,
	}, nil
}

func storeUserCredentials(email, password, ethereumAddress string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	_, err = db.Exec("INSERT INTO users (email, password_hash, ethereum_address) VALUES ($1, $2, $3)",
		email, string(hashedPassword), ethereumAddress)
	if err != nil {
		return fmt.Errorf("failed to store user credentials: %v", err)
	}

	return nil
}

func verifyUserCredentials(email, password string) (string, error) {
	var passwordHash, ethereumAddress string
	err := db.QueryRow("SELECT password_hash, ethereum_address FROM users WHERE email = $1", email).Scan(&passwordHash, &ethereumAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", fmt.Errorf("failed to retrieve user credentials: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	return ethereumAddress, nil
}

func getEthereumAccount(email string) (*EthereumAccount, error) {
	var ethereumAddress string
	err := db.QueryRow("SELECT ethereum_address FROM users WHERE email = $1", email).Scan(&ethereumAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to retrieve Ethereum address: %v", err)
	}

	// In a real-world scenario, you would securely retrieve the private key
	// This is just a placeholder and should not be used in production
	privateKey, err := crypto.HexToECDSA("your_private_key_here")
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return &EthereumAccount{
		PrivateKey: privateKey,
		Address:    common.HexToAddress(ethereumAddress),
	}, nil
}
