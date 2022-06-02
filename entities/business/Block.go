package business

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Height        int64
	PrevBlockHash []byte
	Txs           []*Transaction
	TimeStamp     int64
	Hash          []byte
	Nuance        int64
}

func NewBlock(txs []*Transaction, prevBlock []byte, height int64) *Block {
	currentTime := time.Now().Unix()
	block := &Block{height, prevBlock, txs, currentTime, nil, 0}
	pow := NewProofOfWork(block)
	hash, nuance := pow.run()
	block.Hash = hash
	block.Nuance = nuance
	return block
}

func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, make([]byte, 32, 32), 0)
}

func (block *Block) HashTransactions() []byte {
	// combine transactions' ids, which are hashes of that transaction
	var txHashes [][]byte
	for _, tx := range block.Txs {
		txHashes = append(txHashes, tx.TxId)
	}
	allBytes := bytes.Join(txHashes, []byte{})
	hash := sha256.Sum256(allBytes)
	return hash[:]
}

func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

func DeSerializeBlock(blockBytes []byte) *Block {
	var block Block
	var reader = bytes.NewReader(blockBytes)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
