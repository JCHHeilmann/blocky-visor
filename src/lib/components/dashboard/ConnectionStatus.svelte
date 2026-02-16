<script lang="ts">
  import { blockingStore } from "$lib/stores/blocking.svelte";
  import { settingsStore } from "$lib/stores/settings.svelte";
  import { formatRelativeTime } from "$lib/utils/formatters";

  // Update the relative time display periodically
  let now = $state(Date.now());
  $effect(() => {
    const timer = setInterval(() => (now = Date.now()), 1000);
    return () => clearInterval(timer);
  });

  let displayTime = $derived(
    blockingStore.lastUpdated
      ? formatRelativeTime(blockingStore.lastUpdated)
      : "Never",
  );
</script>

<div class="flex items-center gap-4 text-sm text-text-muted">
  <span>
    API: <span class="font-mono text-text-secondary"
      >{settingsStore.apiUrl}</span
    >
  </span>
  <span class="text-surface-border">|</span>
  <span>Last updated: {displayTime}</span>
  <span class="text-surface-border">|</span>
  <span>Polling: every {settingsStore.refreshInterval}s</span>
</div>
