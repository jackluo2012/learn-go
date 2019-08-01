package blockschain

import "bytes"

//交易输入
type TXInput struct {
	Txid      []byte //交易的 ID
	Vout      int    // 存储输出的序列号(一个交易可以包括多个 TXO)

	Signature []byte //签名
	PubKey    []byte //公共的 key
}
/**
 * 是否可以使用 特写的未哈希的公钥来解锁一个 TXO
 */
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {

	lockingHash := HashPubkey(in.PubKey)
	//比较两个 bytes 是否一致
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
