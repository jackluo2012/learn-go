package blockschain

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp    int64          //区块生成时间
	Transactions []*Transaction //交易地址
	PreBlockHash []byte         //前一块的 hash
	Hash         []byte         //当前 hash
	Nonce        int            //计算出来的值
	Height       int
}

func NewBlock(transactions []*Transaction, preBlockHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), transactions, preBlockHash, []byte{}, 0, height}
	//加入工作证明
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	//	block.SetHash()
	return block
}

/**
 *  获取 交易信息的 sha265
 */
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)
	return mTree.RootNode.Data
}

/**
 * 创世纪链
 */

func NewGenesisBlock(coinbas *Transaction) *Block {
	return NewBlock([]*Transaction{coinbas}, []byte{},0)
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
