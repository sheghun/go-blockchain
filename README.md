# Go-BlockChain

I've been playing around with bitcoin and etheruem for the past few weeks, so I decided to build my
own blockchain protocol, basically just a test of knowledge to know what I've learnt so far

This a blockchain protocol written in Golang,
It is inspired by the bitcoin blockchain protocol, although not exactly like it.

The wallets addresses are encoded in the same base as the bitcoin protocol (Base58Check)
blocks are mined only when transactions takes place, it is currently limited as there can only
one transaction per block it's not a full fledged blockchain protocol
but a very simple version of it, it implements the basic ideas of the bitcoin protocol

##Usage
### Initialisation
When creating the first transaction or mining the first block we have to supply an address so
that the 100 coins reward for mining the genesis block is transferred to the address 

All addresses either displayed or supplied are in the Base58Check
  
 
To Initialize the blockchain we have to first generate an address to which we will transfer the coinbase
transaction reward to generate first wallet address in your terminal type
 
    go run main.go createWallet
 
copy the generated address and run the command 
    
    go run main.go createBlockchain -address <GENERATED ADDRESS e.g "1PfwwjUi97CHHbTQF4W1XzTRDSoaJbM3qG">
        
        
to create/initialise the blockchain, it also transfers the coinbase transaction of 100 coins to the address.

#### Commands
To get the address balance

    go run main.go getBalance -address <ADDRESS e.g "1PfwwjUi97CHHbTQF4W1XzTRDSoaJbM3qG">

To print all the transactions on the blockchain
    
    go run main.go printChain
    
To get a list of all the addresses

    go run main.go listAddresses
    
To transfer from one address to another
 
**Note** All addresses (i.e sender and receiver) should have been created and tied to a wallet using the createWallet command

    go run main.go send -from <SENDER_ADDRESS e.g "1DYLi62NLDQwkey8roEWAap5Xdm3zX7BHd"> -to <RECEIVER_ADDRESS e.g "14waQN7En5QJ6C2iSJukhVVKBhMmovsNWq"> -amount <AMOUNT e.g 30>
    