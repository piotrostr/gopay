package client

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func TestRunNode(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := RunNode()
	if err != nil {
		t.Error(err)
	}
	<-ctx.Done()
}

func TestGetsClientRight(t *testing.T) {
	client := Get()
	if client.client == nil {
		t.Error("client.client is nil")
	}
	if client.privateKey == nil {
		t.Error("client.privateKey is nil")
	}
	if client.address == (common.Address{}) {
		t.Error("client.address is empty")
	}
}
