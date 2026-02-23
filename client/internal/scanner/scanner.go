package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"tml-sync/client/internal/models"
	"tml-sync/shared/pkg"

	"golang.org/x/mod/semver"
)

// GetSteamWorkshopPath attempts to find the tModLoader workshop directory on different OSs.
func GetSteamWorkshopPath() string {
	home, _ := os.UserHomeDir()
	var paths []string

	if pkg.System == pkg.Windows {
		paths = []string{
			`C:\Program Files (x86)\Steam\steamapps\workshop\content\1281930`,
			filepath.Join(home, "AppData/Local/Steam/steamapps/workshop/content/1281930"), // Fallback/Custom
		}
	} else {
		// Common Linux Steam paths
		paths = []string{
			filepath.Join(home, ".steam/debian-installation/steamapps/workshop/content/1281930"),
			filepath.Join(home, ".local/share/Steam/steamapps/workshop/content/1281930"),
			filepath.Join(home, ".steam/steam/steamapps/workshop/content/1281930"),
			filepath.Join(home, ".var/app/com.valvesoftware.Steam/.steam/steam/steamapps/workshop/content/1281930"), // Flatpak
		}
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// GetEnabledJSONPath finds the local enabled.json for tModLoader.
func GetEnabledJSONPath() string {
	home, _ := os.UserHomeDir()
	var paths []string

	if pkg.System == pkg.Windows {
		paths = []string{
			filepath.Join(home, "Documents", "My Games", "Terraria", "tModLoader", "Mods", "enabled.json"),
			filepath.Join(home, "OneDrive", "Documents", "My Games", "Terraria", "tModLoader", "Mods", "enabled.json"), // OneDrive fallback
		}
	} else {
		paths = []string{
			filepath.Join(home, ".local", "share", "Terraria", "tModLoader", "Mods", "enabled.json"),
		}
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

type FoundMod struct {
	Path     string
	Metadata models.ModMetadata
}

// ScanMods finds the best .tmod for each mod based on the server's tML version.
func ScanMods(workshopPath string, serverTMLVersion string) ([]FoundMod, error) {
	var found []FoundMod

	modDirs, err := os.ReadDir(workshopPath)
	if err != nil {
		return nil, err
	}

	for _, modDir := range modDirs {
		if !modDir.IsDir() {
			continue
		}

		modIDPath := filepath.Join(workshopPath, modDir.Name())
		versionDirs, err := os.ReadDir(modIDPath)
		if err != nil {
			continue
		}

		var versions []string
		for _, vDir := range versionDirs {
			if vDir.IsDir() {
				versions = append(versions, vDir.Name())
			}
		}

		if len(versions) == 0 {
			continue
		}

		// Sort versions descending to find the highest compatible one first
		sort.Slice(versions, func(i, j int) bool {
			return compareVersions(versions[i], versions[j]) > 0
		})

		selectedVersion := ""
		
		// Logic similar to start.sh: find first version <= server version
		if serverTMLVersion != "" && serverTMLVersion != "unknown" {
			for _, v := range versions {
				if compareVersions(v, serverTMLVersion) <= 0 {
					selectedVersion = v
					break
				}
			}
		}

		// Fallback to most recent if no compatible found
		if selectedVersion == "" {
			selectedVersion = versions[0]
		}

		versionPath := filepath.Join(modIDPath, selectedVersion)
		files, err := os.ReadDir(versionPath)
		if err != nil {
			continue
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".tmod") {
				path := filepath.Join(versionPath, file.Name())
				modName := strings.TrimSuffix(file.Name(), ".tmod")
				
				hash, err := pkg.CalculateSHA256(path)
				if err != nil {
					continue
				}

				found = append(found, FoundMod{
					Path: path,
					Metadata: models.ModMetadata{
						Name:    modName,
						Version: selectedVersion,
						Hash:    hash,
					},
				})
				break // Only one .tmod per mod ID
			}
		}
	}

	return found, nil
}

// compareVersions mimics basic version comparison logic. 
// It handles tML formats like 2024.12.
func compareVersions(v1, v2 string) int {
	// Simple semver wrapper if prepended with 'v'
	sv1 := v1
	if !strings.HasPrefix(sv1, "v") {
		sv1 = "v" + sv1
	}
	sv2 := v2
	if !strings.HasPrefix(sv2, "v") {
		sv2 = "v" + sv2
	}

	if semver.IsValid(sv1) && semver.IsValid(sv2) {
		return semver.Compare(sv1, sv2)
	}

	// Fallback to basic string comparison if semver fails
	if v1 < v2 {
		return -1
	}
	if v1 > v2 {
		return 1
	}
	return 0
}
