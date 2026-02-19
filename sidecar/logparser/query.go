package logparser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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

// MatchesFilter tests whether a single entry passes the filter criteria.
func MatchesFilter(e *LogEntry, filter LogFilter) bool {
	if filter.Client != "" {
		cl := strings.ToLower(filter.Client)
		if !strings.Contains(strings.ToLower(e.ClientIP), cl) &&
			!strings.Contains(strings.ToLower(e.ClientName), cl) &&
			!strings.Contains(strings.ToLower(e.ResolvedName), cl) {
			return false
		}
	}
	if filter.Domain != "" && !strings.Contains(strings.ToLower(e.Domain), strings.ToLower(filter.Domain)) {
		return false
	}
	if filter.Type != "" {
		switch filter.Type {
		case "blocked":
			if !e.IsBlocked() {
				return false
			}
		case "cached":
			if !e.IsCached() {
				return false
			}
		case "resolved":
			if e.IsBlocked() || e.IsCached() {
				return false
			}
		}
	}
	return true
}

// FilterEntries applies a LogFilter to a set of entries.
func FilterEntries(entries []*LogEntry, filter LogFilter) []*LogEntry {
	if filter.Client == "" && filter.Domain == "" && filter.Type == "" {
		return entries
	}
	var result []*LogEntry
	for _, e := range entries {
		if MatchesFilter(e, filter) {
			result = append(result, e)
		}
	}
	return result
}
