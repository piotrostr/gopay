package schema

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Type         string `json:"type,omitempty"`
	Nonce        string `json:"nonce,omitempty"`
	GasPrice     string `json:"gas_price,omitempty"`
	MaxFeePerGas string `json:"max_fee_per_gas,omitempty"`
	Gas          string `json:"gas,omitempty"`
	Value        string `json:"value,omitempty"`
	Input        string `json:"input,omitempty"`
	V            string `json:"v,omitempty"`
	R            string `json:"r,omitempty"`
	S            string `json:"s,omitempty"`
	To           string `json:"to,omitempty"`
	Hash         string `json:"hash,omitempty"`
}
