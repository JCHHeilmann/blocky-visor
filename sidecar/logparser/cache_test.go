package logparser

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestStatsAccumulatorMerge(t *testing.T) {
	start := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)

	// Build two accumulators from separate entry sets
	a := NewStatsAccumulator(start, end)
	a.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 0, 30, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "RESOLVED", ReturnCode: "NOERROR", DurationMs: 10, ResponseCategory: "RESOLVED"})
	a.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 0, 31, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "CACHED", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "CACHED"})

	b := NewStatsAccumulator(start, end)
	b.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 1, 0, 0, 0, time.UTC), ClientIP: "10.0.0.2", ClientName: "Phone", Domain: "ad.doubleclick.net.", QueryType: "A", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "BLOCKED (ads)"})
	b.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 1, 5, 0, 0, time.UTC), ClientIP: "10.0.0.2", ClientName: "Phone", Domain: "tracker.example.com.", QueryType: "AAAA", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "BLOCKED (ads)"})

	// Merge b into a
	a.Merge(b)
	merged := a.Finalize(2)

	// Compare against computing all entries at once
	all := NewStatsAccumulator(start, end)
	all.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 0, 30, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "RESOLVED", ReturnCode: "NOERROR", DurationMs: 10, ResponseCategory: "RESOLVED"})
	all.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 0, 31, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "CACHED", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "CACHED"})
	all.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 1, 0, 0, 0, time.UTC), ClientIP: "10.0.0.2", ClientName: "Phone", Domain: "ad.doubleclick.net.", QueryType: "A", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "BLOCKED (ads)"})
	all.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 1, 5, 0, 0, time.UTC), ClientIP: "10.0.0.2", ClientName: "Phone", Domain: "tracker.example.com.", QueryType: "AAAA", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "BLOCKED (ads)"})
	direct := all.Finalize(2)

	if merged.Summary != direct.Summary {
		t.Errorf("Summary mismatch:\n  merged: %+v\n  direct: %+v", merged.Summary, direct.Summary)
	}
	for i := range 24 {
		if merged.Hourly[i] != direct.Hourly[i] {
			t.Errorf("Hourly[%d] mismatch: merged=%+v direct=%+v", i, merged.Hourly[i], direct.Hourly[i])
		}
	}
	if len(merged.TopBlocked) != len(direct.TopBlocked) {
		t.Errorf("TopBlocked length: merged=%d direct=%d", len(merged.TopBlocked), len(direct.TopBlocked))
	}
	if len(merged.Clients) != len(direct.Clients) {
		t.Errorf("Clients length: merged=%d direct=%d", len(merged.Clients), len(direct.Clients))
	}
	if len(merged.TopDomains) != len(direct.TopDomains) {
		t.Errorf("TopDomains length: merged=%d direct=%d", len(merged.TopDomains), len(direct.TopDomains))
	}
}

func TestStatsAccumulatorMergeOverlappingKeys(t *testing.T) {
	start := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)

	a := NewStatsAccumulator(start, end)
	a.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 10, 0, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "RESOLVED", ReturnCode: "NOERROR", DurationMs: 5, ResponseCategory: "RESOLVED"})
	a.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 10, 1, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "ad.tracker.net.", QueryType: "A", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 1, ResponseCategory: "BLOCKED (ads)"})

	b := NewStatsAccumulator(start, end)
	// Same client, same domain, same blocked domain — should merge counts
	b.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 10, 2, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "example.com.", QueryType: "A", ResponseReason: "RESOLVED", ReturnCode: "NOERROR", DurationMs: 8, ResponseCategory: "RESOLVED"})
	b.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 10, 3, 0, 0, time.UTC), ClientIP: "10.0.0.1", ClientName: "PC", Domain: "ad.tracker.net.", QueryType: "A", ResponseReason: "BLOCKED (ads)", ReturnCode: "NOERROR", DurationMs: 0, ResponseCategory: "BLOCKED (ads)"})

	a.Merge(b)
	stats := a.Finalize(2)

	if stats.Summary.TotalQueries != 4 {
		t.Errorf("TotalQueries = %d, want 4", stats.Summary.TotalQueries)
	}
	if stats.Summary.BlockedQueries != 2 {
		t.Errorf("BlockedQueries = %d, want 2", stats.Summary.BlockedQueries)
	}
	if stats.Summary.UniqueClients != 1 {
		t.Errorf("UniqueClients = %d, want 1", stats.Summary.UniqueClients)
	}
	if stats.Summary.UniqueDomains != 2 {
		t.Errorf("UniqueDomains = %d, want 2", stats.Summary.UniqueDomains)
	}
	// Client should have merged totals
	if len(stats.Clients) != 1 {
		t.Fatalf("expected 1 client, got %d", len(stats.Clients))
	}
	if stats.Clients[0].Total != 4 {
		t.Errorf("Client total = %d, want 4", stats.Clients[0].Total)
	}
	if stats.Clients[0].Blocked != 2 {
		t.Errorf("Client blocked = %d, want 2", stats.Clients[0].Blocked)
	}
	// Blocked domain count should be merged
	if len(stats.TopBlocked) != 1 || stats.TopBlocked[0].Count != 2 {
		t.Errorf("TopBlocked = %+v, want 1 entry with count 2", stats.TopBlocked)
	}
	// Hourly: all 4 in hour 10
	if stats.Hourly[10].Total != 4 {
		t.Errorf("Hourly[10].Total = %d, want 4", stats.Hourly[10].Total)
	}
}

func TestTimelineAccumulatorMerge(t *testing.T) {
	a := NewTimelineAccumulator(time.Hour)
	a.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 0, 10, 0, 0, time.UTC), ResponseReason: "RESOLVED"})
	a.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 0, 20, 0, 0, time.UTC), ResponseReason: "CACHED"})

	b := NewTimelineAccumulator(time.Hour)
	b.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 0, 30, 0, 0, time.UTC), ResponseReason: "RESOLVED"})
	b.Add(&LogEntry{Timestamp: time.Date(2026, 2, 14, 1, 5, 0, 0, time.UTC), ResponseReason: "BLOCKED (ads)"})

	a.Merge(b)
	buckets := a.Finalize()

	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
	// Hour 0: 3 total (2 from a + 1 from b), 1 cached
	if buckets[0].Total != 3 {
		t.Errorf("bucket[0].Total = %d, want 3", buckets[0].Total)
	}
	if buckets[0].Cached != 1 {
		t.Errorf("bucket[0].Cached = %d, want 1", buckets[0].Cached)
	}
	// Hour 1: 1 total, 1 blocked
	if buckets[1].Total != 1 {
		t.Errorf("bucket[1].Total = %d, want 1", buckets[1].Total)
	}
	if buckets[1].Blocked != 1 {
		t.Errorf("bucket[1].Blocked = %d, want 1", buckets[1].Blocked)
	}
}

func TestReaggregateTo(t *testing.T) {
	acc := NewTimelineAccumulator(time.Hour)

	// Add entries across two days, multiple hours
	entries := []*LogEntry{
		{Timestamp: time.Date(2026, 2, 14, 0, 10, 0, 0, time.UTC), ResponseReason: "RESOLVED"},
		{Timestamp: time.Date(2026, 2, 14, 1, 10, 0, 0, time.UTC), ResponseReason: "RESOLVED"},
		{Timestamp: time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC), ResponseReason: "CACHED"},
		{Timestamp: time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC), ResponseReason: "BLOCKED (ads)"},
		{Timestamp: time.Date(2026, 2, 15, 3, 0, 0, 0, time.UTC), ResponseReason: "RESOLVED"},
	}
	for _, e := range entries {
		acc.Add(e)
	}

	daily := acc.ReaggregateTo(24 * time.Hour)
	buckets := daily.Finalize()

	if len(buckets) != 2 {
		t.Fatalf("expected 2 daily buckets, got %d", len(buckets))
	}
	// Day 1: 3 entries
	if buckets[0].Total != 3 {
		t.Errorf("day1.Total = %d, want 3", buckets[0].Total)
	}
	if buckets[0].Cached != 1 {
		t.Errorf("day1.Cached = %d, want 1", buckets[0].Cached)
	}
	// Day 2: 2 entries
	if buckets[1].Total != 2 {
		t.Errorf("day2.Total = %d, want 2", buckets[1].Total)
	}
	if buckets[1].Blocked != 1 {
		t.Errorf("day2.Blocked = %d, want 1", buckets[1].Blocked)
	}
}

func TestStatsCacheWithFiles(t *testing.T) {
	// Create temp dir with a fake log file
	dir := t.TempDir()

	// Write a log file
	logFile := dir + "/2026-02-14_ALL.log"
	writeTestLogFile(t, logFile, []string{
		"2026-02-14 00:00:37\t10.0.0.1\tPC\t0\tCACHED\texample.com.\tA (1.2.3.4)\tNOERROR\tCACHED\tA\tblocky",
		"2026-02-14 01:00:00\t10.0.0.2\tPhone\t1\tBLOCKED (ads)\tad.tracker.net.\t\tNOERROR\tBLOCKED (ads)\tA\tblocky",
		"2026-02-14 02:00:00\t10.0.0.1\tPC\t5\tRESOLVED\tgoogle.com.\tA (142.250.80.46)\tNOERROR\tRESOLVED\tA\tblocky",
	})

	cache := NewStatsCache()
	start := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 14, 23, 59, 59, 0, time.UTC)

	// First call — cache miss, parses file
	stats1 := cache.ComputeStats(dir, start, end)
	if stats1.Summary.TotalQueries != 3 {
		t.Errorf("First call: TotalQueries = %d, want 3", stats1.Summary.TotalQueries)
	}
	if stats1.Summary.BlockedQueries != 1 {
		t.Errorf("First call: BlockedQueries = %d, want 1", stats1.Summary.BlockedQueries)
	}
	if stats1.Summary.CachedQueries != 1 {
		t.Errorf("First call: CachedQueries = %d, want 1", stats1.Summary.CachedQueries)
	}

	// Second call — should use cache, same results
	stats2 := cache.ComputeStats(dir, start, end)
	if stats2.Summary != stats1.Summary {
		t.Errorf("Cached result differs:\n  first:  %+v\n  second: %+v", stats1.Summary, stats2.Summary)
	}

	// Verify cache is populated
	cache.mu.RLock()
	if len(cache.files) != 1 {
		t.Errorf("Expected 1 cached file, got %d", len(cache.files))
	}
	cache.mu.RUnlock()
}

func TestStatsCacheTimelineWithFiles(t *testing.T) {
	dir := t.TempDir()
	logFile := dir + "/2026-02-14_ALL.log"
	writeTestLogFile(t, logFile, []string{
		"2026-02-14 00:10:00\t10.0.0.1\tPC\t0\tRESOLVED\texample.com.\tA (1.2.3.4)\tNOERROR\tRESOLVED\tA\tblocky",
		"2026-02-14 00:20:00\t10.0.0.1\tPC\t0\tCACHED\texample.com.\tA (1.2.3.4)\tNOERROR\tCACHED\tA\tblocky",
		"2026-02-14 01:05:00\t10.0.0.2\tPhone\t1\tBLOCKED (ads)\tad.tracker.net.\t\tNOERROR\tBLOCKED (ads)\tA\tblocky",
	})

	cache := NewStatsCache()
	start := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 14, 23, 59, 59, 0, time.UTC)

	// Hourly timeline
	tl := cache.ComputeTimeline(dir, start, end, time.Hour)
	if len(tl) != 2 {
		t.Fatalf("Expected 2 hourly buckets, got %d", len(tl))
	}
	if tl[0].Total != 2 {
		t.Errorf("Hour 0 total = %d, want 2", tl[0].Total)
	}
	if tl[1].Blocked != 1 {
		t.Errorf("Hour 1 blocked = %d, want 1", tl[1].Blocked)
	}

	// Daily re-aggregation
	daily := cache.ComputeTimeline(dir, start, end, 24*time.Hour)
	if len(daily) != 1 {
		t.Fatalf("Expected 1 daily bucket, got %d", len(daily))
	}
	if daily[0].Total != 3 {
		t.Errorf("Daily total = %d, want 3", daily[0].Total)
	}
}

func TestStatsCacheInvalidation(t *testing.T) {
	dir := t.TempDir()
	logFile := dir + "/2026-02-14_ALL.log"

	// Initial file with 2 entries
	writeTestLogFile(t, logFile, []string{
		"2026-02-14 00:00:00\t10.0.0.1\tPC\t0\tRESOLVED\texample.com.\tA (1.2.3.4)\tNOERROR\tRESOLVED\tA\tblocky",
		"2026-02-14 01:00:00\t10.0.0.1\tPC\t0\tRESOLVED\texample.com.\tA (1.2.3.4)\tNOERROR\tRESOLVED\tA\tblocky",
	})

	cache := NewStatsCache()
	start := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 14, 23, 59, 59, 0, time.UTC)

	stats1 := cache.ComputeStats(dir, start, end)
	if stats1.Summary.TotalQueries != 2 {
		t.Fatalf("First: TotalQueries = %d, want 2", stats1.Summary.TotalQueries)
	}

	// Append a line (simulates today's file growing)
	writeTestLogFile(t, logFile, []string{
		"2026-02-14 00:00:00\t10.0.0.1\tPC\t0\tRESOLVED\texample.com.\tA (1.2.3.4)\tNOERROR\tRESOLVED\tA\tblocky",
		"2026-02-14 01:00:00\t10.0.0.1\tPC\t0\tRESOLVED\texample.com.\tA (1.2.3.4)\tNOERROR\tRESOLVED\tA\tblocky",
		"2026-02-14 02:00:00\t10.0.0.1\tPC\t0\tCACHED\texample.com.\tA (1.2.3.4)\tNOERROR\tCACHED\tA\tblocky",
	})

	// Cache should be invalidated because file size changed
	stats2 := cache.ComputeStats(dir, start, end)
	if stats2.Summary.TotalQueries != 3 {
		t.Errorf("After append: TotalQueries = %d, want 3", stats2.Summary.TotalQueries)
	}
	if stats2.Summary.CachedQueries != 1 {
		t.Errorf("After append: CachedQueries = %d, want 1", stats2.Summary.CachedQueries)
	}
}

func writeTestLogFile(t *testing.T, path string, lines []string) {
	t.Helper()
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write test log file: %v", err)
	}
}
