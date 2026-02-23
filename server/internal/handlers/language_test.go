package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"tml-sync/server/internal/i18n"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/language", GetLanguage)

	// Set language
	_ = os.MkdirAll("logs", 0755)
	_ = os.WriteFile(i18n.LangFile, []byte("pt"), 0644)
	defer os.RemoveAll("logs")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/language", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pt", w.Body.String())
}
