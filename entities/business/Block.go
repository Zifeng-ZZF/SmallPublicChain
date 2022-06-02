package business

import (
	"bytes"
	"encoding/gob"
	"log"
	"smallPublicChain/entities/business/txBuzzi"
	"time"
)

type Block struct {
	Height        int64
	PrevBlockHash []byte
	Txs           []*txBuzzi.Transaction
	TimeStamp     int64
	Hash          []byte
	Nuance        int64
}

func NewBlock(txs []*txBuzzi.Transaction, prevBlock []byte, height int64) *Block {
	currentTime := time.Now().Unix()
	block := &Block{height, prevBlock, txs, currentTime, nil, 0}
	pow := NewProofOfWork(block)
	hash, nuance := pow.run()
	block.Hash = hash
	block.Nuance = nuance
	return block
}

func CreateGenesisBlock(txs []*txBuzzi.Transaction) *Block {
	return NewBlock(txs, make([]byte, 32, 32), 0)
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
