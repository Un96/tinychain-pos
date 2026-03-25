package merkle

import "testing"

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
	t.Logf("Root hash: %x", tree.Root.Hash)
}
