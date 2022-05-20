package business

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"smallPublicChain/entities/Persistence"
)

const DBFullName = Persistence.DBName

type Blockchain struct {
	Tip     []byte // last block's Hash
	BlockDB *bolt.DB
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
		log.Fatal(err)
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
		log.Fatal(err)
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
		log.Panic(err)
	}

	lastBlock := DeSerializeBlock(lastBlockBytes)
	return &Blockchain{lastBlock.Hash, db}
}

func (bc *Blockchain) AddNewBlock(data string) {
	err := bc.BlockDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Persistence.TableName))
		lastBlockBytes := bucket.Get(bc.Tip)
		lastBlock := DeSerializeBlock(lastBlockBytes)
		block := NewBlock(data, lastBlock.Hash, lastBlock.Height+1)
		if bucket != nil {
			blockBytes := block.Serialize()
			e := bucket.Put(block.Hash, blockBytes)
			if e != nil {
				log.Panic("Put error")
			}
			e = bucket.Put([]byte("l"), blockBytes)
			bc.Tip = block.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (bc *Blockchain) GetIterator() *BlockChainIterator {
	return &BlockChainIterator{CurrentHash: bc.Tip, BlockDB: bc.BlockDB}
}

func dbExist() bool {
	_, err := os.Stat(DBFullName)
	return !os.IsNotExist(err)
}

func (bc *Blockchain) PrintChain() {
	it := bc.GetIterator()
	var hashInt = new(big.Int)
	var zeroInt = big.NewInt(0)
	for {
		cur := it.Next()
		fmt.Printf("\nPrevHash=%x\n", cur.PrevBlockHash)
		fmt.Printf("CurrHash=%x\n", cur.Hash)
		fmt.Printf("Height=%d\n", cur.Height)
		fmt.Printf("TimeStamp=%d\n", cur.TimeStamp)
		fmt.Printf("Nuance=%d\n", cur.Nuance)
		fmt.Printf("Data=%s\n", cur.Data)

		hashInt.SetBytes(cur.PrevBlockHash)
		if hashInt.Cmp(zeroInt) == 0 {
			break
		}
	}
}

func GetBlockChain() *Blockchain {
	if !dbExist() {
		fmt.Printf("No BlockChain Object On Disk")
		return nil
	}
	return getBlockChainFromDB()
}
