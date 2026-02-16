<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import { refreshLists } from "$lib/api/lists";
  import { fetchMetrics } from "$lib/api/metrics";
  import { toastStore } from "$lib/stores/toasts.svelte";
  import { settingsStore } from "$lib/stores/settings.svelte";
  import { formatNumber, formatDate } from "$lib/utils/formatters";
  import type { ParsedMetrics } from "$lib/types/api";

  let loading = $state(false);
  let showConfirm = $state(false);
  let lastRefresh = $state<number | null>(null);
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

  async function handleRefresh() {
    loading = true;
    showConfirm = false;
    try {
      await refreshLists();
      lastRefresh = Date.now();
      await loadMetrics();
      toastStore.success("Lists refreshed successfully");
    } catch (err) {
      toastStore.error(
        `Failed to refresh lists: ${err instanceof Error ? err.message : "Unknown error"}`,
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
      <div>
        <h2 class="text-lg font-semibold text-text-primary">Blocking Lists</h2>
        <p class="text-sm text-text-muted">
          Refresh lists from configured sources
        </p>
        {#if lastRefresh}
          <p class="mt-1 text-xs text-text-faint">
            Last refreshed: {formatDate(lastRefresh)}
          </p>
        {/if}
      </div>
      <Button onclick={() => (showConfirm = true)} {loading}>
        Refresh Lists
      </Button>
    </div>
  </Card>

  <!-- List entries from metrics -->
  {#if metrics?.listEntries && Object.keys(metrics.listEntries).length > 0}
    <Card>
      <h3 class="mb-4 text-sm font-medium text-text-secondary">
        List Entry Counts
      </h3>
      <div class="space-y-3">
        {#each Object.entries(metrics.listEntries) as [name, count]}
          <div class="flex items-center justify-between">
            <span class="truncate font-mono text-sm text-text-secondary"
              >{name}</span
            >
            <span class="ml-4 shrink-0 text-sm font-medium text-text-primary"
              >{formatNumber(count)}</span
            >
          </div>
        {/each}
      </div>
    </Card>
  {:else}
    <Card>
      <p class="text-center text-sm text-text-muted">
        {metrics
          ? "No list entries found in metrics"
          : "Prometheus metrics unavailable — enable metrics in Blocky for entry counts"}
      </p>
    </Card>
  {/if}
</div>

<Modal
  open={showConfirm}
  onclose={() => (showConfirm = false)}
  title="Refresh Lists"
>
  <p class="text-sm text-text-secondary">
    This will re-download and reload all configured blocking lists. This may
    take a moment.
  </p>
  {#snippet actions()}
    <Button variant="ghost" onclick={() => (showConfirm = false)}>Cancel</Button
    >
    <Button onclick={handleRefresh} {loading}>Refresh</Button>
  {/snippet}
</Modal>
