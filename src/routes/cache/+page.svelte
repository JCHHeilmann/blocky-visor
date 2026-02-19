<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import StatsCard from "$lib/components/dashboard/StatsCard.svelte";
  import { flushCache } from "$lib/api/cache";
  import { fetchMetrics } from "$lib/api/metrics";
  import { toastStore } from "$lib/stores/toasts.svelte";
  import { settingsStore } from "$lib/stores/settings.svelte";
  import { formatNumber, formatPercentage } from "$lib/utils/formatters";
  import type { ParsedMetrics } from "$lib/types/api";

  let loading = $state(false);
  let showConfirm = $state(false);
  let metrics = $state<ParsedMetrics | null>(null);

  $effect(() => {
    loadMetrics();
    const interval = setInterval(() => {
      if (!document.hidden) loadMetrics();
    }, settingsStore.metricsInterval * 1000);
    return () => clearInterval(interval);
  });

  async function loadMetrics() {
    metrics = await fetchMetrics();
  }

  async function handleFlush() {
    loading = true;
    showConfirm = false;
    try {
      await flushCache();
      await loadMetrics();
      toastStore.success("Cache flushed successfully");
    } catch (err) {
      toastStore.error(
        `Failed to flush cache: ${err instanceof Error ? err.message : "Unknown error"}`,
      );
    } finally {
      loading = false;
    }
  }
</script>

<div class="space-y-6">
  <Card>
    <div
      class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
    >
      <p class="text-sm text-text-muted">Manage the DNS response cache</p>
      <Button variant="danger" onclick={() => (showConfirm = true)} {loading}>
        Flush Cache
      </Button>
    </div>
  </Card>

  <div class="grid gap-4 sm:grid-cols-3">
    <StatsCard
      label="Cache Entries"
      value={formatNumber(metrics?.cacheEntryCount)}
      icon="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375"
    />
    <StatsCard
      label="Cache Hits"
      value={formatNumber(metrics?.cacheHits)}
      icon="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
      trend={metrics?.cacheHits !== undefined &&
      metrics?.cacheMisses !== undefined
        ? `${formatPercentage(metrics.cacheHits, metrics.cacheHits + metrics.cacheMisses)} hit rate`
        : undefined}
    />
    <StatsCard
      label="Cache Misses"
      value={formatNumber(metrics?.cacheMisses)}
      icon="M9.75 9.75l4.5 4.5m0-4.5-4.5 4.5M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
    />
  </div>

  {#if !metrics}
    <Card>
      <p class="text-center text-sm text-text-muted">
        Prometheus metrics unavailable — enable metrics in Blocky for cache
        statistics
      </p>
    </Card>
  {/if}
</div>

<Modal
  open={showConfirm}
  onclose={() => (showConfirm = false)}
  title="Flush Cache"
>
  <p class="text-sm text-text-secondary">
    This will clear all cached DNS responses. Subsequent queries will be
    resolved fresh from upstream servers.
  </p>
  {#snippet actions()}
    <Button variant="ghost" onclick={() => (showConfirm = false)}>Cancel</Button
    >
    <Button variant="danger" onclick={handleFlush} {loading}>Flush</Button>
  {/snippet}
</Modal>
