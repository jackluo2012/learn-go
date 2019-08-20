package blockschain

import (
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, amount int) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain()

	UTXOSet := UTXOSet{bc}

	defer bc.Db.Close()

	//生成交易 ，并进行一系列的检查
	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)

	//产生输入交易
	cbTx := NewCoinbaseTX(from, "")



	txs := []*Transaction{cbTx, tx}

	//进挖矿产生区块
	newBlock := bc.MineBlock(txs)

	//更新 UTXO
	UTXOSet.Update(newBlock)

	fmt.Println("Success!")

}
