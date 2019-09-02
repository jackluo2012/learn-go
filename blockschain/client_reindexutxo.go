package blockschain

import "fmt"

func (cli *CLI) reindexUTXO(nodeID string) {
	bc := NewBlockchain(nodeID)
	utxoset := UTXOSet{bc}

	utxoset.Reindex()

	count := utxoset.CountTransactions()
	fmt.Print("Done! There are %d transactions in the UTXO set.\n", count)
}
