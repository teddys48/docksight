package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"docker-monitoring/dockerclient"
)

type ActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func GetContainersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if dockerclient.Instance == nil {
		http.Error(w, `{"error":"Docker client not initialized"}`, http.StatusInternalServerError)
		return
	}

	containers, err := dockerclient.Instance.ListContainers(r.Context())
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(containers)
}

func ContainerActionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// URL format: /api/containers/{id}/{action}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.Error(w, `{"error":"Invalid URL path"}`, http.StatusBadRequest)
		return
	}

	containerID := parts[2]
	action := parts[3]

	if dockerclient.Instance == nil {
		http.Error(w, `{"error":"Docker client not initialized"}`, http.StatusInternalServerError)
		return
	}

	var err error
	switch action {
	case "start":
		err = dockerclient.Instance.StartContainer(r.Context(), containerID)
	case "stop":
		err = dockerclient.Instance.StopContainer(r.Context(), containerID)
	case "restart":
		err = dockerclient.Instance.RestartContainer(r.Context(), containerID)
	case "remove":
		force := r.URL.Query().Get("force") == "true"
		err = dockerclient.Instance.RemoveContainer(r.Context(), containerID, force)
	default:
		http.Error(w, `{"error":"Unknown action"}`, http.StatusBadRequest)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ActionResponse{Success: false, Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(ActionResponse{Success: true, Message: "Action " + action + " executed successfully"})
}

func ContainerInspectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, `{"error":"Container ID required"}`, http.StatusBadRequest)
		return
	}
	containerID := parts[2]

	if dockerclient.Instance == nil {
		http.Error(w, `{"error":"Docker client not initialized"}`, http.StatusInternalServerError)
		return
	}

	inspectData, err := dockerclient.Instance.InspectContainer(r.Context(), containerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(inspectData)
}
