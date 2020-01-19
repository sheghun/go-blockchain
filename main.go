package main

import (
	"fmt"
	"github.com/sheghun/blockchain/blockchain"
	"strconv"
)

func main() {
	// Init the chain
	chain := blockchain.InitBlockChain()

	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block after Genesis")
	chain.AddBlock("Third Block after Genesis")

	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		p := blockchain.NewProof(block)
		fmt.Printf("Proof of work: %s\n\n", strconv.FormatBool(p.Validate()))

	}

}
