package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JCHHeilmann/blocky-visor/sidecar/blocky"
)

func ServiceStatus(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := blocky.Status(serviceName)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, blocky.ErrSystemdNotAvailable) {
				status = http.StatusNotImplemented
			}
			http.Error(w, jsonErr(err.Error()), status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

func ServiceRestart(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := blocky.Restart(serviceName); err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, blocky.ErrSystemdNotAvailable) {
				status = http.StatusNotImplemented
			}
			http.Error(w, jsonErr(err.Error()), status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "restarted"})
	}
}
