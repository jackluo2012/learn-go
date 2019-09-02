package blockschain

import (
	"fmt"
	"log"
)

func (cli *CLI) startNode(nodeID, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if ValidateAddress(minerAddress) {
			fmt.Println("Mining is on.Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong niner address!")
		}
	}
	StartServer(nodeID, minerAddress)
}
