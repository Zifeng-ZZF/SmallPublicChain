package business

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"smallPublicChain/entities/Persistence"
	"smallPublicChain/entities/business/txBuzzi"
)

const DBFullName = Persistence.DBName

type Blockchain struct {
	Tip     []byte // last block's Hash
	BlockDB *bolt.DB
}

func CreateBlockChain(address string) *Blockchain {
	if !dbExist() {
		fmt.Printf("Creating genesis block.. DB not exist\n\n")
		return initializeBlockChain(address)
	}
	fmt.Printf("Db exists. Getting block from DB...\n\n")
	return getBlockChainFromDB()
}

func initializeBlockChain(address string) *Blockchain {
	coinbase := txBuzzi.NewCoinbaseTransaction(address)
	genesisBlock := CreateGenesisBlock([]*txBuzzi.Transaction{coinbase})
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

func (bc *Blockchain) AddNewBlock(txs []*txBuzzi.Transaction) {
	err := bc.BlockDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(Persistence.TableName))
		lastBlockBytes := bucket.Get(bc.Tip)
		lastBlock := DeSerializeBlock(lastBlockBytes)
		block := NewBlock(txs, lastBlock.Hash, lastBlock.Height+1)
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
		fmt.Printf("Transactions:\n")
		for _, tx := range cur.Txs {
			fmt.Printf("\t\tTransaction ID=%x\n", tx.TxId)
			fmt.Printf("\t\tInputs:\n")
			for _, in := range tx.TxIns {
				fmt.Printf("\t\t\ttxId=%x\n", in.TxId)
				fmt.Printf("\t\t\toutIdx=%d\n", in.Output)
				fmt.Printf("\t\t\tscriptSig=%s\n", in.ScriptSig)
			}
			fmt.Printf("\t\tOutputs:\n")
			for _, out := range tx.TxOuts {
				fmt.Printf("\t\t\tvalue=%d\n", out.Value)
				fmt.Printf("\t\t\tscriptPubKey=%s\n", out.ScriptPubKey)
			}
		}

		// iterate until prev hash is zero (which is genesis)
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
