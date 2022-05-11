package objects

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Height        int64
	PrevBlockHash []byte
	Data          []byte
	TimeStamp     int64
	Hash          []byte
}

func NewBlock(data string, prevBlock []byte, height int64) *Block {
	currentTime := time.Now().Unix()
	block := &Block{height, prevBlock, []byte(data), currentTime, nil}
	block.SetHash()
	return block
}

func (b *Block) SetHash() {
	heightBytes := getIntBytes(b.Height)
	timeString := strconv.FormatInt(b.TimeStamp, 2)
	timeBytes := []byte(timeString)
	bytesList := [][]byte{
		heightBytes,
		timeBytes,
		b.PrevBlockHash,
		b.Data,
	}
	blockBytes := bytes.Join(bytesList, []byte{})
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:] // deep copy
}

func CreateGenesisBlock(data string) *Block {
	return NewBlock(data, make([]byte, 32, 32), 0)
}
