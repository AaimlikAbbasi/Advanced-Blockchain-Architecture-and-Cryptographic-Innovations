// /verification/amq_filter.go
package verification

import (
	"crypto/sha256"
	"encoding/hex"
	"hash/fnv"
)

type AMQFilter struct {
	items   map[string]bool
	size    int
	buckets []bool
}

func NewAMQFilter() *AMQFilter {
	return &AMQFilter{
		items:   make(map[string]bool),
		size:    1000, // Default size
		buckets: make([]bool, 1000),
	}
}

func (f *AMQFilter) Add(data []byte) {
	h := sha256.Sum256(data)
	f.items[hex.EncodeToString(h[:])] = true
}

func (f *AMQFilter) Contains(data []byte) bool {
	h := sha256.Sum256(data)
	_, ok := f.items[hex.EncodeToString(h[:])]
	return ok
}

func (f *AMQFilter) Exists(item []byte) bool {
	h := fnv.New32a()
	h.Write(item)
	index := int(h.Sum32()) % f.size
	return f.buckets[index]
}

func (f *AMQFilter) PossiblyContains(data []byte) bool {
	return f.Contains(data) || f.Exists(data)
}
