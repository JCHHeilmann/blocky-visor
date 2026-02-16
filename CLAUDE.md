# Blocky Visor

Open-source static SPA dashboard for [Blocky](https://github.com/0xERR0R/blocky), a DNS proxy and ad blocker. Includes an optional Go sidecar for analytics, logs, and config management.

## Tech Stack

- **SvelteKit** with `adapter-static` (SPA mode, `fallback: 'index.html'`)
- **Svelte 5** runes (`$state`, `$derived`, `$effect`, `$props`, `$bindable`)
- **TailwindCSS v4** via `@tailwindcss/vite` plugin — config in `src/app.css` using `@theme`
- **TypeScript** throughout
- **Go** sidecar (`sidecar/`) using `go-chi/chi` router

## Commands

- `npm run dev` — start dev server
- `npm run build` — produce static files in `build/`
- `npm run preview` — preview production build
- `npm run check` — run svelte-check

## Architecture

```
src/
├── lib/
│   ├── api/           # Fetch wrapper + endpoint modules (blocking, query, lists, cache, metrics, sidecar)
│   ├── stores/        # Svelte 5 rune stores (settings, blocking, toasts, metricsHistory)
│   ├── components/
│   │   ├── layout/    # AppShell, Sidebar, Header
│   │   ├── ui/        # Reusable primitives (Button, Card, Toggle, Modal, Input, etc.)
│   │   ├── dashboard/ # BlockingToggle, StatsCard, ConnectionStatus
│   │   ├── query/     # QueryForm, QueryResult
│   │   ├── analytics/ # Charts, LogViewer, CardSkeleton
│   │   └── settings/  # ApiUrlForm
│   ├── utils/         # Prometheus metrics parser, formatters
│   └── types/         # TypeScript interfaces (BlockingStatus, DnsQueryResponse, ParsedMetrics, error classes)
├── routes/            # 8 pages: / (dashboard), /query, /lists, /cache, /settings,
│                      #           /analytics, /logs, /config (last 3 require sidecar)
└── app.css            # TailwindCSS v4 theme config (@theme block)
sidecar/               # Optional Go sidecar service
├── main.go            # Entry point, chi router setup
├── config.go          # Config loading (config.yaml)
├── config.example.yaml
├── handler/           # HTTP handlers (config, service, stats, logs)
├── logparser/         # Blocky log file parser
├── blocky/            # Blocky service interaction (systemctl)
└── middleware/         # Auth (X-API-Key) and CORS middleware
```

## Key Patterns

- **Runtime API URL** — stored in `localStorage` under `blocky-api-url`, configurable in Settings. No rebuild needed to change it. Default: `http://localhost:4000`
- **Sidecar URL** — also in `localStorage`, configured in Settings. Unlocks Analytics, Logs, Config pages in the sidebar when set
- **Stores** use Svelte 5 runes (`$state` + getters/setters), not legacy Svelte stores
- **Polling** — blocking status polls every N seconds (default 5), metrics every 30s. Pauses when tab is hidden via `document.hidden`
- **Metrics history** — `metricsHistoryStore` keeps an in-memory rolling history of metric snapshots for the Query Activity chart
- **Error types** — `ConnectionError` (network/timeout) vs `ApiError` (HTTP status errors), both in `src/lib/types/api.ts`
- **Theme** — System/Light/Dark toggle (oklch-based cool gray palette with teal/cyan accent, defined in `src/app.css` `@theme` block)
- **Fonts** — Inter Variable (sans) and JetBrains Mono Variable (mono), imported via `@fontsource-variable`
- **SPA mode** — `ssr = false`, `prerender = false` in `+layout.ts`
- **Query history** — DNS query page persists up to 20 entries in `localStorage` (key: `blocky-query-history`)

## Blocky API

Blocky REST API (no auth) at port 4000:

- `GET /api/blocking/status` — returns `{ enabled, disabledGroups?, autoEnableInSec? }`
- `GET /api/blocking/enable` — enable blocking
- `GET /api/blocking/disable?duration=5m&groups=ads` — disable blocking
- `POST /api/query` — body `{ query, type }` → DNS query result
- `POST /api/lists/refresh` — re-download blocking lists
- `PUT /api/cache/flush` — flush DNS cache
- `GET /metrics` — Prometheus metrics (enabled in Blocky config)

## Sidecar API

Go sidecar (default port 8550), all authenticated routes require `X-API-Key` header:

- `GET /api/health` — unauthenticated health check
- `GET /api/config` — read Blocky's `config.yml`
- `PUT /api/config` — write Blocky's config (creates timestamped backup)
- `GET /api/service/status` — systemd service status
- `POST /api/service/restart` — restart blocky systemd service
- `GET /api/stats?range=today|yesterday|7d|30d` — aggregated query analytics from log files
- `GET /api/stats/timeline?range=...&interval=5m|15m|1h` — time-bucketed query timeline
- `GET /api/logs?range=...&limit=...&offset=...&client=...&domain=...&type=...` — paginated, filtered log entries

## Blocky Prometheus Metrics

Metric names used by the parser (Blocky v0.28):

- `blocky_cache_hits_total`, `blocky_cache_misses_total`, `blocky_cache_entries`
- `blocky_denylist_cache_entries{group="..."}` — per-group denylist entry counts
- `blocky_prefetch_hits_total`, `blocky_prefetches_total`
- `blocky_blocking_enabled`, `blocky_error_total`

## Conventions

- Components use Svelte 5 `interface Props` + `$props()` pattern
- UI components accept `children` as `Snippet` type
- Form inputs use `$bindable()` for two-way binding
- Confirmation modals for destructive actions (disable blocking, refresh lists, flush cache)
- Toast notifications via `toastStore` for action feedback
- Sidebar conditionally shows Analytics/Logs/Config links only when sidecar is configured
- Loading skeletons (`CardSkeleton`) used on analytics page while data loads
