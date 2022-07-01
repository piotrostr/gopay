package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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
	godotenv.Load("../.env")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
	)
	fmt.Println(dsn)
	dialector := postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true})
	db, err := gorm.Open(dialector, &gorm.Config{})
	db.AutoMigrate(&TransactionClientside{}, &StateOnChain{})
	if err != nil {
		fmt.Printf("err connecting to db: %s\n", err.Error())
		return nil
	}
	return db
}
