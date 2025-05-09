package bft

import (
	"math/big"
)

// Compute the sum of private inputs from multiple parties
func ComputeSum(inputs []*big.Int) *big.Int {
	var sum big.Int
	for _, input := range inputs {
		sum.Add(&sum, input)
	}
	return &sum
}
