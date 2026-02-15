package logparser

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type cachedFile struct {
	modTime  time.Time
	size     int64
	stats    *StatsAccumulator
	timeline *TimelineAccumulator // always at 1-hour granularity
}

// StatsCache caches per-file accumulator state to avoid re-parsing immutable
// historical log files. Today's file is validated by mtime+size on each request.
type StatsCache struct {
	mu    sync.RWMutex
	files map[string]*cachedFile
}

func NewStatsCache() *StatsCache {
	return &StatsCache{
		files: make(map[string]*cachedFile),
	}
}

func (c *StatsCache) get(path string, info os.FileInfo) *cachedFile {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cf, ok := c.files[path]
	if !ok {
		return nil
	}
	if cf.modTime != info.ModTime() || cf.size != info.Size() {
		return nil
	}
	return cf
}

func (c *StatsCache) put(path string, info os.FileInfo, stats *StatsAccumulator, timeline *TimelineAccumulator) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.files[path] = &cachedFile{
		modTime:  info.ModTime(),
		size:     info.Size(),
		stats:    stats,
		timeline: timeline,
	}
}

// processFile parses a single log file, calling fn for each entry.
func processFile(path string, fn func(*LogEntry)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		entry, err := ParseLine(line)
		if err != nil {
			continue
		}
		fn(entry)
	}
	return scanner.Err()
}

// ComputeStats builds stats for a date range using cached per-file accumulators.
func (c *StatsCache) ComputeStats(logDir string, start, end time.Time) *StatsResponse {
	files := LogFilesForRange(logDir, start, end)
	combined := NewStatsAccumulator(start, end)

	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if cached := c.get(path, info); cached != nil {
			combined.Merge(cached.stats)
			continue
		}

		// Cache miss — parse file, cache result
		fileAcc := NewStatsAccumulator(start, end)
		fileTl := NewTimelineAccumulator(time.Hour)
		processFile(path, func(e *LogEntry) {
			fileAcc.Add(e)
			fileTl.Add(e)
		})
		c.put(path, info, fileAcc, fileTl)
		combined.Merge(fileAcc)
	}

	return combined.Finalize(len(files))
}

// ComputeTimeline builds timeline for a date range, re-aggregating cached
// hourly buckets to the requested interval.
func (c *StatsCache) ComputeTimeline(logDir string, start, end time.Time, interval time.Duration) []TimelineBucket {
	files := LogFilesForRange(logDir, start, end)
	combined := NewTimelineAccumulator(time.Hour)

	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if cached := c.get(path, info); cached != nil {
			combined.Merge(cached.timeline)
			continue
		}

		// Cache miss — parse, cache both stats + timeline
		fileAcc := NewStatsAccumulator(start, end)
		fileTl := NewTimelineAccumulator(time.Hour)
		processFile(path, func(e *LogEntry) {
			fileAcc.Add(e)
			fileTl.Add(e)
		})
		c.put(path, info, fileAcc, fileTl)
		combined.Merge(fileTl)
	}

	if interval == time.Hour {
		return combined.Finalize()
	}
	return combined.ReaggregateTo(interval).Finalize()
}
