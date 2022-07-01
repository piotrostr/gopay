package db

import (
	"fmt"
	"testing"
)

func TestStore(t *testing.T) {
	db := Get()
	if db == nil {
		t.Errorf("db is nil")
	}

	db.Create(&Transaction{TxHash: "asdf", Mined: true})

	var tx Transaction
	db.First(&tx)
	fmt.Printf("%+v\n", tx)

	if tx.TxHash != "asdf" {
		t.Errorf("tx hash is not asdf")
	}
}
