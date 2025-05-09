package amf

import (
	"blockchain_A3/core"
	"blockchain_A3/merkle"
	"fmt"
	"math"
	"sort"
)

const maxShardLoad = 10

type Shard struct {
	Tree         *MerkleTree
	Load         int
	Blocks       []*core.Block
	Transactions [][]byte
	RootHash     []byte
	States       [][]byte
}

type ShardManager struct {
	Shards []*Shard
	// Track all transactions to prevent duplicates
	allTransactions map[string]struct{}
}

// NewShard creates a new empty shard
func NewShard() *Shard {
	return &Shard{
		Load:     0,
		RootHash: []byte("initialRootHash"),
		States:   [][]byte{},
	}
}

// NewShardManager initializes a new ShardManager
func NewShardManager() *ShardManager {
	return &ShardManager{
		Shards: []*Shard{
			{
				Load:     0,
				RootHash: []byte("initialRootHash"),
				States:   [][]byte{},
			},
		},
		allTransactions: make(map[string]struct{}),
	}
}

// AddTransaction adds a transaction to the appropriate shard
func (sm *ShardManager) AddTransaction(tx []byte) error {
	// Check for duplicate transaction
	if _, exists := sm.allTransactions[string(tx)]; exists {
		return fmt.Errorf("transaction already exists")
	}

	// Add to all transactions map
	sm.allTransactions[string(tx)] = struct{}{}

	// If no shards exist, create the first one
	if len(sm.Shards) == 0 {
		sm.Shards = append(sm.Shards, NewShard())
	}

	// Find shard with lowest load
	lowestLoadShard := sm.Shards[0]
	for _, shard := range sm.Shards {
		if shard.Load < lowestLoadShard.Load {
			lowestLoadShard = shard
		}
	}

	// Add transaction to shard with lowest load
	lowestLoadShard.Transactions = append(lowestLoadShard.Transactions, tx)
	lowestLoadShard.States = append(lowestLoadShard.States, tx)
	lowestLoadShard.Load = len(lowestLoadShard.Transactions)
	lowestLoadShard.RecalculateRootHash()

	// Check if split is needed
	if lowestLoadShard.Load > maxShardLoad {
		return sm.SplitShard()
	}

	return nil
}

// Helper to get leaves as data
func (mt *MerkleTree) LeavesData() [][]byte {
	var data [][]byte
	for _, leaf := range mt.Leaves {
		data = append(data, leaf.Hash)
	}
	return data
}

// SplitShard splits the shard with highest load
func (sm *ShardManager) SplitShard() error {
	// Find shard with highest load
	var highestLoadShard *Shard
	maxLoad := 0

	for _, shard := range sm.Shards {
		if shard.Load > maxLoad {
			maxLoad = shard.Load
			highestLoadShard = shard
		}
	}

	if highestLoadShard == nil {
		return fmt.Errorf("no shards available to split")
	}

	if highestLoadShard.Load < maxShardLoad {
		return fmt.Errorf("shard load %d is below threshold %d", highestLoadShard.Load, maxShardLoad)
	}

	// Create new shard
	newShard := NewShard()

	// Move half of transactions to new shard
	mid := len(highestLoadShard.Transactions) / 2
	newShard.Transactions = append(newShard.Transactions, highestLoadShard.Transactions[mid:]...)
	newShard.States = append(newShard.States, highestLoadShard.States[mid:]...)

	// Update original shard
	highestLoadShard.Transactions = highestLoadShard.Transactions[:mid]
	highestLoadShard.States = highestLoadShard.States[:mid]

	// Update loads and root hashes
	highestLoadShard.Load = len(highestLoadShard.Transactions)
	newShard.Load = len(newShard.Transactions)
	highestLoadShard.RecalculateRootHash()
	newShard.RecalculateRootHash()

	// Add new shard to manager
	sm.Shards = append(sm.Shards, newShard)

	return nil
}

// ShouldSplit checks if any shard needs splitting
func (sm *ShardManager) ShouldSplit() bool {
	for _, shard := range sm.Shards {
		if shard.Load >= maxShardLoad {
			fmt.Printf("Shard with load %d needs splitting (threshold: %d)\n", shard.Load, maxShardLoad)
			return true
		}
	}
	return false
}

// PrintShards prints the current state of all shards
func (manager *ShardManager) PrintShards() {
	if len(manager.Shards) == 0 {
		fmt.Println("No shards available.")
		return
	}

	fmt.Println("\nShard Status Report:")
	fmt.Println("===================")
	for i, shard := range manager.Shards {
		fmt.Printf("\nShard %d Status:\n", i)
		fmt.Printf("-----------------\n")
		fmt.Printf("Current Load: %d/%d transactions\n", shard.Load, maxShardLoad)
		fmt.Printf("Root Hash: %x\n", shard.RootHash)

		if shard.Load == 0 {
			fmt.Println("Transactions: [Empty]")
		} else {
			fmt.Println("Transactions:")
			for j, tx := range shard.Transactions {
				fmt.Printf("  %d. %s\n", j+1, string(tx))
			}
		}

		// Show shard condition
		if shard.Load >= maxShardLoad {
			fmt.Println("Condition: Overloaded - Needs splitting")
		} else if shard.Load == 0 {
			fmt.Println("Condition: Empty")
		} else {
			fmt.Printf("Condition: Normal (%.1f%% capacity)\n", float64(shard.Load)/float64(maxShardLoad)*100)
		}
		fmt.Println("-----------------")
	}
	fmt.Println("===================")
}

// GetShard retrieves a shard by index
func (manager *ShardManager) GetShard(index int) *Shard {
	if index >= 0 && index < len(manager.Shards) {
		return manager.Shards[index]
	}
	return nil
}
func (sm *ShardManager) ShouldMerge() bool {
	if len(sm.Shards) < 2 {
		return false
	}

	// Find two shards with lowest load
	var lowestLoadShards [2]*Shard
	lowestLoadShards[0] = sm.Shards[0]
	lowestLoadShards[1] = sm.Shards[1]

	for _, shard := range sm.Shards {
		if shard.Load < lowestLoadShards[0].Load {
			lowestLoadShards[1] = lowestLoadShards[0]
			lowestLoadShards[0] = shard
		} else if shard.Load < lowestLoadShards[1].Load {
			lowestLoadShards[1] = shard
		}
	}

	// Check if merge is needed (if combined load is below threshold)
	combinedLoad := lowestLoadShards[0].Load + lowestLoadShards[1].Load
	if combinedLoad <= maxShardLoad {
		fmt.Printf("\nMerge possible: Shards with loads %d and %d (combined: %d, threshold: %d)\n",
			lowestLoadShards[0].Load, lowestLoadShards[1].Load, combinedLoad, maxShardLoad)
		return true
	}
	return false
}

func (sm *ShardManager) MergeShards() error {
	if len(sm.Shards) < 2 {
		return fmt.Errorf("not enough shards to merge")
	}

	// Find two shards with lowest load
	var lowestLoadShards [2]*Shard
	var lowestLoadIndices [2]int
	lowestLoadShards[0] = sm.Shards[0]
	lowestLoadShards[1] = sm.Shards[1]

	for i, shard := range sm.Shards {
		if shard.Load < lowestLoadShards[0].Load {
			lowestLoadShards[1] = lowestLoadShards[0]
			lowestLoadIndices[1] = lowestLoadIndices[0]
			lowestLoadShards[0] = shard
			lowestLoadIndices[0] = i
		} else if shard.Load < lowestLoadShards[1].Load {
			lowestLoadShards[1] = shard
			lowestLoadIndices[1] = i
		}
	}

	// Check if merge is needed
	combinedLoad := lowestLoadShards[0].Load + lowestLoadShards[1].Load
	if combinedLoad > maxShardLoad {
		return fmt.Errorf("combined load %d exceeds threshold %d", combinedLoad, maxShardLoad)
	}

	fmt.Printf("\nMerging shards with loads %d and %d\n", lowestLoadShards[0].Load, lowestLoadShards[1].Load)

	// Create merged shard
	mergedShard := &Shard{
		Transactions: append(lowestLoadShards[0].Transactions, lowestLoadShards[1].Transactions...),
		States:       append(lowestLoadShards[0].States, lowestLoadShards[1].States...),
		Load:         combinedLoad,
	}
	mergedShard.RecalculateRootHash()

	// Remove the two merged shards and add the new one
	sm.Shards = append(sm.Shards[:lowestLoadIndices[0]], sm.Shards[lowestLoadIndices[0]+1:]...)
	if lowestLoadIndices[1] > lowestLoadIndices[0] {
		lowestLoadIndices[1]--
	}
	sm.Shards = append(sm.Shards[:lowestLoadIndices[1]], sm.Shards[lowestLoadIndices[1]+1:]...)
	sm.Shards = append(sm.Shards, mergedShard)

	fmt.Printf("Merge complete. New shard has %d transactions\n", mergedShard.Load)

	// Check if the merged shard needs splitting
	if mergedShard.Load >= maxShardLoad {
		fmt.Println("\nMerged shard exceeds threshold, performing split...")
		return sm.SplitShard()
	}

	return nil
}

func (s *Shard) RecalculateRootHash() {
	// Sort transactions to ensure consistent ordering
	sortedTxs := make([][]byte, len(s.Transactions))
	copy(sortedTxs, s.Transactions)

	// Sort transactions by their string representation
	sort.Slice(sortedTxs, func(i, j int) bool {
		return string(sortedTxs[i]) < string(sortedTxs[j])
	})

	// Build Merkle tree with sorted transactions
	s.RootHash = calculateMerkleRoot(sortedTxs)
}

func calculateMerkleRoot(transactions [][]byte) []byte {
	if len(transactions) == 0 {
		return []byte{}
	}

	// Create a new Merkle tree with the sorted transactions
	tree := merkle.NewMerkleTree(transactions)
	return tree.Root.Hash
}

// ForceReduceLoad reduces the load of all shards to simulate low network conditions
func (sm *ShardManager) ForceReduceLoad() {
	fmt.Println("\nForcing load reduction on all shards...")
	for i, shard := range sm.Shards {
		if len(shard.Transactions) > 3 {
			// Keep only the first 3 transactions
			shard.Transactions = shard.Transactions[:3]
			shard.States = shard.States[:3]
			shard.Load = 3
			shard.RecalculateRootHash()
			fmt.Printf("Reduced load on shard %d to 3 transactions\n", i)
		}
	}
}

func CalculateShardEntropy(shard *Shard) float64 {
	// Create a map to count the frequency of each transaction
	txFrequency := make(map[string]int)

	for _, tx := range shard.Transactions {
		txFrequency[string(tx)]++
	}

	// Calculate entropy using Shannon's entropy formula
	var entropy float64
	totalTransactions := float64(len(shard.Transactions))

	for _, freq := range txFrequency {
		probability := float64(freq) / totalTransactions
		entropy -= probability * math.Log2(probability)
	}

	return entropy
}

// Add a conflict detection method that compares the entropy of two shards
func (sm *ShardManager) DetectConflictBasedOnEntropy(shardIndex1, shardIndex2 int) bool {
	shard1 := sm.Shards[shardIndex1]
	shard2 := sm.Shards[shardIndex2]

	entropy1 := CalculateShardEntropy(shard1)
	entropy2 := CalculateShardEntropy(shard2)

	// Compare the entropy values. If the difference is greater than a threshold, it's a conflict.
	entropyThreshold := 1.5 // You can adjust this threshold based on your requirements

	fmt.Printf("Shard 1 entropy: %.2f, Shard 2 entropy: %.2f\n", entropy1, entropy2)

	if math.Abs(entropy1-entropy2) > entropyThreshold {
		fmt.Println("Conflict detected based on entropy!")
		return true
	}

	return false
}
