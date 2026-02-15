package handler

import (
	"encoding/json"
	"net/http"

	"github.com/JCHHeilmann/blocky-visor/sidecar/blocky"
)

func ServiceStatus(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := blocky.Status(serviceName)
		if err != nil {
			http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

func ServiceRestart(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := blocky.Restart(serviceName); err != nil {
			http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "restarted"})
	}
}
