package blockschain

import (
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, amount int) {
	if !ValidateAddress(from){
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to){
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain(from)
	defer bc.Db.Close()

	//生成交易
	tx := NewUTXOTransaction(from, to, amount, bc)
	//进挖矿产生区块
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")

}
