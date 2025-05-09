package auth

import "math/rand"

func GetTrustScore(nodeID string) float64 {
	// Random trust score generator (demo only)
	return 0.5 + rand.Float64()/2 // returns between 0.5 and 1.0
}
