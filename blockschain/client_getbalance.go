package blockschain

import (
	"fmt"
	"log"
)

func (cli *CLI) getBalance(address, nodeID string) {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := NewBlockchain(nodeID)

	utxoSet := UTXOSet{bc}
	defer bc.Db.Close()

	balance := 0

	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	//*
	//log.Print(utxoSet)
	UTXOs := utxoSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}
	//*/
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
