package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"tml-sync/server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetSyncStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/sync", GetSyncStatus)

	// Ensure Mods directory exists for the test and clean it up afterward
	_ = os.MkdirAll(filepath.Dir(models.SyncFile), 0755)
	defer os.RemoveAll(filepath.Dir(models.SyncFile))

	// Create dummy sync.json
	data := models.SyncData{
		Mods: []models.ModMetadata{
			{Name: "Calamity", Version: "1.0", Hash: "abc"},
		},
		EnabledJSONHash: "xyz",
	}
	content, _ := json.Marshal(data)
	_ = os.WriteFile(models.SyncFile, content, 0644)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sync", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SyncData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Mods, 1)
	assert.Equal(t, "Calamity", response.Mods[0].Name)
}
