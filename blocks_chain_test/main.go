package main

import (
	"gopcp.v2/chapter7/blockschain"
)

func main() {
	//初始化一个 ClI
	cli := blockschain.CLI{}
	cli.Run()
}
