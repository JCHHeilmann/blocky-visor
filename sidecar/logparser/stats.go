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
	TotalQueries  int     `json:"total_queries"`
	BlockedQueries int    `json:"blocked_queries"`
	CachedQueries  int    `json:"cached_queries"`
	UniqueDomains  int    `json:"unique_domains"`
	UniqueClients  int    `json:"unique_clients"`
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

// ComputeStats aggregates a slice of LogEntry into stats.
func ComputeStats(entries []*LogEntry, start, end time.Time, filesParsed int) *StatsResponse {
	stats := &StatsResponse{
		Period: Period{Start: start, End: end, FilesParsed: filesParsed},
		QueryTypes:         make(map[string]int),
		ResponseCategories: make(map[string]int),
		ReturnCodes:        make(map[string]int),
	}

	// Initialize 24 hourly buckets
	stats.Hourly = make([]HourlyBucket, 24)
	for i := range 24 {
		stats.Hourly[i].Hour = i
	}

	domainCounts := make(map[string]int)
	blockedDomains := make(map[string]*BlockedDomain)
	clientMap := make(map[string]*ClientStats)
	uniqueDomains := make(map[string]bool)
	var durations []float64
	var durationSum float64

	for _, e := range entries {
		stats.Summary.TotalQueries++

		if e.IsBlocked() {
			stats.Summary.BlockedQueries++
		}
		if e.IsCached() {
			stats.Summary.CachedQueries++
		}

		// Hourly
		h := e.Timestamp.Hour()
		stats.Hourly[h].Total++
		if e.IsBlocked() {
			stats.Hourly[h].Blocked++
		}
		if e.IsCached() {
			stats.Hourly[h].Cached++
		}

		// Domains
		domainCounts[e.Domain]++
		uniqueDomains[e.Domain] = true

		// Blocked domains
		if e.IsBlocked() {
			if bd, ok := blockedDomains[e.Domain]; ok {
				bd.Count++
			} else {
				blockedDomains[e.Domain] = &BlockedDomain{Domain: e.Domain, Count: 1, Reason: e.ResponseReason}
			}
		}

		// Clients
		if cs, ok := clientMap[e.ClientIP]; ok {
			cs.Total++
			if e.IsBlocked() {
				cs.Blocked++
			}
		} else {
			blocked := 0
			if e.IsBlocked() {
				blocked = 1
			}
			clientMap[e.ClientIP] = &ClientStats{IP: e.ClientIP, Name: e.ClientName, Total: 1, Blocked: blocked}
		}

		// Query types
		stats.QueryTypes[e.QueryType]++

		// Response categories
		stats.ResponseCategories[e.ResponseCategory]++

		// Return codes
		stats.ReturnCodes[e.ReturnCode]++

		// Durations
		durations = append(durations, e.DurationMs)
		durationSum += e.DurationMs
	}

	stats.Summary.UniqueDomains = len(uniqueDomains)
	stats.Summary.UniqueClients = len(clientMap)

	if len(durations) > 0 {
		stats.Summary.AvgDurationMs = math.Round(durationSum/float64(len(durations))*10) / 10
		sort.Float64s(durations)
		p95idx := int(float64(len(durations)) * 0.95)
		if p95idx >= len(durations) {
			p95idx = len(durations) - 1
		}
		stats.Summary.P95DurationMs = durations[p95idx]
	}

	// Top domains (top 20)
	stats.TopDomains = topN(domainCounts, 20)

	// Top blocked (top 20)
	blockedList := make([]BlockedDomain, 0, len(blockedDomains))
	for _, bd := range blockedDomains {
		blockedList = append(blockedList, *bd)
	}
	sort.Slice(blockedList, func(i, j int) bool { return blockedList[i].Count > blockedList[j].Count })
	if len(blockedList) > 20 {
		blockedList = blockedList[:20]
	}
	stats.TopBlocked = blockedList

	// Clients sorted by total
	clients := make([]ClientStats, 0, len(clientMap))
	for _, cs := range clientMap {
		clients = append(clients, *cs)
	}
	sort.Slice(clients, func(i, j int) bool { return clients[i].Total > clients[j].Total })
	stats.Clients = clients

	return stats
}

// ComputeTimeline groups entries into time buckets of the given interval.
func ComputeTimeline(entries []*LogEntry, interval time.Duration) []TimelineBucket {
	if len(entries) == 0 {
		return nil
	}

	bucketMap := make(map[int64]*TimelineBucket)
	for _, e := range entries {
		key := e.Timestamp.Truncate(interval).Unix()
		b, ok := bucketMap[key]
		if !ok {
			b = &TimelineBucket{Timestamp: time.Unix(key, 0).UTC()}
			bucketMap[key] = b
		}
		b.Total++
		if e.IsBlocked() {
			b.Blocked++
		}
		if e.IsCached() {
			b.Cached++
		}
	}

	result := make([]TimelineBucket, 0, len(bucketMap))
	for _, b := range bucketMap {
		result = append(result, *b)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Timestamp.Before(result[j].Timestamp) })
	return result
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
