package blockchain

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// UserStorageContract represents a binding to the UserStorage smart contract
type UserStorageContract struct {
	*bind.BoundContract
}

type UserInfo struct {
	Name              string
	Email             string
	Role              string
	WalletBalance     *big.Int
	CompletedSessions *big.Int
	UpcomingSessions  *big.Int
	TotalHours        *big.Int
	AverageRating     *big.Int
}

type UserByEmail struct {
	Name              string
	Role              string
	WalletBalance     *big.Int
	CompletedSessions *big.Int
	UpcomingSessions  *big.Int
	TotalHours        *big.Int
	AverageRating     *big.Int
}

const UserStorageABI = `[{"inputs":[{"internalType":"string","name":"name","type":"string"},
{"internalType":"string","name":"email","type":"string"},
{"internalType":"string","name":"role","type":"string"}],"name":"createUser","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

// create eth accounts
func generateEthereumAccount() (*ecdsa.PrivateKey, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Printf("Failed to generate Ethereum private key: %v", err)
		return nil, "", err
	}

	publicKey := privateKey.Public()
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	return privateKey, address.Hex(), nil
}

func createUserAccount(email, password string) (string, error) {
	_, address, err := generateEthereumAccount()
	if err != nil {
		return "", err
	}

	// Check for empty credentials
	if email == "" || password == "" || address == "" {
		return "", errors.New("all fields (email, password, address) are required")
	}

	// Implement your logic to store user credentials
	// For example, you can use a database or a file to store the credentials
	// Here, I'm just logging the credentials for demonstration purposes
	log.Printf("Storing user credentials: email=%s, password=%s, address=%s", email, password, address)

	return address, nil
}

// NewUserStorageContract creates a new instance of UserStorageContract, bound to a specific deployed contract.
func NewUserStorageContract(address common.Address, client *ethclient.Client) (*UserStorageContract, error) {
	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(UserStorageABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Create a new instance of the contract binding
	contract := bind.NewBoundContract(address, parsedABI, client, client, client)

	return &UserStorageContract{contract}, nil
}

// CreateUser calls the createUser function of the smart contract
func (c *UserStorageContract) CreateUser(opts *bind.TransactOpts, name, email, role string) (*types.Transaction, error) {
	if name == "" || email == "" || role == "" {
		return nil, errors.New("name, email, and role are required")
	}
	return c.Transact(opts, "createUser", name, email, role)
}

// GetUser calls the getUser function of the smart contract
// func (c *UserStorageContract) GetUserDetails(opts *bind.CallOpts, userAddress common.Address) (UserInfo, error) {
// 	if opts == nil {
// 		return UserInfo{}, errors.New("opts cannot be nil; please provide valid call options")
// 	}
// 	if userAddress == (common.Address{}) {
// 		return UserInfo{}, errors.New("userAddress cannot be empty; please provide a valid user address")
// 	}

// 	var out UserInfo
// 	var iface []interface{}
// 	iface = append(iface, &out)
// 	err := c.Call(opts, &iface, "getUser", userAddress)
// 	if err != nil {
// 		log.WithError(err).WithField("address", userAddress).Error("failed to call getUserDetails")
// 		return UserInfo{}, fmt.Errorf("failed to call getUserDetails for address %s: %w", userAddress, err)
// 	}
// 	return out, nil
// }

// GetUser calls the getUser function of the smart contract
func (c *UserStorageContract) GetUserDetails(opts *bind.CallOpts, userAddress common.Address) (UserInfo, error) {
	if opts == nil {
		return UserInfo{}, errors.New("opts cannot be nil; please provide valid call options")
	}
	if userAddress == (common.Address{}) {
		return UserInfo{}, errors.New("userAddress cannot be empty; please provide a valid user address")
	}

	var out UserInfo
	var name string
	var email string
	var role string
	var walletBalance *big.Int
	var completedSessions *big.Int
	var upcomingSessions *big.Int
	var totalHours *big.Int
	var averageRating *big.Int
	err := c.Call(opts, &[]interface{}{
		&name,
		&email,
		&role,
		&walletBalance,
		&completedSessions,
		&upcomingSessions,
		&totalHours,
		&averageRating,
	}, "getUser", userAddress)
	if err != nil {
		log.WithError(err).WithField("address", userAddress).Error("failed to call getUserDetails")
		return UserInfo{}, fmt.Errorf("failed to call getUserDetails for address %s: %w", userAddress, err)
	}
	out.Name = name
	out.Email = email
	out.Role = role
	out.WalletBalance = walletBalance
	out.CompletedSessions = completedSessions
	out.UpcomingSessions = upcomingSessions
	out.TotalHours = totalHours
	out.AverageRating = averageRating
	return out, nil
}

// GetUserByEmail calls the getUserByEmail function of the smart contract
// GetUserByEmail calls the getUserByEmail function of the smart contract
func (c *UserStorageContract) GetUserByEmail(opts *bind.CallOpts, email string) (UserByEmail, error) {
	if opts == nil {
		return UserByEmail{}, errors.New("opts cannot be nil; please provide valid call options")
	}
	if email == "" {
		return UserByEmail{}, errors.New("email cannot be empty; please provide a valid email")
	}

	var out UserByEmail
	var name string
	var role string
	var walletBalance *big.Int
	var completedSessions *big.Int
	var upcomingSessions *big.Int
	var totalHours *big.Int
	var averageRating *big.Int
	err := c.Call(opts, &[]interface{}{
		&name,
		&role,
		&walletBalance,
		&completedSessions,
		&upcomingSessions,
		&totalHours,
		&averageRating,
	}, "getUserByEmail", email)
	if err != nil {
		log.WithError(err).WithField("email", email).Error("failed to call getUserByEmail")
		return UserByEmail{}, fmt.Errorf("failed to call getUserByEmail for email %s: %w", email, err)
	}
	out.Name = name
	out.Role = role
	out.WalletBalance = walletBalance
	out.CompletedSessions = completedSessions
	out.UpcomingSessions = upcomingSessions
	out.TotalHours = totalHours
	out.AverageRating = averageRating
	return out, nil
}

// UpdateUser calls the updateUser function of the smart contract
func (c *UserStorageContract) UpdateUser(opts *bind.TransactOpts, name, email, role string) (*types.Transaction, error) {
	return c.Transact(opts, "updateUser", name, email, role)
}

// UpdateWalletBalance calls the updateWalletBalance function of the smart contract
func (c *UserStorageContract) UpdateWalletBalance(opts *bind.TransactOpts, newBalance *big.Int) (*types.Transaction, error) {
	return c.Transact(opts, "updateWalletBalance", newBalance)
}

// UpdateLearningStats calls the updateLearningStats function of the smart contract
func (c *UserStorageContract) UpdateLearningStats(opts *bind.TransactOpts, completedSessions, upcomingSessions, totalHours, averageRating *big.Int) (*types.Transaction, error) {
	return c.Transact(opts, "updateLearningStats", completedSessions, upcomingSessions, totalHours, averageRating)
}
