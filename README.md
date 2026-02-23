<div align="center">

# tModLoader sync

![License](https://img.shields.io/github/license/Ashu11-A/tModLoader-sync?style=for-the-badge&color=302D41&labelColor=f9e2af&logoColor=302D41)
![Stars](https://img.shields.io/github/stars/Ashu11-A/tModLoader-sync?style=for-the-badge&color=302D41&labelColor=f9e2af&logoColor=302D41)
![Last Commit](https://img.shields.io/github/last-commit/Ashu11-A/tModLoader-sync?style=for-the-badge&color=302D41&labelColor=b4befe&logoColor=302D41)
![Repo Size](https://img.shields.io/github/repo-size/Ashu11-A/tModLoader-sync?style=for-the-badge&color=302D41&labelColor=90dceb&logoColor=302D41)

<br>

<p align="center">
  <strong>Automatically synchronizes your Steam mods with the server, making upload and maintenance easier.</strong>
  <br>
  <sub>
    This tool was created to be used with the tModLoader egg found here: <a href="https://github.com/Ashu11-A/Ashu_eggs">Ashu_eggs</a>
  </sub>
</p>

<p align="center">
  <a href="https://github.com/Ashu11-A/tModLoader-sync/stargazers">
    <img src="https://img.shields.io/badge/Leave%20a%20Star%20🌟-302D41?style=for-the-badge&color=302D41&labelColor=302D41" alt="Star Repo">
  </a>
</p>

</div>

---

## 🚀 Overview

**tModLoader-sync** is a synchronization tool designed to bridge the gap between tModLoader servers and mod uploading. It ensures that all mods from your Steam workshop are synchronized with the server, using SHA256 hashing for precise verification.

### Key Features

* **Automatic Synchronization:** Detects missing or outdated mods and uploads them to the server.
* **Multi-platform:** Support for Linux and Windows (x64 and ARM64).
* **Interactive TUI:** Friendly terminal interface built with [Bubbletea](https://github.com/charmbracelet/bubbletea).
* **Pterodactyl Ready:** Specifically optimized for Pterodactyl-based deployments.
* **Multilingual:** Support for English and Portuguese (PT-BR).

---

## 🏗️ Architecture

The project is divided into three main components:

1. **Server (`/server`):** A Gin-based REST API that manages mod storage and synchronization status.
2. **Client (`/client`):** A TUI application that scans local Steam Workshop mods and interacts with the server.
3. **Shared (`/shared`):** Common logic for operating system detection, hashing, and versioning.

---

## 🛠️ Building

The project uses a Go workspace. You can build all targets using the provided script:

```bash
./build.sh
```

Binaries will be available in the `build/` directory for:

* Linux (x64, ARM64)
* Windows (x64, ARM64)

---

## 📖 How to Use

### Server

Run the server with the mandatory `--port` flag:

```bash
./server --port 8080
```

### Client

Run the client by providing the server's host and port:

```bash
./client --host <SERVER_IP> --port 8080
```

The client will automatically detect your Steam Workshop path and start the synchronization process.

---

## 📥 Installation (Quick Script)

### Linux (Bash)

```bash
curl -fsSL https://raw.githubusercontent.com/Ashu11-A/tModLoader-sync/main/sync.sh | h=<SERVER_IP> p=:<PORT> bash
```

### Windows (PowerShell)

```powershell
powershell -c "$h='<SERVER_IP>';$p=':<PORT>';irm https://raw.githubusercontent.com/Ashu11-A/tModLoader-sync/main/sync.ps1|iex"
```
