package schema

import (
	"gorm.io/gorm"
)

type StateOnChain struct {
	gorm.Model
	TxHash       string `json:"tx_hash,omitempty"`
	ContentHash  string `json:"content_hash,omitempty"`
	PayeeAddress string `json:"payee_address,omitempty"`
	PayeePaid    string `json:"payee_paid,omitempty"`
}
