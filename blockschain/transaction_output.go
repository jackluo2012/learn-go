package blockschain

import "bytes"

//交易输出
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
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
