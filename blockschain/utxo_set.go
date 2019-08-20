package blockschain

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

const utxoBucket = "chainstate"

// 设置交易信息
// 缓存 未消费的
type UTXOSet struct {
	Blockchain *Blockchain
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
/**
 * 基于 UTXO 集合 进行计算
 * 查找未被消费的交易区块和金额
 */
func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.Db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			// 获取id
			txID := hex.EncodeToString(k)
			//解码
			outs := DeserializeOutputs(v)

			for outIdx, out := range outs.Outputs {
				//查找剩于 的 金额
				if out.IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

// FindUTXO finds UTXO for a public key hash
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {

	var UTXOs []TXOutput

	db := u.Blockchain.Db

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoBucket))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// CountTransactions returns the number of transactions in the UTXO set
func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.Db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

// Reindex rebuilds the UTXO set
/**
 *  将Blockchain.FindUTXO结果存储到数据库中
 */
func (u UTXOSet) Reindex() {
	db := u.Blockchain.Db

	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		//若存在UTXO集合，将其删除,
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}
		//然后再创建 名称
		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//再查找 未消费的
	UTXO := u.Blockchain.FindUTXO()
	//存入数据中
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			//序列化存入
			err = b.Put(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
}

// Update updates the UTXO set with transactions from the Block
// The Block is considered to be the tip of a blockchain
/**
 * UTXO集合的出现意味着 交易信息分为两部分存储:实际交易存在了blockchain中
 * UTXO数据存储在UTXO集中,为了保持数据一致性的同时避免每一次挖到 一个新的block 都 重建 UTXO集合
 * 需要一个数据同步机制
 *
 * 当新的block生成时,需要更新UTXO集合,
 * 删除已消费的TXO,同时将新生成的交易中的UTXO添加进来
 */
func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.Db

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			//如果不是创世区块链
			if tx.IsCoinbase() == false {
				//检查 输入
				for _, vin := range tx.Vin {
					updatedOuts := TXOutputs{}
					outsBytes := b.Get(vin.Txid)
					outs := DeserializeOutputs(outsBytes)
					//将别个转入的钱 ,加入
					for outIdx, out := range outs.Outputs {
						if outIdx != vin.Vout {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}
					//如果为0 就删除
					if len(updatedOuts.Outputs) == 0 {
						err := b.Delete(vin.Txid)
						if err != nil {
							log.Panic(err)
						}
					} else {
						//更新一下
						err := b.Put(vin.Txid, updatedOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
					}

				}
			}

			newOutputs := TXOutputs{}
			//检查 输入 ,这个应该是 挖矿奖励
			for _, out := range tx.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			err := b.Put(tx.ID, newOutputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
