package config

import (
	"os"

	"github.com/agez0s/todoGo/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	logger := GetLogger("sqlite")
	// pathDb := "./db/main.db"
	pathDb := "./db/" + DBFILE

	_, err := os.Stat(pathDb)
	if os.IsNotExist(err) {
		logger.Info("Criando banco de dados...")
		err = os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			return nil, err
		}
		file, err := os.Create(pathDb)
		if err != nil {
			return nil, err
		}
		file.Close()
	}
	db, err := gorm.Open(sqlite.Open(pathDb), &gorm.Config{})
	if err != nil {
		logger.ErrorF("SQLite error: %v", err)
		return nil, err
	}
	err = db.AutoMigrate(&schema.User{}, &schema.Todo{})

	if err != nil {
		logger.ErrorF("SQLite migration error: %v", err)
		return nil, err
	}
	return db, nil

}
