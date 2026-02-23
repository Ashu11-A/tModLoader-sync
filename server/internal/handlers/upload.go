package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"tml-sync/server/internal/models"
	"tml-sync/shared/pkg"

	"github.com/gin-gonic/gin"
)

// UploadMod receives a mod file and performs version/hash checks
func UploadMod(ctx *gin.Context) {
	file, err := ctx.FormFile("mod")
	if err != nil {
		fmt.Printf("Error receiving mod file: %v\n", err)
		// Try to see if there are ANY files
		form, _ := ctx.MultipartForm()
		if form != nil {
			fmt.Printf("Received files: %v\n", form.File)
			fmt.Printf("Received values: %v\n", form.Value)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("No mod file provided: %v", err),
			"details": "Check server logs for more info",
		})
		return
	}

	modName := ctx.PostForm("name")
	modVersion := ctx.PostForm("version")
	clientHash := ctx.PostForm("hash")

	if modName == "" || modVersion == "" || clientHash == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Mod name, version, and hash are required"})
		return
	}

	// 1. Check if we already have this exact mod (same hash)
	if isAlreadySynced(modName, clientHash) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Mod is already up to date", "status": "skipped"})
		return
	}

	// 2. Ensure Mods directory exists
	if _, err := os.Stat(models.ModsDir); os.IsNotExist(err) {
		os.MkdirAll(models.ModsDir, 0755)
	}

	// 3. Save the file temporarily to verify hash
	tempPath := filepath.Join(models.ModsDir, file.Filename+".tmp")
	if err := ctx.SaveUploadedFile(file, tempPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temporary file"})
		return
	}
	defer os.Remove(tempPath)

	// 4. Verify hash of the uploaded file
	serverHash, err := pkg.CalculateSHA256(tempPath)
	if err != nil || serverHash != clientHash {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Hash mismatch or verification failed"})
		return
	}

	// 5. Finalize the file
	finalPath := filepath.Join(models.ModsDir, file.Filename)
	if err := os.Rename(tempPath, finalPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize mod file"})
		return
	}

	// 6. Update metadata
	updateSyncMetadata(modName, modVersion, serverHash)

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Mod %s uploaded and verified", modName),
		"status":  "success",
	})
}

// UploadEnabledJSON receives the enabled.json file
func UploadEnabledJSON(ctx *gin.Context) {
	file, err := ctx.FormFile("enabled")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	clientHash := ctx.PostForm("hash")
	if clientHash == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Hash is required"})
		return
	}

	// 1. Ensure Mods directory exists
	if _, err := os.Stat(models.ModsDir); os.IsNotExist(err) {
		os.MkdirAll(models.ModsDir, 0755)
	}

	// 2. Save the file temporarily to verify hash
	tempPath := filepath.Join(models.ModsDir, "enabled.json.tmp")
	if err := ctx.SaveUploadedFile(file, tempPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temporary file"})
		return
	}
	defer os.Remove(tempPath)

	// 3. Verify hash of the uploaded file
	serverHash, err := pkg.CalculateSHA256(tempPath)
	if err != nil || serverHash != clientHash {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Hash mismatch or verification failed"})
		return
	}

	// 4. Finalize the file
	finalPath := filepath.Join(models.ModsDir, "enabled.json")
	if err := os.Rename(tempPath, finalPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize enabled.json"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "enabled.json uploaded and verified",
		"status":  "success",
	})
}

func isAlreadySynced(name, hash string) bool {
	// Ensure the directory exists
	if _, err := os.Stat(filepath.Dir(models.SyncFile)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(models.SyncFile), 0755)
	}

	var data models.SyncData
	content, err := os.ReadFile(models.SyncFile)
	if err != nil {
		return false
	}
	json.Unmarshal(content, &data)

	for _, mod := range data.Mods {
		if mod.Name == name && mod.Hash == hash {
			return true
		}
	}
	return false
}

func updateSyncMetadata(name, version, hash string) {
	// Ensure the directory exists
	if _, err := os.Stat(filepath.Dir(models.SyncFile)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(models.SyncFile), 0755)
	}

	var data models.SyncData
	content, err := os.ReadFile(models.SyncFile)
	if err == nil {
		json.Unmarshal(content, &data)
	}

	found := false
	for i, mod := range data.Mods {
		if mod.Name == name {
			data.Mods[i].Version = version
			data.Mods[i].Hash = hash
			found = true
			break
		}
	}

	if !found {
		data.Mods = append(data.Mods, models.ModMetadata{Name: name, Version: version, Hash: hash})
	}

	updatedContent, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(models.SyncFile, updatedContent, 0644)
}
