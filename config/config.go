package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	logger     *Logger
	DBFILE     string
	JWT_SECRET string
)

func Init() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBFILE = os.Getenv("DBFILE")
	JWT_SECRET = os.Getenv("JWT_SECRET")

	fmt.Printf("DBFILE: %s\n", DBFILE)
	fmt.Printf("JWT_SECRET: %s\n", JWT_SECRET)

	db, err = InitializeDB()
	if err != nil {
		return fmt.Errorf("Error initializing database: %v", err)
	}
	return nil
}

func GetSQLite() *gorm.DB {
	return db
}

func GetLogger(p string) *Logger {
	logger = NewLogger(p)
	return logger
}
