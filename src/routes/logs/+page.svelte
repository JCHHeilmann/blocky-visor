<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import DateRangeSelector from "$lib/components/analytics/DateRangeSelector.svelte";
  import LogViewer from "$lib/components/analytics/LogViewer.svelte";
  import { sidecarStore } from "$lib/stores/sidecar.svelte";
  import type { StatsRange } from "$lib/api/sidecar-stats";

  let range = $state<StatsRange>("today");
  let live = $state(false);
</script>

<div class="flex flex-col h-full">
  <div class="shrink-0 flex justify-end items-center gap-1.5 pb-4">
    <div
      class="flex gap-1.5"
      class:opacity-40={live}
      class:pointer-events-none={live}
    >
      <DateRangeSelector value={range} onchange={(r) => (range = r)} />
    </div>

    <span class="mx-1 h-5 w-px bg-surface-border"></span>

    <button
      onclick={() => (live = !live)}
      class="flex items-center gap-2 rounded-lg px-3 py-1.5 text-sm font-medium transition-colors cursor-pointer
        {live
        ? 'bg-emerald-600/15 text-emerald-400 border border-emerald-500/30'
        : 'bg-surface-secondary text-text-secondary border border-surface-border hover:border-text-muted'}"
    >
      {#if live}
        <span class="relative flex h-2 w-2">
          <span
            class="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75"
          ></span>
          <span class="relative inline-flex h-2 w-2 rounded-full bg-emerald-400"
          ></span>
        </span>
      {:else}
        <span class="inline-flex h-2 w-2 rounded-full bg-text-muted"></span>
      {/if}
      Live
    </button>
  </div>

  {#if !sidecarStore.configured}
    <Card>
      <p class="text-sm text-text-muted">
        Configure a sidecar connection in Settings to view query logs.
      </p>
    </Card>
  {:else}
    <Card class="flex-1 min-h-0 flex flex-col overflow-hidden !p-4">
      <LogViewer {range} {live} />
    </Card>
  {/if}
</div>
