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
  <div class="shrink-0 flex justify-end pb-4">
    {#if live}
      <div class="flex gap-1.5">
        <span
          class="rounded-lg px-3 py-1.5 text-sm font-medium bg-surface-secondary text-text-muted border border-surface-border cursor-not-allowed opacity-50"
        >
          Live mode active
        </span>
      </div>
    {:else}
      <DateRangeSelector value={range} onchange={(r) => (range = r)} />
    {/if}
  </div>

  {#if !sidecarStore.configured}
    <Card>
      <p class="text-sm text-text-muted">
        Configure a sidecar connection in Settings to view query logs.
      </p>
    </Card>
  {:else}
    <Card class="flex-1 min-h-0 flex flex-col overflow-hidden !p-4">
      <LogViewer {range} {live} onlivechange={(v) => (live = v)} />
    </Card>
  {/if}
</div>
