package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
	"time"
)

// Difficulty level
const Difficulty = 12

// ProofOfWork struct
type ProofOfWork struct {
	Block  *Block   // Block to verify
	Target *big.Int // Target to proof
}

// NewProofOfWork cares a new ProofOfWork instance
func NewProof(b *Block) *ProofOfWork {
	t := big.NewInt(1)
	t.Lsh(t, uint(256-Difficulty))

	p := &ProofOfWork{b, t}
	return p
}

// InitData returns the bytes of the proof of work to hash
// Takes the nonce as an argument
// Combines the bytes of the block(prevHash + blockData), the nonce and returns the byte
func (p *ProofOfWork) InitData(n int) []byte {
	b := bytes.Join(
		[][]byte{
			p.Block.PrevHash,
			p.Block.HashTransactions(),
			ToHex(int64(n)),
			ToHex(int64(Difficulty))},
		[]byte{},
	)

	return b
}

// Run does the ProofOfWork calculation
// starts the nonce from zero(bytes) appends the nonce to the block bytes
// hashes the bytes check if it's lower than the target
func (p *ProofOfWork) Run() (int, []byte) {
	nonce := 0 // nonce to increment

	c := make(chan struct {
		nonce int
		hash  []byte
	})

	cl := make(chan bool)

	for nonce < math.MaxInt64 {
		select {
		// Check if nonce has been found and sent to channel
		case d := <-c:
			// Close the other channels
			close(cl)
			return d.nonce, d.hash[:]
		default:
			go calNonce(c, nonce, nonce+3000, p, cl)
			nonce += 3000
		}
		time.Sleep(100 * time.Millisecond)
	}
	// Wait for the nonce to be supplied
	d := <-c
	close(cl)
	return d.nonce, d.hash[:]
}

// Validate verifies the block by running the proof of work algorithm
// to check if the block none satisfies the target
func (p *ProofOfWork) Validate() bool {

	var bigH big.Int

	d := p.InitData(p.Block.Nonce)
	h := sha256.Sum256(d)

	bigH.SetBytes(h[:])

	return bigH.Cmp(p.Target) == -1

}

// ToHex serializes the supplied num into bytes and returns the byte
func ToHex(num int64) []byte {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.BigEndian, num)

	if err != nil {
		log.Panic(err)
	}

	return b.Bytes()
}

// Calculates the supplied nonce input and check if it satisfy the target
func calNonce(c chan struct {
	nonce int    // This is the nonce
	hash  []byte // The hash byte
}, n int, e int, p *ProofOfWork, closeC <-chan bool) {
	var bigH big.Int

	for n < e {
		select {
		// Check if channel is closed
		case _, cl := <-closeC:
			if !cl {
				// Close the goroutine
				break
			}
		// Run the default case
		default:
			d := p.InitData(n)    // proof of work byt es
			h := sha256.Sum256(d) // hash

			bigH.SetBytes(h[:]) // Read the bytes into the int

			if bigH.Cmp(p.Target) == -1 {
				// Send in the struct
				c <- struct {
					nonce int
					hash  []byte
				}{
					nonce: n,
					hash:  h[:],
				}

				break
			}

			n++ // Increment the nonce
		}

	}

}
