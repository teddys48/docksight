package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"docker-monitoring/db"
	"docker-monitoring/metrics"
)

func GetSystemStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if metrics.GlobalCollector == nil {
		http.Error(w, `{"error":"Collector not initialized"}`, http.StatusInternalServerError)
		return
	}

	stats := metrics.GlobalCollector.GetCurrentStats()
	json.NewEncoder(w).Encode(stats)
}

func GetSystemHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limitStr := r.URL.Query().Get("limit")
	limit := 60
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	snapshots, err := db.GetMetricsHistory(limit)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(snapshots)
}
