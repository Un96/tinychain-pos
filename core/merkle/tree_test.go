package merkle

import (
	"testing"
)

func TestNewTree(t *testing.T) {
	txs := [][]byte{
		[]byte("alice->bob:10"),
		[]byte("bob->carol:5"),
		[]byte("carol->alice:3"),
		[]byte("dave->eve:20"),
	}

	tree := NewTree(txs)
	if tree == nil {
		t.Fatal("tree is nil")
	}
	if tree.Root == nil {
		t.Fatal("root is nil")
	}

	t.Logf("✅ Root hash: %s", tree.RootHash())
	t.Logf("   Leaf count: %d", len(tree.Leafs))
}

func TestMerkleProof(t *testing.T) {
	txs := [][]byte{
		[]byte("alice->bob:10"),  // index 0
		[]byte("bob->carol:5"),   // index 1
		[]byte("carol->alice:3"), // index 2
		[]byte("dave->eve:20"),   // index 3
	}

	tree := NewTree(txs)
	root := tree.Root.Hash

	// 测试每个叶子的 proof
	for i := 0; i < len(txs); i++ {
		proof := tree.GetProof(i)
		if proof == nil {
			t.Fatalf("GetProof(%d) returned nil", i)
		}

		// 验证 proof
		if !Verify(root, proof) {
			t.Fatalf("Verify failed for index %d", i)
		}

		t.Logf("✅ Index %d proof verified", i)
	}

	// 测试篡改数据后验证失败
	proof := tree.GetProof(0)
	proof.TargetHash[0] ^= 0xFF // 篡改哈希
	if Verify(root, proof) {
		t.Fatal("Verify should fail for tampered data")
	}
	t.Log("✅ Tampered data correctly rejected")
}
