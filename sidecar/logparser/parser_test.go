package logparser

import (
	"testing"
	"time"
)

func TestParseLine(t *testing.T) {
	line := "2026-02-14 00:00:37\t10.0.0.101\t10.0.0.101\t0\tCACHED\tbag-cdn.itunes-apple.com.akadns.net.\tCNAME (...), A (151.101.131.6)\tNOERROR\tCACHED\tA\tblocky"

	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("ParseLine() error: %v", err)
	}

	expected := time.Date(2026, 2, 14, 0, 0, 37, 0, time.UTC)
	if !entry.Timestamp.Equal(expected) {
		t.Errorf("Timestamp = %v, want %v", entry.Timestamp, expected)
	}
	if entry.ClientIP != "10.0.0.101" {
		t.Errorf("ClientIP = %q, want %q", entry.ClientIP, "10.0.0.101")
	}
	if entry.ClientName != "10.0.0.101" {
		t.Errorf("ClientName = %q, want %q", entry.ClientName, "10.0.0.101")
	}
	if entry.DurationMs != 0 {
		t.Errorf("DurationMs = %f, want 0", entry.DurationMs)
	}
	if entry.ResponseReason != "CACHED" {
		t.Errorf("ResponseReason = %q, want %q", entry.ResponseReason, "CACHED")
	}
	if entry.Domain != "bag-cdn.itunes-apple.com.akadns.net." {
		t.Errorf("Domain = %q", entry.Domain)
	}
	if entry.ReturnCode != "NOERROR" {
		t.Errorf("ReturnCode = %q, want %q", entry.ReturnCode, "NOERROR")
	}
	if entry.QueryType != "A" {
		t.Errorf("QueryType = %q, want %q", entry.QueryType, "A")
	}
	if entry.Source != "blocky" {
		t.Errorf("Source = %q, want %q", entry.Source, "blocky")
	}
}

func TestParseLineBlocked(t *testing.T) {
	line := "2026-02-14 12:30:00\t10.0.0.50\tdesktop.local\t1\tBLOCKED (ads)\tad.doubleclick.net.\t\tNOERROR\tBLOCKED (ads)\tA\tblocky"
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("ParseLine() error: %v", err)
	}
	if !entry.IsBlocked() {
		t.Error("expected IsBlocked() = true")
	}
	if entry.IsCached() {
		t.Error("expected IsCached() = false")
	}
}

func TestParseLineCached(t *testing.T) {
	line := "2026-02-14 12:30:00\t10.0.0.50\tdesktop.local\t0\tCACHED\texample.com.\tA (1.2.3.4)\tNOERROR\tCACHED\tAAAA\tblocky"
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("ParseLine() error: %v", err)
	}
	if entry.IsBlocked() {
		t.Error("expected IsBlocked() = false")
	}
	if !entry.IsCached() {
		t.Error("expected IsCached() = true")
	}
}

func TestParseLineInvalidFields(t *testing.T) {
	_, err := ParseLine("not enough\tfields")
	if err == nil {
		t.Error("expected error for insufficient fields")
	}
}

func TestParseLineInvalidTimestamp(t *testing.T) {
	line := "not-a-date\t1\t2\t3\t4\t5\t6\t7\t8\t9\t10"
	_, err := ParseLine(line)
	if err == nil {
		t.Error("expected error for invalid timestamp")
	}
}

func TestStatsAccumulator(t *testing.T) {
	start := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)
	acc := NewStatsAccumulator(start, end)

	entries := []*LogEntry{
		{Timestamp: time.Date(2026, 2, 14, 0, 30, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "RESOLVED", ReturnCode: "NOERROR", DurationMs: 10, ResponseCategory: "RESOLVED"},
		{Timestamp: time.Date(2026, 2, 14, 0, 31, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "CACHED", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "CACHED"},
		{Timestamp: time.Date(2026, 2, 14, 1, 0, 0, 0, time.UTC), ClientIP: "10.0.0.2", ClientName: "Phone", Domain: "ad.doubleclick.net.", QueryType: "A", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "BLOCKED (ads)"},
	}

	for _, e := range entries {
		acc.Add(e)
	}
	stats := acc.Finalize(1)

	if stats.Summary.TotalQueries != 3 {
		t.Errorf("TotalQueries = %d, want 3", stats.Summary.TotalQueries)
	}
	if stats.Summary.BlockedQueries != 1 {
		t.Errorf("BlockedQueries = %d, want 1", stats.Summary.BlockedQueries)
	}
	if stats.Summary.CachedQueries != 1 {
		t.Errorf("CachedQueries = %d, want 1", stats.Summary.CachedQueries)
	}
	if stats.Summary.UniqueClients != 2 {
		t.Errorf("UniqueClients = %d, want 2", stats.Summary.UniqueClients)
	}
	if stats.Summary.UniqueDomains != 2 {
		t.Errorf("UniqueDomains = %d, want 2", stats.Summary.UniqueDomains)
	}
	if stats.Period.FilesParsed != 1 {
		t.Errorf("FilesParsed = %d, want 1", stats.Period.FilesParsed)
	}
	if stats.Hourly[0].Total != 2 {
		t.Errorf("Hour 0 total = %d, want 2", stats.Hourly[0].Total)
	}
	if stats.Hourly[1].Blocked != 1 {
		t.Errorf("Hour 1 blocked = %d, want 1", stats.Hourly[1].Blocked)
	}
	if len(stats.TopBlocked) != 1 || stats.TopBlocked[0].Domain != "ad.doubleclick.net." {
		t.Errorf("TopBlocked = %v, want ad.doubleclick.net.", stats.TopBlocked)
	}
}

func TestTimelineAccumulator(t *testing.T) {
	acc := NewTimelineAccumulator(time.Hour)

	entries := []*LogEntry{
		{Timestamp: time.Date(2026, 2, 14, 0, 10, 0, 0, time.UTC), ResponseReason: "RESOLVED"},
		{Timestamp: time.Date(2026, 2, 14, 0, 20, 0, 0, time.UTC), ResponseReason: "CACHED"},
		{Timestamp: time.Date(2026, 2, 14, 1, 5, 0, 0, time.UTC), ResponseReason: "BLOCKED (ads)"},
	}

	for _, e := range entries {
		acc.Add(e)
	}
	buckets := acc.Finalize()

	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
	if buckets[0].Total != 2 {
		t.Errorf("bucket[0].Total = %d, want 2", buckets[0].Total)
	}
	if buckets[0].Cached != 1 {
		t.Errorf("bucket[0].Cached = %d, want 1", buckets[0].Cached)
	}
	if buckets[1].Total != 1 {
		t.Errorf("bucket[1].Total = %d, want 1", buckets[1].Total)
	}
	if buckets[1].Blocked != 1 {
		t.Errorf("bucket[1].Blocked = %d, want 1", buckets[1].Blocked)
	}
}

func TestStatsAccumulatorMatchesComputeStats(t *testing.T) {
	entries := []*LogEntry{
		{Timestamp: time.Date(2026, 2, 14, 0, 30, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "RESOLVED", ReturnCode: "NOERROR", DurationMs: 10, ResponseCategory: "RESOLVED"},
		{Timestamp: time.Date(2026, 2, 14, 0, 31, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "CACHED", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "CACHED"},
		{Timestamp: time.Date(2026, 2, 14, 1, 0, 0, 0, time.UTC), ClientIP: "10.0.0.2", ClientName: "Phone", Domain: "ad.doubleclick.net.", QueryType: "A", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "BLOCKED (ads)"},
		{Timestamp: time.Date(2026, 2, 14, 1, 5, 0, 0, time.UTC), ClientIP: "10.0.0.2", ClientName: "Phone", Domain: "tracker.example.com.", QueryType: "AAAA", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "BLOCKED (ads)"},
	}

	start := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)

	// Use ComputeStats (wrapper)
	wrapper := ComputeStats(entries, start, end, 1)

	// Use accumulator directly
	acc := NewStatsAccumulator(start, end)
	for _, e := range entries {
		acc.Add(e)
	}
	direct := acc.Finalize(1)

	if wrapper.Summary != direct.Summary {
		t.Errorf("Summary mismatch:\n  wrapper: %+v\n  direct:  %+v", wrapper.Summary, direct.Summary)
	}
	if wrapper.Period != direct.Period {
		t.Errorf("Period mismatch:\n  wrapper: %+v\n  direct:  %+v", wrapper.Period, direct.Period)
	}
}

func TestStatsAggregation(t *testing.T) {
	entries := []*LogEntry{
		{Timestamp: time.Date(2026, 2, 14, 0, 30, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "RESOLVED", ReturnCode: "NOERROR", DurationMs: 10, ResponseCategory: "RESOLVED"},
		{Timestamp: time.Date(2026, 2, 14, 0, 31, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "CACHED", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "CACHED"},
		{Timestamp: time.Date(2026, 2, 14, 1, 0, 0, 0, time.UTC), ClientIP: "10.0.0.2", ClientName: "Phone", Domain: "ad.doubleclick.net.", QueryType: "A", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "BLOCKED (ads)"},
		{Timestamp: time.Date(2026, 2, 14, 1, 5, 0, 0, time.UTC), ClientIP: "10.0.0.2", ClientName: "Phone", Domain: "tracker.example.com.", QueryType: "AAAA", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "BLOCKED (ads)"},
	}

	stats := ComputeStats(entries, time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC), time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC), 1)

	if stats.Summary.TotalQueries != 4 {
		t.Errorf("TotalQueries = %d, want 4", stats.Summary.TotalQueries)
	}
	if stats.Summary.BlockedQueries != 2 {
		t.Errorf("BlockedQueries = %d, want 2", stats.Summary.BlockedQueries)
	}
	if stats.Summary.CachedQueries != 1 {
		t.Errorf("CachedQueries = %d, want 1", stats.Summary.CachedQueries)
	}
	if stats.Summary.UniqueClients != 2 {
		t.Errorf("UniqueClients = %d, want 2", stats.Summary.UniqueClients)
	}
	if stats.Summary.UniqueDomains != 3 {
		t.Errorf("UniqueDomains = %d, want 3", stats.Summary.UniqueDomains)
	}
	if len(stats.Hourly) != 24 {
		t.Errorf("Hourly length = %d, want 24", len(stats.Hourly))
	}
	// Hour 0 should have 2 queries
	if stats.Hourly[0].Total != 2 {
		t.Errorf("Hour 0 total = %d, want 2", stats.Hourly[0].Total)
	}
	// Hour 1 should have 2 queries
	if stats.Hourly[1].Total != 2 {
		t.Errorf("Hour 1 total = %d, want 2", stats.Hourly[1].Total)
	}
	if len(stats.TopDomains) == 0 {
		t.Error("expected TopDomains to be non-empty")
	}
	if stats.TopDomains[0].Domain != "example.com." {
		t.Errorf("top domain = %q, want %q", stats.TopDomains[0].Domain, "example.com.")
	}
}
