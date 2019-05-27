package blockschain

import (
	"bytes"
	"encoding/binary"
	"log"
)

/**
 * 将一个 int64 转化为一个字节数组
 */
func IntToHex2(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func IntToHex(i int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, i)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
