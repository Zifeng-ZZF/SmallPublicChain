package business

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// ProofOfWork 经济学上认为，理性的人都是逐利的，PoW抑制了节点的恶意动机
type ProofOfWork struct {
	Block  *Block   // block to validate
	Target *big.Int // target hash
}

// TargetBit difficulty
const TargetBit = 16

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target, 256-TargetBit)
	return &ProofOfWork{block, target}
}

func (pow *ProofOfWork) run() ([]byte, int64) {
	nuance := 0
	hashInt := new(big.Int)
	var hash [32]byte
	for true {
		dataBytes := pow.prepareBlockDataWithNuance(int64(nuance))
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		if pow.Target.Cmp(hashInt) > 0 {
			fmt.Printf("nuance=%d, hash=%x\n", nuance, hash)
			break
		}
		nuance++
	}
	return hash[:], int64(nuance)
}

func (pow *ProofOfWork) prepareBlockDataWithNuance(nuance int64) []byte {
	b := pow.Block
	heightBytes := GetIntBytes(b.Height)
	timeBytes := GetIntBytes(b.TimeStamp)
	nuanceBytes := GetIntBytes(nuance)
	data := bytes.Join([][]byte{
		heightBytes,
		timeBytes,
		nuanceBytes,
		b.PrevBlockHash,
		b.HashTransactions(),
	}, []byte{})
	return data
}

func (pow *ProofOfWork) isValid() bool {
	hashInt := new(big.Int)
	hashInt.SetBytes(pow.Block.Hash)
	return pow.Target.Cmp(hashInt) > 0
}
