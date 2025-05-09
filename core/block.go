package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func GenerateMerkleRoot(transactions []*Transaction) string {
	if len(transactions) == 0 {
		return Hash("")
	}

	var hashes []string
	for _, tx := range transactions {
		hashes = append(hashes, Hash(tx.Sender+tx.Receiver+fmt.Sprintf("%f", tx.Amount)))
	}

	for len(hashes) > 1 {
		var newLevel []string
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				newLevel = append(newLevel, Hash(hashes[i]+hashes[i+1]))
			} else {
				newLevel = append(newLevel, Hash(hashes[i]))
			}
		}
		hashes = newLevel
	}

	return hashes[0]
}

type Block struct {
	Index        int
	Timestamp    time.Time
	Transactions []*Transaction
	PrevHash     string
	Hash         string
	MerkleRoot   string
	Data         []byte
}

func CreateBlock(index int, transactions []*Transaction, prevHash string) Block {
	timestamp := time.Now()
	merkleRoot := GenerateMerkleRoot(transactions)
	blockData := fmt.Sprintf("%d%s%s%s", index, timestamp.String(), merkleRoot, prevHash)
	hash := Hash(blockData)

	return Block{
		Index:        index,
		Timestamp:    timestamp,
		Transactions: transactions,
		PrevHash:     prevHash,
		Hash:         hash,
		MerkleRoot:   merkleRoot,
	}
}

func NewBlock(index int, data []byte) *Block {
	return &Block{
		Index:     index,
		Timestamp: time.Now(),
		Data:      data,
		PrevHash:  "",
		Hash:      "",
	}
}
