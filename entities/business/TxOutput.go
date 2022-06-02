package business

// TxOutput is indivisible
type TxOutput struct {
	Value        int64  // Value of the UTXO
	ScriptPubKey string // bitcoin address is the encodings of ScriptPubKey
}

// UnlockScript temporary: comparing address to pubKey script (no script yet)
func (to *TxOutput) UnlockScript(address string) bool {
	return address == to.ScriptPubKey
}
