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
</script>

<div class="flex flex-wrap items-center justify-between gap-3">
  <div class="flex gap-1">
    <button
      onclick={() => onmodechange("live")}
      class="rounded-lg px-3 py-1.5 text-sm font-medium transition-colors cursor-pointer
				{mode === 'live'
        ? 'bg-accent-600/15 text-accent-600 dark:text-accent-400 border border-accent-500/30'
        : 'bg-surface-secondary text-text-secondary border border-surface-border hover:border-text-muted'}"
      >Live</button
    >
    <button
      onclick={() => onmodechange("historical")}
      class="rounded-lg px-3 py-1.5 text-sm font-medium transition-colors cursor-pointer
				{mode === 'historical'
        ? 'bg-accent-600/15 text-accent-600 dark:text-accent-400 border border-accent-500/30'
        : 'bg-surface-secondary text-text-secondary border border-surface-border hover:border-text-muted'}"
      >Historical</button
    >
  </div>

  <div class="flex items-center gap-3">
    {#if mode === "live" && uptimeLabel}
      <span class="text-xs text-text-faint">Blocky started {uptimeLabel}</span>
    {:else if mode === "historical"}
      <DateRangeSelector value={range} onchange={onrangechange} />
      {#if periodLabel}
        <span class="hidden sm:inline text-xs text-text-faint"
          >{periodLabel}</span
        >
      {/if}
    {/if}
  </div>
</div>
