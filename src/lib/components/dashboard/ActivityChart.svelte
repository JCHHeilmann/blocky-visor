<script lang="ts">
  import type { ActivityPoint } from "$lib/stores/metrics-history.svelte";
  import { tooltipStore } from "$lib/stores/tooltip.svelte";
  import { formatNumber } from "$lib/utils/formatters";

  interface Props {
    data: ActivityPoint[];
  }

  let { data }: Props = $props();

  const HEIGHT = 120;
  const BAR_GAP = 2;

  let chartWidth = $state(600);
  let containerEl = $state<HTMLDivElement>();

  $effect(() => {
    if (!containerEl) return;
    const observer = new ResizeObserver((entries) => {
      chartWidth = entries[0].contentRect.width;
    });
    observer.observe(containerEl);
    return () => observer.disconnect();
  });

  let maxValue = $derived(Math.max(1, ...data.map((d) => d.total)));

  let barWidth = $derived(
    data.length > 0
      ? Math.max(4, (chartWidth - (data.length - 1) * BAR_GAP) / data.length)
      : 20,
  );

  let bars = $derived(
    data.map((point, i) => {
      const x = i * (barWidth + BAR_GAP);
      const totalH = (point.total / maxValue) * (HEIGHT - 4);
      const blockedH = (point.blocked / maxValue) * (HEIGHT - 4);
      const allowedH = totalH - blockedH;
      return {
        x,
        totalH,
        blockedH,
        allowedH,
        blockedY: HEIGHT - blockedH,
        allowedY: HEIGHT - totalH,
        point,
      };
    }),
  );

  function formatTime(ts: number): string {
    return new Date(ts).toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  }

  function showTooltip(e: MouseEvent, bar: (typeof bars)[0]) {
    const allowed = bar.point.total - bar.point.blocked;
    const html = `
			<div class="font-medium text-text-primary">${formatTime(bar.point.timestamp)}</div>
			<div class="mt-1 space-y-0.5 text-xs">
				<div class="flex items-center gap-1.5">
					<span class="inline-block h-2 w-2 rounded-sm bg-accent-500/70"></span>
					<span class="text-text-secondary">Allowed</span>
					<span class="ml-auto font-mono text-text-primary">${formatNumber(allowed)}</span>
				</div>
				<div class="flex items-center gap-1.5">
					<span class="inline-block h-2 w-2 rounded-sm bg-error/70"></span>
					<span class="text-text-secondary">Blocked</span>
					<span class="ml-auto font-mono text-text-primary">${formatNumber(bar.point.blocked)}</span>
				</div>
				<div class="border-t border-surface-border mt-1 pt-1 flex justify-between">
					<span class="text-text-muted">Total</span>
					<span class="font-mono text-text-primary">${formatNumber(bar.point.total)}</span>
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

<div bind:this={containerEl} class="w-full">
  {#if data.length === 0}
    <div
      class="flex h-[120px] items-center justify-center text-sm text-text-muted"
    >
      Collecting data... chart will appear after a few polling intervals.
    </div>
  {:else}
    <svg width={chartWidth} height={HEIGHT} class="w-full">
      <!-- Horizontal grid lines -->
      {#each [0.25, 0.5, 0.75] as pct}
        <line
          x1="0"
          y1={HEIGHT * (1 - pct)}
          x2={chartWidth}
          y2={HEIGHT * (1 - pct)}
          stroke="currentColor"
          class="text-surface-border"
          stroke-dasharray="4 4"
          stroke-width="1"
        />
      {/each}
      <defs>
        {#each bars as bar, i}
          {#if bar.totalH > 0}
            <clipPath id="bar-clip-{i}">
              <rect
                x={bar.x}
                y={bar.allowedY}
                width={barWidth}
                height={bar.totalH}
                rx="3"
              />
            </clipPath>
          {/if}
        {/each}
      </defs>
      {#each bars as bar, i}
        <g>
          {#if bar.totalH > 0}
            <g clip-path="url(#bar-clip-{i})">
              <!-- Allowed portion (teal) -->
              {#if bar.allowedH > 0}
                <rect
                  x={bar.x}
                  y={bar.allowedY}
                  width={barWidth}
                  height={bar.allowedH}
                  class="fill-accent-500/70"
                />
              {/if}
              <!-- Blocked portion (stacked at bottom) -->
              {#if bar.blockedH > 0}
                <rect
                  x={bar.x}
                  y={bar.blockedY}
                  width={barWidth}
                  height={bar.blockedH}
                  class="fill-error/70"
                />
              {/if}
            </g>
          {/if}
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <rect
            x={bar.x}
            y="0"
            width={barWidth}
            height={HEIGHT}
            fill="transparent"
            onmouseenter={(e) => showTooltip(e, bar)}
            onmousemove={moveTooltip}
            onmouseleave={() => tooltipStore.hide()}
          />
        </g>
      {/each}
    </svg>
    <div class="mt-2 flex items-center gap-4 text-xs text-text-muted">
      <span class="flex items-center gap-1.5">
        <span class="inline-block h-2.5 w-2.5 rounded-sm bg-accent-500/70"
        ></span>
        Allowed
      </span>
      <span class="flex items-center gap-1.5">
        <span class="inline-block h-2.5 w-2.5 rounded-sm bg-error/70"></span>
        Blocked
      </span>
    </div>
  {/if}
</div>
