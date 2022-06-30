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

func Get() *Client {
	client, err := ethclient.Dial(os.Getenv("RPC_URL"))
	checkErr(err)

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	checkErr(err)

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &Client{client, privateKey, address}
}

func (c *Client) RunNode() error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd := exec.CommandContext(
		ctx, "geth", "--dev", "--http", "--http.api", "eth,web3,personal,net", "--http.corsdomain", "http://remix.ethereum.org")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	// TODO make the node shut down when the context is cancelled
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Balance() {
	balance, err := c.client.BalanceAt(ctx, c.address, nil)
	checkErr(err)
	fmt.Println("Balance: ", balance)
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