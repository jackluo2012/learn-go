package blockschain

import (
	"bytes"
	"encoding/gob"
	"log"
)

//交易输出
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

/**
 * 对TXO进行加锁。
 * 当发给某人货币时，仅仅知道他们的地址，
 * 因此该方法唯一入参就是地址信息。从地址中从解码出哈希后的公钥，将其保存到PubKeyHash中
 */
func (out *TXOutput) Lock(address []byte) {
	//解码
	pubKeyHash := Base58Decode(address)
	//获取公钥
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	//设置输出的公钥
	out.PubKeyHash = pubKeyHash
}

/**
 * 核查 TXO 是否被特定的 hash后的公钥锁定
 *
 */
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	//对比 bytes 是否是一致
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

/**
 * 初始化一个输出交易块
 */
func NewTXOutput(value int, address string) *TXOutput {

	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

type TXOutputs struct {
	Outputs []TXOutput
}

func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}
	return outputs
}
