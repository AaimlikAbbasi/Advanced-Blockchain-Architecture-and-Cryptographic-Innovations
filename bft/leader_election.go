package bft

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// Generate a VRF value for leader election
func GenerateVRF() ([]byte, []byte) {
	randomData := make([]byte, 32)
	_, err := rand.Read(randomData)
	if err != nil {
		fmt.Println("Error generating random data:", err)
	}

	hash := sha256.Sum256(randomData)
	return randomData, hash[:]
}
