package models

type ModMetadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

type SyncData struct {
	Mods            []ModMetadata `json:"mods"`
	EnabledJSONHash string        `json:"enabled_json_hash"`
}

type ServerVersionResponse struct {
	ServerVersion string `json:"server_version"`
	TMLVersion    string `json:"tml_version"`
}
