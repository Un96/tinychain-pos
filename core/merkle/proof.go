package merkle

// Proof Merkle证明
type Proof struct {
	TargetHash  []byte
	ProofHashes [][]byte
	IsLeft      []bool
}

// GetProof 获取第 index 个叶子的证明
func (t *Tree) GetProof(index int) *Proof {
	if t == nil || index < 0 || index >= len(t.Leafs) {
		return nil
	}

	target := t.Leafs[index]
	proof := &Proof{
		TargetHash:  target.Hash,
		ProofHashes: [][]byte{},
		IsLeft:      []bool{},
	}

	// 从叶子向上遍历到根
	current := target
	for current.Parent != nil {
		parent := current.Parent

		if parent.Left == current {
			proof.ProofHashes = append(proof.ProofHashes, parent.Right.Hash)
			proof.IsLeft = append(proof.IsLeft, false)
		} else {
			proof.ProofHashes = append(proof.ProofHashes, parent.Left.Hash)
			proof.IsLeft = append(proof.IsLeft, true)
		}

		current = parent
	}

	return proof
}

// Verify 验证 proof 是否匹配 root
func Verify(root []byte, proof *Proof) bool {
	if proof == nil || len(proof.TargetHash) == 0 {
		return false
	}

	currentHash := proof.TargetHash

	// 逐层向上计算
	for i, sibling := range proof.ProofHashes {
		if proof.IsLeft[i] {
			// 兄弟在左：兄弟 + 我
			currentHash = hash(append(sibling, currentHash...))
		} else {
			// 兄弟在右：我 + 兄弟
			currentHash = hash(append(currentHash, sibling...))
		}
	}

	// 对比根哈希
	return string(currentHash) == string(root)
}
