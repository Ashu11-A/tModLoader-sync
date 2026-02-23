package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"tml-sync/shared/pkg"
)

// TriggerServerUpdate downloads the new version of the server and restarts it.
func TriggerServerUpdate(targetVersion string) error {
	repo := "Ashu11-A/tModLoader-sync"
	osName := pkg.GetOSName()
	arch := pkg.GetArch()
	
	binaryName := fmt.Sprintf("server-%s-%s", osName, arch)
	if osName == "windows" {
		binaryName += ".exe"
	}

	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/v%s/%s", repo, targetVersion, binaryName)

	fmt.Printf("Updating server to v%s...", targetVersion)
	fmt.Printf("Downloading: %s", downloadURL)

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Create a temporary file for the new binary
	tempPath := exePath + ".new"
	out, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save update: %w", err)
	}
	out.Close()

	// Set executable permissions
	err = os.Chmod(tempPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}

	// Swap binaries
	oldPath := exePath + ".old"
	_ = os.Remove(oldPath) // Remove if exists
	
	err = os.Rename(exePath, oldPath)
	if err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	err = os.Rename(tempPath, exePath)
	if err != nil {
		// Try to restore backup
		_ = os.Rename(oldPath, exePath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	fmt.Println("Server updated! Restarting...")

	// Restart the process
	if pkg.System == pkg.Windows {
		// On Windows, we start a new process via cmd to allow a small delay
		// so the current process can exit and release its port.
		args := ""
		if len(os.Args) > 1 {
			args = " " + strings.Join(os.Args[1:], " ")
		}

		restartCmd := fmt.Sprintf("ping 127.0.0.1 -n 3 > nul && \"%s\"%s", exePath, args)
		cmd := exec.Command("cmd", "/C", restartCmd)
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to launch restarter: %w", err)
		}
		os.Exit(0)
	} else {
		err = syscall.Exec(exePath, os.Args, os.Environ())
		if err != nil {
			return fmt.Errorf("failed to restart server: %w", err)
		}
	}

	return nil
}
