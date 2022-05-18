package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"smallPublicChain/entities/business"
)

const DBPath = "entities/Persistence/blockChain"

func main() {
	//testDB()
	testSerialize()
}

func testBlockChainBasic() {
	blockChain := business.CreateBlockChain("Genesis block")
	blocks := &blockChain.Chain

	blockChain.AddNewBlock("a sends 500 to b", (*blocks)[len(*blocks)-1].Hash,
		(*blocks)[len(*blocks)-1].Height+1)
	blockChain.AddNewBlock("c sends 10000 to b", (*blocks)[len(*blocks)-1].Hash,
		(*blocks)[len(*blocks)-1].Height+1)
	blockChain.AddNewBlock("a sends 100 to c", (*blocks)[len(*blocks)-1].Hash,
		(*blocks)[len(*blocks)-1].Height+1)

	for i, v := range *blocks {
		fmt.Printf("\nBlockNum: %d\n", i)
		fmt.Printf("PrevHash=%x\n", v.PrevBlockHash)
		fmt.Printf("Hash=%x\n", v.Hash)
		fmt.Printf("Height=%d\n", v.Height)
		fmt.Printf("Data=%s\n", v.Data)
	}
}

func testDB() {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// writing to DB
	err = db.Update(func(tx *bolt.Tx) error {
		bName := []byte("testBucket")
		b, e := tx.CreateBucket(bName)
		if e != nil {
			if e != bolt.ErrBucketExists {
				return fmt.Errorf("create bucket:%s", e)
			}
			b = tx.Bucket(bName)
		}

		if b != nil {
			key := []byte("l") // key & val are bytes array
			value := []byte("a sends 10 yuan to b")
			e = b.Put(key, value)
			if e != nil {
				log.Panicf("put kv pair:%s", e)
			}
		} else {
			log.Panic("bucket nil")
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// reading from DB
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("testBucket"))
		if b != nil {
			value := b.Get([]byte("l"))
			fmt.Println(value)
			fmt.Printf("value=%s\n", value)
			value2 := b.Get([]byte("NotExist"))
			fmt.Println(value2)
			fmt.Printf("value=%s\n", value2)
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	// cursor, iterating from DB
	err = db.View(func(tx *bolt.Tx) error {
		fmt.Println("Test cursor:")
		b := tx.Bucket([]byte("testBucket"))
		cursor := b.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			fmt.Printf("Key=%s, Value=%s\n", k, v)
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func testSerialize() {
	block := business.NewBlock("test", make([]byte, 32, 32), 0)
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	// adding serialized block
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("s"))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("s"))
			if err != nil {
				log.Panic("Fail creating bucket")
				return err
			}
		}
		blockBytes := block.Serialize()
		err = bucket.Put([]byte("l"), blockBytes)
		if err != nil {
			log.Panic("Fail putting")
			return err
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
		return
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("s"))
		if bucket == nil {
			log.Panic("No such bucket")
			return bolt.ErrBucketExists
		}
		blockBytes := bucket.Get([]byte("l"))
		if blockBytes == nil {
			log.Panic("No such key")
			return bolt.ErrInvalid
		}
		block := business.DeSerializeBlock(blockBytes)
		fmt.Printf("deserialized block hash= %x\n", block.Hash)
		return nil
	})
}
