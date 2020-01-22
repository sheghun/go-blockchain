package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
)

// Transaction struct
type Transaction struct {
	ID      []byte // ID of the transaction
	Inputs  []TxInput
	Outputs []TxOutput
}

// SetID derives and sets the transaction hash
func (t *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(t)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	t.ID = hash[:]

}

// CoinbaseTx initiates the first transaction in the genesis block
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

// NewTransaction initiates a new transaction
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	var txn *Transaction

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		Handle(errors.New("not enough funds"))
	}

	for txid, outs := range validOutputs {
		txid, err := hex.DecodeString(txid) // Decode id back to a string
		Handle(err)

		for _, out := range outs {
			input := TxInput{txid, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	txn = &Transaction{nil, inputs, outputs}
	txn.SetID()

	return txn
}

// IsCoinbase checks if the current transactions is a coinbase transaction
func (t *Transaction) IsCoinbase() bool {
	return len(t.Inputs) == 1 && len(t.Inputs[0].ID) == 0 && t.Inputs[0].Out == -1
}
