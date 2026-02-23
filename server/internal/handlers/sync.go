package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"tml-sync/server/internal/models"
	"tml-sync/shared/pkg"

	"github.com/gin-gonic/gin"
)

// GetSyncStatus returns the list of mods and the hash of enabled.json on the server.
func GetSyncStatus(ctx *gin.Context) {
	// Ensure the directory exists
	if _, err := os.Stat(filepath.Dir(models.SyncFile)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(models.SyncFile), 0755)
	}

	var data models.SyncData

	content, err := os.ReadFile(models.SyncFile)
	if err == nil {
		json.Unmarshal(content, &data)
	} else {
		data.Mods = []models.ModMetadata{}
	}

	// Hash enabled.json if it exists
	enabledPath := filepath.Join(models.ModsDir, "enabled.json")
	if hash, err := pkg.CalculateSHA256(enabledPath); err == nil {
		data.EnabledJSONHash = hash
	}

	ctx.JSON(http.StatusOK, data)
}
