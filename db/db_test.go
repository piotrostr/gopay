package db

import (
	"testing"
)

func TestStore(t *testing.T) {
	db := Get()
	if db == nil {
		t.Errorf("db is nil")
	}
}
