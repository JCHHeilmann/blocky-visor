<script lang="ts">
    import Card from "$lib/components/ui/Card.svelte";
    import BlockingToggle from "$lib/components/dashboard/BlockingToggle.svelte";
    import StatsCard from "$lib/components/dashboard/StatsCard.svelte";
    import ConnectionStatus from "$lib/components/dashboard/ConnectionStatus.svelte";
    import ActivityChart from "$lib/components/dashboard/ActivityChart.svelte";
    import ResponseBreakdown from "$lib/components/dashboard/ResponseBreakdown.svelte";
    import ClientTable from "$lib/components/dashboard/ClientTable.svelte";
    import { settingsStore } from "$lib/stores/settings.svelte";
    import { metricsHistoryStore } from "$lib/stores/metrics-history.svelte";
    import { fetchMetrics } from "$lib/api/metrics";
    import { formatNumber, formatPercentage } from "$lib/utils/formatters";
    import type { ParsedMetrics } from "$lib/types/api";

    let metrics = $state<ParsedMetrics | null>(null);

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

    let totalDenylistEntries = $derived(
        metrics?.listEntries
            ? Object.values(metrics.listEntries).reduce((a, b) => a + b, 0)
            : undefined,
    );

    let blockedCount = $derived.by(() => {
        if (!metrics?.responsesByReason) return undefined;
        let total = 0;
        for (const [reason, count] of Object.entries(
            metrics.responsesByReason,
        )) {
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
</script>

<div class="space-y-6">
    <!-- Blocking toggle -->
    <Card>
        <div
            class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
        >
            <div>
                <h2 class="text-lg font-semibold text-text-primary mb-1">
                    DNS Blocking
                </h2>
                <p class="text-sm text-text-muted">
                    Control ad and tracker blocking
                </p>
            </div>
            <BlockingToggle />
        </div>
    </Card>

    <!-- Key metrics -->
    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatsCard
            label="Total Queries"
            value={formatNumber(metrics?.totalQueries)}
            icon="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z"
            trend={metricsHistoryStore.queriesPerInterval !== undefined
                ? `+${metricsHistoryStore.queriesPerInterval} last interval`
                : undefined}
        />
        <StatsCard
            label="Blocked"
            value={formatNumber(blockedCount)}
            icon="M18.364 18.364A9 9 0 0 0 5.636 5.636m12.728 12.728A9 9 0 0 1 5.636 5.636m12.728 12.728L5.636 5.636"
            trend={blockedPct ? `${blockedPct} of all responses` : undefined}
        />
        <StatsCard
            label="Cache Hit Rate"
            value={cacheHitRate ?? "N/A"}
            icon="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
            trend={metrics?.cacheHits !== undefined
                ? `${formatNumber(metrics.cacheHits)} hits / ${formatNumber(metrics.cacheMisses)} misses`
                : undefined}
        />
        <StatsCard
            label="Avg Latency"
            value={avgLatency ?? "N/A"}
            icon="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
            trend={metrics?.requestDurationCount !== undefined
                ? `over ${formatNumber(metrics.requestDurationCount)} requests`
                : undefined}
        />
    </div>

    <!-- Activity chart -->
    <Card>
        <h3 class="text-sm font-medium text-text-secondary mb-3">
            Query Activity
        </h3>
        <ActivityChart data={metricsHistoryStore.activity} />
    </Card>

    <!-- Breakdown row -->
    <div class="grid gap-4 lg:grid-cols-2">
        <Card>
            <h3 class="text-sm font-medium text-text-secondary mb-3">
                Response Breakdown
            </h3>
            <ResponseBreakdown
                responsesByReason={metrics?.responsesByReason}
                total={metrics?.totalResponses}
            />
        </Card>
        <Card>
            <h3 class="text-sm font-medium text-text-secondary mb-3">
                Queries by Client
            </h3>
            <ClientTable
                queriesByClient={metrics?.queriesByClient}
                total={metrics?.totalQueries}
            />
        </Card>
    </div>

    <!-- Secondary stats -->
    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatsCard
            label="Denylist Entries"
            value={formatNumber(totalDenylistEntries)}
            icon="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 0 1 0 3.75H5.625a1.875 1.875 0 0 1 0-3.75Z"
            trend={metrics?.listEntries
                ? `${Object.keys(metrics.listEntries).length} groups`
                : undefined}
        />
        <StatsCard
            label="Cache Entries"
            value={formatNumber(metrics?.cacheEntryCount)}
            icon="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375"
        />
        <StatsCard
            label="Prefetch Hits"
            value={formatNumber(metrics?.prefetchHits)}
            icon="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75Z"
            trend={metrics?.prefetches !== undefined
                ? `${formatNumber(metrics.prefetches)} total prefetches`
                : undefined}
        />
        <StatsCard
            label="Errors"
            value={formatNumber(metrics?.errors ?? 0)}
            icon="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"
            trend={metrics?.failedDownloads !== undefined &&
            metrics.failedDownloads > 0
                ? `${metrics.failedDownloads} failed downloads`
                : undefined}
        />
    </div>

    <!-- Connection info -->
    <div class="flex justify-center">
        <ConnectionStatus />
    </div>

    {#if !metrics}
        <p class="text-center text-sm text-text-faint">
            Prometheus metrics unavailable — stats will show N/A. Enable metrics
            in your Blocky config for full stats.
        </p>
    {/if}
</div>
