package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/piotrostr/gopay/schema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Get() *gorm.DB {
	godotenv.Load("../.env")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
	)
	fmt.Println(dsn)
	dialector := postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true})
	db, err := gorm.Open(dialector, &gorm.Config{})
	db.AutoMigrate(&schema.Transaction{})
	if err != nil {
		fmt.Printf("err connecting to db: %s\n", err.Error())
		return nil
	}
	return db
}
