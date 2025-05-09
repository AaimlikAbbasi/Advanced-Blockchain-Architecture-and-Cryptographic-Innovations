package auth

import (
	"crypto/sha256"
	"fmt"
)

func MultiFactorValidate(nodeID string) bool {
	// Mock MFA: check token + signature
	token := getNodeToken(nodeID)
	signature := getNodeSignature(nodeID)

	// Simulate validation
	expected := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", expected[:]) == signature
}

func getNodeToken(nodeID string) string {
	return "secure_token_123" // In real case, dynamically issued
}

func getNodeSignature(nodeID string) string {
	// Fake correct signature
	token := getNodeToken(nodeID)
	sig := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", sig[:])
}
