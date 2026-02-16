<script lang="ts">
  import { formatNumber, formatPercentage } from "$lib/utils/formatters";
  import { tooltipStore } from "$lib/stores/tooltip.svelte";

  interface Props {
    data: Record<string, number>;
    colorFn?: (key: string) => string;
    maxItems?: number;
  }

  let { data, colorFn, maxItems }: Props = $props();

  let sorted = $derived(Object.entries(data).sort(([, a], [, b]) => b - a));

  let displaySorted = $derived.by(() => {
    if (!maxItems || sorted.length <= maxItems) return sorted;
    const top = sorted.slice(0, maxItems);
    const restCount = sorted
      .slice(maxItems)
      .reduce((sum, [, count]) => sum + count, 0);
    if (restCount > 0) {
      return [...top, ["Other", restCount] as [string, number]];
    }
    return top;
  });

  let total = $derived(sorted.reduce((sum, [, count]) => sum + count, 0));
  let maxCount = $derived(displaySorted.length > 0 ? displaySorted[0][1] : 1);

  const defaultColors = [
    "bg-accent-600/70",
    "bg-blue-500/70",
    "bg-purple-500/70",
    "bg-amber-500/70",
    "bg-emerald-500/70",
    "bg-rose-500/70",
    "bg-cyan-500/70",
    "bg-orange-500/70",
  ];

  const defaultHexColors = [
    "#0d9488",
    "#3b82f6",
    "#a855f7",
    "#f59e0b",
    "#10b981",
    "#f43f5e",
    "#06b6d4",
    "#f97316",
  ];

  function getColor(key: string, index: number): string {
    if (key === "Other") return "bg-gray-500/70";
    if (colorFn) return colorFn(key);
    return defaultColors[index % defaultColors.length];
  }

  function getHexColor(key: string, index: number): string {
    if (key === "Other") return "#6b7280";
    if (colorFn) {
      const cls = colorFn(key);
      if (cls.includes("emerald")) return "#10b981";
      if (cls.includes("amber")) return "#f59e0b";
      if (cls.includes("red")) return "#ef4444";
      if (cls.includes("blue")) return "#3b82f6";
      return defaultHexColors[index % defaultHexColors.length];
    }
    return defaultHexColors[index % defaultHexColors.length];
  }

  function showBarTooltip(
    e: MouseEvent,
    key: string,
    count: number,
    index: number,
  ) {
    const pct = formatPercentage(count, total);
    const hex = getHexColor(key, index);
    const html = `
			<div class="flex items-center gap-2">
				<span class="inline-block h-2.5 w-2.5 rounded-sm" style="background:${hex}"></span>
				<span class="font-medium text-text-primary">${key}</span>
			</div>
			<div class="mt-1 text-xs space-y-0.5">
				<div class="flex justify-between gap-4">
					<span class="text-text-muted">Count</span>
					<span class="font-mono text-text-primary">${formatNumber(count)}</span>
				</div>
				<div class="flex justify-between gap-4">
					<span class="text-text-muted">Share</span>
					<span class="font-mono text-text-primary">${pct}</span>
				</div>
			</div>`;
    tooltipStore.show(html, e.clientX, e.clientY);
  }

  function moveTooltip(e: MouseEvent) {
    if (tooltipStore.current) {
      tooltipStore.show(tooltipStore.current.html, e.clientX, e.clientY);
    }
  }
</script>

{#if displaySorted.length === 0}
  <p class="text-sm text-text-muted py-4 text-center">No data available</p>
{:else}
  <!-- Stacked overview bar -->
  <div class="flex h-4 rounded-full overflow-hidden bg-surface-secondary mb-4">
    {#each displaySorted as [key, count], i}
      {@const pct = (count / total) * 100}
      {#if pct > 0.3}
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
          class="h-full {getColor(key, i)} cursor-default"
          style="width: {pct}%"
          onmouseenter={(e) => showBarTooltip(e, key, count, i)}
          onmousemove={moveTooltip}
          onmouseleave={() => tooltipStore.hide()}
        ></div>
      {/if}
    {/each}
  </div>

  <!-- Individual rows -->
  <div class="space-y-2.5">
    {#each displaySorted as [key, count], i}
      {@const pct = formatPercentage(count, total)}
      {@const barPct = (count / maxCount) * 100}
      <div>
        <div class="flex items-center justify-between mb-1">
          <span class="flex items-center gap-2 text-sm">
            <span class="h-3 w-3 rounded-sm shrink-0 {getColor(key, i)}"></span>
            <span class="font-mono text-text-primary">{key}</span>
          </span>
          <span class="text-sm tabular-nums">
            <span class="font-mono text-text-secondary"
              >{formatNumber(count)}</span
            >
            <span class="font-mono text-text-muted ml-1">{pct}</span>
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
