<script lang="ts">
	interface HourlyData {
		hour: number;
		total: number;
		blocked: number;
		cached: number;
	}

	interface Props {
		data: HourlyData[];
	}

	let { data }: Props = $props();

	let maxVal = $derived(Math.max(...data.map((d) => d.total), 1));
</script>

<div class="flex items-end gap-1 h-48">
	{#each data as bucket}
		{@const totalPct = (bucket.total / maxVal) * 100}
		{@const blockedPct = (bucket.blocked / maxVal) * 100}
		{@const cachedPct = (bucket.cached / maxVal) * 100}
		{@const resolvedPct = totalPct - blockedPct - cachedPct}
		<div class="group relative flex-1 flex flex-col justify-end h-full" title="Hour {bucket.hour}: {bucket.total} queries">
			<!-- Stacked bar -->
			<div class="flex flex-col justify-end h-full gap-px">
				{#if resolvedPct > 0}
					<div
						class="w-full rounded-t-sm bg-accent-600/60"
						style="height: {resolvedPct}%"
					></div>
				{/if}
				{#if cachedPct > 0}
					<div
						class="w-full bg-blue-500/60"
						style="height: {cachedPct}%"
					></div>
				{/if}
				{#if blockedPct > 0}
					<div
						class="w-full rounded-b-sm bg-red-500/60"
						style="height: {blockedPct}%"
					></div>
				{/if}
			</div>
			<!-- Tooltip -->
			<div class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 hidden group-hover:block z-10
				whitespace-nowrap rounded-lg border border-surface-border bg-surface-primary px-3 py-2 text-xs shadow-lg">
				<p class="font-semibold text-text-primary">{bucket.hour}:00</p>
				<p class="text-text-muted">{bucket.total} total</p>
				<p class="text-red-400">{bucket.blocked} blocked</p>
				<p class="text-blue-400">{bucket.cached} cached</p>
			</div>
		</div>
	{/each}
</div>
<!-- X-axis labels -->
<div class="flex mt-1">
	{#each data as bucket}
		<div class="flex-1 text-center text-[10px] text-text-faint">
			{#if bucket.hour % 3 === 0}
				{bucket.hour}
			{/if}
		</div>
	{/each}
</div>
<!-- Legend -->
<div class="flex gap-4 mt-3 text-xs text-text-muted">
	<span class="flex items-center gap-1.5"><span class="h-2.5 w-2.5 rounded-sm bg-accent-600/60"></span> Resolved</span>
	<span class="flex items-center gap-1.5"><span class="h-2.5 w-2.5 rounded-sm bg-blue-500/60"></span> Cached</span>
	<span class="flex items-center gap-1.5"><span class="h-2.5 w-2.5 rounded-sm bg-red-500/60"></span> Blocked</span>
</div>
