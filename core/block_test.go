package core

import (
	"testing"
)

func TestNewBlockchain(t *testing.T) {
	bc := NewBlockchain([]byte("validator1"))

	if bc == nil {
		t.Fatal("blockchain is nil")
	}
	if len(bc.Blocks) != 1 {
		t.Fatalf("should have genesis block only, got %d", len(bc.Blocks))
	}
	if bc.LatestHeight != 0 {
		t.Fatalf("latest height should be 0, got %d", bc.LatestHeight)
	}

	genesis := bc.GetLatestBlock()
	t.Logf("✅ Genesis Hash: %x", genesis.Hash)
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockchain([]byte("validator1"))

	// 添加区块1
	txs := []*Transaction{
		{From: "alice", To: "bob", Amount: 10},
	}

	err := bc.AddBlock(txs, []byte("validator2"))
	if err != nil {
		t.Fatalf("add block failed: %v", err)
	}

	if bc.LatestHeight != 1 {
		t.Fatalf("height should be 1, got %d", bc.LatestHeight)
	}

	block1 := bc.GetLatestBlock()
	if block1.Header.Height != 1 {
		t.Fatal("block height mismatch")
	}
	if string(block1.Header.PrevHash) != string(bc.Blocks[0].Hash) {
		t.Fatal("prevHash should match genesis hash")
	}

	t.Logf("✅ Block #1 added")
	t.Logf("   Hash: %x", block1.Hash)
	t.Logf("   PrevHash: %x", block1.Header.PrevHash)
	t.Logf("   Tx Count: %d", len(block1.Transactions))
}

func TestGetBlockByHeight(t *testing.T) {
	bc := NewBlockchain([]byte("validator1"))
	bc.AddBlock([]*Transaction{{From: "a", To: "b", Amount: 5}}, []byte("v2"))

	block, err := bc.GetBlockByHeight(1)
	if err != nil {
		t.Fatal(err)
	}
	if block.Header.Height != 1 {
		t.Fatal("wrong block")
	}

	// 查询不存在的区块
	_, err = bc.GetBlockByHeight(999)
	if err == nil {
		t.Fatal("should return error for non-existent block")
	}

	t.Log("✅ GetBlockByHeight OK")
}
