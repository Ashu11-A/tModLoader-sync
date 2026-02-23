package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Stop(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Server is shutting down...",
	})

	// Give it a tiny bit of time to send the response before triggering shutdown
	go func() {
		time.Sleep(500 * time.Millisecond)
		close(ShutdownChan)
	}()
}
