package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockchainOperations interface {
	StoreSessionID(sessionID string) (*types.Transaction, error)
	GetSessionID() (string, error)
	GetChainID() (*big.Int, error)
	GetContractABI() abi.ABI
	GetTransactionHistory() []*types.Transaction
	GetLastError() error
}
