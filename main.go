package main

import (
	"github.com/agez0s/todoGo/config"
	"github.com/agez0s/todoGo/router"
)

var (
	logger *config.Logger
)

func main() {
	logger = config.GetLogger("main")

	err := config.Init()
	if err != nil {
		logger.Error("Error initializing config: %s", err)
		return
	}

	router.Initialize()
}
