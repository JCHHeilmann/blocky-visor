package logparser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// LogEntry represents a single parsed Blocky query log line.
type LogEntry struct {
	Timestamp        time.Time `json:"timestamp"`
	ClientIP         string    `json:"client_ip"`
	ClientName       string    `json:"client_name"`
	ResolvedName     string    `json:"resolved_name,omitempty"`
	DurationMs       float64   `json:"duration_ms"`
	ResponseReason   string    `json:"response_reason"`
	Domain           string    `json:"domain"`
	ResponseAnswer   string    `json:"response_answer"`
	ReturnCode       string    `json:"return_code"`
	ResponseCategory string    `json:"response_category"`
	QueryType        string    `json:"query_type"`
	Source           string    `json:"source"`
}

// ParseLine parses a single TSV log line into a LogEntry.
// Format: timestamp\tclientIP\tclientName\tduration\tresponseReason\tdomain\tresponseAnswer\treturnCode\tresponseCategory\tqueryType\tsource
func ParseLine(line string) (*LogEntry, error) {
	fields := strings.Split(line, "\t")
	if len(fields) < 11 {
		return nil, fmt.Errorf("expected 11 tab-separated fields, got %d", len(fields))
	}

	ts, err := time.ParseInLocation("2006-01-02 15:04:05", fields[0], time.Local)
	if err != nil {
		return nil, fmt.Errorf("parse timestamp: %w", err)
	}

	duration, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		duration = 0
	}

	return &LogEntry{
		Timestamp:        ts,
		ClientIP:         fields[1],
		ClientName:       fields[2],
		DurationMs:       duration,
		ResponseReason:   fields[4],
		Domain:           fields[5],
		ResponseAnswer:   fields[6],
		ReturnCode:       fields[7],
		ResponseCategory: fields[8],
		QueryType:        fields[9],
		Source:           fields[10],
	}, nil
}

// IsBlocked returns true if the entry was blocked.
func (e *LogEntry) IsBlocked() bool {
	return strings.HasPrefix(strings.ToUpper(e.ResponseReason), "BLOCKED")
}

// IsCached returns true if the entry was served from cache.
func (e *LogEntry) IsCached() bool {
	upper := strings.ToUpper(e.ResponseReason)
	return upper == "CACHED" || strings.HasPrefix(upper, "CACHED ")
}
