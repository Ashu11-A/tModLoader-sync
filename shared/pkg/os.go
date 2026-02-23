package pkg

import (
	"fmt"
	"runtime"
	"os"
)

type OperatingSystem int

const (
	Windows OperatingSystem = iota
	Linux
)

func checkSystem() OperatingSystem {
	system := runtime.GOOS

	switch system {
	case "windows": return Windows
	case "linux": return Linux
	}

	fmt.Printf("Warning: Operating system ‘%s’ not supported. Shutting down the application...\n", system)
	os.Exit(0)

	return 0
}

var System OperatingSystem = checkSystem()

// GetArch returns the architecture in the format used by the release script (x64, arm64)
func GetArch() string {
	switch runtime.GOARCH {
	case "amd64":
		return "x64"
	case "arm64":
		return "arm64"
	default:
		return runtime.GOARCH
	}
}

// GetOSName returns the OS name in the format used by the release script (windows, linux)
func GetOSName() string {
	return runtime.GOOS
}