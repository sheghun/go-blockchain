package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
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

func DBExits() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

//InitBlockChain starts the blockchain system
func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBExits() {
		fmt.Println("Database already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)

		gen := Genesis(cbtx)

		fmt.Println("Genesis proved")

		err = txn.Set(gen.Hash, gen.Serialize())
		Handle(err)

		err = txn.Set([]byte("lh"), gen.Hash)

		lastHash = gen.Hash

		return err
	})
	Handle(err)
	// Return blockchain instance
	return &BlockChain{lastHash, db}
}

// ContinueBlockChain retrieves the last hash ID on the database
func ContinueBlockChain() *BlockChain {
	if DBExits() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})
	Handle(err)
	chain := BlockChain{lastHash, db}

	return &chain
}

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTxos := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)

			// Outputs label
		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTxos[txId] != nil {
					for _, spentOut := range spentTxos[txId] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// Check if PubKey is address
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTxos[inTxID] = append(spentTxos[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTxs
}

// FindUTXO returns the list of unspent transactions output
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs

}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unSpentTxs := chain.FindUnspentTransactions(address)
	accumalated := 0

Work:
	for _, tx := range unSpentTxs {
		txId := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumalated < out.Value {
				accumalated += out.Value
				unspentOuts[txId] = append(unspentOuts[txId], outIdx)

				if accumalated >= amount {
					break Work
				}
			}
		}
	}

	return accumalated, unspentOuts
}

// AddBlock adds a new block to the chain
func (chain *BlockChain) AddBlock(txn []*Transaction) {
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

	newBlock := CreateBlock(txn, lastHash)

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
