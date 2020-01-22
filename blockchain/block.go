/*
	Package blockchain does all the block chain work
	Adding blocks, constructing the block chain,
	Signing the blocks
*/
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// The block
// Contains bytes of the hash, data and the previous block hash
type Block struct {
	Hash         []byte         // Block hash
	Transactions []*Transaction // Block data
	PrevHash     []byte         // Previous Block hash
	Nonce        int            // Nonce that qualifies the target
}

// HashTransactions hashes and returns the sha256 hash of all the transaction
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// CreateBlock creates a new block
func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	b := &Block{Transactions: txs, Hash: []byte{}, PrevHash: prevHash, Nonce: 0}
	p := NewProof(b) // Generate new proof of work
	n, h := p.Run()  // Returns the block hash and nonce

	b.Hash = h[:]
	b.Nonce = n
	return b
}

// Genesis creates the first block in the blockchain
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

// Serialize converts the block into a gob byte
func (b *Block) Serialize() []byte {
	buffer := new(bytes.Buffer)

	encoder := gob.NewEncoder(buffer)

	err := encoder.Encode(b)
	Handle(err)

	return buffer.Bytes()
}

// Deserialize converts the supplied byte into a block
func Deserialize(data []byte) *Block {
	var b Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&b)
	Handle(err)

	return &b

}

// Handle takes the error and prints it out
func Handle(err error) {
	if err != nil {
		log.Panic(err)

	}
}
