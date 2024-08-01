package blockchain

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type UserStorage struct {
	contract *UserStorageContract
	client   *ethclient.Client
}

type User struct {
	Name              string
	Email             string
	Role              string
	WalletBalance     *big.Int
	CompletedSessions *big.Int
	UpcomingSessions  *big.Int
	TotalHours        *big.Int
	AverageRating     *big.Int
}

func NewUserStorage(address common.Address, client *ethclient.Client) (*UserStorage, error) {
	contract, err := NewUserStorageContract(address, client)
	if err != nil {
		return nil, err
	}
	if contract == nil || client == nil {
		return nil, fmt.Errorf("failed to initialize UserStorage: contract or client is nil")
	}
	return &UserStorage{
		contract: contract,
		client:   client,
	}, nil
}

func (u *UserStorage) CreateUser(auth *bind.TransactOpts, name, email, role string) (*types.Transaction, error) {
	// Create a new transaction
	tx, err := u.contract.Transact(auth, "createUser", name, email, role)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (us *UserStorage) GetUser(address common.Address) (User, error) {
	var user User
	var result []interface{}
	result = make([]interface{}, 8)
	result[0] = &user.Name
	result[1] = &user.Email
	result[2] = &user.Role
	result[3] = &user.WalletBalance
	result[4] = &user.CompletedSessions
	result[5] = &user.UpcomingSessions
	result[6] = &user.TotalHours
	result[7] = &user.AverageRating
	err := us.contract.Call(nil, &result, "getUser", address)
	return user, err
}

func (us *UserStorage) GetUserByEmail(email string) (User, error) {
	var user User
	var result []interface{}
	result = make([]interface{}, 8)
	result[0] = &user.Name
	result[1] = &user.Email
	result[2] = &user.Role
	result[3] = &user.WalletBalance
	result[4] = &user.CompletedSessions
	result[5] = &user.UpcomingSessions
	result[6] = &user.TotalHours
	result[7] = &user.AverageRating
	err := us.contract.Call(nil, &result, "getUserByEmail", email)
	return user, err
}

func (us *UserStorage) UpdateUser(auth *bind.TransactOpts, name, email, role string) error {
	_, err := us.contract.Transact(auth, "updateUser", name, email, role)
	return err
}

func (us *UserStorage) UpdateWalletBalance(auth *bind.TransactOpts, newBalance *big.Int) error {
	_, err := us.contract.Transact(auth, "updateWalletBalance", newBalance)
	return err
}

func (us *UserStorage) UpdateLearningStats(auth *bind.TransactOpts, completedSessions, upcomingSessions, totalHours, averageRating *big.Int) error {
	_, err := us.contract.Transact(auth, "updateLearningStats", completedSessions, upcomingSessions, totalHours, averageRating)
	return err
}
