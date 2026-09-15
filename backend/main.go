package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"docker-monitoring/db"
	"docker-monitoring/dockerclient"
	"docker-monitoring/handlers"
	"docker-monitoring/metrics"
)

// staticFiles embeds static frontend build if present
//
//go:embed dist/*
var staticFiles embed.FS

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting docksight Backend...")

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./monitoring.db"
	}

	// 1. Init DB
	if err := db.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 2. Init Docker Client
	_, err := dockerclient.InitDockerClient()
	if err != nil {
		log.Printf("Warning: Docker client init error: %v", err)
	}

	// 3. Start Collector
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collector := metrics.InitCollector()
	go collector.Start(ctx)

	// 4. Mux & Handlers
	mux := http.NewServeMux()

	// REST APIs
	mux.HandleFunc("/api/system/stats", handlers.GetSystemStatsHandler)
	mux.HandleFunc("/api/system/history", handlers.GetSystemHistoryHandler)
	mux.HandleFunc("/api/containers", handlers.GetContainersHandler)
	mux.HandleFunc("/api/containers/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/inspect") {
			handlers.ContainerInspectHandler(w, r)
		} else {
			handlers.ContainerActionHandler(w, r)
		}
	})
	mux.HandleFunc("/api/images", handlers.GetImagesHandler)
	mux.HandleFunc("/api/volumes", handlers.GetVolumesHandler)

	// SSE Streams
	mux.HandleFunc("/api/sse/stats", handlers.StatsSSEHandler)
	mux.HandleFunc("/api/sse/logs", handlers.LogsSSEHandler)

	// Frontend Static File Server with SPA Fallback
	distFS, err := fs.Sub(staticFiles, "dist")
	var fileServer http.Handler
	if err == nil {
		fileServer = http.FileServer(http.FS(distFS))
	} else {
		// Fallback to local ./dist folder if available
		fileServer = http.FileServer(http.Dir("./dist"))
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If path starts with /api/, return 404 if not matched
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Check if file exists in static FS
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path != "" {
			if err == nil {
				if f, openErr := distFS.Open(path); openErr == nil {
					f.Close()
					fileServer.ServeHTTP(w, r)
					return
				}
			} else {
				if _, statErr := os.Stat(filepath.Join("./dist", path)); statErr == nil {
					fileServer.ServeHTTP(w, r)
					return
				}
			}
		}

		// Fallback to index.html for SPA client-side routing
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: enableCORS(mux),
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		server.Shutdown(shutdownCtx)
		cancel()
	}()

	log.Printf("Server listening on http://0.0.0.0:%s\n", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}
