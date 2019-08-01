package blockschain

import (
	"fmt"
)

func (cli *CLI) createWallet()  {
	wallets, _ :=NewWallets()
	address :=wallets.CresteWallet()

	wallets.SaveToFile()

	fmt.Printf("You new address: %s\n",address)
}