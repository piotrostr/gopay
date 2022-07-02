package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const DEAD = "0x0000000000000000000000000000000000000000"

var ctx = context.Background()

type Client struct {
	Eth        *ethclient.Client
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

func PrintJson(data any) {
	barr, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(barr))
}

func Get() *Client {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial(os.Getenv("RPC_URL"))
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &Client{client, privateKey, address}
}

func (c *Client) Balance() (*big.Int, error) {
	balance, err := c.Eth.BalanceAt(ctx, c.address, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

/* this method will be a mock to run on ganache */
func (c *Client) SendTx() (*types.Transaction, error) {
	nonce, err := c.Eth.NonceAt(ctx, c.address, nil)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(DEAD),
		big.NewInt(1),
		100000,
		big.NewInt(1),
		nil,
	)
	signedTx, err := types.SignTx(
		tx,
		types.NewEIP155Signer(big.NewInt(1337)),
		c.privateKey,
	)
	if err != nil {
		return nil, err
	}

	err = c.Eth.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (c *Client) GetTx(txAddress string) (*types.Transaction, error) {
	common.HexToAddress(txAddress)
	tx, isPending, err := c.Eth.TransactionByHash(
		ctx,
		common.HexToHash(txAddress),
	)
	if err != nil {
		return nil, err
	}

	if isPending {
		time.Sleep(time.Second * 1)
		return c.GetTx(txAddress)
	}

	return tx, nil
}
