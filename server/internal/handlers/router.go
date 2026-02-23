package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ShutdownChan = make(chan struct{})
)

type Router struct {
	engine *gin.Engine
	server *http.Server
}

func New() *Router {
	return &Router{
		engine: gin.Default(),
	}
}

func (router *Router) Register() {
	router.engine.GET("/version", GetVersion)

	group := router.engine.Group("/v1")
	
	group.GET("/language", GetLanguage)
	group.GET("/update", Update)
	group.GET("/sync", GetSyncStatus)
	group.POST("/upload", UploadMod)
	group.POST("/enabled", UploadEnabledJSON)
	group.POST("/stop", Stop)
}

func (router *Router) Start(address string) {
	router.server = &http.Server{
		Addr:    address,
		Handler: router.engine,
	}

	go func() {
		if err := router.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for shutdown signal
	<-ShutdownChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := router.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
