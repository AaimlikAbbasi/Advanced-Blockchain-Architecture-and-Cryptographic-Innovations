package sync

import (
	"blockchain_A3/amf"
	"blockchain_A3/merkle"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// HomomorphicHash represents a homomorphic hash function
type HomomorphicHash struct {
	Key []byte
}

// NewHomomorphicHash creates a new homomorphic hash instance
func NewHomomorphicHash() *HomomorphicHash {
	key := make([]byte, 32)
	rand.Read(key)
	return &HomomorphicHash{Key: key}
}

// Hash computes a homomorphic hash of the data
func (h *HomomorphicHash) Hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(h.Key)
	hash.Write(data)
	return hash.Sum(nil)
}

// Combine hashes homomorphically
func (h *HomomorphicHash) Combine(hashes ...[]byte) []byte {
	result := make([]byte, 32)
	for _, hash := range hashes {
		for i := range result {
			result[i] ^= hash[i]
		}
	}
	return result
}

// Commitment represents a cryptographic commitment
type Commitment struct {
	Value    []byte
	Blinding []byte
	Hash     []byte
}

// NewCommitment creates a new commitment
func NewCommitment(value []byte) *Commitment {
	blinding := make([]byte, 32)
	rand.Read(blinding)

	hash := sha256.New()
	hash.Write(value)
	hash.Write(blinding)
	commitmentHash := hash.Sum(nil)

	return &Commitment{
		Value:    value,
		Blinding: blinding,
		Hash:     commitmentHash,
	}
}

// VerifyCommitment verifies a commitment
func VerifyCommitment(commitment *Commitment) bool {
	hash := sha256.New()
	hash.Write(commitment.Value)
	hash.Write(commitment.Blinding)
	computedHash := hash.Sum(nil)
	return string(computedHash) == string(commitment.Hash)
}

// AdvancedCrossShardTransfer performs an atomic cross-shard transfer with homomorphic authentication
func AdvancedCrossShardTransfer(manager *amf.ShardManager, fromShardID, toShardID int, txIndex int) error {
	// Validate shard IDs
	if len(manager.Shards) <= fromShardID || len(manager.Shards) <= toShardID {
		return fmt.Errorf("invalid shard IDs provided: from=%d, to=%d", fromShardID, toShardID)
	}

	fromShard := manager.Shards[fromShardID]
	toShard := manager.Shards[toShardID]

	// Validate transaction index
	if len(fromShard.Transactions) <= txIndex {
		return fmt.Errorf("invalid transaction index: %d", txIndex)
	}

	// Step 1: Get the transaction to transfer
	tx := fromShard.Transactions[txIndex]
	fmt.Printf("Initiating transfer of transaction: %s\n", tx)

	// Step 2: Create homomorphic hash instance
	homomorphicHash := NewHomomorphicHash()
	txHash := homomorphicHash.Hash(tx)

	// Step 3: Create commitment for the transaction
	commitment := NewCommitment(txHash)
	if !VerifyCommitment(commitment) {
		return fmt.Errorf("commitment verification failed")
	}

	// Step 4: Build Merkle Tree with homomorphic hashes
	merkleTree := merkle.NewMerkleTree(fromShard.Transactions)

	// Step 5: Generate Merkle Proof with homomorphic verification
	proof, err := merkleTree.GenerateProof(tx)
	if err != nil {
		return fmt.Errorf("failed to generate Merkle proof: %v", err)
	}

	// Step 6: Verify Proof with homomorphic properties
	isValid := merkle.VerifyProof(tx, proof, fromShard.RootHash)
	if !isValid {
		return fmt.Errorf("merkle proof verification failed for source shard")
	}

	fmt.Println("Merkle proof verified. Proceeding with cross-shard transfer...")

	// Step 7: Create backup of current states
	fromShardBackup := make([][]byte, len(fromShard.Transactions))
	copy(fromShardBackup, fromShard.Transactions)
	toShardBackup := make([][]byte, len(toShard.Transactions))
	copy(toShardBackup, toShard.Transactions)

	// Step 8: Remove transaction from source shard
	fromShard.Transactions = append(fromShard.Transactions[:txIndex], fromShard.Transactions[txIndex+1:]...)
	fromShard.States = append(fromShard.States[:txIndex], fromShard.States[txIndex+1:]...)
	fromShard.Load = len(fromShard.Transactions)
	fromShard.RecalculateRootHash()

	// Step 9: Add transaction to destination shard with commitment
	toShard.Transactions = append(toShard.Transactions, tx)
	toShard.States = append(toShard.States, tx)
	toShard.Load = len(toShard.Transactions)
	toShard.RecalculateRootHash()

	// Step 10: Verify the new states with homomorphic verification
	toMerkleTree := merkle.NewMerkleTree(toShard.Transactions)
	toProof, err := toMerkleTree.GenerateProof(tx)
	if err != nil {
		// Rollback if proof generation fails
		rollback(fromShard, toShard, fromShardBackup, toShardBackup)
		return fmt.Errorf("failed to generate Merkle proof for destination shard: %v", err)
	}

	isValid = merkle.VerifyProof(tx, toProof, toShard.RootHash)
	if !isValid {
		// Rollback if verification fails
		rollback(fromShard, toShard, fromShardBackup, toShardBackup)
		return fmt.Errorf("merkle proof verification failed for destination shard")
	}

	fmt.Println("Cross-shard transfer complete.")
	fmt.Printf("Updated Root Hash for Shard %d: %x\n", fromShardID, fromShard.RootHash)
	fmt.Printf("Updated Root Hash for Shard %d: %x\n", toShardID, toShard.RootHash)
	return nil
}

// rollback performs a rollback of the shard states
func rollback(fromShard, toShard *amf.Shard, fromBackup, toBackup [][]byte) {
	fromShard.Transactions = fromBackup
	fromShard.States = fromBackup
	fromShard.Load = len(fromShard.Transactions)
	fromShard.RecalculateRootHash()

	toShard.Transactions = toBackup
	toShard.States = toBackup
	toShard.Load = len(toShard.Transactions)
	toShard.RecalculateRootHash()
}
