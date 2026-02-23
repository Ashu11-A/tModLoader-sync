# tModLoader-sync Project Context

tModLoader-sync is a Go-based project designed to synchronize mods and configuration for tModLoader (Terraria) between clients and servers. It includes a server component for Pterodactyl-based deployments and a cross-platform client (Windows/Linux).

## Project Overview

- **Purpose:** Synchronize `.tmod` files and `enabled.json` between a central server and clients.
- **Main Technologies:** Go (Gin for API, Bubbletea for TUI), Bash for server-side automation.
- **Architecture:** Client-Server architecture with shared logic in `shared/`.
- **Target Platforms:** Linux (Server/Client), Windows (Client).

## Architecture & Components

### 1. Server (`server/`)
- **Web Framework:** [Gin](https://github.com/gin-gonic/gin).
- **Functionality:**
    - Provides APIs for mod synchronization (`/v1/sync`, `/v1/upload`) and version checks (`/version`).
    - Tracks synced mods in `sync.json`.
    - Stores `.tmod` files in the `Mods/` directory relative to the executable.
- **Configuration:** Managed via `server/configs/config.go` (uses `pflag`). Requires a `--port` flag.
- **Deployment:** Specifically designed for **Pterodactyl Docker containers** with read-only filesystems and restricted access to `/home/container`.

### 2. Client (`client/`)
- **UI Framework:** [Bubbletea](https://github.com/charmbracelet/bubbletea) / [Lipgloss](https://github.com/charmbracelet/lipgloss).
- **Functionality:**
    - Scans local Steam Workshop directories for tModLoader mods.
    - Connects to the server to check synchronization status.
    - Uploads missing or updated mods (using SHA256 hashes for verification).
- **Configuration:** Managed via `client/configs/config.go` (uses `pflag`). Requires `--host` and `--port`.
- **Paths:**
    - **Linux:** `~/.steam/debian-installation/steamapps/workshop/content/1281930`
    - **Windows:** `C:\Program Files (x86)\Steam\steamapps\workshop\content\1281930`

### 3. Shared Logic (`shared/`)
- Contains cross-platform logic (e.g., OS detection in `pkg/os.go`, Hashing in `pkg/hash.go`).

### 4. Automation (`start.sh`)
- A bash script for the Linux server that:
    - Detects the tModLoader version.
    - Manages symbolic links for Steam Workshop mods.
    - Automatically generates `enabled.json`.
    - Launches the tModLoader server via `dotnet`.

## Building and Running

### Build Script
Use `build.sh` to build for all target platforms:
```bash
chmod +x build.sh
./build.sh
```
Outputs are placed in the `build/` directory.

### Manual Build
To build individual components:
```bash
# Server
cd server && go build -o server cmd/main.go
# Client
cd client && go build -o client cmd/main.go
```

### Running the Server
```bash
./server --port 8000
```
*Note: The `--port` flag is mandatory.*

### Running the Client
```bash
./client --host <server-ip> --port 8000
```
*Note: The `--host` and `--port` flags are mandatory.*

## Development Conventions

- **Go Workspace:** The project uses a Go workspace (managed via `go.work`, which is ignored by git).
- **API Versioning:**
    - `/version`: Root-level endpoint (no prefix) for version compatibility checks.
    - `/v1/*`: All other functional endpoints (e.g., `/v1/sync`, `/v1/upload`).
- **Configuration:** Use `configs.Load()` in `cmd/main.go` for centralized flag parsing.
- **Hashing:** Use `pkg.CalculateSHA256` from `shared/pkg/hash.go` for file verification.
- **I18n:** The client uses an internal i18n package for localized strings.
- **Formatting:** Use standard `go fmt` and idiomatic Go practices.

## Directory Structure

- `build/`: Contains compiled binaries and local test data (e.g., `sync.json`, `Mods/`).
- `client/cmd/`: Client entry point.
- `client/configs/config.go`: Client configuration and flag handling.
- `client/internal/api/`: REST API client for interacting with the server.
- `client/internal/ui/`: Bubbletea-based terminal UI components.
- `server/cmd/`: Server entry point.
- `server/configs/config.go`: Server configuration and flag handling.
- `server/internal/handlers/`: API route handlers (sync, upload, version).
- `shared/pkg/os.go`: OS detection logic.
- `shared/pkg/hash.go`: Shared SHA256 hashing utility.
- `shared/pkg/`: Common utilities shared between client and server.

## Sync Logic Workflow

1. **Client** scans the local Workshop directory for `.tmod` files.
2. **Client** calls `GetVersion` (`/version`) to verify compatibility.
3. **Client** requests `GetSyncStatus` (`/v1/sync`) from the server to get a list of currently synced mods.
4. **Client** compares local mods with the server list.
5. If a mod is missing or has a different hash, **Client** calculates the SHA256 hash (using `pkg.CalculateSHA256`) and calls `UploadMod` (`/v1/upload`).
6. **Server** receives the file, verifies the hash, and saves it to the `Mods/` directory.
7. **Server** updates `sync.json` with the new metadata.
8. **Server** (via `start.sh`) regenerates `enabled.json` to include the new/updated mods.
