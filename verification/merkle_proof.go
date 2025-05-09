package verification

import (
	"bytes"
	"crypto/sha256"
)

type MerkleProof struct {
	Proof [][]byte
}

func VerifyMerkleProof(proof *MerkleProof, leafHash, rootHash []byte) bool {
	computed := leafHash

	for _, siblingHash := range proof.Proof {
		var combined []byte
		if bytes.Compare(computed, siblingHash) < 0 {
			combined = append(computed, siblingHash...)
		} else {
			combined = append(siblingHash, computed...)
		}
		hash := sha256.Sum256(combined)
		computed = hash[:]
	}

	return bytes.Equal(computed, rootHash)
}
