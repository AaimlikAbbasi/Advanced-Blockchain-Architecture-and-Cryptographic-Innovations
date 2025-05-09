package bft

// Adaptive consensus threshold based on reputation
func (node *Node) ConsensusThreshold() float64 {
	return 1.0 - node.Reputation // Higher reputation reduces the threshold
}
