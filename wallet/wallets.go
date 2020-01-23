package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "./tmp/blocks/wallets.data"

// Wallets struct
type Wallets struct {
	Wallets map[string]*Wallet
}

// CreateWallets creates and returns wallet
func CreateWallets() *Wallets {
	wallets := &Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	wallets.LoadFile()

	return wallets
}

// AddWallet creates and adds a wallet to a user's address
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

// GetWallet returns the wallet details of the address
func (ws Wallets) GetWallet(addr string) Wallet {
	return *ws.Wallets[addr]
}

// GetAllAddresses returns all the addresses in the wallet
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// Opens and loads the wallet file
func (ws *Wallets) LoadFile() {
	var wallets Wallets

	// If no data has been saved
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		_, _ = os.Create(walletFile)
		return // exit the function without modifying the wallet struct
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	Handle(err)

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)

	// If no wallets exists/nil
	if len(wallets.Wallets) == 0 {
		return
	}

	ws.Wallets = wallets.Wallets

}

// SaveFile saves the wallets to a file
func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	Handle(err)

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	Handle(err)

}

// Handle takes the error and prints it out
func Handle(err error) {
	if err != nil {
		log.Panic(err)

	}
}
