package merkle

import (
	"crypto/sha256"
	"fmt"
	"sort"
)

type Node struct {
	Left  *Node
	Right *Node
	Hash  []byte
}

type MerkleTree struct {
	Root   *Node
	Leaves []*Node
}

type ProofStep struct {
	Hash   []byte
	IsLeft bool // true if this hash should be on the left when combining
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	if len(data) == 0 {
		return &MerkleTree{Root: &Node{Hash: hash([]byte{})}}
	}

	// Sort the data to ensure consistent ordering
	sortedData := make([][]byte, len(data))
	copy(sortedData, data)
	sort.Slice(sortedData, func(i, j int) bool {
		return string(sortedData[i]) < string(sortedData[j])
	})

	var leaves []*Node
	for _, d := range sortedData {
		leaves = append(leaves, &Node{Hash: hash(d)})
	}

	root := buildTree(leaves)
	return &MerkleTree{Root: root, Leaves: leaves}
}

func buildTree(nodes []*Node) *Node {
	if len(nodes) == 1 {
		return nodes[0]
	}

	var newLevel []*Node
	for i := 0; i < len(nodes); i += 2 {
		if i+1 < len(nodes) {
			// Use nodes in their original order
			left, right := nodes[i], nodes[i+1]

			// Always combine hashes in the same order (left.Hash followed by right.Hash)
			combinedHash := make([]byte, len(left.Hash)+len(right.Hash))
			copy(combinedHash, left.Hash)
			copy(combinedHash[len(left.Hash):], right.Hash)

			newNode := &Node{
				Left:  left,
				Right: right,
				Hash:  hash(combinedHash),
			}
			newLevel = append(newLevel, newNode)
		} else {
			// For odd number of nodes, duplicate the last node
			combinedHash := make([]byte, len(nodes[i].Hash)*2)
			copy(combinedHash, nodes[i].Hash)
			copy(combinedHash[len(nodes[i].Hash):], nodes[i].Hash)

			newNode := &Node{
				Left:  nodes[i],
				Right: nodes[i],
				Hash:  hash(combinedHash),
			}
			newLevel = append(newLevel, newNode)
		}
	}

	return buildTree(newLevel)
}

func hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func findParent(root, target *Node) *Node {
	if root == nil || root == target {
		return nil
	}

	if root.Left == target || root.Right == target {
		return root
	}

	if parent := findParent(root.Left, target); parent != nil {
		return parent
	}

	return findParent(root.Right, target)
}

func (mt *MerkleTree) GenerateProof(data []byte) ([]ProofStep, error) {
	dataHash := hash(data)
	var targetNode *Node
	var targetIndex int
	for i, leaf := range mt.Leaves {
		if string(leaf.Hash) == string(dataHash) {
			targetNode = leaf
			targetIndex = i
			break
		}
	}

	if targetNode == nil {
		return nil, fmt.Errorf("data not found in tree")
	}

	var proof []ProofStep
	current := targetNode
	parent := findParent(mt.Root, current)

	for parent != nil {
		// Determine if current node is left or right child
		isLeft := parent.Left == current

		if isLeft {
			// When current is left child, we need the right sibling
			if parent.Right != nil {
				proof = append(proof, ProofStep{Hash: parent.Right.Hash, IsLeft: false})
			} else {
				// This shouldn't happen with our balanced tree
				proof = append(proof, ProofStep{Hash: parent.Left.Hash, IsLeft: false})
			}
		} else {
			// When current is right child, we need the left sibling
			proof = append(proof, ProofStep{Hash: parent.Left.Hash, IsLeft: true})
		}
		current = parent
		parent = findParent(mt.Root, current)
	}

	// Print debug information about the proof
	fmt.Printf("Debug - Generated proof for transaction: %s (index: %d)\n", string(data), targetIndex)
	for i, step := range proof {
		fmt.Printf("Debug - Proof[%d]: %x (IsLeft: %v)\n", i, step.Hash, step.IsLeft)
	}

	return proof, nil
}

func VerifyProof(data []byte, proof []ProofStep, rootHash []byte) bool {
	currentHash := hash(data)
	fmt.Printf("Debug - Initial hash: %x\n", currentHash)

	for i, step := range proof {
		fmt.Printf("Debug - Step %d - Current hash: %x\n", i, currentHash)
		fmt.Printf("Debug - Step %d - Sibling hash: %x (IsLeft: %v)\n", i, step.Hash, step.IsLeft)

		// Combine hashes based on the order specified in the proof
		var combinedHash []byte
		if step.IsLeft {
			// If sibling is left, it goes first
			combinedHash = make([]byte, len(step.Hash)+len(currentHash))
			copy(combinedHash, step.Hash)
			copy(combinedHash[len(step.Hash):], currentHash)
		} else {
			// If sibling is right, current hash goes first
			combinedHash = make([]byte, len(currentHash)+len(step.Hash))
			copy(combinedHash, currentHash)
			copy(combinedHash[len(currentHash):], step.Hash)
		}
		currentHash = hash(combinedHash)

		fmt.Printf("Debug - Step %d - Combined hash: %x\n", i, currentHash)
	}

	fmt.Printf("Debug - Final hash: %x\n", currentHash)
	fmt.Printf("Debug - Expected root: %x\n", rootHash)

	return string(currentHash) == string(rootHash)
}
