package blockschain

import (
	"fmt"
	"log"
)
/**
 * 创建一个连接
 */
func (cli *CLI) createBlockchain(address string) {

	if !ValidateAddress(address){
		log.Panic("ERROR: Address is not valid")
	}

	bc := CreateBlockchain(address)
	bc.Db.Close()
	fmt.Println("Done!")
}