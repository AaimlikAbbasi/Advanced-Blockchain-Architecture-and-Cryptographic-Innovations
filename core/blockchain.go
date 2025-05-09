package core

import (
	"fmt"
)

type Blockchain struct {
	Blocks []Block
}

func NewBlockchain() *Blockchain {
	genesis := CreateBlock(0, []*Transaction{}, "0")
	return &Blockchain{
		Blocks: []Block{genesis},
	}
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := CreateBlock(lastBlock.Index+1, transactions, lastBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
	fmt.Println("Block added:", newBlock.Index)
}

func (bc *Blockchain) GetMerkleRoot() []byte {
	if len(bc.Blocks) == 0 {
		return []byte{}
	}
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	return []byte(lastBlock.MerkleRoot)
}

func (bc *Blockchain) VerifyMerkleProof(rootHash []byte) bool {
	if len(bc.Blocks) == 0 {
		return false
	}
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	return string(rootHash) == lastBlock.MerkleRoot
}
