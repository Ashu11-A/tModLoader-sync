package handlers

import (
	"net/http"
	"tml-sync/server/internal/update"

	"github.com/gin-gonic/gin"
)

func Update(ctx *gin.Context) {
	targetVersion := ctx.Query("version")
	if targetVersion == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing version parameter"})
		return
	}

	// This is a blocking/restarting operation
	// We run it in a goroutine to respond to the client first
	go func() {
		err := update.TriggerServerUpdate(targetVersion)
		if err != nil {
			// Since it's in a goroutine, we can only log it here
			// In a real scenario, we might want to notify someone
			println("Update error:", err.Error())
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Updating server to v" + targetVersion + "...",
	})
}
