package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/JCHHeilmann/blocky-visor/sidecar/logparser"
	"github.com/JCHHeilmann/blocky-visor/sidecar/resolver"
)

// fileTracker holds the offset for a tailed log file.
type fileTracker struct {
	path   string
	offset int64
}

// todayLogFiles returns the log file paths for today's date using LogFilesForDate.
func todayLogFiles(logDir string) []string {
	return logparser.LogFilesForDate(logDir, time.Now())
}

func StreamLogs(logDir string, hr *resolver.HostResolver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, `{"error":"streaming not supported"}`, http.StatusInternalServerError)
			return
		}

		filter := logparser.LogFilter{
			Client: r.URL.Query().Get("client"),
			Domain: r.URL.Query().Get("domain"),
			Type:   r.URL.Query().Get("type"),
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")
		flusher.Flush()

		ctx := r.Context()

		currentDate := time.Now().Format("2006-01-02")
		logPaths := todayLogFiles(logDir)

		// Send recent historical entries on connect
		const backfillCount = 50
		var allFiltered []*logparser.LogEntry

		// Build initial trackers (one per file) and read existing entries
		trackers := make([]fileTracker, 0, len(logPaths))
		for _, logPath := range logPaths {
			f, err := os.Open(logPath)
			if err != nil {
				continue
			}
			scanner := bufio.NewScanner(f)
			scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}
				entry, err := logparser.ParseLine(line)
				if err != nil {
					continue
				}
				if name := hr.Lookup(entry.ClientIP); name != "" {
					entry.ResolvedName = name
				}
				if logparser.MatchesFilter(entry, filter) {
					allFiltered = append(allFiltered, entry)
				}
			}
			info, _ := f.Stat()
			var offset int64
			if info != nil {
				offset = info.Size()
			}
			f.Close()
			trackers = append(trackers, fileTracker{path: logPath, offset: offset})
		}

		// Sort by timestamp and send last N as backfill
		sort.Slice(allFiltered, func(i, j int) bool {
			return allFiltered[i].Timestamp.Before(allFiltered[j].Timestamp)
		})
		start := 0
		if len(allFiltered) > backfillCount {
			start = len(allFiltered) - backfillCount
		}
		if len(allFiltered[start:]) > 0 {
			data, _ := json.Marshal(allFiltered[start:])
			fmt.Fprintf(w, "event: backfill\ndata: %s\n\n", data)
			flusher.Flush()
		}

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Check for day rollover
				newDate := time.Now().Format("2006-01-02")
				if newDate != currentDate {
					currentDate = newDate
					logPaths = todayLogFiles(logDir)
					trackers = make([]fileTracker, 0, len(logPaths))
					for _, p := range logPaths {
						trackers = append(trackers, fileTracker{path: p, offset: 0})
					}
				}

				// Check for new per-client files that appeared since last tick
				currentPaths := todayLogFiles(logDir)
				knownPaths := make(map[string]bool, len(trackers))
				for _, t := range trackers {
					knownPaths[t.path] = true
				}
				for _, p := range currentPaths {
					if !knownPaths[p] {
						trackers = append(trackers, fileTracker{path: p, offset: 0})
					}
				}

				// Poll each tracked file for new data
				for i := range trackers {
					entries := pollFile(&trackers[i], hr, filter)
					for _, entry := range entries {
						data, err := json.Marshal(entry)
						if err != nil {
							continue
						}
						fmt.Fprintf(w, "data: %s\n\n", data)
					}
				}
				flusher.Flush()
			}
		}
	}
}

// pollFile reads new data from a single tracked file and returns matching entries.
func pollFile(t *fileTracker, hr *resolver.HostResolver, filter logparser.LogFilter) []*logparser.LogEntry {
	f, err := os.Open(t.path)
	if err != nil {
		return nil
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil
	}

	// File was truncated or rotated
	if info.Size() < t.offset {
		t.offset = 0
	}

	if info.Size() <= t.offset {
		return nil
	}

	if _, err := f.Seek(t.offset, io.SeekStart); err != nil {
		return nil
	}

	buf := make([]byte, info.Size()-t.offset)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return nil
	}
	t.offset += int64(n)

	var entries []*logparser.LogEntry
	lines := splitLines(buf[:n])
	for _, line := range lines {
		if line == "" {
			continue
		}
		entry, err := logparser.ParseLine(line)
		if err != nil {
			continue
		}
		if name := hr.Lookup(entry.ClientIP); name != "" {
			entry.ResolvedName = name
		}
		if !logparser.MatchesFilter(entry, filter) {
			continue
		}
		entries = append(entries, entry)
	}
	return entries
}

func splitLines(data []byte) []string {
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
