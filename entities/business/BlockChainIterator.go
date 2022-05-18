package business

import (
	"github.com/boltdb/bolt"
	"log"
	"smallPublicChain/entities/Persistence"
)

type BlockChainIterator struct {
	CurrentHash []byte
	BlockDB     *bolt.DB
}

func (it *BlockChainIterator) Next() *Block {
	var curBlock *Block
	err := it.BlockDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Persistence.TableName))
		if bucket != nil {
			blockBytes := bucket.Get(it.CurrentHash)
			curBlock = DeSerializeBlock(blockBytes)
			it.CurrentHash = curBlock.PrevBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panicf("Next():%s\n", err)
	}
	return curBlock
}
