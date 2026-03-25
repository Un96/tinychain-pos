package merkle

import (
	"crypto/sha256"
	"encoding/hex"
)

// Node Merkle 树节点
type Node struct {
	Hash   []byte
	Left   *Node
	Right  *Node
	Parent *Node
	IsLeaf bool //判断是不是叶子节点
}

// Tree Merkle
type Tree struct {
	Root  *Node
	Leafs []*Node
}

// 计算hash值
func hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// NewTree  构建merkle 树
func NewTree(data [][]byte) *Tree {
	/*思路：
	1.创建叶子节点
	2.节点创建完后计算其hash值
	3.逐层向上计算父节点（奇数复制最后一个的节点数据）
	*/
	var nodes []*Node

	//⭐数据判断，确定输入的数据不为空
	if len(data) == 0 {
		return nil
	}

	//叶子节点
	for _, d := range data {
		node := &Node{
			Hash:   hash(d),
			IsLeaf: true,
		}
		nodes = append(nodes, node)
	}
	//保存原始叶子引用，用于回溯
	leafs := make([]*Node, len(nodes))
	copy(leafs, nodes)

	//构建merkelTree
	for len(nodes) > 1 {
		var newLevel []*Node
		//奇数节点/偶数节点处理
		if len(nodes)%2 == 1 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}
		//创建父节点
		for i := 0; i < len(nodes); i += 2 {
			parent := &Node{
				Left:  nodes[i],
				Right: nodes[i+1],
				Hash:  hash(append(nodes[i].Hash, nodes[i+1].Hash...)),
			}
			// 子节点指向父节点
			nodes[i].Parent = parent
			nodes[i+1].Parent = parent

			newLevel = append(newLevel, parent)
		}
		nodes = newLevel
	}
	return &Tree{
		Root:  nodes[0],
		Leafs: leafs,
	}
}

// RootHash 根哈希转换为十六进制字符串
func (t *Tree) RootHash() string {
	if t == nil || t.Root == nil {
		return ""
	}
	return hex.EncodeToString(t.Root.Hash)
}
