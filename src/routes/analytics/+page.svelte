<script lang="ts">
	import Card from '$lib/components/ui/Card.svelte';
	import StatsCard from '$lib/components/dashboard/StatsCard.svelte';
	import DateRangeSelector from '$lib/components/analytics/DateRangeSelector.svelte';
	import HourlyChart from '$lib/components/analytics/HourlyChart.svelte';
	import TopDomainsTable from '$lib/components/analytics/TopDomainsTable.svelte';
	import ResponseBreakdownChart from '$lib/components/analytics/ResponseBreakdownChart.svelte';
	import ClientBreakdown from '$lib/components/analytics/ClientBreakdown.svelte';
	import DistributionChart from '$lib/components/analytics/DistributionChart.svelte';
	import CardSkeleton from '$lib/components/analytics/CardSkeleton.svelte';
	import { sidecarStore } from '$lib/stores/sidecar.svelte';
	import { fetchStats, type StatsRange } from '$lib/api/sidecar-stats';
	import { formatNumber, formatPercentage } from '$lib/utils/formatters';
	import type { SidecarStatsResponse } from '$lib/types/api';

	let range = $state<StatsRange>('today');
	let stats = $state<SidecarStatsResponse | null>(null);
	let loading = $state(false);
	let error = $state<string | null>(null);

	async function loadStats() {
		loading = true;
		error = null;
		try {
			stats = await fetchStats(range);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load stats';
			stats = null;
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		range; // track
		loadStats();
	});

	let blockedPct = $derived(
		stats ? formatPercentage(stats.summary.blocked_queries, stats.summary.total_queries) : undefined
	);

	let cachedPct = $derived(
		stats ? formatPercentage(stats.summary.cached_queries, stats.summary.total_queries) : undefined
	);

	function returnCodeColor(key: string): string {
		if (key === 'NOERROR') return 'bg-emerald-500/70';
		if (key === 'NXDOMAIN') return 'bg-amber-500/70';
		if (key === 'SERVFAIL') return 'bg-red-500/70';
		return 'bg-blue-500/70';
	}
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<h1 class="text-xl font-bold text-text-primary">Query Analytics</h1>
		<DateRangeSelector value={range} onchange={(r) => (range = r)} />
	</div>

	{#if !sidecarStore.configured}
		<Card>
			<p class="text-sm text-text-muted">
				Configure a sidecar connection in Settings to view query log analytics.
			</p>
		</Card>
	{:else if error}
		<div class="rounded-xl border border-red-600/30 bg-red-600/10 px-5 py-4 text-sm text-red-400">
			<p class="font-medium">Failed to load analytics</p>
			<p class="mt-1 text-text-muted">{error}</p>
		</div>
	{:else}
		<!-- Summary cards -->
		<div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
			<StatsCard
				label="Total Queries"
				value={loading && !stats ? '...' : formatNumber(stats?.summary.total_queries)}
				icon="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z"
				trend={stats ? `${stats.summary.unique_domains.toLocaleString()} unique domains` : undefined}
			/>
			<StatsCard
				label="Blocked"
				value={loading && !stats ? '...' : formatNumber(stats?.summary.blocked_queries)}
				icon="M18.364 18.364A9 9 0 0 0 5.636 5.636m12.728 12.728A9 9 0 0 1 5.636 5.636m12.728 12.728L5.636 5.636"
				trend={blockedPct ? `${blockedPct} of total` : undefined}
			/>
			<StatsCard
				label="Cached"
				value={loading && !stats ? '...' : formatNumber(stats?.summary.cached_queries)}
				icon="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
				trend={cachedPct ? `${cachedPct} cache rate` : undefined}
			/>
			<StatsCard
				label="Avg Latency"
				value={loading && !stats ? '...' : stats ? `${stats.summary.avg_duration_ms}ms` : 'N/A'}
				icon="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
				trend={stats ? `P95: ${stats.summary.p95_duration_ms}ms` : undefined}
			/>
		</div>

		<!-- Hourly chart -->
		<Card>
			<h3 class="text-sm font-medium text-text-secondary mb-3">Hourly Breakdown</h3>
			{#if loading && !stats}
				<div class="h-48 animate-pulse rounded-lg bg-surface-hover"></div>
			{:else if stats}
				<HourlyChart data={stats.hourly} />
			{/if}
		</Card>

		<!-- Top domains + top blocked -->
		<div class="grid gap-4 lg:grid-cols-2">
			<Card>
				<h3 class="text-sm font-medium text-text-secondary mb-3">Top Queried Domains</h3>
				{#if loading && !stats}
					<CardSkeleton lines={8} />
				{:else if stats}
					<TopDomainsTable domains={stats.top_domains} />
				{/if}
			</Card>
			<Card>
				<h3 class="text-sm font-medium text-text-secondary mb-3">Top Blocked Domains</h3>
				{#if loading && !stats}
					<CardSkeleton lines={8} />
				{:else if stats}
					<TopDomainsTable domains={stats.top_blocked} showReason />
				{/if}
			</Card>
		</div>

		<!-- Response breakdown + clients -->
		<div class="grid gap-4 lg:grid-cols-2">
			<Card>
				<h3 class="text-sm font-medium text-text-secondary mb-3">Response Categories</h3>
				{#if loading && !stats}
					<CardSkeleton lines={4} />
				{:else if stats}
					<ResponseBreakdownChart categories={stats.response_categories} />
				{/if}
			</Card>
			<Card>
				<h3 class="text-sm font-medium text-text-secondary mb-3">Clients</h3>
				{#if loading && !stats}
					<CardSkeleton lines={6} />
				{:else if stats}
					<ClientBreakdown clients={stats.clients} />
				{/if}
			</Card>
		</div>

		<!-- Query types + return codes -->
		<div class="grid gap-4 lg:grid-cols-2">
			<Card>
				<h3 class="text-sm font-medium text-text-secondary mb-3">Query Types</h3>
				{#if loading && !stats}
					<CardSkeleton lines={5} />
				{:else if stats}
					<DistributionChart data={stats.query_types} />
				{/if}
			</Card>
			<Card>
				<h3 class="text-sm font-medium text-text-secondary mb-3">Return Codes</h3>
				{#if loading && !stats}
					<CardSkeleton lines={3} />
				{:else if stats}
					<DistributionChart data={stats.return_codes} colorFn={returnCodeColor} />
				{/if}
			</Card>
		</div>
	{/if}
</div>
