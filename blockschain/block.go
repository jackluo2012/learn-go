package blockschain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Block struct {
	Timestamp    int64          //区块生成时间
	Transactions []*Transaction //交易地址
	PreBlockHash []byte         //前一块的 hash
	Hash         []byte         //当前 hash
	Nonce        int            //计算出来的值
}

func NewBlock(transactions []*Transaction, preBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, preBlockHash, []byte{}, 0}
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
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	//将半岛晨报
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}


/**
 * 创世纪链
 */

func NewGenesisBlock(coinbas *Transaction) *Block {
	return NewBlock([]*Transaction{coinbas}, []byte{})
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
