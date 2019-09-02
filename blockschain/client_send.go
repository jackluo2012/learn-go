package blockschain

import (
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, amount int, nodeID string, mineNow bool) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain(nodeID)

	utxoset := UTXOSet{bc}

	defer bc.Db.Close()

	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	//生成交易 ，并进行一系列的检查
	tx := NewUTXOTransaction(&wallet, to, amount, &utxoset)
	if mineNow {
		//产生输入交易
		cbTx := NewCoinbaseTX(from, "")

		txs := []*Transaction{cbTx, tx}

		//进挖矿产生区块
		newBlock := bc.MineBlock(txs)

		//更新 UTXO
		utxoset.Update(newBlock)
	} else {
		sendTx(knownNodes[0], tx)
	}

	fmt.Println("Success!")

}
