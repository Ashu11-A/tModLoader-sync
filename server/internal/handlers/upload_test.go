package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"tml-sync/server/internal/models"
	"tml-sync/shared/pkg"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUploadMod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/upload", UploadMod)

	// Create a dummy mod file
	modName := "TestMod"
	modVersion := "1.0"
	modContent := []byte("dummy mod data")
	_ = os.MkdirAll(models.ModsDir, 0755)
	modPath := filepath.Join(models.ModsDir, modName+".tmod")
	_ = os.WriteFile(modPath, modContent, 0644)
	
	hash, _ := pkg.CalculateSHA256(modPath)
	defer os.RemoveAll(models.ModsDir)
	defer os.Remove(models.SyncFile)

	// Prepare multipart request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	
	part, _ := writer.CreateFormFile("mod", modName+".tmod")
	_, _ = io.Copy(part, bytes.NewReader(modContent))
	
	_ = writer.WriteField("name", modName)
	_ = writer.WriteField("version", modVersion)
	_ = writer.WriteField("hash", hash)
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	// Verify metadata update
	content, _ := os.ReadFile(models.SyncFile)
	var data models.SyncData
	json.Unmarshal(content, &data)
	assert.Len(t, data.Mods, 1)
	assert.Equal(t, modName, data.Mods[0].Name)
	assert.Equal(t, hash, data.Mods[0].Hash)
}
