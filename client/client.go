package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
		panic(err)
	}
}

func printJson(data any) {
	barr, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(barr))
}

func CheckEnv() {
	if os.Getenv("PRIVATE_KEY") == "" {
		panic("PRIVATE_KEY is not set")
	}
	if os.Getenv("RPC_URL") == "" {
		fmt.Println("RPC_URL is not set, will dial http://localhost:8545")
	}
}

func Get(url ...string) *Client {
	CheckEnv()

	if len(url) == 0 {
		if os.Getenv("RPC_URL") == "" {
			url = append(url, "http://localhost:8545")
		} else {
			url = append(url, os.Getenv("RPC_URL"))
		}
	}
	client, err := ethclient.Dial(url[0])
	checkErr(err)

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	checkErr(err)

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &Client{client, privateKey, address}
}

func RunNode() error {
	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error)
	defer cancel()
	cmd := exec.CommandContext(
		ctx, "geth", "--dev", "--http", "--http.api", "eth,web3,personal,net", "--http.corsdomain", "http://remix.ethereum.org")
	go func() {
		cmd.Stdout = os.Stdout
		fmt.Printf("running node")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}()
	return <-errChan
}

func (c *Client) Balance() *big.Int {
	balance, err := c.client.BalanceAt(ctx, c.address, nil)
	checkErr(err)
	fmt.Println("Balance: ", balance)
	return balance
}

/* this method will be a mock to run on ganache */
func (c *Client) SendTx() *types.Transaction {
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
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, c.privateKey)
	checkErr(err)

	err = c.client.SendTransaction(ctx, signedTx)
	checkErr(err)
	return signedTx
}

func (c *Client) GetTx(txAddress string) {
	common.HexToAddress(txAddress)
	tx, isPending, err := c.client.TransactionByHash(ctx, common.HexToHash(txAddress))
	checkErr(err)
	if isPending {
		fmt.Println("Pending")
		time.Sleep(time.Second * 1)
		c.GetTx(txAddress)
	} else {
		fmt.Println("Not pending")
		printJson(tx)
	}
}
