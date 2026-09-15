package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"docker-monitoring/dockerclient"
	"docker-monitoring/metrics"
)

func StatsSSEHandler(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Initial send
	if metrics.GlobalCollector != nil {
		stats := metrics.GlobalCollector.GetCurrentStats()
		data, err := json.Marshal(stats)
		if err == nil {
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			if metrics.GlobalCollector != nil {
				stats := metrics.GlobalCollector.GetCurrentStats()
				data, err := json.Marshal(stats)
				if err == nil {
					fmt.Fprintf(w, "data: %s\n\n", data)
					flusher.Flush()
				}
			}
		}
	}
}

func LogsSSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", `{"error":"Container ID required"}`)
		flusher.Flush()
		return
	}

	tail := r.URL.Query().Get("tail")
	if tail == "" {
		tail = "100"
	}

	stdout := r.URL.Query().Get("stdout") != "false"
	stderr := r.URL.Query().Get("stderr") != "false"
	timestamps := r.URL.Query().Get("timestamps") == "true"

	if dockerclient.Instance == nil {
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", `{"error":"Docker client not initialized"}`)
		flusher.Flush()
		return
	}

	logChan := make(chan string, 100)
	ctx := r.Context()

	go func() {
		defer close(logChan)
		_ = dockerclient.Instance.StreamLogs(ctx, containerID, tail, stdout, stderr, timestamps, true, logChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-logChan:
			if !ok {
				fmt.Fprintf(w, "event: end\ndata: Log stream ended\n\n")
				flusher.Flush()
				return
			}

			// Format SSE JSON frame
			payload, _ := json.Marshal(map[string]string{"log": line})
			fmt.Fprintf(w, "data: %s\n\n", payload)
			flusher.Flush()
		}
	}
}
