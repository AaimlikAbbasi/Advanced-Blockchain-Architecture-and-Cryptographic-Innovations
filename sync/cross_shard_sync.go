package sync

import (
	"blockchain_A3/amf"
	"blockchain_A3/merkle"
	"fmt"
)

func CrossShardTransfer(manager *amf.ShardManager, fromShardID, toShardID int, txIndex int) error {
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

	// Step 2: Build Merkle Tree of fromShard transactions
	merkleTree := merkle.NewMerkleTree(fromShard.Transactions)

	// Step 3: Generate Merkle Proof for the transaction
	proof, err := merkleTree.GenerateProof(tx)
	if err != nil {
		return fmt.Errorf("failed to generate Merkle proof: %v", err)
	}

	// Step 4: Verify Proof against the shard's root hash
	isValid := merkle.VerifyProof(tx, proof, fromShard.RootHash)
	if !isValid {
		fmt.Printf("Debug - Transaction: %s\n", string(tx))
		fmt.Printf("Debug - Root Hash: %x\n", fromShard.RootHash)
		fmt.Printf("Debug - Proof Length: %d\n", len(proof))
		for i, p := range proof {
			fmt.Printf("Debug - Proof[%d]: %x (IsLeft: %v)\n", i, p.Hash, p.IsLeft)
		}
		return fmt.Errorf("merkle proof verification failed for source shard")
	}

	fmt.Println("Merkle proof verified. Proceeding with cross-shard transfer...")

	// Step 5: Create backup of current states
	fromShardBackup := make([][]byte, len(fromShard.Transactions))
	copy(fromShardBackup, fromShard.Transactions)
	toShardBackup := make([][]byte, len(toShard.Transactions))
	copy(toShardBackup, toShard.Transactions)

	// Step 6: Remove transaction from source shard first
	fromShard.Transactions = append(fromShard.Transactions[:txIndex], fromShard.Transactions[txIndex+1:]...)
	fromShard.States = append(fromShard.States[:txIndex], fromShard.States[txIndex+1:]...)
	fromShard.Load = len(fromShard.Transactions)
	fromShard.RecalculateRootHash()

	// Step 7: Add transaction to destination shard in a consistent position
	// Always append to the end to maintain consistent ordering
	toShard.Transactions = append(toShard.Transactions, tx)
	toShard.States = append(toShard.States, tx)
	toShard.Load = len(toShard.Transactions)
	toShard.RecalculateRootHash()

	// Step 8: Verify the new states
	fromMerkleTree := merkle.NewMerkleTree(fromShard.Transactions)
	_, err = fromMerkleTree.GenerateProof(tx)
	if err == nil {
		// If proof generation succeeds, it means the transaction wasn't properly removed
		rollbackShardState(fromShard, toShard, fromShardBackup, toShardBackup)
		return fmt.Errorf("transaction still present in source shard after removal")
	}

	// Step 9: Verify destination shard state
	toMerkleTree := merkle.NewMerkleTree(toShard.Transactions)
	toProof, err := toMerkleTree.GenerateProof(tx)
	if err != nil {
		// Rollback if proof generation fails
		rollbackShardState(fromShard, toShard, fromShardBackup, toShardBackup)
		return fmt.Errorf("failed to generate Merkle proof for destination shard: %v", err)
	}

	// Step 10: Verify the proof against the new root hash
	isValid = merkle.VerifyProof(tx, toProof, toShard.RootHash)
	if !isValid {
		fmt.Printf("Debug - Destination verification failed:\n")
		fmt.Printf("Transaction: %s\n", string(tx))
		fmt.Printf("Root Hash: %x\n", toShard.RootHash)
		fmt.Printf("Proof Length: %d\n", len(toProof))
		for i, p := range toProof {
			fmt.Printf("Proof[%d]: %x (IsLeft: %v)\n", i, p.Hash, p.IsLeft)
		}
		// Rollback if verification fails
		rollbackShardState(fromShard, toShard, fromShardBackup, toShardBackup)
		return fmt.Errorf("merkle proof verification failed for destination shard")
	}

	fmt.Println("Cross-shard transfer complete.")
	fmt.Printf("Updated Root Hash for Shard %d: %x\n", fromShardID, fromShard.RootHash)
	fmt.Printf("Updated Root Hash for Shard %d: %x\n", toShardID, toShard.RootHash)
	return nil
}

// rollbackShardState performs a rollback of the shard states
func rollbackShardState(fromShard, toShard *amf.Shard, fromBackup, toBackup [][]byte) {
	fromShard.Transactions = fromBackup
	fromShard.States = fromBackup
	fromShard.Load = len(fromShard.Transactions)
	fromShard.RecalculateRootHash()

	toShard.Transactions = toBackup
	toShard.States = toBackup
	toShard.Load = len(toShard.Transactions)
	toShard.RecalculateRootHash()
}
