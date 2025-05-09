package verification

import "crypto/sha256"

// Simple cryptographic accumulator
type Accumulator struct {
	Value []byte
}

func NewAccumulator() *Accumulator {
	return &Accumulator{Value: []byte{}}
}

func (a *Accumulator) Add(data []byte) {
	hash := sha256.Sum256(append(a.Value, data...))
	a.Value = hash[:]
}

func (a *Accumulator) GetValue() []byte {
	return a.Value
}
