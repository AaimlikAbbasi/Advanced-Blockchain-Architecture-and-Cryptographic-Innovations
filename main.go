package main

import (
	"blockchain_A3/amf"
	"blockchain_A3/bft"
	"blockchain_A3/core"
	"blockchain_A3/sync"
	"blockchain_A3/verification"
	"fmt"
	"math/big"
)

func main() {
	fmt.Println("Initializing Blockchain...")

	// Initialize the ConsistencyOrchestrator
	co := core.NewConsistencyOrchestrator()

	// Simulate network telemetry (this can be dynamically updated in a real system)
	networkTelemetry := core.NetworkTelemetry{
		Latency:    120,  // in ms
		PacketLoss: 0.15, // 15% packet loss
		Throughput: 50.0, // 50 Mbps throughput
	}

	// Set network telemetry to the ConsistencyOrchestrator
	co.SetNetworkTelemetry(networkTelemetry)

	// Calculate network partition risk based on telemetry
	partitionRisk := calculatePartitionRisk(networkTelemetry)
	fmt.Printf("Predicted Network Partition Risk: %.2f%%\n", partitionRisk)

	// Adjust consistency and network settings based on telemetry
	co.AdjustConsistency()
	fmt.Printf("Current Consistency Level: %s\n", getConsistencyLevelName(co.ConsistencyLevel))
	fmt.Printf("Timeout: %ds, Retries: %d\n", 5, 3)

	// Create new blockchain
	bc := core.NewBlockchain()

	// Add dummy transactions
	tx1 := core.NewTransaction("Alice", "Bob", 5)
	tx2 := core.NewTransaction("Bob", "Charlie", 2)

	// Add block
	bc.AddBlock([]*core.Transaction{tx1, tx2})

	fmt.Println("\nBlockchain created with genesis block.")

	// ------------------------
	// Merkle Tree Demonstration
	// ------------------------
	transactions := [][]byte{
		[]byte("Alice -> Bob:5"),
		[]byte("Bob -> Charlie:2"),
	}

	tree := amf.NewMerkleTree(transactions)

	fmt.Printf("\nMerkle Root: %x\n", tree.Root.Hash)

	// Generate Merkle Proof for first transaction
	proof := tree.GenerateMerkleProof(0)

	// Verify the proof
	isValid := verification.VerifyMerkleProof(proof, amf.GetTransactionHash(transactions[0]), tree.Root.Hash)
	fmt.Println("Merkle Proof Valid:", isValid)

	// ------------------------
	// Shard Management
	// ------------------------
	fmt.Println("\nInitializing Shard Manager...")
	manager := amf.NewShardManager()

	// Simulate adding load to shards
	fmt.Println("\nAdding transactions to shards...")

	// Add initial transactions (more than maxShardLoad to force split)
	initialTransactions := []string{
		"User0 -> User1: 0",
		"User1 -> User2: 10",
		"User2 -> User3: 20",
		"User3 -> User4: 30",
		"User4 -> User5: 40",
		"User5 -> User6: 50",
		"User6 -> User7: 60",
		"User7 -> User8: 70",
		"User8 -> User9: 80",
		"User9 -> User10: 90",
		"User10 -> User11: 100", // This should trigger split
	}

	// Add initial transactions to a single shard
	fmt.Println("\n=== INITIAL STATE - SINGLE SHARD ===")
	for _, tx := range initialTransactions {
		if err := manager.AddTransaction([]byte(tx)); err != nil {
			fmt.Printf("Error adding transaction: %v\n", err)
		}
	}
	manager.PrintShards()

	// Force split when load exceeds threshold
	if manager.ShouldSplit() {
		fmt.Println("\n=== PERFORMING SHARD SPLIT ===")
		if err := manager.SplitShard(); err != nil {
			fmt.Printf("Error splitting shard: %v\n", err)
		}
		fmt.Println("\n=== AFTER SPLIT ===")
		manager.PrintShards()
	}

	// Add more transactions
	fmt.Println("\n=== ADDING MORE TRANSACTIONS ===")
	extraTransactions := []string{
		"User11 -> User12: 110",
		"User12 -> User13: 120",
		"User13 -> User14: 130",
		"User14 -> User15: 140",
		"User15 -> User16: 150",
		"User16 -> User17: 160",
		"User17 -> User18: 170",
		"User18 -> User19: 180",
		"User19 -> User20: 190",
		"User20 -> User21: 200",
	}

	for _, tx := range extraTransactions {
		if err := manager.AddTransaction([]byte(tx)); err != nil {
			fmt.Printf("Error adding transaction: %v\n", err)
		}
	}
	manager.PrintShards()

	// Check and split overloaded shards
	for manager.ShouldSplit() {
		fmt.Println("\n=== PERFORMING SHARD SPLIT ===")
		if err := manager.SplitShard(); err != nil {
			fmt.Printf("Error splitting shard: %v\n", err)
			break
		}
		fmt.Println("\n=== AFTER SPLIT ===")
		manager.PrintShards()
	}

	// Check and perform merge if needed
	if manager.ShouldMerge() {
		fmt.Println("\n=== PERFORMING SHARD MERGE ===")
		if err := manager.MergeShards(); err != nil {
			fmt.Printf("Error in merge/split operation: %v\n", err)
		} else {
			fmt.Println("\n=== AFTER MERGE/SPLIT ===")
			manager.PrintShards()
		}
	} else {
		fmt.Println("\nNo shards available for merging at this time")
	}

	// ------------------------
	// Cross-Shard Transfer
	// ------------------------
	fmt.Println("\n--- Simulating cross-shard transfer ---")

	// Ensure we have at least 2 shards for transfer
	if len(manager.Shards) < 2 {
		fmt.Println("Need at least 2 shards for cross-shard transfer")
		return
	}

	// Perform a cross-shard transfer
	if err := sync.CrossShardTransfer(manager, 0, 1, 0); err != nil {
		fmt.Printf("Error in cross-shard transfer: %v\n", err)
	}

	// Show final shard states
	fmt.Println("\nFinal Shard States:")
	manager.PrintShards()

	// ------------------------
	// AMQ Filter Membership Check
	// ------------------------
	fmt.Println("\nChecking AMQ Filter membership...")
	amq := verification.NewAMQFilter()

	// Insert some items into AMQ filter
	amq.Add([]byte("Alice"))
	amq.Add([]byte("Bob"))

	// Query items in the AMQ filter
	fmt.Println("Is 'Alice' possibly in the set?", amq.PossiblyContains([]byte("Alice")))
	fmt.Println("Is 'Charlie' possibly in the set?", amq.PossiblyContains([]byte("Charlie")))

	// ------------------------
	// Additional Splitting Case
	// ------------------------
	fmt.Println("\n--- Additional Case: Forcing Split ---")

	for i := 15; i < 30; i++ { // More transactions to force a split
		tx := []byte(fmt.Sprintf("ExtraUser%d -> ExtraUser%d: %d", i, i+1, i*10))
		if err := manager.AddTransaction(tx); err != nil {
			fmt.Printf("Error adding transaction: %v\n", err)
		}
	}

	fmt.Println("\nShard States AFTER adding more transactions (to force split):")
	manager.PrintShards()

	if manager.ShouldSplit() {
		fmt.Println("\nShard splitting triggered (Extra Case).")
		if err := manager.SplitShard(); err != nil {
			fmt.Printf("Error splitting shard: %v\n", err)
		}
		manager.PrintShards()
	}

	// ------------------------
	// Additional Merging Case
	// ------------------------
	fmt.Println("\n--- Additional Case: Forcing Merge ---")

	// Simulate load reduction by clearing out some shards manually
	manager.ForceReduceLoad()

	if manager.ShouldMerge() {
		fmt.Println("\nShard merging triggered (Extra Case).")
		if err := manager.MergeShards(); err != nil {
			fmt.Printf("Error merging shards: %v\n", err)
		}
		manager.PrintShards()
	}
	// ------------------------
	// Additional BFT Operations
	// ------------------------

	// Initialize a new Node with reputation system from the BFT package

	node := bft.NewNode("Node1")
	node.UpdateReputation(true) // Simulate success

	// Print out node's reputation
	fmt.Printf("Node Reputation: %.2f\n", node.GetNodeReputation())

	// Check consensus threshold
	fmt.Printf("Consensus Threshold: %.2f\n", node.ConsensusThreshold())

	// Cryptographic validation
	priv, pub, err := bft.GenerateKeyPair(2048)
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return
	}

	message := []byte("Hello Blockchain")
	signature, err := bft.SignMessage(priv, message)
	if err != nil {
		fmt.Println("Error signing message:", err)
		return
	}

	err = bft.VerifySignature(pub, message, signature)
	if err != nil {
		fmt.Println("Signature verification failed:", err)
	} else {
		fmt.Println("Signature verified successfully.")
	}

	// VRF Leader Election
	randomData, vrfValue := bft.GenerateVRF()
	fmt.Printf("VRF Random Data: %x\n", randomData)
	fmt.Printf("VRF Leader Value: %x\n", vrfValue)

	// MPC Computation example
	inputs := []*big.Int{big.NewInt(5), big.NewInt(10)}
	result := bft.ComputeSum(inputs)
	fmt.Printf("MPC Computed Sum: %s\n", result.String())
}

// calculatePartitionRisk calculates the risk of network partition based on telemetry
func calculatePartitionRisk(telemetry core.NetworkTelemetry) float64 {
	// Simple risk calculation based on latency and packet loss
	latencyRisk := float64(telemetry.Latency) / 1000.0 // Normalize to 0-1 range
	packetLossRisk := telemetry.PacketLoss
	throughputRisk := 1.0 - (telemetry.Throughput / 1000.0) // Normalize to 0-1 range

	// Weighted average of risk factors
	risk := (latencyRisk*0.4 + packetLossRisk*0.4 + throughputRisk*0.2) * 100
	if risk > 100 {
		risk = 100
	}
	return risk
}

// getConsistencyLevelName returns a human-readable name for the consistency level
func getConsistencyLevelName(level int) string {
	switch level {
	case 1:
		return "eventual"
	case 2:
		return "causal"
	case 3:
		return "sequential"
	case 4:
		return "linearizable"
	case 5:
		return "strict"
	default:
		return "unknown"
	}
}
