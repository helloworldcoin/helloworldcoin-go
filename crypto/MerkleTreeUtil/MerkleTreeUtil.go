package MerkleTreeUtil

import (
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/crypto/Sha256Util"
)

func CalculateMerkleTreeRoot(dataList [][]byte) []byte {
	tree := dataList[:]
	levelOffset := 0
	for levelSize := len(tree); levelSize > 1; levelSize = (levelSize + 1) / 2 {
		for left := 0; left < levelSize; left += 2 {
			right := min(left+1, levelSize-1)
			leftBytes := tree[levelOffset+left]
			rightBytes := tree[levelOffset+right]
			tree = append(tree, Sha256Util.DoubleDigest(ByteUtil.Concat(leftBytes, rightBytes)))
		}
		levelOffset += levelSize
	}
	return tree[len(tree)-1]
}

func min(num1 int, num2 int) int {
	if num1 < num2 {
		return num1
	} else {
		return num2
	}
}
