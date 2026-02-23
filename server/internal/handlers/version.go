package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tml-sync/shared/pkg"

	"github.com/gin-gonic/gin"
)

// GetVersion returns the current version of the server and the tModLoader.
func GetVersion(ctx *gin.Context) {
	var tmlVersion string
	
	directory, err := os.Getwd()
	if err != nil {
		directory = "."
	}

	path := filepath.Join(directory, "logs", "tml_version.conf")
	data, err := os.ReadFile(path)

	if err != nil {
		tmlVersion = "unknown"
	} else {
		tmlVersion = strings.TrimSpace(string(data))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"server_version": pkg.Version,
		"tml_version":    tmlVersion,
	})
}
