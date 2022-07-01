package db

import (
	"testing"

	"github.com/piotrostr/gopay/schema"
)

func TestStore(t *testing.T) {
	db := Get()
	if db == nil {
		t.Errorf("db is nil")
	}

	db.Create(&schema.Transaction{Hash: "asdf", Type: "asdf"})
}
