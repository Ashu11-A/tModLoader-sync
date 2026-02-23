// Package handlers contains the HTTP handler functions for the tModLoader-sync server.
package handlers

import (
	"net/http"
	"tml-sync/server/internal/i18n"

	"github.com/gin-gonic/gin"
)

// GetLanguage returns the server's configured language.
func GetLanguage(ctx *gin.Context) {
	ctx.String(http.StatusOK, i18n.GetLanguage())
}
