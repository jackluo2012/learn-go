package blockschain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const (
	Subsidy = 50
)

//定义交易结构
type Transaction struct {
	ID   []byte
	Vin  []TXInput  //多个交易输入
	Vout []TXOutput //交易输出
}

//交易输入
type TXInput struct {
	Txid      []byte //交易的 ID
	Vout      int    // 存储输出的序列号(一个交易可以包括多个 TXO)
	ScriptSig string //与之关联的	TXO 的 解锁 并生成新的 TXO
}

// 输入 加锁 / 解锁
func (input TXInput) CanUnlockOutputWith(address string) bool {
	return input.ScriptSig == address
}

// 输出 - 加锁 / 解锁
func (output TXOutput) CanBeUnlockedWith(address string) bool {
	return output.ScriptPubKey == address
}

//交易输出
type TXOutput struct {
	Value        int
	ScriptPubKey string //用户定义的字符串
}


/**
 * 新建一个转账交易
 */
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	//查找未被消费的交易区块和金额
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	//检查够转账的金额
	if acc < amount {
		log.Panic("错误: 余额不足!!!")
	}

	// 构建输入列表
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// 构建输出列表
	outputs = append(outputs, TXOutput{amount, to})
	//如果是 大于,再加入
	if acc > amount {
		//剩下的
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

/**
 * 创建一个新交易
 */
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{Subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

/**
 * 设置 交易的 ID 值
 */
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	//将交易编码
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	//算出sha256的值
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}
/**
 * 判断是否是 创世区块
 */
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}
