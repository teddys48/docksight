# docksight 👁️ - Modern Real-time Docker & Podman Monitoring

**docksight** is a lightweight, high-performance, real-time container monitoring and management application built for **Docker** and **Podman** using **Golang**, **Svelte 5**, **Vite**, **Bun**, **SQLite**, and **TailwindCSS**.

![docksight Tech Stack](https://img.shields.io/badge/Stack-Golang%20%7C%20Svelte%205%20%7C%20Bun%20%7C%20SQLite%20%7C%20Tailwind-blue)
![Container Engine API](https://img.shields.io/badge/Engine-Docker%20%7C%20Podman-sky)
![License](https://img.shields.io/badge/License-MIT-emerald)

---

## ✨ Features

### 📊 1. System & Host Metrics Dashboard

- **Host Resource Tracking**: Real-time monitoring of Host CPU usage %, Host RAM (Used/Total), and Host Disk (Used/Total).
- **Container Overview**: Total containers count, active running containers, and stopped/exited containers.
- **Real-time Metrics Charts**: Interactive timeline line charts for Host CPU and RAM usage powered by Chart.js.
- **Container I/O Summary**: Live Network I/O (RX/TX bytes) and Block I/O (Disk Read/Write) metrics.

### 🐳 2. Container Management & Pagination (Docker & Podman)

- **Container Directory**: Detailed list showing container name, status badges, CPU %, RAM usage, Network I/O, and Uptime with **Pagination controls** (5, 10, 20, 50 rows per page).
- **Search & Filters**: Instant live text search (by container name, ID, or image name) and tabbed status filters (All, Running, Stopped).
- **Lifecycle Control Operations**:
  - ▶️ **Start** container
  - ⏹️ **Stop** container
  - 🔄 **Restart** container
  - 🗑️ **Remove** container (with optional Force flag)
  - 🔍 **Inspect JSON View**

### 💾 3. Volume Management & Pagination

- **Volume Directory**: View host volumes, driver, mount points, scope, and creation dates.
- **Paginated Volume Table**: Full pagination support (5, 10, 20, 50 rows per page) and instant search filtering.

### 📜 4. Real-time Container Log Viewer (SSE Event Stream)

- **Server-Sent Events Tailing**: Live streaming stdout and stderr log outputs without refreshing the page.
- **Interactive Log Controls**:
  - ⏸️ **Pause / Resume** live log stream
  - 🧹 **Clear** log view
  - ⬇️ **Auto-scroll** toggle to keep focus on latest log output
  - 🕒 **Timestamps** display toggle
  - ⚙️ **STDOUT / STDERR** channel filters
  - 🔍 **Live Log Search**: Filter log lines instantly with query matching
  - 🔢 **Configurable Tail Lines**: Select 50, 100, 200, or 500 initial tail lines

### 🖼️ 5. Container Images Explorer & Pagination

- View all local Docker/Podman images with Repository Tags, Image ID, File Size (MB/GB), and Created Date.
- **Paginated Image Directory**: Full pagination support (5, 10, 20, 50 rows per page) and instant search filtering.
- Automatic detection and filtering of **Dangling Images** (`<none>:<none>`).

### 🎨 6. Modern UI & Ergonomics

- **Dark & Light Mode Switcher**: Seamless theme toggle with persistent `localStorage` user preferences.
- **Skeleton Loader UI**: Smooth placeholder loading states during data fetch.
- **Toast Notifications**: Interactive alert toasts for container start, stop, restart, and deletion operations.
- **Fully Responsive**: Optimized for desktop, tablet, and mobile displays.

---

## 🏗️ Architecture & Technology Stack

```
 ┌─────────────────────────────────────────────────────────┐
 │                   Svelte 5 + Vite Frontend              │
 │  - Dashboard (Host Metrics, Interactive Line Charts)    │
 │  - Container Manager (Paginated, Actions & Live Search) │
 │  - Volume Explorer (Paginated, Search & Mount points)   │
 │  - SSE Log Viewer (Filter STDOUT/STDERR, Search, Pause) │
 │  - Image Explorer (Dangling badges & Tags)              │
 └────────────────────────────┬────────────────────────────┘
                              │ REST APIs / SSE Streams
 ┌────────────────────────────▼────────────────────────────┘
 │                     Go Backend API                      │
 │  - Docker/Podman Client (github.com/docker/docker)      │
 │  - Auto-detection for Docker & Podman Unix Sockets      │
 │  - System Metrics Collector (github.com/shirou/gopsutil) │
 │  - CGO-Free SQLite Driver (modernc.org/sqlite)          │
 │  - Go Embed (Single self-contained binary deployment)   │
 └────────────────────────────┬────────────────────────────┘
                              │ Docker / Podman Socket
 ┌────────────────────────────▼────────────────────────────┐
 │              Container Engine (Docker / Podman)         │
 │  (/var/run/docker.sock or /run/user/<UID>/podman.sock)  │
 └─────────────────────────────────────────────────────────┘
```

---

## 🦭 Podman Support & Setup Guide

**docksight** fully supports monitoring and managing containers, volumes, and images managed by **Podman** (both Rootless and Rootful modes) via Podman's Docker-compatible API socket.

### 1. Auto-Detection Mechanism

When launched, **docksight** automatically checks candidate socket endpoints in the following order:

1. `DOCKER_HOST` environment variable (if explicitly set)
2. `/var/run/docker.sock` (Standard Docker socket)
3. `unix://$XDG_RUNTIME_DIR/podman/podman.sock` (Podman rootless socket)
4. `unix:///run/user/<UID>/podman/podman.sock` (Podman rootless socket fallback)
5. `unix:///run/podman/podman.sock` (Podman rootful socket)
6. `unix:///var/run/podman/podman.sock` (Alternative Podman rootful socket)

### 2. Enabling Podman System Service Socket

Before running docksight with Podman, ensure the Podman API service socket is active:

#### Rootless Podman (Recommended):

```bash
# Enable and start the user podman socket service
systemctl --user enable --now podman.socket
```

Verify the socket is active:

```bash
ls -la $XDG_RUNTIME_DIR/podman/podman.sock
```

#### Rootful Podman:

```bash
# Enable and start the system podman socket service
sudo systemctl enable --now podman.socket
```

Verify the socket is active:

```bash
sudo ls -la /run/podman/podman.sock
```

### 3. Running docksight with Podman

#### Option A: Running Standalone Binary with Podman

If running the compiled binary directly on your system:

```bash
# docksight automatically auto-detects the Podman socket!
./docksight

# Or explicitly set DOCKER_HOST if using a non-standard socket:
export DOCKER_HOST="unix://$XDG_RUNTIME_DIR/podman/podman.sock"
./docksight
```

#### Option B: Running in Container using `podman run`

**Rootless Podman Container:**

```bash
podman run -d \
  --name docksight \
  -p 8080:8080 \
  -v $XDG_RUNTIME_DIR/podman/podman.sock:/var/run/docker.sock:ro \
  -v monitoring_data:/data \
  docksight
```

**Rootful Podman Container:**

```bash
sudo podman run -d \
  --name docksight \
  -p 8080:8080 \
  -v /run/podman/podman.sock:/var/run/docker.sock:ro \
  -v monitoring_data:/data \
  docksight
```

#### Option C: Running with `podman-compose`

If using `podman-compose`, you can specify the Podman socket mount in `docker-compose.yml`:

```yaml
services:
  docksight:
    build: .
    container_name: docksight
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      # Podman rootless socket mapping example:
      - ${XDG_RUNTIME_DIR}/podman/podman.sock:/var/run/docker.sock:ro
      - monitoring_data:/data
    environment:
      - PORT=8080
      - DB_PATH=/data/monitoring.db

volumes:
  monitoring_data:
    driver: local
```

Then start using:

```bash
podman-compose up -d
```

---

## 🚀 Quick Start Guide (Docker)

### Option 1: Run with Docker Compose (Recommended)

To start **docksight** instantly using Docker Compose:

```bash
docker-compose up -d --build
```

Access the dashboard in your web browser:
👉 **`http://localhost:8080`**

---

### Option 2: Local Development Mode

#### 1. Start Go Backend

```bash
cd backend
go run main.go
```

The backend server will start on `http://localhost:8080`.

#### 2. Start Frontend (Bun + Vite)

In a separate terminal window:

```bash
cd frontend
bun install
bun run dev
```

Open your browser at `http://localhost:5173`. Vite will proxy API requests to `http://localhost:8080`.

---

## 📦 Building a Single Executable Binary

You can compile **docksight** into a single self-contained executable binary with embedded static frontend assets:

```bash
# 1. Build frontend dist
cd frontend
bun run build

# 2. Copy static build assets to backend
cd ..
rm -rf backend/dist
cp -r frontend/dist backend/dist

# 3. Build standalone Go binary
cd backend
CGO_ENABLED=0 go build -o docksight main.go
```

Run `./docksight` anywhere without external frontend files or Node/Bun runtime dependencies!

---

## ⚙️ Environment Variables

| Variable      | Default           | Description                                                                                                   |
| :------------ | :---------------- | :------------------------------------------------------------------------------------------------------------ |
| `PORT`        | `8080`            | Port for the docksight HTTP web server                                                                        |
| `DB_PATH`     | `./monitoring.db` | Path to the SQLite metrics database file                                                                      |
| `DOCKER_HOST` | _(Auto-detected)_ | Custom socket endpoint URI (e.g. `unix:///run/user/1000/podman/podman.sock` or `unix:///var/run/docker.sock`) |

---

## 📄 License

MIT License. Developed with high-performance coding standards.
