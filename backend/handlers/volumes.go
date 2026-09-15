package handlers

import (
	"encoding/json"
	"net/http"

	"docker-monitoring/dockerclient"
)

func GetVolumesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if dockerclient.Instance == nil {
		http.Error(w, `{"error":"Docker client not initialized"}`, http.StatusInternalServerError)
		return
	}

	volumes, err := dockerclient.Instance.ListVolumes(r.Context())
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(volumes)
}
