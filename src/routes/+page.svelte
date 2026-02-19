<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import BlockingToggle from "$lib/components/dashboard/BlockingToggle.svelte";
  import StatsCard from "$lib/components/dashboard/StatsCard.svelte";
  import ConnectionStatus from "$lib/components/dashboard/ConnectionStatus.svelte";
  import ActivityChart from "$lib/components/dashboard/ActivityChart.svelte";
  import ResponseBreakdown from "$lib/components/dashboard/ResponseBreakdown.svelte";
  import ClientTable from "$lib/components/dashboard/ClientTable.svelte";
  import ModeBar from "$lib/components/dashboard/ModeBar.svelte";
  import ResponseBreakdownChart from "$lib/components/analytics/ResponseBreakdownChart.svelte";
  import ClientBreakdown from "$lib/components/analytics/ClientBreakdown.svelte";
  import HourlyChart from "$lib/components/analytics/HourlyChart.svelte";
  import DailyChart from "$lib/components/analytics/DailyChart.svelte";
  import TopDomainsTable from "$lib/components/analytics/TopDomainsTable.svelte";
  import DistributionChart from "$lib/components/analytics/DistributionChart.svelte";
  import CardSkeleton from "$lib/components/analytics/CardSkeleton.svelte";
  import { settingsStore } from "$lib/stores/settings.svelte";
  import { sidecarStore } from "$lib/stores/sidecar.svelte";
  import { metricsHistoryStore } from "$lib/stores/metrics-history.svelte";
  import { fetchMetrics } from "$lib/api/metrics";
  import {
    fetchStats,
    fetchTimeline,
    type StatsRange,
  } from "$lib/api/sidecar-stats";
  import { formatNumber, formatPercentage } from "$lib/utils/formatters";
  import type {
    ParsedMetrics,
    SidecarStatsResponse,
    SidecarTimelineBucket,
  } from "$lib/types/api";

  // --- State ---
  let metrics = $state<ParsedMetrics | null>(null);
  let mode = $state<"live" | "historical">("live");
  let range = $state<StatsRange>("today");
  let sidecarStats = $state<SidecarStatsResponse | null>(null);
  let timeline = $state<SidecarTimelineBucket[] | null>(null);
  let sidecarLoading = $state(false);
  let sidecarError = $state<string | null>(null);

  // --- Prometheus polling (always runs) ---
  $effect(() => {
    loadMetrics();
    const interval = setInterval(() => {
      if (!document.hidden) loadMetrics();
    }, settingsStore.metricsInterval * 1000);
    return () => clearInterval(interval);
  });

  async function loadMetrics() {
    const m = await fetchMetrics();
    if (m) {
      metrics = m;
      metricsHistoryStore.push(m);
    }
  }

  // --- Sidecar data fetching ---
  async function loadSidecarData() {
    sidecarLoading = true;
    sidecarError = null;
    try {
      const useDaily = range === "7d" || range === "30d";
      const [statsResult, timelineResult] = await Promise.all([
        fetchStats(range),
        useDaily ? fetchTimeline(range, "1d") : null,
      ]);
      sidecarStats = statsResult;
      timeline = timelineResult;
    } catch (err) {
      sidecarError = err instanceof Error ? err.message : "Failed to load";
      sidecarStats = null;
      timeline = null;
    } finally {
      sidecarLoading = false;
    }
  }

  // Fetch sidecar data when in historical mode or range changes
  $effect(() => {
    if (mode === "historical") {
      range; // track
      loadSidecarData();
    }
  });

  // Force back to live if sidecar becomes unconfigured
  $effect(() => {
    if (!sidecarStore.configured && mode === "historical") {
      mode = "live";
    }
  });

  // --- Live mode derived values ---
  let totalDenylistEntries = $derived(
    metrics?.listEntries
      ? Object.values(metrics.listEntries).reduce((a, b) => a + b, 0)
      : undefined,
  );

  let blockedCount = $derived.by(() => {
    if (!metrics?.responsesByReason) return undefined;
    let total = 0;
    for (const [reason, count] of Object.entries(metrics.responsesByReason)) {
      if (reason.toUpperCase().includes("BLOCKED")) total += count;
    }
    return total;
  });

  let blockedPct = $derived(
    blockedCount !== undefined && metrics?.totalResponses
      ? formatPercentage(blockedCount, metrics.totalResponses)
      : undefined,
  );

  let cacheHitRate = $derived(
    metrics?.cacheHits !== undefined && metrics?.cacheMisses !== undefined
      ? formatPercentage(
          metrics.cacheHits,
          metrics.cacheHits + metrics.cacheMisses,
        )
      : undefined,
  );

  let avgLatency = $derived.by(() => {
    if (!metrics?.requestDurationSum || !metrics?.requestDurationCount)
      return undefined;
    const avgMs =
      (metrics.requestDurationSum / metrics.requestDurationCount) * 1000;
    if (avgMs < 1) return `${(avgMs * 1000).toFixed(0)}µs`;
    if (avgMs < 1000) return `${avgMs.toFixed(1)}ms`;
    return `${(avgMs / 1000).toFixed(2)}s`;
  });

  // --- Historical mode derived values ---
  let showDailyChart = $derived(range === "7d" || range === "30d");

  let dailyData = $derived.by(() => {
    if (!timeline || !showDailyChart) return null;
    return timeline.map((b) => ({
      label: b.timestamp.slice(0, 10),
      total: b.total,
      blocked: b.blocked,
      cached: b.cached,
    }));
  });

  const CHART_TITLES: Record<StatsRange, string> = {
    today: "Hourly Breakdown (Today)",
    yesterday: "Hourly Breakdown (Yesterday)",
    "7d": "Daily Breakdown (7 Days)",
    "30d": "Daily Breakdown (30 Days)",
  };
  let chartTitle = $derived(CHART_TITLES[range]);

  const RETURN_CODE_COLORS: Record<string, string> = {
    NOERROR: "bg-emerald-500/70",
    NXDOMAIN: "bg-amber-500/70",
    SERVFAIL: "bg-red-500/70",
  };
  function returnCodeColor(key: string): string {
    return RETURN_CODE_COLORS[key] ?? "bg-blue-500/70";
  }

  function formatPeriodDate(iso: string): string {
    const d = new Date(iso);
    return d.toLocaleDateString(undefined, { month: "short", day: "numeric" });
  }

  let periodLabel = $derived.by(() => {
    if (!sidecarStats?.period) return null;
    const start = formatPeriodDate(sidecarStats.period.start);
    const end = formatPeriodDate(sidecarStats.period.end);
    return `${start} – ${end}`;
  });

  // --- Shared derived values ---
  let uptimeLabel = $derived.by(() => {
    if (!metrics?.processStartTime) return null;
    const startMs = metrics.processStartTime * 1000;
    const diffMs = Date.now() - startMs;
    const mins = Math.floor(diffMs / 60_000);
    if (mins < 1) return "just now";
    if (mins < 60) return `${mins}m ago`;
    const hours = Math.floor(mins / 60);
    if (hours < 24) return `${hours}h ${mins % 60}m ago`;
    const days = Math.floor(hours / 24);
    return `${days}d ${hours % 24}h ago`;
  });

  // Stats card values that switch on mode
  // In historical mode without data (error/loading), show N/A rather than live values
  let isHistorical = $derived(mode === "historical");

  let statsTotal = $derived(
    isHistorical
      ? sidecarStats
        ? formatNumber(sidecarStats.summary.total_queries)
        : "N/A"
      : formatNumber(metrics?.totalQueries),
  );
  let statsTotalTrend = $derived(
    isHistorical
      ? sidecarStats
        ? `${sidecarStats.summary.unique_domains.toLocaleString()} unique domains`
        : undefined
      : metricsHistoryStore.queriesPerInterval !== undefined
        ? `+${metricsHistoryStore.queriesPerInterval} last interval`
        : undefined,
  );

  let statsBlocked = $derived(
    isHistorical
      ? sidecarStats
        ? formatNumber(sidecarStats.summary.blocked_queries)
        : "N/A"
      : formatNumber(blockedCount),
  );
  let statsBlockedTrend = $derived.by(() => {
    if (isHistorical) {
      if (!sidecarStats) return undefined;
      const pct =
        sidecarStats.summary.total_queries > 0
          ? (
              (sidecarStats.summary.blocked_queries /
                sidecarStats.summary.total_queries) *
              100
            ).toFixed(1)
          : "0";
      return `${pct}% of total`;
    }
    return blockedPct ? `${blockedPct} of total` : undefined;
  });

  let statsCard3Label = $derived(
    isHistorical ? "Cached Queries" : "Cache Hit Rate",
  );
  let statsCard3Value = $derived(
    isHistorical
      ? sidecarStats
        ? formatNumber(sidecarStats.summary.cached_queries)
        : "N/A"
      : (cacheHitRate ?? "N/A"),
  );
  let statsCard3Trend = $derived.by(() => {
    if (isHistorical) {
      if (!sidecarStats) return undefined;
      const pct =
        sidecarStats.summary.total_queries > 0
          ? (
              (sidecarStats.summary.cached_queries /
                sidecarStats.summary.total_queries) *
              100
            ).toFixed(1)
          : "0";
      return `${pct}% cache rate`;
    }
    if (metrics?.cacheHits !== undefined) {
      return `${formatNumber(metrics.cacheHits)} hits / ${formatNumber(metrics.cacheMisses)} misses`;
    }
    return undefined;
  });

  let statsLatency = $derived(
    isHistorical
      ? sidecarStats
        ? `${sidecarStats.summary.avg_duration_ms.toFixed(1)}ms`
        : "N/A"
      : (avgLatency ?? "N/A"),
  );
  let statsLatencyTrend = $derived(
    isHistorical
      ? sidecarStats
        ? `P95: ${sidecarStats.summary.p95_duration_ms}ms`
        : undefined
      : metrics?.requestDurationCount !== undefined
        ? `over ${formatNumber(metrics.requestDurationCount)} requests`
        : undefined,
  );
</script>

{#snippet loadingOverlay()}
  {#if sidecarLoading}
    <div class="absolute inset-0 z-10 rounded-lg bg-surface-primary/60"></div>
  {/if}
{/snippet}

<div class="space-y-6">
  <!-- Blocking toggle -->
  <Card>
    <div
      class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
    >
      <p class="text-sm text-text-muted">Control ad and tracker blocking</p>
      <BlockingToggle />
    </div>
  </Card>

  <!-- Mode toggle bar (only if sidecar configured) -->
  {#if sidecarStore.configured}
    <ModeBar
      {mode}
      onmodechange={(m) => (mode = m)}
      {range}
      onrangechange={(r) => (range = r)}
      {uptimeLabel}
      {periodLabel}
    />
  {/if}

  <!-- Error banner for historical mode -->
  {#if mode === "historical" && sidecarError}
    <div
      class="rounded-xl border border-red-600/30 bg-red-600/10 px-5 py-4 text-sm text-red-400"
    >
      <p class="font-medium">Failed to load analytics</p>
      <p class="mt-1 text-text-muted">{sidecarError}</p>
    </div>
  {/if}

  <!-- Key metrics -->
  <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
    <StatsCard
      label="Total Queries"
      value={statsTotal}
      icon="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z"
      trend={statsTotalTrend}
      color="accent"
    />
    <StatsCard
      label="Blocked"
      value={statsBlocked}
      icon="M18.364 18.364A9 9 0 0 0 5.636 5.636m12.728 12.728A9 9 0 0 1 5.636 5.636m12.728 12.728L5.636 5.636"
      trend={statsBlockedTrend}
      color="rose"
    />
    <StatsCard
      label={statsCard3Label}
      value={statsCard3Value}
      icon="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
      trend={statsCard3Trend}
      color="emerald"
    />
    <StatsCard
      label="Avg Latency"
      value={statsLatency}
      icon="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
      trend={statsLatencyTrend}
      color="amber"
    />
  </div>

  <!-- Chart card -->
  {#if mode === "live"}
    <Card>
      <div class="flex items-baseline justify-between mb-3">
        <h3 class="text-sm font-medium text-text-secondary">Query Activity</h3>
        {#if !sidecarStore.configured && uptimeLabel}
          <span class="text-xs text-text-faint"
            >Blocky started {uptimeLabel}</span
          >
        {/if}
      </div>
      <ActivityChart data={metricsHistoryStore.activity} />
    </Card>
  {:else}
    <Card>
      <h3 class="text-sm font-medium text-text-secondary mb-3">{chartTitle}</h3>
      {#if sidecarLoading && !sidecarStats}
        <div class="h-48 animate-pulse rounded-lg bg-surface-hover"></div>
      {:else if sidecarStats}
        <div class="relative">
          {#if sidecarLoading}
            <div
              class="absolute inset-0 z-10 flex items-center justify-center rounded-lg bg-surface-primary/60"
            >
              <div
                class="h-5 w-5 animate-spin rounded-full border-2 border-accent-600 border-t-transparent"
              ></div>
            </div>
          {/if}
          {#if showDailyChart && dailyData}
            <DailyChart data={dailyData} />
          {:else}
            <HourlyChart data={sidecarStats.hourly} />
          {/if}
        </div>
      {/if}
    </Card>
  {/if}

  <!-- Breakdown row -->
  <div class="grid gap-4 lg:grid-cols-2">
    {#if mode === "live"}
      <Card>
        <h3 class="text-sm font-medium text-text-secondary mb-3">
          Response Breakdown
        </h3>
        <div class="min-h-[200px]">
          <ResponseBreakdown
            responsesByReason={metrics?.responsesByReason}
            total={metrics?.totalResponses}
          />
        </div>
      </Card>
      <Card>
        <h3 class="text-sm font-medium text-text-secondary mb-3">
          Queries by Client
        </h3>
        <div class="min-h-[200px]">
          <ClientTable
            queriesByClient={metrics?.queriesByClient}
            total={metrics?.totalQueries}
          />
        </div>
      </Card>
    {:else}
      <Card>
        <h3 class="text-sm font-medium text-text-secondary mb-3">
          Response Breakdown
        </h3>
        <div class="min-h-[200px]">
          {#if sidecarLoading && !sidecarStats}
            <CardSkeleton lines={5} />
          {:else if sidecarStats}
            <div class="relative">
              {@render loadingOverlay()}
              <ResponseBreakdownChart
                categories={sidecarStats.response_categories}
              />
            </div>
          {/if}
        </div>
      </Card>
      <Card>
        <h3 class="text-sm font-medium text-text-secondary mb-3">
          Queries by Client
        </h3>
        <div class="min-h-[200px]">
          {#if sidecarLoading && !sidecarStats}
            <CardSkeleton lines={5} />
          {:else if sidecarStats}
            <div class="relative">
              {@render loadingOverlay()}
              <ClientBreakdown clients={sidecarStats.clients} />
            </div>
          {/if}
        </div>
      </Card>
    {/if}
  </div>

  <!-- Live-only: secondary stats -->
  {#if mode === "live"}
    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <StatsCard
        label="Denylist Entries"
        value={formatNumber(totalDenylistEntries)}
        icon="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 0 1 0 3.75H5.625a1.875 1.875 0 0 1 0-3.75Z"
        trend={metrics?.listEntries
          ? `${Object.keys(metrics.listEntries).length} groups`
          : undefined}
        color="accent"
      />
      <StatsCard
        label="Cache Entries"
        value={formatNumber(metrics?.cacheEntryCount)}
        icon="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375"
        color="emerald"
      />
      <StatsCard
        label="Prefetch Hits"
        value={formatNumber(metrics?.prefetchHits)}
        icon="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75Z"
        trend={metrics?.prefetches !== undefined
          ? `${formatNumber(metrics.prefetches)} total prefetches`
          : undefined}
        color="amber"
      />
      <StatsCard
        label="Errors"
        value={formatNumber(metrics?.errors ?? 0)}
        icon="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"
        trend={metrics?.failedDownloads !== undefined &&
        metrics.failedDownloads > 0
          ? `${metrics.failedDownloads} failed downloads`
          : undefined}
        color="rose"
      />
    </div>
  {/if}

  <!-- Historical-only: top domains -->
  {#if mode === "historical"}
    <div class="grid gap-4 lg:grid-cols-2">
      <Card>
        <h3 class="text-sm font-medium text-text-secondary mb-3">
          Top Queried Domains
        </h3>
        {#if sidecarLoading && !sidecarStats}
          <CardSkeleton lines={8} />
        {:else if sidecarStats}
          <div class="relative">
            {#if sidecarLoading}<div
                class="absolute inset-0 z-10 rounded-lg bg-surface-primary/60"
              ></div>{/if}
            <TopDomainsTable domains={sidecarStats.top_domains} />
          </div>
        {/if}
      </Card>
      <Card>
        <h3 class="text-sm font-medium text-text-secondary mb-3">
          Top Blocked Domains
        </h3>
        {#if sidecarLoading && !sidecarStats}
          <CardSkeleton lines={8} />
        {:else if sidecarStats}
          <div class="relative">
            {#if sidecarLoading}<div
                class="absolute inset-0 z-10 rounded-lg bg-surface-primary/60"
              ></div>{/if}
            <TopDomainsTable domains={sidecarStats.top_blocked} showReason />
          </div>
        {/if}
      </Card>
    </div>

    <!-- Query types + return codes -->
    <div class="grid gap-4 lg:grid-cols-2">
      <Card>
        <h3 class="text-sm font-medium text-text-secondary mb-3">
          Query Types
        </h3>
        {#if sidecarLoading && !sidecarStats}
          <CardSkeleton lines={5} />
        {:else if sidecarStats}
          <div class="relative">
            {#if sidecarLoading}<div
                class="absolute inset-0 z-10 rounded-lg bg-surface-primary/60"
              ></div>{/if}
            <DistributionChart data={sidecarStats.query_types} />
          </div>
        {/if}
      </Card>
      <Card>
        <h3 class="text-sm font-medium text-text-secondary mb-3">
          Return Codes
        </h3>
        {#if sidecarLoading && !sidecarStats}
          <CardSkeleton lines={3} />
        {:else if sidecarStats}
          <div class="relative">
            {#if sidecarLoading}<div
                class="absolute inset-0 z-10 rounded-lg bg-surface-primary/60"
              ></div>{/if}
            <DistributionChart
              data={sidecarStats.return_codes}
              colorFn={returnCodeColor}
              maxItems={5}
            />
          </div>
        {/if}
      </Card>
    </div>
  {/if}

  <!-- Connection info -->
  <div class="flex justify-center">
    <ConnectionStatus />
  </div>

  {#if !metrics}
    <p class="text-center text-sm text-text-faint">
      Prometheus metrics unavailable — stats will show N/A. Enable metrics in
      your Blocky config for full stats.
    </p>
  {/if}
</div>
