<script lang="ts">
	import { tooltipStore } from '$lib/stores/tooltip.svelte';

	let el = $state<HTMLDivElement>();

	// Offset from cursor
	const OFFSET_X = 12;
	const OFFSET_Y = -8;

	let style = $derived.by(() => {
		const t = tooltipStore.current;
		if (!t || !el) return '';

		// Position tooltip, flip if near viewport edge
		let x = t.x + OFFSET_X;
		let y = t.y + OFFSET_Y;

		const rect = el.getBoundingClientRect();
		const vw = window.innerWidth;
		const vh = window.innerHeight;

		if (x + rect.width > vw - 8) x = t.x - rect.width - OFFSET_X;
		if (y + rect.height > vh - 8) y = t.y - rect.height - OFFSET_X;
		if (y < 8) y = 8;

		return `left:${x}px;top:${y}px`;
	});
</script>

{#if tooltipStore.current}
	<div
		bind:this={el}
		class="fixed z-[200] pointer-events-none rounded-lg border border-surface-border bg-surface-primary px-3 py-2 text-sm shadow-lg"
		style={style}
	>
		{@html tooltipStore.current.html}
	</div>
{/if}
