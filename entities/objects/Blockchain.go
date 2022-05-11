package objects

type Blockchain struct {
	Chain []*Block
}

func CreateBlockChain(genesisDataStr string) *Blockchain {
	genesisBlock := CreateGenesisBlock(genesisDataStr)
	blockchain := &Blockchain{[]*Block{genesisBlock}}
	return blockchain
}

func (bc *Blockchain) AddNewBlock(data string, prevBlock []byte, height int64) {
	block := NewBlock(data, prevBlock, height)
	bc.Chain = append(bc.Chain, block)
}
