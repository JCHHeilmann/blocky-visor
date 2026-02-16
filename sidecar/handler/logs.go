package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/JCHHeilmann/blocky-visor/sidecar/logparser"
	"github.com/JCHHeilmann/blocky-visor/sidecar/resolver"
)

func GetLogs(logDir string, hr *resolver.HostResolver) http.HandlerFunc {
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

		entries, _, err := logparser.LoadEntriesForRange(logDir, start, end)
		if err != nil {
			http.Error(w, jsonErr(err.Error()), http.StatusInternalServerError)
			return
		}

		// Enrich with resolved hostnames before filtering
		enrichEntries(entries, hr)

		// Filter
		filtered := logparser.FilterEntries(entries, filter)

		// Sort reverse chronological
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Timestamp.After(filtered[j].Timestamp)
		})

		// Paginate
		total := len(filtered)
		if offset >= total {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(logparser.LogsResponse{Total: total, Offset: offset, Limit: limit, Entries: []*logparser.LogEntry{}})
			return
		}
		endIdx := offset + limit
		if endIdx > total {
			endIdx = total
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logparser.LogsResponse{
			Total:   total,
			Offset:  offset,
			Limit:   limit,
			Entries: filtered[offset:endIdx],
		})
	}
}

func enrichEntries(entries []*logparser.LogEntry, hr *resolver.HostResolver) {
	for _, e := range entries {
		if name := hr.Lookup(e.ClientIP); name != "" {
			e.ResolvedName = name
		}
	}
}
