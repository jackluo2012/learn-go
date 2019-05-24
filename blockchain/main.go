package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"math"
	"math/big"
	"os"
	"strconv"
	"time"
)

//设置挖矿难度  前24位为 0
const (
	targetBits   = 24
	maxNonce     = math.MaxInt64
	dbFile       = "blockchian.db"
	blocksBucket = "blocks"
)

//加入工作量证明  pow
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	//进行位移操作
	target.Lsh(target, uint(256-targetBits))
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
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{})
	return data
}

/**
 * Pow 核心算法 :
 */

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0 //初始化为0
	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
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

/**
 * 将一个 int64 转化为一个字节数组
 */
func IntToHex(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

type Block struct {
	Timestamp    int64  //区块生成时间
	Data         []byte //记录的值
	PreBlockHash []byte //前一块的 hash
	Hash         []byte //当前 hash
	Nonce        int    //计算出来的值
}

func NewBlock(data string, preBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), preBlockHash, []byte{}, 0}
	//加入工作证明
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	//	block.SetHash()
	return block
}

type Blockchain struct {
	tip []byte   //最新的那个值
	db  *bolt.DB //数组长度
}

/**
 * 加入区块链
 */
func (bc *Blockchain) AddBlock(data string) {

	var lastHsh []byte
	//使用只读事务获取 当前数据库中最新 block 的 hash 值
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHsh = b.Get([]byte("1"))
		return nil
	})
	//创建一个最新的区块
	newBlock := NewBlock(data, lastHsh)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))             //拿到表
		err = b.Put(newBlock.Hash, newBlock.Serialize()) //存入值
		err = b.Put([]byte("1"), newBlock.Hash)          //设置为最新的
		bc.tip = newBlock.Hash                           //指向最新的 hash 地址
		return nil
	})

}

/**
 * 实现迭代器来逐个遍历 block
 */
type BlockchainItertor struct {
	currentHash []byte   //当前 最新 hash 的值,也就是最近一个
	db          *bolt.DB //拿到操作的数据库
}

/**
 * 迭代遍历
 */

func (bc *Blockchain) Iterator() *BlockchainItertor {
	bci := &BlockchainItertor{bc.tip, bc.db}
	return bci
}

/**
 * 返回下一块区块
 */
func (i *BlockchainItertor) Next() *Block {
	var block *Block

	//读操作
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))   //获取表
		encodedBlock := b.Get(i.currentHash)   //获取最新的一个
		block = DeserializeBlock(encodedBlock) //解码得到
		return nil
	})
	if err != nil {

	}
	//将当前的,指向上一块区块的链接地址
	i.currentHash = block.PreBlockHash
	return block
}

/**
 * 创世纪链
 */

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

/**
 * 将创世链加入区块链中
 */
func NewBlockchain() *Blockchain {
	var tip []byte

	//打开文件标准
	db, err := bolt.Open(dbFile, 0600, nil)
	//使用读写事务
	err = db.Update(func(tx *bolt.Tx) error {
		//尝试获取 一个 bucket
		b := tx.Bucket([]byte(blocksBucket))
		//不存在
		if b == nil {
			//创建创业区块
			genesis := NewGenesisBlock()
			//
			b, err = tx.CreateBucket([]byte(blocksBucket))
			//往里面写入数据
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("1"), genesis.Hash)
			//将 tip 设置成最新的 区块 hash 值
			tip = genesis.Hash
		} else {
			//获取 1
			tip = b.Get([]byte("1"))
		}
		return err
	})
	bc := Blockchain{tip, db}

	return &bc
}

/**
 * 将 Block 序列化为字节 数组
 */
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {

	}
	return result.Bytes()
}

/**
 * 将 byte字节数组反序列化为 Block 对象
 */
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {

	}
	return &block
}

// 新建一个 Blockchain 客户端操作
type CLI struct {
	bc *Blockchain
}

func (cli *CLI) Run() {
	//检测参数
	cli.validateArgs()

	//创建子命令
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "区块数据不能为空")
	switch os.Args[1] {
	case "addblock":
		fmt.Println("12121")
		addBlockCmd.Parse(os.Args[2:])
	case "printchain":
		printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(0)
	}
	//检析成功是否为函数
	if addBlockCmd.Parsed() {
		fmt.Printf(*addBlockData)

		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()
	for {
		block := bci.Next()
		fmt.Printf("Prev. hash: %x\n", block.PreBlockHash)
		fmt.Printf("Data. hash: %s\n", block.Data)
		fmt.Printf("Hash. hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PreBlockHash) == 0 {
			break
		}
	}
}
func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("success!")
}
func (cli *CLI) printUsage() {

}
func (cli *CLI) validateArgs() {

}

func main() {
	bc := NewBlockchain()

	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}
