package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JCHHeilmann/blocky-visor/sidecar/logparser"
)

func GetStats(logDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start, end := parseRange(r)

		acc := logparser.NewStatsAccumulator(start, end)
		filesParsed, err := logparser.ProcessFilesStreaming(logDir, start, end, acc.Add)
		if err != nil {
			http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
			return
		}

		stats := acc.Finalize(filesParsed)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}

func GetTimeline(logDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start, end := parseRange(r)

		intervalStr := r.URL.Query().Get("interval")
		interval := 15 * time.Minute
		switch intervalStr {
		case "5m":
			interval = 5 * time.Minute
		case "15m":
			interval = 15 * time.Minute
		case "1h":
			interval = time.Hour
		case "1d":
			interval = 24 * time.Hour
		}

		acc := logparser.NewTimelineAccumulator(interval)
		_, err := logparser.ProcessFilesStreaming(logDir, start, end, acc.Add)
		if err != nil {
			http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
			return
		}

		timeline := acc.Finalize()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(timeline)
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
