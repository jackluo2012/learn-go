package blockschain

import (
	"fmt"
	"log"
)

/**
 * 将一个钱包地址加入到 区 块 链
 */
func (cli *CLI) createBlockchain(address string) {

	//检查 地址是否合法
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	// 创建 一个创世链，并奖励50 ,并且把它写入区块链中
	bc := CreateBlockchain(address)

	defer bc.Db.Close()

	UTXOSet := UTXOSet{bc}

	UTXOSet.Reindex()

	fmt.Println("Done!")
}
