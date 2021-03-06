package MerkleTreeUtil

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/crypto/Sha256Util"
	"helloworldcoin-go/util/ByteUtil"
	"helloworldcoin-go/util/MathUtil"
)

func CalculateMerkleTreeRoot(datas [][]byte) []byte {
	tree := datas[:]
	levelOffset := 0
	for levelSize := len(tree); levelSize > 1; levelSize = (levelSize + 1) / 2 {
		for left := 0; left < levelSize; left += 2 {
			right := MathUtil.Min(left+1, levelSize-1)
			leftBytes := tree[levelOffset+left]
			rightBytes := tree[levelOffset+right]
			tree = append(tree, Sha256Util.DoubleDigest(ByteUtil.Concatenate(leftBytes, rightBytes)))
		}
		levelOffset += levelSize
	}
	return tree[len(tree)-1]
}
