package store

type TransactionClientside struct {
	TxHash       string
	Mined        bool
	ContentHash  string
	PayeeAddress string
	OnChain      *StateOnChain
}

type StateOnChain struct {
	TxHash       string
	ContentHash  string
	PayeeAddress string
	PayeePaid    string
}

func GetTx(txHash string) (*TransactionClientside, error) {
	return nil, nil
}

func (t *TransactionClientside) Wait() {}

func (t *TransactionClientside) Sync(txHash string) {}
