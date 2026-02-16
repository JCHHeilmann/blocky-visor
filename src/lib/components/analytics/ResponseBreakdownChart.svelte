<script lang="ts">
  import { formatNumber, formatPercentage } from "$lib/utils/formatters";
  import { tooltipStore } from "$lib/stores/tooltip.svelte";

  interface Props {
    categories: Record<string, number>;
  }

  let { categories }: Props = $props();

  interface Segment {
    label: string;
    count: number;
    pct: string;
    pctNum: number;
    colorBar: string;
    colorDot: string;
    colorHex: string;
  }

  const REASON_COLORS: Record<
    string,
    { bar: string; dot: string; hex: string }
  > = {
    BLOCKED: {
      bar: "bg-red-500/70 dark:bg-red-500/60",
      dot: "bg-red-500",
      hex: "#ef4444",
    },
    CACHED: {
      bar: "bg-accent-500/70 dark:bg-accent-500/60",
      dot: "bg-accent-500",
      hex: "#14b8a6",
    },
    RESOLVED: {
      bar: "bg-emerald-500/70 dark:bg-emerald-500/60",
      dot: "bg-emerald-500",
      hex: "#10b981",
    },
    CONDITIONAL: {
      bar: "bg-amber-500/70 dark:bg-amber-500/60",
      dot: "bg-amber-500",
      hex: "#f59e0b",
    },
    CUSTOM: {
      bar: "bg-violet-500/70 dark:bg-violet-500/60",
      dot: "bg-violet-500",
      hex: "#8b5cf6",
    },
    SPECIAL: {
      bar: "bg-sky-500/70 dark:bg-sky-500/60",
      dot: "bg-sky-500",
      hex: "#0ea5e9",
    },
  };

  const FALLBACK_COLOR = {
    bar: "bg-gray-500/70 dark:bg-gray-500/60",
    dot: "bg-gray-500",
    hex: "#6b7280",
  };

  function getColor(reason: string): { bar: string; dot: string; hex: string } {
    const upper = reason.toUpperCase();
    for (const [key, color] of Object.entries(REASON_COLORS)) {
      if (upper.includes(key)) return color;
    }
    return FALLBACK_COLOR;
  }

  let total = $derived(
    Object.values(categories).reduce((sum, count) => sum + count, 0),
  );

  let segments = $derived.by(() => {
    if (!categories || total === 0) return [];
    return Object.entries(categories)
      .sort(([, a], [, b]) => b - a)
      .map(([category, count]): Segment => {
        const color = getColor(category);
        return {
          label: category,
          count,
          pct: formatPercentage(count, total),
          pctNum: (count / total) * 100,
          colorBar: color.bar,
          colorDot: color.dot,
          colorHex: color.hex,
        };
      });
  });

  function showSegmentTooltip(e: MouseEvent, seg: Segment) {
    const html = `
			<div class="flex items-center gap-2">
				<span class="inline-block h-2.5 w-2.5 rounded-sm" style="background:${seg.colorHex}"></span>
				<span class="font-medium text-text-primary">${seg.label}</span>
			</div>
			<div class="mt-1 text-xs space-y-0.5">
				<div class="flex justify-between gap-4">
					<span class="text-text-muted">Count</span>
					<span class="font-mono text-text-primary">${formatNumber(seg.count)}</span>
				</div>
				<div class="flex justify-between gap-4">
					<span class="text-text-muted">Share</span>
					<span class="font-mono text-text-primary">${seg.pct}</span>
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

{#if segments.length === 0}
  <p class="text-sm text-text-muted">No response data available.</p>
{:else}
  <!-- Stacked bar -->
  <div
    class="flex h-4 w-full overflow-hidden rounded-full bg-surface-secondary"
  >
    {#each segments as seg}
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class="{seg.colorBar} transition-all cursor-default"
        style="width: {seg.pctNum}%"
        onmouseenter={(e) => showSegmentTooltip(e, seg)}
        onmousemove={moveTooltip}
        onmouseleave={() => tooltipStore.hide()}
      ></div>
    {/each}
  </div>

  <!-- Legend list -->
  <div class="mt-4 space-y-2">
    {#each segments as seg}
      <div class="flex items-center gap-3 text-sm">
        <span class="h-3 w-3 shrink-0 rounded-sm {seg.colorDot}"></span>
        <span class="flex-1 text-text-primary">{seg.label}</span>
        <span class="shrink-0 tabular-nums font-mono text-text-secondary"
          >{formatNumber(seg.count)}</span
        >
        <span
          class="shrink-0 w-14 text-right tabular-nums font-mono text-text-muted"
          >{seg.pct}</span
        >
      </div>
    {/each}
  </div>
{/if}
