/*
	Package blockchain does all the block chain work
	Adding blocks, constructing the block chain,
	Signing the blocks
*/
package blockchain

// BlockChain the chain(slice) containing the blocks
type BlockChain struct {
	Blocks []*Block
}

// The block
// Contains bytes of the hash, data and the previous block hash
type Block struct {
	Hash     []byte // Block hash
	Data     []byte // Block data
	PrevHash []byte // Previous Block hash
	Nonce    int    // Nonce that qualifies the target
}

// AddBlock adds a new block to the chain
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)

}

// CreateBlock creates a new block
func CreateBlock(data string, prevHash []byte) *Block {
	b := &Block{Data: []byte(data), Hash: []byte{}, PrevHash: prevHash, Nonce: 0}
	p := NewProof(b) // Generate new proof of work
	n, h := p.Run()  // Returns the block hash and nonce

	b.Hash = h[:]
	b.Nonce = n
	return b
}

// Genesis creates the first block in the blockchain
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

//InitBlockChain starts the blockchain system
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
