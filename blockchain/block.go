/*
	Package blockchain does all the block chain work
	Adding blocks, constructing the block chain,
	Signing the blocks
*/
package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
)

const (
	dbPath = "./tmp/blocks"
)

// BlockChain the chain(slice) containing the blocks
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// Iterator loops through the database and retrieves all blocks
type Iterator struct {
	currentHash []byte
	Database    *badger.DB
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
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			err = txn.Set([]byte("lh"), newBlock.Hash)

			chain.LastHash = newBlock.Hash
			return err
		})

		Handle(err)

		err = txn.Set(newBlock.Hash, newBlock.Serialize())
		return err
	})
	Handle(err)

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
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// Try getting the last has from the database
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")

			gen := Genesis()

			fmt.Println("Genesis proved")

			err = txn.Set(gen.Hash, gen.Serialize())
			Handle(err)

			err = txn.Set([]byte("lh"), gen.Hash)

			lastHash = gen.Hash

			return err
		}

		// Get the last block has
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})
	Handle(err)
	// Return blockchain instance
	return &BlockChain{lastHash, db}
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

// Returns the iterator struct to iterate the blocks in the database
func (chain *BlockChain) Iterator() *Iterator {
	iter := &Iterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *Iterator) Next() *Block {
	var block *Block
	// Retrieve from the database
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.currentHash)
		Handle(err)

		// Get the stored block bytes
		err = item.Value(func(val []byte) error {
			// Deserialize into block struct
			block = Deserialize(val)
			return nil
		})
		return err
	})
	Handle(err)

	iter.currentHash = block.PrevHash

	return block
}

// Handle takes the error and prints it out
func Handle(err error) {
	if err != nil {
		log.Panic(err)

	}
}
