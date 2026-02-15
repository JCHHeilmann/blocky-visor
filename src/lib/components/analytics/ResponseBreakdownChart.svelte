<script lang="ts">
	interface Props {
		categories: Record<string, number>;
	}

	let { categories }: Props = $props();

	let sorted = $derived(
		Object.entries(categories)
			.sort(([, a], [, b]) => b - a)
	);

	let total = $derived(sorted.reduce((sum, [, count]) => sum + count, 0));
</script>

{#if sorted.length === 0}
	<p class="text-sm text-text-faint py-4 text-center">No data available</p>
{:else}
	<!-- Stacked horizontal bar -->
	<div class="flex h-4 rounded-full overflow-hidden mb-4">
		{#each sorted as [category, count]}
			{@const pct = (count / total) * 100}
			{@const color = category.includes('BLOCKED') ? 'bg-red-500/70' : category.includes('CACHED') ? 'bg-blue-500/60' : 'bg-accent-600/60'}
			{#if pct > 0.5}
				<div
					class="h-full {color}"
					style="width: {pct}%"
					title="{category}: {count.toLocaleString()} ({pct.toFixed(1)}%)"
				></div>
			{/if}
		{/each}
	</div>
	<!-- Legend list -->
	<div class="space-y-1.5">
		{#each sorted as [category, count]}
			{@const pct = total > 0 ? ((count / total) * 100).toFixed(1) : '0'}
			{@const color = category.includes('BLOCKED') ? 'bg-red-500/70' : category.includes('CACHED') ? 'bg-blue-500/60' : 'bg-accent-600/60'}
			<div class="flex items-center justify-between text-sm">
				<span class="flex items-center gap-2">
					<span class="h-2.5 w-2.5 rounded-sm {color}"></span>
					<span class="text-text-secondary">{category}</span>
				</span>
				<span class="text-text-muted tabular-nums">{count.toLocaleString()} <span class="text-text-faint">({pct}%)</span></span>
			</div>
		{/each}
	</div>
{/if}
