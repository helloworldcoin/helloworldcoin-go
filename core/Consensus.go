package core

import (
	"helloworldcoin-go/core/Model"
	"math/big"
)

type Consensus struct {
}

func (c *Consensus) checkConsensus(blockchainDataBase *BlockchainDatabase, block *Model.Block) bool {
        difficulty := block.Difficulty
        if StringUtil.isNullOrEmpty(difficulty) {
            difficulty = calculateDifficult(blockchainDataBase,block) 
            block.setDifficulty(difficulty) 
        }

         hash := block. Hash
        if StringUtil.isNullOrEmpty(hash) {
            hash = BlockTool.calculateBlockHash(block)
        }
        return new BigInteger(difficulty,16).compareTo(new BigInteger(hash,16)) > 0
}
