// Package models defines the data structures used by the tModLoader-sync server.
package models

// ModMetadata represents the basic information of a tModLoader mod.
type ModMetadata struct {
	Name    string `json:"name"`    // Internal name of the mod
	Version string `json:"version"` // Version string of the mod
	Hash    string `json:"hash"`    // SHA256 hash of the .tmod file
}

// SyncData represents the current synchronization state of the server.
type SyncData struct {
	Mods            []ModMetadata `json:"mods"`              // List of synced mods
	EnabledJSONHash string        `json:"enabled_json_hash"` // SHA256 hash of the enabled.json file
}

const (
	// SyncFile is the name of the file storing server's sync metadata.
	SyncFile = "Mods/sync.json"
	// ModsDir is the directory where the mod files are stored.
	ModsDir = "Mods"
)
