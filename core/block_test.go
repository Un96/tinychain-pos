package core

import (
	"testing"
)

func TestNewGenesisBlock(t *testing.T) {
	validator := []byte("validator1")
	genesis := NewGenesisBlock(validator)

	if genesis == nil {
		t.Fatal("genesis block is nil")
	}
	if genesis.Header.Height != 0 {
		t.Fatalf("genesis height should be 0, got %d", genesis.Header.Height)
	}
	if len(genesis.Header.PrevHash) != 0 {
		t.Fatal("genesis prevHash should be empty")
	}
	if len(genesis.Hash) == 0 {
		t.Fatal("genesis hash should not be empty")
	}

	t.Logf("✅ Genesis Block Hash: %x", genesis.Hash)
	t.Logf("   Timestamp: %d", genesis.Header.Timestamp)
}

func TestNewBlock(t *testing.T) {
	// 先创建创世区块
	genesis := NewGenesisBlock([]byte("validator1"))

	// 创建交易
	txs := []*Transaction{
		{From: "alice", To: "bob", Amount: 10},
		{From: "bob", To: "carol", Amount: 5},
	}

	// 创建新区块
	block := NewBlock(txs, genesis.Hash, 1, []byte("validator2"))

	if block.Header.Height != 1 {
		t.Fatalf("height should be 1, got %d", block.Header.Height)
	}
	if string(block.Header.PrevHash) != string(genesis.Hash) {
		t.Fatal("prevHash should match genesis hash")
	}
	if len(block.Header.MerkleRoot) == 0 {
		t.Fatal("merkleRoot should not be empty")
	}
	if len(block.Transactions) != 2 {
		t.Fatalf("should have 2 transactions, got %d", len(block.Transactions))
	}

	t.Logf("✅ Block #1 Hash: %x", block.Hash)
	t.Logf("   PrevHash: %x", block.Header.PrevHash)
	t.Logf("   MerkleRoot: %x", block.Header.MerkleRoot)
	t.Logf("   Tx Count: %d", len(block.Transactions))
}

func TestBlockSerialize(t *testing.T) {
	block := NewGenesisBlock([]byte("validator1"))

	data, err := block.Serialize()
	if err != nil {
		t.Fatalf("serialize failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("serialized data is empty")
	}

	// 反序列化
	recovered, err := DeserializeBlock(data)
	if err != nil {
		t.Fatalf("deserialize failed: %v", err)
	}

	if recovered.Header.Height != block.Header.Height {
		t.Fatal("height mismatch after deserialize")
	}
	if string(recovered.Hash) != string(block.Hash) {
		t.Fatal("hash mismatch after deserialize")
	}

	t.Log("✅ Serialize/Deserialize OK")
}
