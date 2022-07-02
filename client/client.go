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
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
	checkErr(err)

	client, err := ethclient.Dial(os.Getenv("RPC_URL"))
	checkErr(err)

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	checkErr(err)

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &Client{client, privateKey, address}
}

func (c *Client) Balance() *big.Int {
	balance, err := c.client.BalanceAt(ctx, c.address, nil)
	checkErr(err)
	return balance
}

/* this method will be a mock to run on ganache */
func (c *Client) SendTx() (*types.Transaction, error) {
	nonce, err := c.client.NonceAt(ctx, c.address, nil)
	checkErr(err)

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
	checkErr(err)

	err = c.client.SendTransaction(ctx, signedTx)
	checkErr(err)
	return signedTx, nil
}

func (c *Client) GetTx(txAddress string) (*types.Transaction, error) {
	common.HexToAddress(txAddress)
	tx, isPending, err := c.client.TransactionByHash(
		ctx,
		common.HexToHash(txAddress),
	)
	checkErr(err)
	if isPending {
		fmt.Println("Pending tx..")
		time.Sleep(time.Second * 1)
		return c.GetTx(txAddress)
	} else {
		PrintJson(tx)
	}
	return tx, nil
}
