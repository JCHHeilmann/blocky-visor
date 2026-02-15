package logparser

import (
	"math"
	"sort"
	"time"
)

type StatsResponse struct {
	Period             Period              `json:"period"`
	Summary            Summary             `json:"summary"`
	Hourly             []HourlyBucket      `json:"hourly"`
	TopDomains         []DomainCount       `json:"top_domains"`
	TopBlocked         []BlockedDomain     `json:"top_blocked"`
	Clients            []ClientStats       `json:"clients"`
	QueryTypes         map[string]int      `json:"query_types"`
	ResponseCategories map[string]int      `json:"response_categories"`
	ReturnCodes        map[string]int      `json:"return_codes"`
}

type Period struct {
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	FilesParsed int       `json:"files_parsed"`
}

type Summary struct {
	TotalQueries   int     `json:"total_queries"`
	BlockedQueries int     `json:"blocked_queries"`
	CachedQueries  int     `json:"cached_queries"`
	UniqueDomains  int     `json:"unique_domains"`
	UniqueClients  int     `json:"unique_clients"`
	AvgDurationMs  float64 `json:"avg_duration_ms"`
	P95DurationMs  float64 `json:"p95_duration_ms"`
}

type HourlyBucket struct {
	Hour    int `json:"hour"`
	Total   int `json:"total"`
	Blocked int `json:"blocked"`
	Cached  int `json:"cached"`
}

type DomainCount struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

type BlockedDomain struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
	Reason string `json:"reason"`
}

type ClientStats struct {
	IP      string `json:"ip"`
	Name    string `json:"name"`
	Total   int    `json:"total"`
	Blocked int    `json:"blocked"`
}

type TimelineBucket struct {
	Timestamp time.Time `json:"timestamp"`
	Total     int       `json:"total"`
	Blocked   int       `json:"blocked"`
	Cached    int       `json:"cached"`
}

// StatsAccumulator incrementally aggregates LogEntry data for stats computation.
// Entries can be fed one at a time via Add(), then Finalize() produces the response.
type StatsAccumulator struct {
	start, end         time.Time
	hourly             [24]HourlyBucket
	domainCounts       map[string]int
	blockedDomains     map[string]*BlockedDomain
	clientMap          map[string]*ClientStats
	uniqueDomains      map[string]struct{}
	queryTypes         map[string]int
	responseCategories map[string]int
	returnCodes        map[string]int
	durations          []float64
	durationSum        float64
	totalQueries       int
	blockedQueries     int
	cachedQueries      int
}

func NewStatsAccumulator(start, end time.Time) *StatsAccumulator {
	a := &StatsAccumulator{
		start:              start,
		end:                end,
		domainCounts:       make(map[string]int),
		blockedDomains:     make(map[string]*BlockedDomain),
		clientMap:          make(map[string]*ClientStats),
		uniqueDomains:      make(map[string]struct{}),
		queryTypes:         make(map[string]int),
		responseCategories: make(map[string]int),
		returnCodes:        make(map[string]int),
	}
	for i := range 24 {
		a.hourly[i].Hour = i
	}
	return a
}

// Add processes a single log entry into the accumulator.
func (a *StatsAccumulator) Add(e *LogEntry) {
	a.totalQueries++

	blocked := e.IsBlocked()
	cached := e.IsCached()

	if blocked {
		a.blockedQueries++
	}
	if cached {
		a.cachedQueries++
	}

	// Hourly
	h := e.Timestamp.Hour()
	a.hourly[h].Total++
	if blocked {
		a.hourly[h].Blocked++
	}
	if cached {
		a.hourly[h].Cached++
	}

	// Domains
	a.domainCounts[e.Domain]++
	a.uniqueDomains[e.Domain] = struct{}{}

	// Blocked domains
	if blocked {
		if bd, ok := a.blockedDomains[e.Domain]; ok {
			bd.Count++
		} else {
			a.blockedDomains[e.Domain] = &BlockedDomain{Domain: e.Domain, Count: 1, Reason: e.ResponseReason}
		}
	}

	// Clients
	if cs, ok := a.clientMap[e.ClientIP]; ok {
		cs.Total++
		if blocked {
			cs.Blocked++
		}
	} else {
		b := 0
		if blocked {
			b = 1
		}
		a.clientMap[e.ClientIP] = &ClientStats{IP: e.ClientIP, Name: e.ClientName, Total: 1, Blocked: b}
	}

	// Query types
	a.queryTypes[e.QueryType]++

	// Response categories
	a.responseCategories[e.ResponseCategory]++

	// Return codes
	a.returnCodes[e.ReturnCode]++

	// Durations
	a.durations = append(a.durations, e.DurationMs)
	a.durationSum += e.DurationMs
}

// Finalize computes the final StatsResponse from accumulated data.
func (a *StatsAccumulator) Finalize(filesParsed int) *StatsResponse {
	stats := &StatsResponse{
		Period:             Period{Start: a.start, End: a.end, FilesParsed: filesParsed},
		QueryTypes:         a.queryTypes,
		ResponseCategories: a.responseCategories,
		ReturnCodes:        a.returnCodes,
	}

	stats.Hourly = make([]HourlyBucket, 24)
	copy(stats.Hourly, a.hourly[:])

	stats.Summary = Summary{
		TotalQueries:   a.totalQueries,
		BlockedQueries: a.blockedQueries,
		CachedQueries:  a.cachedQueries,
		UniqueDomains:  len(a.uniqueDomains),
		UniqueClients:  len(a.clientMap),
	}

	if len(a.durations) > 0 {
		stats.Summary.AvgDurationMs = math.Round(a.durationSum/float64(len(a.durations))*10) / 10
		sort.Float64s(a.durations)
		p95idx := int(float64(len(a.durations)) * 0.95)
		if p95idx >= len(a.durations) {
			p95idx = len(a.durations) - 1
		}
		stats.Summary.P95DurationMs = a.durations[p95idx]
	}

	// Top domains (top 20)
	stats.TopDomains = topN(a.domainCounts, 20)

	// Top blocked (top 20)
	blockedList := make([]BlockedDomain, 0, len(a.blockedDomains))
	for _, bd := range a.blockedDomains {
		blockedList = append(blockedList, *bd)
	}
	sort.Slice(blockedList, func(i, j int) bool { return blockedList[i].Count > blockedList[j].Count })
	if len(blockedList) > 20 {
		blockedList = blockedList[:20]
	}
	stats.TopBlocked = blockedList

	// Clients sorted by total
	clients := make([]ClientStats, 0, len(a.clientMap))
	for _, cs := range a.clientMap {
		clients = append(clients, *cs)
	}
	sort.Slice(clients, func(i, j int) bool { return clients[i].Total > clients[j].Total })
	stats.Clients = clients

	return stats
}

// ComputeStats aggregates a slice of LogEntry into stats.
func ComputeStats(entries []*LogEntry, start, end time.Time, filesParsed int) *StatsResponse {
	acc := NewStatsAccumulator(start, end)
	for _, e := range entries {
		acc.Add(e)
	}
	return acc.Finalize(filesParsed)
}

// TimelineAccumulator incrementally aggregates LogEntry data for timeline computation.
type TimelineAccumulator struct {
	interval  time.Duration
	bucketMap map[int64]*TimelineBucket
}

func NewTimelineAccumulator(interval time.Duration) *TimelineAccumulator {
	return &TimelineAccumulator{
		interval:  interval,
		bucketMap: make(map[int64]*TimelineBucket),
	}
}

// Add processes a single log entry into the timeline accumulator.
func (a *TimelineAccumulator) Add(e *LogEntry) {
	key := e.Timestamp.Truncate(a.interval).Unix()
	b, ok := a.bucketMap[key]
	if !ok {
		b = &TimelineBucket{Timestamp: time.Unix(key, 0).UTC()}
		a.bucketMap[key] = b
	}
	b.Total++
	if e.IsBlocked() {
		b.Blocked++
	}
	if e.IsCached() {
		b.Cached++
	}
}

// Finalize returns sorted timeline buckets.
func (a *TimelineAccumulator) Finalize() []TimelineBucket {
	if len(a.bucketMap) == 0 {
		return nil
	}
	result := make([]TimelineBucket, 0, len(a.bucketMap))
	for _, b := range a.bucketMap {
		result = append(result, *b)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Timestamp.Before(result[j].Timestamp) })
	return result
}

// ComputeTimeline groups entries into time buckets of the given interval.
func ComputeTimeline(entries []*LogEntry, interval time.Duration) []TimelineBucket {
	if len(entries) == 0 {
		return nil
	}
	acc := NewTimelineAccumulator(interval)
	for _, e := range entries {
		acc.Add(e)
	}
	return acc.Finalize()
}

func topN(counts map[string]int, n int) []DomainCount {
	list := make([]DomainCount, 0, len(counts))
	for domain, count := range counts {
		list = append(list, DomainCount{Domain: domain, Count: count})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Count > list[j].Count })
	if len(list) > n {
		list = list[:n]
	}
	return list
}
