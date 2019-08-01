package blockschain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

type Blockchain struct {
	Tip []byte   //最新的那个值
	Db  *bolt.DB //数组长度
}

/**
 * 实现迭代器来逐个遍历 block
 */
type BlockchainItertor struct {
	CurrentHash []byte   //当前 最新 hash 的值,也就是最近一个
	Db          *bolt.DB //拿到操作的数据库
}

// 使用提供的事务挖掘新块
func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	//进行挖矿
	newBlock := NewBlock(transactions, lastHash)

	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		//设置最新hash 链为 tip
		bc.Tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

}

/**
 * 迭代遍历
 */

func (bc *Blockchain) Iterator() *BlockchainItertor {

	bci := &BlockchainItertor{bc.Tip, bc.Db}

	return bci
}

// FindUTXO finds and returns all unspent transaction outputs
func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	//获取 未被交易的 utxo
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	//未返回的
	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

/**
 *  查找未被消费的交易区块和金额
 */
func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	//获取款被使用的金额
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		//遍历所有的交易记录
		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs

}

/**
 * 查找 未消费 TXO(UTXO)
 * 未消费表式TXO 未被任何的 TXI 所引用
 */
func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {

	var unspentTXs []Transaction

	spentTXOs := make(map[string][]int)
	// 返回 blockchain 封装 接口
	bci := bc.Iterator()

	//查找
	for {
		// 交易存在 block 中,我们要遍历blockchain 中的每个 block
		block := bci.Next()

		//遍历 blockchain 中的交易
		for _, tx := range block.Transactions {

			txID := hex.EncodeToString(tx.ID)

		Outputs:
			//遍历每个  交易 中的 输出块
			for outIdx, out := range tx.Vout {
				// 如果是有使用有记录的不进行
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				//返回未交易的 UTXO
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			//如果不是创世区块
			if tx.IsCoinbase() == false {
				//遍历交易的所有输出
				for _, in := range tx.Vin {
					//验证交易是否正确
					if in.UsesKey(pubKeyHash) {
						//获取 输入的id 并转换成字符串
						inTxID := hex.EncodeToString(in.Txid)
						//将 UTXI 的 UTXO 加入
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		//查找到头结点
		if len(block.PreBlockHash) == 0 {
			break
		}
	}

	//返回未交易的 UTXO
	return unspentTXs
}

func (bc *Blockchain) SignTransaction(tx *Transaction, privkey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privkey, prevTXs)
}

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PreBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

/**
 * 返回下一块区块
 */
func (i *BlockchainItertor) Next() *Block {
	var block *Block

	//读操作
	err := i.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))   //获取表
		encodedBlock := b.Get(i.CurrentHash)   //获取最新的一个
		block = DeserializeBlock(encodedBlock) //解码得到
		return nil
	})
	if err != nil {

	}
	//将当前的,指向上一块区块的链接地址
	i.CurrentHash = block.PreBlockHash
	return block
}

/**
 * 将创世链加入区块链中
 */
func NewBlockchain(address string) *Blockchain {
	//检查 文件是否存在,如果 不存,就返回
	if dbExists() == false {
		fmt.Println("No existing blockchain found,Create on first")
		os.Exit(1)
	}

	var tip []byte

	//打开文件标准
	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	//使用读写事务
	err = db.Update(func(tx *bolt.Tx) error {
		//获取 最新 区块的地址
		b := tx.Bucket([]byte(BlocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

//检查区块链文件是否存在
func dbExists() bool {
	if _, err := os.Stat(DbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

/**
 * 创建一个  区块链 数据库
 */
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	cbtx := NewCoinbaseTX(address, GenesisCoinbaseData)
	//创建创世区块
	genesis := NewGenesisBlock(cbtx)

	db, err := bolt.Open(DbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		//创建一个表
		b, err := tx.CreateBucket([]byte(BlocksBucket))
		if err != nil {
			log.Panic(err)
		}

		//将写第一个区块写进去
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		//将 hash 存入最新的区块
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{Tip: tip, Db: db}

	return &bc
}
