package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/sheghun/blockchain/wallet"
	"log"
	"math/big"
	"strings"
)

// Transaction struct
type Transaction struct {
	ID      []byte // ID of the transaction
	Inputs  []TxInput
	Outputs []TxOutput
}

// SetID derives axnd sets the transaction hash
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

	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTxOutput(100, to)

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx.SetID()

	return &tx
}

// NewTransaction initiates a new transaction
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	var txn *Transaction

	wallets := wallet.CreateWallets()

	w := wallets.GetWallet(from)
	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc, validOutputs := chain.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		Handle(errors.New("not enough funds"))
	}

	for txid, outs := range validOutputs {
		txid, err := hex.DecodeString(txid) // Decode id back to a string
		Handle(err)

		for _, out := range outs {
			input := TxInput{txid, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, *NewTxOutput(amount, to))

	if acc > amount {
		outputs = append(outputs, *NewTxOutput(acc-amount, from))
	}

	txn = &Transaction{nil, inputs, outputs}
	txn.ID = txn.Hash()
	chain.SignTransaction(txn, w.PrivateKey)

	return txn
}

func (t *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *t
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// IsCoinbase checks if the current transactions is a coinbase transaction
func (t *Transaction) IsCoinbase() bool {
	return len(t.Inputs) == 1 && len(t.Inputs[0].ID) == 0 && t.Inputs[0].Out == -1
}

func (t *Transaction) Sign(privKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if t.IsCoinbase() {
		return
	}

	for _, in := range t.Inputs {
		if prevTxs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic(errors.New("previous transaction does not exists"))
		}
	}

	tCopy := t.TrimmedCopy()

	// Loop through and sign all inputs
	for inId, in := range tCopy.Inputs {
		prevTx := prevTxs[hex.EncodeToString(in.ID)]
		tCopy.Inputs[inId].Signature = nil
		tCopy.Inputs[inId].PubKey = prevTx.Outputs[in.Out].PubKeyHash
		tCopy.ID = tCopy.Hash()
		tCopy.Inputs[inId].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, tCopy.ID)
		Handle(err)

		signature := append(r.Bytes(), s.Bytes()...)

		t.Inputs[inId].Signature = signature

	}
}

// Serialize encodes and returns the byte representation of the transaction
func (t Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(t)
	Handle(err)

	return encoded.Bytes()

}

func (t *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range t.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, in.PubKey})
	}

	for _, out := range t.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{t.ID, inputs, outputs}

	return txCopy
}

// Verify checks if the transaction input signatures are valid
func (t *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if t.IsCoinbase() {
		return true
	}

	for _, in := range t.Inputs {
		if prevTxs[hex.EncodeToString(in.ID)].ID == nil {
			Handle(errors.New("previous transaction does not exists"))
		}
	}

	tCopy := t.TrimmedCopy()
	curve := elliptic.P256()

	// Loop through and sign all inputs
	for inId, in := range tCopy.Inputs {
		prevTx := prevTxs[hex.EncodeToString(in.ID)]
		tCopy.Inputs[inId].Signature = nil
		tCopy.Inputs[inId].PubKey = prevTx.Outputs[in.Out].PubKeyHash
		tCopy.ID = tCopy.Hash()
		tCopy.Inputs[inId].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		// Create a new public key
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

		// return verification result either true/false
		return ecdsa.Verify(&rawPubKey, tCopy.ID, &r, &s)

	}

	return true
}

// String converts transaction to string
func (t *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("------ Transaction %x:", t.ID))


	// Write out inputs strings
	for i, input := range t.Inputs {
		lines = append(lines, fmt.Sprintf("		Input %d:", i))
		lines = append(lines, fmt.Sprintf("			TXID:		%x", input.ID))
		lines = append(lines, fmt.Sprintf("			Out:		%d", input.Out))
		lines = append(lines, fmt.Sprintf("			Signature:	%x", input.Signature))
		lines = append(lines, fmt.Sprintf("			PubKey:		%x", input.PubKey))

	}

	for i, output := range t.Outputs {
		lines = append(lines, fmt.Sprintf("		Output %d:", i))
		lines = append(lines, fmt.Sprintf("			Value:		%d", output.Value))
		lines = append(lines, fmt.Sprintf("			Script: 	%x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
