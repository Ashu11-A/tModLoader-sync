package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"tml-sync/shared/pkg"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/version", GetVersion)

	// Create dummy version file
	_ = os.MkdirAll("logs", 0755)
	versionPath := filepath.Join("logs", "tml_version.conf")
	_ = os.WriteFile(versionPath, []byte("1.4.4.9"), 0644)
	defer os.RemoveAll("logs")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/version", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, pkg.Version, response["server_version"])
	assert.Equal(t, "1.4.4.9", response["tml_version"])
}
