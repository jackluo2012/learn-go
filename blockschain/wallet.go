package blockschain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const (
	version            = byte(0x00)
	addressChecksumLen = 4
	walletFile         = "wallet.dat"
)

//钱包 是一个密钥对
type Wallet struct {
	PrivateKey ecdsa.PrivateKey //私钥
	PublicKey  []byte           //公钥
}

/**
 * 用于生成一个新钱包
 */
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

/**
 * 使用椭圆曲线算法生成私钥
 */
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	cure := elliptic.P256()
	//使用椭圆曲线算法生成私钥
	private, err := ecdsa.GenerateKey(cure, rand.Reader)
	if err != nil {
		log.Panic("解密 加密出错", err)
	}
	//通过私钥 生成公钥 椭圆曲线算法中,公钥是曲线上的点集合.因些,公钥由 X,Y 坐标混合而成.
	//坐标组合在一起,生成公钥
	pubkey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubkey
}

/**
 * 钱包地址 =  Base58 对 version + pubKeyHash + checksum
 */
func (w Wallet) GetAddress() []byte {

	pubKeyHash := HashPubkey(w.PublicKey)
	// 版本信息 zhui 加到 pubkeyhas 前
	versionedPayload := append([]byte{version}, pubKeyHash...)
	//二次 sha256 取前 几位
	checksum := checksum(versionedPayload)
	//加 加入到后面
	fullPayload := append(versionedPayload, checksum...)
	//
	address := Base58Encode(fullPayload)

	return address
}

/**
 * 使用 ripemd160 ( sha256 (Pubkey) ) 	对公钥进行两次哈希,生成 pub KeyHash
 */
func HashPubkey(pubKey []byte) []byte {

	//先进行 sha256
	publicSHA256 := sha256.Sum256(pubKey)
	//再进行
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic("", err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

/**
 * 进行两次哈希得到一个 hash值,取该值的前 n 个字节最终生成 checksum
 */
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
