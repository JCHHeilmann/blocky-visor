package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JCHHeilmann/blocky-visor/sidecar/blocky"
)

func systemdErrorStatus(err error) int {
	if errors.Is(err, blocky.ErrSystemdNotAvailable) {
		return http.StatusNotImplemented
	}
	return http.StatusInternalServerError
}

func ServiceStatus(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := blocky.Status(serviceName)
		if err != nil {
			http.Error(w, jsonErr(err.Error()), systemdErrorStatus(err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

func ServiceRestart(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := blocky.Restart(serviceName); err != nil {
			http.Error(w, jsonErr(err.Error()), systemdErrorStatus(err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "restarted"})
	}
}
