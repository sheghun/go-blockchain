package blockchain

import (
	"bytes"
	"github.com/sheghun/blockchain/wallet"
)

// Transaction output struct
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// Transaction input struct
type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

// NewTxOutput creates and returns a new utxo locked to the supplied address
func NewTxOutput(value int, addr string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(addr))

	return txo
}

// UsesKey check if the transaction input uses this key
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock locks the transactions output
func (out *TxOutput) Lock(address []byte) error {
	pubKeyHash, _, err := wallet.Base58Decode(address)
	out.PubKeyHash = pubKeyHash

	return err
}

// IsLockedWithKey checks if the utxo is locked with key
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
