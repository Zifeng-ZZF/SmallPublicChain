package txBuzzi

type TxInput struct {
	TxId      []byte // transaction id
	Output    int    // the index to previous tx's output, used by this input
	ScriptSig string // signature and pub keys to satisfy ScriptPubKey
}

// UnlockScript temporary: comparing address to sig script (no script yet)
func (ti *TxInput) UnlockScript(address string) bool {
	return address == ti.ScriptSig
}
