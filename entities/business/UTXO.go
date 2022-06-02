package business

type UTXO struct {
	TxId  []byte    // txId from which
	Index int       // index of this output out of multiple output
	TxOut *TxOutput // corresponding output
}
