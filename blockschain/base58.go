package blockschain

import (
	"bytes"
	"math/big"
)
/**
 * Base58算法将一些字符移了,
 * 因此没有了0(数值0)O(o 的大写形式) I(i 的大写) l(L 的小写形式)
 * 此外没有+ 和 /
 */
var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base58Encode encodes a byte array to Base58
func Base58Encode(input []byte) []byte {

	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	ReverseBytes(result)
	for b := range input {
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}
	return result
}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0
	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}
	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)
	return decoded
}
