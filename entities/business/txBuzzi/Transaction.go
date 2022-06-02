package txBuzzi

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"log"
	"smallPublicChain/entities/business"
	"time"
)

type Transaction struct {
	TxId   []byte
	TxIns  []*TxInput
	TxOuts []*TxOutput
}

// There are two types of transactions
// one is created with mining - coinbase transaction (no input) or genesis block
// ont is created when actual transactions happen

func NewCoinbaseTransaction(address string) *Transaction {
	txInput := TxInput{[]byte{}, -1, "Coinbase"}
	txOutput := TxOutput{10, address}
	txCoinbase := Transaction{[]byte{}, []*TxInput{&txInput}, []*TxOutput{&txOutput}}
	txCoinbase.HashTransactionId() // set transaction id
	return &txCoinbase
}

// HashTransactionId Set the id of transaction, which is the hash of tx and current time
func (tx *Transaction) HashTransactionId() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	buffBytes := bytes.Join([][]byte{
		business.GetIntBytes(time.Now().Unix()),
		buffer.Bytes(),
	}, []byte{})
	hash := sha1.Sum(buffBytes)
	tx.TxId = hash[:]
}
