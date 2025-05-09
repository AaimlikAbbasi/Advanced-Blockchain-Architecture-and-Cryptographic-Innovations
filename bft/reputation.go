package bft

// Node structure for tracking node reputation
type Node struct {
	ID               string
	Reputation       float64 // Reputation score
	FailedBlocks     int     // Count of failed blocks by this node
	SuccessfulBlocks int     // Count of successful blocks by this node
}

// UpdateReputation adjusts reputation based on node performance
func (node *Node) UpdateReputation(success bool) {
	if success {
		node.SuccessfulBlocks++
		node.Reputation += 0.1 // Increase reputation for successful behavior
	} else {
		node.FailedBlocks++
		node.Reputation -= 0.2 // Decrease reputation for failed behavior
	}

	// Ensure reputation stays within bounds [0, 1]
	if node.Reputation < 0 {
		node.Reputation = 0
	}
	if node.Reputation > 1 {
		node.Reputation = 1
	}
}

// GetNodeReputation returns the current reputation of the node
func (node *Node) GetNodeReputation() float64 {
	return node.Reputation
}

func NewNode(id string) *Node {
	return &Node{
		ID:         id,
		Reputation: 0.5, // Default neutral reputation
	}
}
