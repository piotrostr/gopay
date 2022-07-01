package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/piotrostr/gopay/client"
	_ "github.com/piotrostr/gopay/db"
)

// goal is to replicate stripe functionality but more easily

type Address string

type Payment struct {
	Amount      big.Int `json:"amount,omitempty"`
	From        Address `json:"from,omitempty"`
	To          Address `json:"to,omitempty"`
	ContentHash string  `json:"content_hash,omitempty"`
	TxHash      string  `json:"tx_hash,omitempty"`
	Successful  bool    `json:"successful,omitempty"`
}

type Payments map[string]Payment

func (p *Payment) CheckAndUpdate() {
	// TODO check the transaction here and validate tx exists
	// make an idempotent function that
	// checks if tx exists
	// if not finished do nothing mark as pending
	// if finished mark as completed
	// if confirmed, mark as paid
	// if err include the error (make error enum and add codes and docs)
}

func SetupRouter() (r *gin.Engine, err error) {
	gin.SetMode(gin.ReleaseMode)
	payments := Payments{}

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

		// store the payment
		id := uuid.New().String()
		payments[id] = payload

		// return the payment Id
		c.JSON(http.StatusOK, gin.H{
			"id":     id,
			"status": "created",
			"tx":     payload.TxHash,
		})
	})

	r.GET("/:paymentId", func(c *gin.Context) {
		txHash := c.Param("paymentId")
		fmt.Println(txHash)
		payment, ok := payments[txHash]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
			return
		}
		payment.CheckAndUpdate()
		if payment.Successful {
			c.JSON(http.StatusOK, gin.H{
				"id":     txHash,
				"status": "paid",
			})
		}
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
