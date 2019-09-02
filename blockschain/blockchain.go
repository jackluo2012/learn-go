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
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {

	var lastHash []byte

	var lastHeight int
	// 在交易加入到block前进行交易验证
	// 验证交易是否合法

	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}
	//查询 出最后一块交易的id
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		lastHash = b.Get([]byte("l"))

		blockData := b.Get(lastHash)
		block := DeserializeBlock(blockData)

		lastHeight = block.Height

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//进行挖矿
	newBlock := NewBlock(transactions, lastHash, lastHeight+1)
	//写入区块记录
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
	return newBlock
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
/**
 * 遍历block 返回所有的UTXO
 * 和 FIndUnspentTransactions 完全一致
 * 把回的 是 交易ID和 输出值
 *
 */
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)

	bci := bc.Iterator()

	for {
		//迭代处理
		block := bci.Next()
		//获取所有的 区块 交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							//最后一块就跳出吗？
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PreBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

/**
 * 迭代遍历
 */

func (bc *Blockchain) Iterator() *BlockchainItertor {

	bci := &BlockchainItertor{bc.Tip, bc.Db}

	return bci
}

// GetBestHeight returns the height of the latest block
func (bc *Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
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
 * 1.遍历所有的block 返回UTXO ,2.通过交易满足要求的可以供消费的交易，3.返回指定公钥hash值的所有UTXO,用于计算yu额
 * 4.根据
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

/**
 * SignTransaction 对于一个交易找到其所有引用的交易后，进行签名
 */
func (bc *Blockchain) SignTransaction(tx *Transaction, privkey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		//按交易ID做键值对
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privkey, prevTXs)
}

/**
 * 现在，需要一个根据交易ID获取交易的函数，
 * 由于需要访问blockchain,
 * 因此作为blockchain的一个方法实现:
 * 根据ID查找并返回交易；
 */
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

//流程: 1.进行查找交易 ，2.引用交易,进行签名,3.进行验证
/**
 * 对于一个交易到其所有引用的交易后,进行验证。
 */
func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	//挖矿奖励
	if tx.IsCoinbase() {
		return true
	}

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
func NewBlockchain(nodeID string) *Blockchain {
	//检查 文件是否存在,如果 不存,就返回
	dbFile := fmt.Sprintf(DbFile, nodeID)
	if dbExists(dbFile) == false {
		fmt.Println("No existing blockchain found,Create on first")
		os.Exit(1)
	}

	var tip []byte

	//打开文件标准
	db, err := bolt.Open(dbFile, 0600, nil)
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

// AddBlock saves the block into the blockchain
func (bc *Blockchain) AddBlock(block *Block) {
	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		blockInDb := b.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte("l"))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			err = b.Put([]byte("l"), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.Tip = block.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// GetBlock finds a block by its hash and returns it
func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

// GetBlockHashes returns a list of hashes of all the blocks in the chain
func (bc *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		block := bci.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PreBlockHash) == 0 {
			break
		}
	}

	return blocks
}

//检查区块链文件是否存在
func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

/**
 * 创建一个  区块链 数据库
 */
func CreateBlockchain(address string, nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(DbFile, nodeID)
	if dbExists(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	//组织 交易信息
	cbtx := NewCoinbaseTX(address, GenesisCoinbaseData)
	//创建创世区块
	genesis := NewGenesisBlock(cbtx)

	db, err := bolt.Open(dbFile, 0600, nil)
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
