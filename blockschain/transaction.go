package blockschain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
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

/**
 * 将交易进行序列化
 */
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

/**
 * 交易的 hash 值
 */
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

/**
 * privKey 私钥
 *
 * 以及所引用交易的之前的集合
 */
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//检查是否创世区块链
	// 没有真实的TXI,因此此交易不进行签名。
	if tx.IsCoinbase() {
		return
	}
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previos transaction is not correct")
		}
	}
	//基于TrimmedCopy 交易进行签名
	txCopy := tx.TrimmedCopy()
	// 对交易进行签名时，
	// 需要获取该交易所有TXI引用的TXO列表
	//
	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		// 再次把Signaure字段设置为 nil
		// 此时，除了当前被处理的交易外，其他所有交易都是 "空" 交易
		txCopy.Vin[inID].Signature = nil
		// 将PubKey 设置为其所引用的TXO的 PubKeyHash
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		//由于比特币允许交易包含来自于不同地址的TXI
		//单独地对每个TXI进行签名,（虽然对于我们的应用来说没必要这么做，一个交易包中的TXI均来自同一地址）
		txCopy.ID = txCopy.Hash() //Hash方法将交易
		//我们将PubKey字段重新设置为nil避免影响后续的迭代
		txCopy.Vin[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].Signature = signature
	}
}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("		Input %d:", i))
		lines = append(lines, fmt.Sprintf("		TXID %x:", input.Txid))
		lines = append(lines, fmt.Sprintf("		Out %d:", input.Vout))
		lines = append(lines, fmt.Sprintf("		Signature %x:", input.Signature))
		lines = append(lines, fmt.Sprintf("		PubKey %x:", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("		Output %d:", i))
		lines = append(lines, fmt.Sprintf("		Value: %d", output.Value))
		lines = append(lines, fmt.Sprintf("		PubKeyHash: %x", output.PubKeyHash))
	}
	return strings.Join(lines, "\n")
}

/**
 * 基于TrimmedCopay交易进行签名
 */
func (tx *Transaction) TrimmedCopy() Transaction {

	var inputs []TXInput
	var outputs []TXOutput
	//复制原始效果所有的TXI 和TXO
	for _, vin := range tx.Vin {
		//  将Signature 和 PubKey 被设置为了nil
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}
	for _, vout := range tx.Vout {

		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}
	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

/**
 * 我们使用ECDSA签名算法通过私钥privKey对txCopy.ID进行签名,
 * 生成一对数字序列,生成一对数字序列。
 */
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Rrevious transaction is not correct")
		}
	}

	//生成TrimmedCopy交易
	txCopy := tx.TrimmedCopy()
	//创建椭圆曲线用于生成键值对
	curve := elliptic.P256()

	//对于每个TXI的签名进行验证 :
	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		//这个过程和Sign方法是一致的,因为验证的数据需要和签名的数据是一致的
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		//在这里，之前，将签名生成的两个字节序列组合生成TXI的Signature;将椭圆 的X,Y坐标点集合(其实也是两个字节)组合生成TXI的PubKey。
		//现在，我们需要将TXI的Signature和Pubkey中的数据进行"拆包",用于crypto/ecdsa库进行验证使用。

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}

/*
// 输入 加锁 / 解锁
func (input TXInput) CanUnlockOutputWith(address string) bool {
	return input.ScriptSig == address
}

// 输出 - 加锁 / 解锁
func (output TXOutput) CanBeUnlockedWith(address string) bool {
	return output.ScriptPubKey == address
}
*/
/**
 * 创建一个创世 新交易 信息
 */

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	//输入为空
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	//添加输入 进行奖励
	txout := NewTXOutput(Subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}

	tx.ID = tx.Hash()
	return &tx
}

/**
 * 新建一个转账交易
 *
 */
func NewUTXOTransaction(wallet *Wallet, to string, amount int, utxoset *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	pubKeyHash := HashPubkey(wallet.PublicKey)
	//查找未被消费的交易区块和金额
	acc, validOutputs := utxoset.FindSpendableOutputs(pubKeyHash, amount)

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
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	from := fmt.Sprintf("%s", wallet.GetAddress())
	// 构建输出列表
	outputs = append(outputs, *NewTXOutput(amount, to))
	//如果是 大于,再加入
	if acc > amount {
		//剩下的
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}
	//组织交易
	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()

	//进行签名
	utxoset.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

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

// DeserializeTransaction deserializes a transaction
func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}
