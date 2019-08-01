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
 *
 */
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previos transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
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

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}
	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}
	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Rrevious transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

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
 * 创建一个新交易
 */

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(Subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}

	tx.ID = tx.Hash()
	return &tx
}

/**
 * 新建一个转账交易
 */
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}

	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubkey(wallet.PublicKey)
	//查找未被消费的交易区块和金额
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

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

	// 构建输出列表
	outputs = append(outputs, *NewTXOutput(amount, to))
	//如果是 大于,再加入
	if acc > amount {
		//剩下的
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, wallet.PrivateKey)

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
