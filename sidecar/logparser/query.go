package logparser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LogFilter defines filter criteria for log queries.
type LogFilter struct {
	Client string // filter by client IP substring
	Domain string // filter by domain substring
	Type   string // "blocked", "cached", "resolved", or "" for all
}

// LogsResponse is the paginated logs response.
type LogsResponse struct {
	Total   int         `json:"total"`
	Offset  int         `json:"offset"`
	Limit   int         `json:"limit"`
	Entries []*LogEntry `json:"entries"`
}

// ParseFile reads a Blocky log file and returns all parsed entries.
func ParseFile(path string) ([]*LogEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}
	defer f.Close()

	var entries []*LogEntry
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		entry, err := ParseLine(line)
		if err != nil {
			continue // skip unparseable lines
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

// LogFilesForRange returns the log file paths for dates within the given range.
func LogFilesForRange(logDir string, start, end time.Time) []string {
	var files []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		name := fmt.Sprintf("%s_ALL.log", d.Format("2006-01-02"))
		path := filepath.Join(logDir, name)
		if _, err := os.Stat(path); err == nil {
			files = append(files, path)
		}
	}
	return files
}

// LoadEntriesForRange parses all log files in the date range and returns entries + file count.
func LoadEntriesForRange(logDir string, start, end time.Time) ([]*LogEntry, int, error) {
	files := LogFilesForRange(logDir, start, end)
	var allEntries []*LogEntry
	for _, f := range files {
		entries, err := ParseFile(f)
		if err != nil {
			continue // skip unreadable files
		}
		allEntries = append(allEntries, entries...)
	}
	return allEntries, len(files), nil
}

// FilterEntries applies a LogFilter to a set of entries.
func FilterEntries(entries []*LogEntry, filter LogFilter) []*LogEntry {
	if filter.Client == "" && filter.Domain == "" && filter.Type == "" {
		return entries
	}

	var result []*LogEntry
	for _, e := range entries {
		if filter.Client != "" && !strings.Contains(e.ClientIP, filter.Client) && !strings.Contains(e.ClientName, filter.Client) {
			continue
		}
		if filter.Domain != "" && !strings.Contains(strings.ToLower(e.Domain), strings.ToLower(filter.Domain)) {
			continue
		}
		if filter.Type != "" {
			switch filter.Type {
			case "blocked":
				if !e.IsBlocked() {
					continue
				}
			case "cached":
				if !e.IsCached() {
					continue
				}
			case "resolved":
				if e.IsBlocked() || e.IsCached() {
					continue
				}
			}
		}
		result = append(result, e)
	}
	return result
}

// QueryLogs loads, filters, and paginates log entries in reverse chronological order.
func QueryLogs(logDir string, start, end time.Time, filter LogFilter, limit, offset int) (*LogsResponse, error) {
	entries, _, err := LoadEntriesForRange(logDir, start, end)
	if err != nil {
		return nil, err
	}

	// Apply filters
	filtered := FilterEntries(entries, filter)

	// Sort reverse chronological
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.After(filtered[j].Timestamp)
	})

	total := len(filtered)

	// Paginate
	if offset >= total {
		return &LogsResponse{Total: total, Offset: offset, Limit: limit, Entries: []*LogEntry{}}, nil
	}
	end2 := offset + limit
	if end2 > total {
		end2 = total
	}

	return &LogsResponse{
		Total:   total,
		Offset:  offset,
		Limit:   limit,
		Entries: filtered[offset:end2],
	}, nil
}
