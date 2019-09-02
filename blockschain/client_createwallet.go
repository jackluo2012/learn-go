package blockschain

import (
	"fmt"
)

func (cli *CLI) createWallet(nodeID string) {
	wallets, _ := NewWallets(nodeID)
	address := wallets.CresteWallet()

	wallets.SaveToFile(nodeID)

	fmt.Printf("You new address: %s\n", address)
}
