package blockschain

/**
 * mo 克尔树
 * 每个block有一个默克树，树中每个叶子节点是一个交易的hash值。
 * 叶子节点为的数量一定是偶数，然后并非每个block都恰好有偶数个交易。
 * 当block有奇数个交易时，最后一个交易会被复制一次(复制仅仅发生在默克尔树中，而不是block中)
 * 默克尔树自下而上的进行组织，叶子节点成对分组后将两个hash值组合后生成新的hash值,形成上层的树节点，重复整个过程直到只有一个树节点为止
 * 也就是所说的根节点.根节点的hash值是整个交易集的唯一标识，保存在block头信息中，用于Pow过程
 * 默克尔树的好处是:验证某个交易不需要现在整个block,而仅仅需要交易hash值，形成上层的树节点，
 * 默克尔树节点hash值以及默克尔路径即刻
 */

import "crypto/sha256"

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right

	return &mNode
}
/**
 * 创建一个默克尔树之前，需要确保有偶数个叶子节点，然后将数据转换为叶子节点,最后生成整个默克树
 */
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}
	//上述实现存在一个问题:当叶子节点的数量是2n时，可以正常创建树 （for j :=0;j <len(nodes);j+=2...该段代码体现
	//如果不是2n时,创建树会发生异常!
	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}
		nodes = newLevel
	}

	mTree := MerkleTree{&nodes[0]}

	return &mTree
}
