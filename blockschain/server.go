package blockschain

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

/**
 * 当一个节点收到version消息后，需要做些什么呢？它将发送自身的version消息给对方。
 * 这是一种握手协议:不握手，不能做其它任何事情。verion用于发现一个更长的blockchain。
 * 当一个节点收到version消息后，如果 发送节点的blockchain更长的话，
 * 接收节点将会发送请求获取缺失的block
 */

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var knownNodes = []string{"localhost:3000"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]Transaction)

type addr struct {
	AddrList []string
}
// 区块 消息
type block struct {
	AddrFrom string
	Block    []byte
}
/**
 * getblocks消息（比特币实现中更加复杂）;
 */
type getblocks struct {
	AddrFrom string
}
// 用于获取某个block或交易信息
type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

/**
 * inv  消息包含发送节点所拥有的block或效果的hash值列表。
 * Type 用于表示是block还是交易。
 */
type inv struct {
	AddrFrom string
	Type     string // 表示是 block 还是交易
	Items    [][]byte
}
// 交易消息
type tx struct {
	AddFrom     string
	Transaction []byte
}

type verzion struct {
	Version    int    //不存储任何重要信息
	BestHeight int    // 存放该节点的blockchain的长度
	AddrFrom   string // 存放消息发送者的地址。
}

/**
 * 消息由字节数组构成： 前12个字节表示命令名（此时是"version"),紧接着是gob编码的消息体
 */
func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}
/**
 * commandToBytes函数创建长度为12的字节数组,
 * 然后将命令名填充到其中。
 *
 */
func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func extractCommand(request []byte) []byte {
	return request[:commandLength]
}

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}

func sendAddr(address string) {
	nodes := addr{knownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(address, request)
}

func sendBlock(addr string, b *Block) {
	data := block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func sendTx(addr string, tnx *Transaction) {
	data := tx{nodeAddress, tnx.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

/**
 * 若当前节点不是中心节点，
 * 需要向中心节点发送version消息,
 * 用于检查当前节点的blockchain是否过时
 */
func sendVersion(addr string, bc *Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(verzion{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(knownNodes))
	requestBlocks()
}
/**
 * 当收到一个新block时，将该block加入blockchain中。
 * 如果还有其它待下载的block,将向相同的节点发送消息来获取block信息.
 * UTXO将被重建。(待完善：目前实现认为新block是有效的，后续 在新block加入blockchain之前，应该对其进行有效性验证)
 * (目前实现使用UTXOSet.Reindex()对UTXO进行重建，后续 应该使用UTXOSet.Update(Block),避免重建整个blockchain带来
 */
func handleBlock(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex()
	}
}

/**
 * 为了跟踪已下载block信息,当block的hash值完成传输后，将被保存到blocksInTransit中,
 * 这样做可以使我们从不同的节点下载block。一旦block置为transit状态后，
 * 将向inv消息的发送节点发送getdata消息并更新blocksInTransit。在真实的P2P网络中,将从不同的节点传输block信息
 * 在我们的实现中,inv消息中仅有一个block或交易的hash值,这也就是当payload.Type ==『tx』时仅仅获取第一个hash值的原因。
 * 当mempool没有找到对应block或交易时，将发送getdata消息获取相关信息
 */
func handleInv(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}
/**
 * 返回当前节点所拥有的block的hash值列表,而不是所有block的详细信息,
 * 出于降低网络负荷的目的,若需要下载block，
 * 可以从多个节点同时下载,没必要从单个节点下载。相应的消息处理函数了很简单
 */
func handleGetBlocks(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}
/**
 * 该消息处理函数根据消息中的数据类型,返回block或交易。
 * 需要注意的是,该函数没有检查对应的block或交易是否存在
 */
func handleGetData(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		sendTx(payload.AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

func handleTx(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == knownNodes[0] {
		for _, node := range knownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(mempool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*Transaction

			for id := range mempool {
				tx := mempool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := NewCoinbaseTX(miningAddress, "")
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)
			UTXOSet := UTXOSet{bc}
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			for _, node := range knownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}
/**
 * version的消息处理函数如下:
 * 首先，对请求中的数据进行解码获取 消息体，
 * 这对所有消息处理函数都是一样的,后续介绍其它消息处理函数时将不再进行介绍
 */
func handleVersion(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload verzion

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if myBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	// sendAddr(payload.AddrFrom)
	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}
  /*
   * bytesToCommand与commandToBytes作用正好相反： 当一个节点接收到一个命令时，通过bytesToCommand函数从中解析出命令名和消息体。
   * handleConnection函数会根据解析出的命令来调用不同的消息处理函数
   *
   * 随后，接收节点将自身的BestHeight与消息体（即送节点）中的BestHeight进行比较。
   * 若接收节点的blockchain更长,将发送version消息;若发送节点的blockchain更长，将发送getblocks消息
   *
 */
func handleConnection(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

/**
 * 创建一个Server用于接收消息:
 * minerAddress 接收挖矿奖励的地址
 */

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)

	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := NewBlockchain(nodeID)

	/**
	 * 首先，将中心节点地址硬维码到程序中： 一天始每个节点必须要知道地址.
	 */
	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}
func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}
