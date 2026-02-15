<script lang="ts">
	interface Props {
		data: Record<string, number>;
		colorFn?: (key: string) => string;
	}

	let { data, colorFn }: Props = $props();

	let sorted = $derived(
		Object.entries(data).sort(([, a], [, b]) => b - a)
	);

	let total = $derived(sorted.reduce((sum, [, count]) => sum + count, 0));
	let maxCount = $derived(sorted.length > 0 ? sorted[0][1] : 1);

	const defaultColors = [
		'bg-accent-600/70', 'bg-blue-500/70', 'bg-purple-500/70', 'bg-amber-500/70',
		'bg-emerald-500/70', 'bg-rose-500/70', 'bg-cyan-500/70', 'bg-orange-500/70'
	];

	function getColor(key: string, index: number): string {
		if (colorFn) return colorFn(key);
		return defaultColors[index % defaultColors.length];
	}
</script>

{#if sorted.length === 0}
	<p class="text-sm text-text-faint py-4 text-center">No data available</p>
{:else}
	<!-- Stacked overview bar -->
	<div class="flex h-3 rounded-full overflow-hidden mb-4">
		{#each sorted as [key, count], i}
			{@const pct = (count / total) * 100}
			{#if pct > 0.3}
				<div
					class="h-full {getColor(key, i)}"
					style="width: {pct}%"
					title="{key}: {count.toLocaleString()} ({pct.toFixed(1)}%)"
				></div>
			{/if}
		{/each}
	</div>

	<!-- Individual rows -->
	<div class="space-y-2.5">
		{#each sorted as [key, count], i}
			{@const pct = total > 0 ? ((count / total) * 100).toFixed(1) : '0'}
			{@const barPct = (count / maxCount) * 100}
			<div>
				<div class="flex items-center justify-between mb-1">
					<span class="flex items-center gap-2 text-sm">
						<span class="h-2.5 w-2.5 rounded-sm shrink-0 {getColor(key, i)}"></span>
						<span class="font-mono text-text-primary">{key}</span>
					</span>
					<span class="text-sm text-text-muted tabular-nums">
						{count.toLocaleString()}
						<span class="text-text-faint">({pct}%)</span>
					</span>
				</div>
				<div class="h-1.5 rounded-full bg-surface-secondary overflow-hidden">
					<div
						class="h-full rounded-full transition-all {getColor(key, i)}"
						style="width: {barPct}%"
					></div>
				</div>
			</div>
		{/each}
	</div>
{/if}
