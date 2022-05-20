package business

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// CLI Command-Line Interface
type CLI struct {
}

const createChainName = "create_blockchain"
const addBlockName = "add_block"
const printChainName = "print_chain"

func (cli *CLI) Run() {
	validateArgs()

	// configuring cmd
	createChainCmd := flag.NewFlagSet(createChainName, flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet(addBlockName, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(printChainName, flag.ExitOnError)

	createDataPtr := createChainCmd.String("data", "Default", "Genesis Block Data")
	addBlockDataPtr := addBlockCmd.String("data", "Default", "Add New Block Data")

	command := os.Args[1]
	switch command {
	case createChainName:
		err := createChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case addBlockName:
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case printChainName:
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if createChainCmd.Parsed() {
		if *createDataPtr == "" {
			printUsage()
			os.Exit(1)
		}
		cli.createGenesisBlockchain(*createDataPtr)
	}

	if addBlockCmd.Parsed() {
		if *addBlockDataPtr == "" {
			printUsage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockDataPtr)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func validateArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\t" + createChainName + " -data DATA -- Create Genesis Block")
	fmt.Println("\t" + addBlockName + " -data Data -- Add Transaction Data")
	fmt.Println("\t" + printChainName + " -- Print Information")
}

func (cli *CLI) createGenesisBlockchain(data string) {
	CreateBlockChain(data)
}

func (cli *CLI) addBlock(data string) {
	bc := GetBlockChain()
	if bc == nil {
		fmt.Printf("No genesis block. Failed.")
		os.Exit(1)
	}
	bc.AddNewBlock(data)
	defer bc.BlockDB.Close()
}

func (cli *CLI) printChain() {
	bc := GetBlockChain()
	if bc == nil {
		fmt.Printf("No block to print.")
		os.Exit(1)
	}
	bc.PrintChain()
}
