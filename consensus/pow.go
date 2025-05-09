package consensus

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Simulated Proof-of-Work to inject randomness
func GeneratePoWRandomness(data string) string {
	nonce := 0
	for {
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%d", data, nonce)))
		if hash[0] == 0x00 {
			return fmt.Sprintf("%x", hash[:])
		}
		nonce++
		if nonce > 1e6 {
			break // prevent infinite loops in demo
		}
	}
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
