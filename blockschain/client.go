package blockschain

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// 新建一个 Blockchain 客户端操作
type CLI struct{}

func (cli *CLI) Run() {
	//检测参数
	cli.validateArgs()

	//创建子命令
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

func (cli *CLI) getBalance(address string) {
	bc := NewBlockchain(address)
	defer bc.Db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockchain(from)
	defer bc.Db.Close()

	//生成交易
	tx := NewUTXOTransaction(from, to, amount, bc)
	//进挖矿产生区块
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bc := NewBlockchain("")
	defer bc.Db.Close()
	bci := bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PreBlockHash)
		fmt.Printf("Hash. hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PreBlockHash) == 0 {
			break
		}
	}
}

/**
 * 创建一个连接
 */
func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockchain(address)
	bc.Db.Close()
	fmt.Println("Done!")
}

/**
 * 打印提示
 */
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

/**
 * 进行检查
 */
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
