package main

import (
	"testing"
)

// All of the methods should use the same context and ethclient, probably will
// make them methods of a struct

func Fixture() {}

// retrieve transaction, store it into the transaction struct
func TestGetTx(t *testing.T) {}

// wait for the tx to finish, put in struct it as finished if it is ready
func TestWaitTx(t *testing.T) {}

func SyncTx(t *testing.T) {}

// check if the payload is corrent and the hash in the smart contract matches
// the payload from the clientside as well
func TestCheckPayload(t *testing.T) {}
