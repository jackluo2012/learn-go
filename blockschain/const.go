package blockschain

import (
	"math"
)

//设置挖矿难度  前24位为 0
const (
	TargetBits          = 24
	MaxNonce            = math.MaxInt64
	DbFile              = "blockchain_%s.db"
	BlocksBucket        = "blocks"
	GenesisCoinbaseData = "创世区块的创立"
)
