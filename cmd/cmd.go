package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sheghun/blockchain/blockchain"
	"github.com/sheghun/blockchain/wallet"
	"log"
	"os"
	"runtime"
	"strconv"
)

// Cmd struct for handling command line related tasks
type Cmd struct {
	blockchain *blockchain.BlockChain
}

// validateAddress validates the supplied address and exit
// the runtime
func (cli Cmd) validateAddress(address string) {
	if wallet.ValidateAddress(address) == false {
		log.Printf("\n\n\n\n ---------- Address is not valid -------------- \n\n\n\n")
		return
	}
}

// validate checks the cmd supplied arguments
func (cli *Cmd) validate() error {
	if len(os.Args) < 2 {
		cli.printUsage()
		return errors.New("add a command")
	}
	return nil
}

// printChain iterates and prints all the blocks in the database
func (cli *Cmd) printChain() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		p := blockchain.NewProof(block)
		fmt.Printf("Proof of work: %s\n\n", strconv.FormatBool(p.Validate()))

		for _, tx := range block.Transactions {
			fmt.Println(tx.String())
		}
		fmt.Println()

		// Check if at last block
		if len(block.PrevHash) == 0 {
			break // Exit functions
		}
	}
}

// createBlockchain creates a new blockchain for an address
func (cli *Cmd) createBlockchain(address string) {
	cli.validateAddress(address)

	chain := blockchain.InitBlockChain(address)
	defer chain.Database.Close()

	fmt.Printf("\n\n\n\n ----- Blockchain was created and 100 coins was transfered to %s in coinbase transaction \n\n\n\n", address)
}

func (cli *Cmd) getBalance(address string) {
	cli.validateAddress(address)

	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	balance := 0
	pubKeyHash, _, _ := wallet.Base58Decode([]byte(address))

	UTXOs := chain.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("\n\n\n\n ------------ Balance of %s: %d ----------------- \n\n\n\n", address, balance)
}

// listAddresses print out all the list to the cmd
func (cli *Cmd) listAddresses() {
	wallets := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	if len(addresses) == 0 {
		fmt.Println()
		fmt.Println()
		fmt.Println("--------------------")
		fmt.Println("------- No Addresses found")
		fmt.Println("--------------------")
		fmt.Println()
		return
	}

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *Cmd) createWallet() {
	wallets := wallet.CreateWallets()
	address := wallets.AddWallet()

	wallets.SaveFile()

	fmt.Printf("\n\n\n\n\n ------ New address is %s -------\n\n\n\n\n", address)
}

// Sends a transaction from one address to another
func (cli *Cmd) send(from, to string, amount int) {
	cli.validateAddress(to)
	cli.validateAddress(from)

	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Printf("\n\n\n\n -------- Transactions successful --------- \n\n\n\n")
}

// printUsage prints the command line possible commands
func (cli *Cmd) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getBalance -address ADDRESS - get the balance of the address")
	fmt.Println(" createBlockchain -address ADDRESS creates a blockchain")
	fmt.Println(" printChain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send the amount ")
	fmt.Println(" listAddresses - Lists the addresses in our wallet file")
	fmt.Println(" createWallet - creates a new wallet")
}

// Run takes in the command line inputs
func (cli *Cmd) Run() {

	if err := cli.validate(); err != nil {
		blockchain.Handle(err)
		return // Exit the function
	}

	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listAddresses", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address balance to retrieve")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "Address to create blockchain for")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	// Listen for the command flags
	switch os.Args[1] {
	case "getBalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "createBlockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "listAddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "createWallet":
		err := createWalletCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		return // Exit the functions
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
			return
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
			return
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == 0 {
			sendCmd.Usage()
			runtime.Goexit()
			return
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
		runtime.Goexit()
		return
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
		runtime.Goexit()
		return
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
		runtime.Goexit()
		return
	}

}
