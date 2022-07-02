package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/piotrostr/gopay/client"
	_ "github.com/piotrostr/gopay/db"
)

type Address string

type Payment struct {
	Id          string  `json:"id,omitempty"`
	Amount      big.Int `json:"amount,omitempty"`
	From        Address `json:"from,omitempty"`
	To          Address `json:"to,omitempty"`
	ContentHash string  `json:"content_hash,omitempty"`
	TxHash      string  `json:"tx_hash,omitempty"`
	Successful  bool    `json:"successful,omitempty"`
	IsPending   bool    `json:"is_pending,omitempty"`
	Error       error   `json:"error,omitempty"`
}

type Payments map[string]Payment

var payments = Payments{}

var ctx = context.Background()

func Verify(p *Payment, tx *types.Transaction) bool {
	if p.ContentHash != tx.Hash().Hex() {
		p.Error = fmt.Errorf("content hash mismatch")
		return false
	}

	// TODO verify sender address as well

	if p.To != Address(tx.To().Hex()) {
		p.Error = fmt.Errorf("to address mismatch")
		return false
	}
	if p.Amount.Cmp(tx.Value()) != 0 {
		p.Error = fmt.Errorf("amount mismatch")
		return false
	}
	return true
}

// make an idempotent function that
func (p *Payment) UpdateStatus() error {
	c := client.Get()

	// checks if tx exists
	tx, isPending, err := c.Eth.TransactionByHash(
		ctx,
		common.HexToHash(p.TxHash),
	)

	// if err include the error in p.Error
	if err == ethereum.NotFound {
		p.Error = err
		return fmt.Errorf("tx not found")
	} else if err != nil {
		p.Error = err
		return err
	}

	valid := Verify(p, tx)
	if !valid {
		return p.Error
	}

	// if not finished do nothing mark as pending
	p.IsPending = isPending

	// if finished mark and everything in tact, mark as successful
	if !isPending {
		p.Successful = true
	}

	return nil
}

func (p *Payment) Commit() error {
	payments[p.Id] = *p
	return nil
}

func (p *Payment) CreateResponse(c *gin.Context) {
	err := p.UpdateStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = p.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if p.Successful {
		c.JSON(http.StatusOK, gin.H{
			"id":     p.Id,
			"status": "paid",
		})
	} else if p.IsPending {
		c.JSON(http.StatusOK, gin.H{
			"id":     p.Id,
			"status": "pending",
		})
	} else if p.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"id":     p.Id,
			"status": "error",
			"error":  p.Error.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"id":     p.Id,
			"status": "error",
			"error":  "unknown error",
		})
	}
}

func SetupRouter() (r *gin.Engine, err error) {
	gin.SetMode(gin.ReleaseMode)

	r.POST("/create", func(c *gin.Context) {
		// get the request body
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var payload Payment

		// validate payload
		err = c.ShouldBindJSON(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// unmarshal payload to json
		err = json.Unmarshal(body, &payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// create the payment struct
		p := Payment{
			Id:          uuid.New().String(),
			Amount:      payload.Amount,
			From:        payload.From,
			To:          payload.To,
			ContentHash: payload.ContentHash,
			IsPending:   true,
		}

		p.CreateResponse(c)
	})

	r.GET("/:paymentId", func(c *gin.Context) {
		txHash := c.Param("paymentId")

		// retrieve the payment, update and return
		p, ok := payments[txHash]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		p.CreateResponse(c)
	})
	return
}

func main() {
	router, err := SetupRouter()
	if err != nil {
		log.Fatalln(err)
	}

	port := flag.String("port", "80", "port to listen on")
	flag.Parse()

	log.Printf("running on :%s\n", *port)
	if err := router.Run(fmt.Sprintf(":%s", *port)); err != nil {
		log.Fatalln(err)
	}
}
