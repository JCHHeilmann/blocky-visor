package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/JCHHeilmann/blocky-visor/sidecar/blocky"
)

func GetConfig(configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := blocky.ReadConfig(configPath)
		if err != nil {
			http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(data)
	}
}

func PutConfig(configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB max
		if err != nil {
			http.Error(w, jsonErr("failed to read body"), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(data) == 0 {
			http.Error(w, jsonErr("empty config body"), http.StatusBadRequest)
			return
		}

		backupPath, err := blocky.WriteConfig(configPath, data)
		if err != nil {
			http.Error(w, jsonErr(err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "saved",
			"backup": backupPath,
		})
	}
}

func jsonErr(msg string) string {
	b, _ := json.Marshal(map[string]string{"error": msg})
	return string(b)
}
