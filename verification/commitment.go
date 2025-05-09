package verification

import (
	"crypto/sha256"
	"fmt"
)

// CreateCommitment generates a cryptographic commitment for a piece of data
func CreateCommitment(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// VerifyCommitment verifies the commitment using the provided data
func VerifyCommitment(commitment, data []byte) bool {
	expectedCommitment := CreateCommitment(data)
	return fmt.Sprintf("%x", expectedCommitment) == fmt.Sprintf("%x", commitment)
}
