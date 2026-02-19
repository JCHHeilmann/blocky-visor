<script lang="ts">
  import DateRangeSelector from "$lib/components/analytics/DateRangeSelector.svelte";
  import type { StatsRange } from "$lib/api/sidecar-stats";

  interface Props {
    mode: "live" | "historical";
    onmodechange: (mode: "live" | "historical") => void;
    range: StatsRange;
    onrangechange: (range: StatsRange) => void;
    uptimeLabel: string | null;
    periodLabel: string | null;
  }

  let {
    mode,
    onmodechange,
    range,
    onrangechange,
    uptimeLabel,
    periodLabel,
  }: Props = $props();

  const BASE_BTN =
    "rounded-lg px-3 py-1.5 text-sm font-medium transition-colors cursor-pointer border";
  const ACTIVE_BTN =
    "bg-accent-600/15 text-accent-600 dark:text-accent-400 border-accent-500/30";
  const INACTIVE_BTN =
    "bg-surface-secondary text-text-secondary border-surface-border hover:border-text-muted";

  function buttonClass(active: boolean): string {
    return `${BASE_BTN} ${active ? ACTIVE_BTN : INACTIVE_BTN}`;
  }
</script>

<div class="flex flex-wrap items-center justify-between gap-3">
  <div class="flex items-center gap-3">
    <div class="flex gap-1">
      <button
        onclick={() => onmodechange("live")}
        class={buttonClass(mode === "live")}>Live</button
      >
      <button
        onclick={() => onmodechange("historical")}
        class={buttonClass(mode === "historical")}>Historical</button
      >
    </div>

    {#if mode === "live" && uptimeLabel}
      <span class="text-xs text-text-faint"
        >Showing data since Blocky started {uptimeLabel}</span
      >
    {:else if mode === "historical" && periodLabel}
      <span class="hidden sm:inline text-xs text-text-faint"
        >Showing data for {periodLabel}</span
      >
    {/if}
  </div>

  {#if mode === "historical"}
    <DateRangeSelector value={range} onchange={onrangechange} />
  {/if}
</div>
