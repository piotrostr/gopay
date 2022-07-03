package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/piotrostr/gopay/schema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Get() *gorm.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
	)
	dialector := postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&schema.Transaction{})
	if err != nil {
		panic(err)
	}

	if err != nil {
		fmt.Printf("err connecting to db: %s\n", err.Error())
		return nil
	}
	return db
}
