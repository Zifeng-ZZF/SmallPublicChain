package business

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
	"smallPublicChain/entities/Persistence"
)

const DBFullName = "entities/Persistence" + Persistence.DBName

type Blockchain struct {
	tip     []byte // last block's Hash
	blockDB *bolt.DB
}

func CreateBlockChain(genesisDataStr string) *Blockchain {
	if !dbExist() {
		return initializeBlockChain(genesisDataStr)
	}
	return getBlockChainFromDB()
}

func initializeBlockChain(genesisDataStr string) *Blockchain {
	genesisBlock := CreateGenesisBlock(genesisDataStr)
	db, err := bolt.Open(DBFullName, 0600, nil)
	if err != nil {
		log.Fatal("Not opening db")
		return nil
	}

	err = db.Update(func(tx *bolt.Tx) error {
		tbNameBytes := []byte(Persistence.TableName)
		bucket, e := tx.CreateBucket(tbNameBytes)
		if e != nil {
			log.Panic("create fail")
		}
		blockBytes := genesisBlock.Serialize()
		e = bucket.Put(genesisBlock.Hash, blockBytes)
		if e != nil {
			log.Panic("put genesis fail")
		}
		e = bucket.Put([]byte("l"), blockBytes)
		if e != nil {
			log.Panic("put genesis last fail")
		}
		return nil
	})

	if err != nil {
		log.Panic("Err in creating genesis")
	}
	return &Blockchain{genesisBlock.Hash, db}
}

func getBlockChainFromDB() *Blockchain {
	db, err := bolt.Open(DBFullName, 0600, nil)
	if err != nil {
		log.Fatal("Not opening db")
		return nil
	}

	var lastBlockBytes []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Persistence.TableName))
		if bucket != nil {
			lastBlockBytes = bucket.Get([]byte("l"))
		}
		return nil
	})

	if err != nil {
		log.Panic("Fail to read block chain")
	}

	return &Blockchain{lastBlockBytes, db}
}

func (bc *Blockchain) AddNewBlock(data string) {
	err := bc.blockDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Persistence.TableName))
		lastBlockBytes := bucket.Get(bc.tip)
		lastBlock := DeSerializeBlock(lastBlockBytes)
		block := NewBlock(data, lastBlock.Hash, lastBlock.Height+1)
		if bucket != nil {
			blockBytes := block.Serialize()
			e := bucket.Put(block.Hash, blockBytes)
			if e != nil {
				log.Panic("Put error")
			}
			e = bucket.Put([]byte("l"), blockBytes)
			bc.tip = block.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func dbExist() bool {
	_, err := os.Stat(Persistence.DBName)
	return !os.IsNotExist(err)
}
