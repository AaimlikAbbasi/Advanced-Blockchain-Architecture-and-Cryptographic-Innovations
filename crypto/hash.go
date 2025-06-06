package crypto

import (
	"blockchain_A3/core"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func GenerateMerkleRoot(transactions []core.Transaction) string {
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
