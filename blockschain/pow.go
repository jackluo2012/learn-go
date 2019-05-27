package blockschain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//加入工作量证明  pow
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	//进行位移操作
	target.Lsh(target, uint(256-TargetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

/**
 * 将 block 各个字段 和 nonce(counter 值)作为输入,计算得到的 hash 值作为输出
 */
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PreBlockHash,
			pow.block.HashTransactions(),//加入交易的改变
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(TargetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{})
	return data
}

/**
 * Pow 工作量计算 核心算法 :
 */

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0 //初始化为0

	fmt.Println("Mining a new block")
	for nonce < MaxNonce {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		//将  byte 转换成 int
		hashInt.SetBytes(hash[:])
		//比较x和y的大小。如x < y返回-1；如x > y返回+1；否则返回0。
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

/**
 * 对工作量时迁验证
 */
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce) //拿到值开始计划
	hash := sha256.Sum256(data)              //是不是256
	hashInt.SetBytes(hash[:])                //将 byte 转换成 int

	isValid := hashInt.Cmp(pow.target) == -1 //判断值是否小于 目标值

	return isValid
}


