package functions

import (
	"crypto/sha256"
	"time"
)

type Block struct {
	Timestamp    int64
	Data         []byte
	PreviousHash []byte
	Hash         []byte
}

// Blockchain block linking structure
func CalculateHash(block Block) []byte {
	record := string(block.Timestamp) + string(block.Data) + string(block.PreviousHash)
	hash := sha256.New()
	hash.Write([]byte(record))
	return hash.Sum(nil)
}

func GenerateBlock(previousBlock Block, data []byte) Block {
	newBlock := Block{
		Timestamp:    time.Now().Unix(),
		Data:         data,
		PreviousHash: previousBlock.Hash,
	}
	newBlock.Hash = CalculateHash(newBlock)
	return newBlock
}
