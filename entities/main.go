package main

import (
	"fmt"
	"smallPublicChain/entities/business"
)

func main() {
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
