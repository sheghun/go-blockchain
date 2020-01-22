package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sheghun/blockchain/blockchain"
	"os"
	"strconv"
)

// CommandLine struct for handling command line related tasks
type CommandLine struct {
	blockchain *blockchain.BlockChain
}

// printUsage prints the command line possible commands
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println("print - Prints the blocks in the chain")
}

// validate checks the cmd supplied arguments
func (cli *CommandLine) validate() error {
	if len(os.Args) < 2 {
		cli.printUsage()
		return errors.New("Add a command")
	}
	return nil
}

// addBlock adds a block to block chain
func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data) // Add block to chain
	fmt.Println("Block has been added")
}

// printChain iterates and prints all the blocks in the database
func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		p := blockchain.NewProof(block)
		fmt.Printf("Proof of work: %s\n\n", strconv.FormatBool(p.Validate()))

		// Check if at last block
		if len(block.PrevHash) == 0 {
			break // Exit functions
		}
	}
}

// run takes in the command line inputs
func (cli *CommandLine) run() {

	err := cli.validate()
	if err != nil {
		blockchain.Handle(err)
		return // Exit the function
	}

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block Data")

	// Listen for the command flags
	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		return // Exit the functions
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			cli.printUsage()
			return // Exit the function
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
		return
	}
}

func main() {
	// Init the chain
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := CommandLine{chain}

	cli.run() // Run the cli

}
