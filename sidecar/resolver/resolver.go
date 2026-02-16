package resolver

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"
)

// HostResolver performs reverse DNS lookups with caching.
// It uses a configurable DNS server to avoid routing PTR queries
// through Blocky (which would pollute the very logs we're displaying).
type HostResolver struct {
	resolver *net.Resolver
	mu       sync.RWMutex
	cache    map[string]cacheEntry
	ttl      time.Duration
}

type cacheEntry struct {
	name      string
	expiresAt time.Time
}

// New creates a HostResolver. If dnsServer is empty, uses the system resolver.
// If dnsServer is set (e.g. "192.168.178.1:53"), queries that server directly.
func New(dnsServer string) *HostResolver {
	var r *net.Resolver
	if dnsServer != "" {
		if !strings.Contains(dnsServer, ":") {
			dnsServer += ":53"
		}
		r = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: 2 * time.Second}
				return d.DialContext(ctx, "udp", dnsServer)
			},
		}
	} else {
		r = net.DefaultResolver
	}

	return &HostResolver{
		resolver: r,
		cache:    make(map[string]cacheEntry),
		ttl:      1 * time.Hour,
	}
}

// Lookup resolves an IP to a hostname. Returns empty string if unresolvable.
// Results are cached.
func (h *HostResolver) Lookup(ip string) string {
	h.mu.RLock()
	if entry, ok := h.cache[ip]; ok && time.Now().Before(entry.expiresAt) {
		h.mu.RUnlock()
		return entry.name
	}
	h.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	names, err := h.resolver.LookupAddr(ctx, ip)
	name := ""
	if err == nil && len(names) > 0 {
		name = strings.TrimSuffix(names[0], ".")
	}

	h.mu.Lock()
	h.cache[ip] = cacheEntry{name: name, expiresAt: time.Now().Add(h.ttl)}
	h.mu.Unlock()

	return name
}
