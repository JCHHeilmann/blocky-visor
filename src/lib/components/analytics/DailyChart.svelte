<script lang="ts">
  import { formatNumber } from "$lib/utils/formatters";
  import { tooltipStore } from "$lib/stores/tooltip.svelte";

  interface DailyData {
    label: string; // "YYYY-MM-DD"
    total: number;
    blocked: number;
    cached: number;
  }

  interface Props {
    data: DailyData[];
  }

  let { data }: Props = $props();

  let maxVal = $derived(Math.max(...data.map((d) => d.total), 1));

  function formatDayLabel(iso: string): string {
    const date = new Date(iso + "T00:00:00");
    return date.toLocaleDateString(undefined, {
      weekday: "short",
      month: "short",
      day: "numeric",
    });
  }

  function formatShortLabel(iso: string): string {
    const date = new Date(iso + "T00:00:00");
    return date.toLocaleDateString(undefined, {
      month: "short",
      day: "numeric",
    });
  }

  function showTooltip(e: MouseEvent, bucket: DailyData) {
    const resolved = bucket.total - bucket.blocked - bucket.cached;
    const dayLabel = formatDayLabel(bucket.label);
    const html = `
			<div class="flex items-center gap-2">
				<span class="font-medium text-text-primary">${dayLabel}</span>
			</div>
			<div class="mt-1 text-xs space-y-0.5">
				<div class="flex justify-between gap-4">
					<span class="text-text-muted">Total</span>
					<span class="font-mono text-text-primary">${formatNumber(bucket.total)}</span>
				</div>
				<div class="flex justify-between gap-4">
					<span class="flex items-center gap-1.5">
						<span class="inline-block h-2.5 w-2.5 rounded-sm" style="background:#ef4444"></span>
						<span class="text-text-muted">Blocked</span>
					</span>
					<span class="font-mono text-text-primary">${formatNumber(bucket.blocked)}</span>
				</div>
				<div class="flex justify-between gap-4">
					<span class="flex items-center gap-1.5">
						<span class="inline-block h-2.5 w-2.5 rounded-sm" style="background:#3b82f6"></span>
						<span class="text-text-muted">Cached</span>
					</span>
					<span class="font-mono text-text-primary">${formatNumber(bucket.cached)}</span>
				</div>
				<div class="flex justify-between gap-4">
					<span class="flex items-center gap-1.5">
						<span class="inline-block h-2.5 w-2.5 rounded-sm" style="background:#0d9488"></span>
						<span class="text-text-muted">Resolved</span>
					</span>
					<span class="font-mono text-text-primary">${formatNumber(resolved)}</span>
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

<div class="flex items-end gap-1.5 h-48">
  {#each data as bucket}
    {@const totalPct = (bucket.total / maxVal) * 100}
    {@const blockedPct = (bucket.blocked / maxVal) * 100}
    {@const cachedPct = (bucket.cached / maxVal) * 100}
    {@const resolvedPct = totalPct - blockedPct - cachedPct}
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      class="relative flex-1 flex flex-col justify-end h-full cursor-default"
      onmouseenter={(e) => showTooltip(e, bucket)}
      onmousemove={moveTooltip}
      onmouseleave={() => tooltipStore.hide()}
    >
      <div class="flex flex-col justify-end h-full gap-px">
        {#if resolvedPct > 0}
          <div
            class="w-full rounded-t bg-accent-600/60"
            style="height: {resolvedPct}%"
          ></div>
        {/if}
        {#if cachedPct > 0}
          <div class="w-full bg-blue-500/60" style="height: {cachedPct}%"></div>
        {/if}
        {#if blockedPct > 0}
          <div
            class="w-full rounded-b bg-red-500/60"
            style="height: {blockedPct}%"
          ></div>
        {/if}
      </div>
    </div>
  {/each}
</div>
<!-- X-axis labels -->
<div class="flex mt-1">
  {#each data as bucket}
    <div class="flex-1 text-center text-[10px] text-text-faint">
      {formatShortLabel(bucket.label)}
    </div>
  {/each}
</div>
<!-- Legend -->
<div class="flex gap-4 mt-3 text-xs text-text-muted">
  <span class="flex items-center gap-1.5"
    ><span class="h-2.5 w-2.5 rounded-sm bg-accent-600/60"></span> Resolved</span
  >
  <span class="flex items-center gap-1.5"
    ><span class="h-2.5 w-2.5 rounded-sm bg-blue-500/60"></span> Cached</span
  >
  <span class="flex items-center gap-1.5"
    ><span class="h-2.5 w-2.5 rounded-sm bg-red-500/60"></span> Blocked</span
  >
</div>
