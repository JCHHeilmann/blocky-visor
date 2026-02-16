package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/JCHHeilmann/blocky-visor/sidecar/logparser"
)

func StreamLogs(logDir string) http.HandlerFunc {
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
		logPath := filepath.Join(logDir, currentDate+"_ALL.log")

		// Send recent historical entries on connect
		const backfillCount = 50
		var fileOffset int64
		if f, err := os.Open(logPath); err == nil {
			var allFiltered []*logparser.LogEntry
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
				if matchesFilter(entry, filter) {
					allFiltered = append(allFiltered, entry)
				}
			}
			// Take last N entries, send as a single "backfill" event
			start := 0
			if len(allFiltered) > backfillCount {
				start = len(allFiltered) - backfillCount
			}
			if len(allFiltered[start:]) > 0 {
				data, _ := json.Marshal(allFiltered[start:])
				fmt.Fprintf(w, "event: backfill\ndata: %s\n\n", data)
				flusher.Flush()
			}

			info, _ := f.Stat()
			if info != nil {
				fileOffset = info.Size()
			}
			f.Close()
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
					logPath = filepath.Join(logDir, currentDate+"_ALL.log")
					fileOffset = 0
				}

				f, err := os.Open(logPath)
				if err != nil {
					continue
				}

				info, err := f.Stat()
				if err != nil {
					f.Close()
					continue
				}

				// File was truncated or rotated
				if info.Size() < fileOffset {
					fileOffset = 0
				}

				if info.Size() <= fileOffset {
					f.Close()
					continue
				}

				if _, err := f.Seek(fileOffset, io.SeekStart); err != nil {
					f.Close()
					continue
				}

				buf := make([]byte, info.Size()-fileOffset)
				n, err := f.Read(buf)
				f.Close()
				if err != nil && err != io.EOF {
					continue
				}
				fileOffset += int64(n)

				// Parse lines from the chunk
				lines := splitLines(buf[:n])
				for _, line := range lines {
					if line == "" {
						continue
					}
					entry, err := logparser.ParseLine(line)
					if err != nil {
						continue
					}
					if !matchesFilter(entry, filter) {
						continue
					}
					data, err := json.Marshal(entry)
					if err != nil {
						continue
					}
					fmt.Fprintf(w, "data: %s\n\n", data)
				}
				flusher.Flush()
			}
		}
	}
}

func splitLines(data []byte) []string {
	var lines []string
	start := 0
	for i, b := range data {
		if b == '\n' {
			line := string(data[start:i])
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	return lines
}

func matchesFilter(entry *logparser.LogEntry, filter logparser.LogFilter) bool {
	filtered := logparser.FilterEntries([]*logparser.LogEntry{entry}, filter)
	return len(filtered) > 0
}
