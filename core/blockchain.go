package core

import (
	"errors"
	"fmt"
	"github.com/Un96/tinychain-pos/core/merkle"
)

// Blockchain 区块链
type Blockchain struct {
	Blocks       map[uint64]*Block // height -> block
	LatestHash   []byte            // 最新区块哈希
	LatestHeight uint64            // 最新高度
}

// NewBlockchain 创建新区块链（包含创世区块）
func NewBlockchain(genesisValidator []byte) *Blockchain {
	genesis := NewGenesisBlock(genesisValidator)

	bc := &Blockchain{
		Blocks:       make(map[uint64]*Block),
		LatestHash:   genesis.Hash,
		LatestHeight: 0,
	}
	bc.Blocks[0] = genesis

	return bc
}

// AddBlock 添加新区块
func (bc *Blockchain) AddBlock(txs []*Transaction, validator []byte) error {
	// TODO:
	// 1. 获取最新区块
	prevBlock := bc.GetLatestBlock()
	// 2. 创建新区块（height = LatestHeight + 1, prevHash = LatestHash）
	newBlock := NewBlock(txs, prevBlock.Hash, prevBlock.Header.Height+1, validator)
	// 3.验证区块
	if err := bc.ValidateBlock(newBlock); err != nil {
		return fmt.Errorf("block validation failed: %v", err)
	}
	// 4. 添加到 Blocks
	bc.Blocks[newBlock.Header.Height] = newBlock
	bc.LatestHash = newBlock.Hash
	bc.LatestHeight = newBlock.Header.Height

	return nil
}

// GetBlockByHeight 通过高度获取区块
func (bc *Blockchain) GetBlockByHeight(height uint64) (*Block, error) {
	block, exists := bc.Blocks[height]
	if !exists {
		return nil, fmt.Errorf("block not found at height %d", height)
	}
	return block, nil
}

// GetLatestBlock 获取最新区块
func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.Blocks[bc.LatestHeight]
}

// ValidateBlock 验证区块（简化版）
func (bc *Blockchain) ValidateBlock(block *Block) error {
	// 1. 检查 PrevHash 是否匹配 LatestHash
	latestBlock := bc.GetLatestBlock()
	if string(block.Header.PrevHash) != string(latestBlock.Hash) {
		return errors.New("prevHash does not match latest block hash")
	}

	// 2. 检查 Height 是否连续
	if block.Header.Height != latestBlock.Header.Height+1 {
		return fmt.Errorf("height should be %d, got %d",
			latestBlock.Header.Height+1, block.Header.Height)
	}

	// 3. 验证 MerkleRoot（重新计算对比）
	var txData [][]byte
	for _, tx := range block.Transactions {
		txData = append(txData, tx.Serialize())
	}
	tree := merkle.NewTree(txData)
	if tree == nil {
		if len(block.Header.MerkleRoot) != 0 {
			return errors.New("merkleRoot mismatch: should be empty")
		}
	} else {
		if string(tree.Root.Hash) != string(block.Header.MerkleRoot) {
			return errors.New("merkleRoot verification failed")
		}
	}

	return nil
}
