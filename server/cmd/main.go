package main

import (
	"os"
	"tml-sync/server/configs"
	"tml-sync/server/internal/handlers"
	"tml-sync/server/internal/i18n"
	"tml-sync/server/internal/ui"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.Load()

	// 1. Language Setup
	i18n.Setup()

	// 2. Gin Mode Configuration
	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. Stylish Banner
	ui.PrintBanner(cfg.Address())

	router := handlers.New()
	router.Register()
	router.Start(cfg.Address())
}
