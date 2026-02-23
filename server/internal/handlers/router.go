package handlers

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
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
}

func (router *Router) Start(address string) {
	router.engine.Run(address)
}
