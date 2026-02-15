package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/JCHHeilmann/blocky-visor/sidecar/logparser"
)

func GetLogs(logDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start, end := parseRange(r)

		limit := 100
		if v := r.URL.Query().Get("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 1000 {
				limit = n
			}
		}

		offset := 0
		if v := r.URL.Query().Get("offset"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				offset = n
			}
		}

		filter := logparser.LogFilter{
			Client: r.URL.Query().Get("client"),
			Domain: r.URL.Query().Get("domain"),
			Type:   r.URL.Query().Get("type"),
		}

		// Default range for logs: today + yesterday for more history
		if r.URL.Query().Get("range") == "" {
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			start = today.AddDate(0, 0, -1)
			end = now
		}

		result, err := logparser.QueryLogs(logDir, start, end, filter, limit, offset)
		if err != nil {
			http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
