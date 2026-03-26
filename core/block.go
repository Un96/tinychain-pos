package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/Un96/tinychain-pos/core/merkle"
	"time"
)

// Transaction 简化的交易
type Transaction struct {
	From   string
	To     string
	Amount int64
	Data   []byte
}

// Header 区块头
type Header struct {
	PrevHash   []byte
	MerkleRoot []byte
	Timestamp  int64
	Height     uint64
	Validator  []byte // PoS: 验证者地址
}

// Block 完整区块
type Block struct {
	Header       *Header
	Transactions []*Transaction
	Hash         []byte // 缓存的区块哈希
}

// Serialize 交易序列化
func (tx *Transaction) Serialize() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(tx)
	return buf.Bytes()
}

// Serialize Header序列化（用于计算哈希）
func (h *Header) Serialize() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(h)
	return buf.Bytes()
}

// NewBlock 创建新区块
func NewBlock(txs []*Transaction, prevHash []byte, height uint64, validator []byte) *Block {
	// TODO:
	// 1. 计算 MerkleRoot（用之前实现的 merkle 包）
	var txData [][]byte
	for _, tx := range txs {
		txData = append(txData, tx.Serialize())
	}
	tree := merkle.NewTree(txData)
	var merkleRoot []byte
	if tree != nil {
		merkleRoot = tree.Root.Hash
	}
	// 2. 组装 Header
	header := &Header{
		PrevHash:   prevHash,
		MerkleRoot: merkleRoot,
		Timestamp:  time.Now().Unix(),
		Height:     height,
		Validator:  validator,
	}
	// 3. 计算区块哈希（Header 的哈希）
	block := &Block{
		Header:       header,
		Transactions: txs,
	}
	block.Hash = block.HashBlock()
	// 4. 返回 Block
	return block
}

// NewGenesisBlock 创建创世区块（第一个区块）
func NewGenesisBlock(validator []byte) *Block {
	return NewBlock([]*Transaction{}, []byte{}, 0, validator)
}

// HashBlock 计算区块哈希
func (b *Block) HashBlock() []byte {
	// TODO: 对 Header 进行 gob 编码后哈希
	headerData := b.Header.Serialize()
	hash := sha256.Sum256(headerData)
	return hash[:]
}

// Serialize 区块序列化
func (b *Block) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(b)
	return buf.Bytes(), err
}

// DeserializeBlock 反序列化区块
func DeserializeBlock(data []byte) (*Block, error) {
	var block Block
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&block)
	return &block, err
}
