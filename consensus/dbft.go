package consensus

func dBFTConsensus(randomInput string) string {
	// Mocked dBFT-style voting (simulated)
	// In real use, validators would sign + vote on block
	quorum := 3
	votes := 0

	// Random logic to simulate validator agreement
	if len(randomInput)%2 == 0 {
		votes = 3
	} else {
		votes = 2
	}

	if votes >= quorum {
		return "Block Approved by dBFT"
	}
	return "Block Rejected by dBFT"
}
