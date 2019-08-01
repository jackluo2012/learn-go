package blockschain

import (
	"fmt"
	"strconv"
)

func (cli *CLI) printChain() {
	bc := NewBlockchain("")
	defer bc.Db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("========Block %x ===========\n", block.Hash)
		fmt.Printf("Prev. block: %x\n", block.PreBlockHash)
		pow := NewProofOfWork(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println("\n\n")

		if len(block.PreBlockHash) == 0 {
			break
		}
	}
}
