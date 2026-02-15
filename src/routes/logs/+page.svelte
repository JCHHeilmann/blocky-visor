<script lang="ts">
	import Card from '$lib/components/ui/Card.svelte';
	import DateRangeSelector from '$lib/components/analytics/DateRangeSelector.svelte';
	import LogViewer from '$lib/components/analytics/LogViewer.svelte';
	import { sidecarStore } from '$lib/stores/sidecar.svelte';
	import type { StatsRange } from '$lib/api/sidecar-stats';

	let range = $state<StatsRange>('today');
</script>

<div class="space-y-6">
	<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<h1 class="text-xl font-bold text-text-primary">Query Logs</h1>
		<DateRangeSelector value={range} onchange={(r) => (range = r)} />
	</div>

	{#if !sidecarStore.configured}
		<Card>
			<p class="text-sm text-text-muted">
				Configure a sidecar connection in Settings to view query logs.
			</p>
		</Card>
	{:else}
		<Card>
			<LogViewer {range} />
		</Card>
	{/if}
</div>
