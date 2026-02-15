package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JCHHeilmann/blocky-visor/sidecar/logparser"
)

func GetStats(logDir string, cache *logparser.StatsCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start, end := parseRange(r)
		stats := cache.ComputeStats(logDir, start, end)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}

func GetTimeline(logDir string, cache *logparser.StatsCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start, end := parseRange(r)
		interval := parseInterval(r)
		timeline := cache.ComputeTimeline(logDir, start, end, interval)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(timeline)
	}
}

func parseInterval(r *http.Request) time.Duration {
	switch r.URL.Query().Get("interval") {
	case "5m":
		return 5 * time.Minute
	case "1h":
		return time.Hour
	case "1d":
		return 24 * time.Hour
	default:
		return 15 * time.Minute
	}
}

func parseRange(r *http.Request) (time.Time, time.Time) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	rangeStr := r.URL.Query().Get("range")
	switch rangeStr {
	case "yesterday":
		return today.AddDate(0, 0, -1), today.Add(-time.Second)
	case "7d":
		return today.AddDate(0, 0, -6), now
	case "30d":
		return today.AddDate(0, 0, -29), now
	default: // "today" or empty
		return today, now
	}
}
