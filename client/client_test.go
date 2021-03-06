package client

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetsClientRight(t *testing.T) {
	client := Get()
	if client.Eth == nil {
		t.Error("client.Eth is nil")
	}
	if client.privateKey == nil {
		t.Error("client.privateKey is nil")
	}
	if client.address == (common.Address{}) {
		t.Error("client.address is empty")
	}
}

func TestGetBalance(t *testing.T) {
	client := Get()

	balance, err := client.Balance()
	if err != nil {
		t.Error(err)
	}

	if balance == big.NewInt(0) {
		t.Error("balance is 0")
	}
}

func TestSendsTx(t *testing.T) {
	client := Get()
	tx, err := client.SendTx()
	if err != nil {
		t.Error(err)
	}
	if tx == nil {
		t.Error("tx is nil")
	}
	PrintJson(tx)
}
