package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type SessionChain struct {
	chainID         *big.Int
	contractABI     abi.ABI
	txHistory       []*types.Transaction
	lastError       error
	client          *ethclient.Client
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
	instance        *SessionStorage
}

type SessionStorage struct {
	contract *bind.BoundContract
}

func NewSessionChain(client *ethclient.Client, contractAddress common.Address, privateKey *ecdsa.PrivateKey) (*SessionChain, error) {
	instance, err := NewSessionStorage(contractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate smart contract: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	return &SessionChain{
		client:          client,
		contractAddress: contractAddress,
		privateKey:      privateKey,
		instance:        instance,
		chainID:         chainID,
	}, nil
}

func (sc *SessionChain) StoreSessionID(sessionID string) (*types.Transaction, error) {
	auth, err := sc.getTransactOpts()
	if err != nil {
		return nil, err
	}

	tx, err := sc.instance.contract.Transact(auth, "StoreSessionID", sessionID)
	if err != nil {
		sc.lastError = fmt.Errorf("failed to store session ID: %w", err)
		return nil, sc.lastError
	}

	sc.txHistory = append(sc.txHistory, tx)
	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())
	return tx, nil
}

func (sc *SessionChain) GetSessionID() (string, error) {
	sessionID, err := sc.instance.GetSessionID(nil)
	if err != nil {
		sc.lastError = fmt.Errorf("failed to retrieve session ID: %w", err)
		return "", sc.lastError
	}
	return sessionID, nil
}

func (sc *SessionChain) GetChainID() (*big.Int, error) {
	return sc.chainID, nil
}

func (sc *SessionChain) GetContractABI() abi.ABI {
	return sc.contractABI
}

func (sc *SessionChain) GetTransactionHistory() []*types.Transaction {
	return sc.txHistory
}

func (sc *SessionChain) GetLastError() error {
	return sc.lastError
}

func (sc *SessionChain) getTransactOpts() (*bind.TransactOpts, error) {
	publicKey := sc.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := sc.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve nonce: %w", err)
	}

	gasPrice, err := sc.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve suggested gas price: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(sc.privateKey, sc.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	return auth, nil
}

func NewSessionStorage(contractAddress common.Address, client *ethclient.Client) (*SessionStorage, error) {
	abiStr := `[{"constant":true,"inputs":[],"name":"GetSessionID","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_sessionID","type":"string"}],"name":"StoreSessionID","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

	abiParsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	contract := bind.NewBoundContract(contractAddress, abiParsed, client, client, client)

	return &SessionStorage{
		contract: contract,
	}, nil
}

func (ss *SessionStorage) GetSessionID(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := ss.contract.Call(opts, &out, "GetSessionID")
	if err != nil {
		return "", err
	}
	return out[0].(string), nil
}
