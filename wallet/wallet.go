package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/sheghun/blockchain/blockchain"
	"golang.org/x/crypto/ripemd160"
)

const (
	checkSumlength = 4
	version        = byte(0x00)
)

// Wallet struct contains the private key and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// Address generates an address for the wallet
func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)
	return address
}

// NewKeyPair generates
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	blockchain.Handle(err)

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

// MakeWallet creates a new wallet
func MakeWallet() *Wallet {
	private, public := NewKeyPair()

	wallet := &Wallet{private, public}

	return wallet
}

// PublicKeyHash hashes and returns the hash of the public key
func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	ripe := ripemd160.New()

	_, err := ripe.Write(pubHash[:])
	blockchain.Handle(err)

	publicRipMd := ripe.Sum(nil)

	return publicRipMd

}

// Returns the checksum of the payload
func Checksum(payload []byte) []byte {
	// Hash twice
	hash := sha256.Sum256(payload)
	hash = sha256.Sum256(hash[:])

	return hash[:checkSumlength]

}
