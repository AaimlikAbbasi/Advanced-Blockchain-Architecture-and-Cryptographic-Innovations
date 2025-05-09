package amf

import (
	"blockchain_A3/verification"
	"crypto/sha256"
)

// MerkleNode structure
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

// MerkleTree structure
type MerkleTree struct {
	Root   *MerkleNode
	Leaves []*MerkleNode
}

// NewMerkleNode creates a new Merkle node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Hash = hash[:]
	} else {
		combined := append(left.Hash, right.Hash...)
		hash := sha256.Sum256(combined)
		node.Hash = hash[:]
	}

	node.Left = left
	node.Right = right

	return node
}

// NewMerkleTree builds a new Merkle Tree
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []*MerkleNode

	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, node)
	}

	leaves := nodes

	for len(nodes) > 1 {
		var newLevel []*MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			if i+1 == len(nodes) {
				newLevel = append(newLevel, NewMerkleNode(nodes[i], nodes[i], nil))
			} else {
				newLevel = append(newLevel, NewMerkleNode(nodes[i], nodes[i+1], nil))
			}
		}

		nodes = newLevel
	}

	return &MerkleTree{Root: nodes[0], Leaves: leaves}
}

// GetTransactionHash hashes a single transaction
func GetTransactionHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func (tree *MerkleTree) GenerateMerkleProof(index int) *verification.MerkleProof {
	proof := &verification.MerkleProof{
		Proof: make([][]byte, 0),
	}

	if index < 0 || index >= len(tree.Leaves) {
		return proof
	}

	// Add sibling hashes to the proof
	for i := 0; i < len(tree.Leaves); i++ {
		if i != index {
			proof.Proof = append(proof.Proof, tree.Leaves[i].Hash)
		}
	}

	return proof
}
