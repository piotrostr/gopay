package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TransactionClientside struct {
	gorm.Model
	TxHash      string        `json:"tx_hash,omitempty"`
	Mined       bool          `json:"mined,omitempty"`
	ContentHash string        `json:"content_hash,omitempty"`
	OnChain     *StateOnChain `json:"on_chain,omitempty"`
}

type StateOnChain struct {
	gorm.Model
	TxHash       string `json:"tx_hash,omitempty"`
	ContentHash  string `json:"content_hash,omitempty"`
	PayeeAddress string `json:"payee_address,omitempty"`
	PayeePaid    string `json:"payee_paid,omitempty"`
}

func Get() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("err connecting to db: %s\n", err.Error())
		return nil
	}
	return db
}
