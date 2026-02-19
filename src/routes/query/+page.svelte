<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import QueryForm from "$lib/components/query/QueryForm.svelte";
  import QueryResult from "$lib/components/query/QueryResult.svelte";
  import EmptyState from "$lib/components/ui/EmptyState.svelte";
  import { dnsQuery } from "$lib/api/query";
  import { toastStore } from "$lib/stores/toasts.svelte";
  import type { DnsQueryResponse, QueryHistoryEntry } from "$lib/types/api";
  import { formatRelativeTime } from "$lib/utils/formatters";

  const HISTORY_KEY = "blocky-query-history";
  const MAX_HISTORY = 20;

  let loading = $state(false);
  let currentResult = $state<{
    domain: string;
    type: string;
    response: DnsQueryResponse;
  } | null>(null);
  let history = $state<QueryHistoryEntry[]>(loadHistory());

  function loadHistory(): QueryHistoryEntry[] {
    if (typeof window === "undefined") return [];
    try {
      return JSON.parse(localStorage.getItem(HISTORY_KEY) || "[]");
    } catch {
      return [];
    }
  }

  function saveHistory() {
    localStorage.setItem(HISTORY_KEY, JSON.stringify(history));
  }

  async function handleQuery(domain: string, type: string) {
    loading = true;
    try {
      const response = await dnsQuery(domain, type);
      currentResult = { domain, type, response };

      history = [
        { domain, type, response, timestamp: Date.now() },
        ...history,
      ].slice(0, MAX_HISTORY);
      saveHistory();
    } catch (err) {
      toastStore.error(
        `Query failed: ${err instanceof Error ? err.message : "Unknown error"}`,
      );
    } finally {
      loading = false;
    }
  }

  function clearHistory() {
    history = [];
    localStorage.removeItem(HISTORY_KEY);
  }

  let tick = $state(0);
  $effect(() => {
    const timer = setInterval(() => tick++, 10000);
    return () => clearInterval(timer);
  });
</script>

<div class="space-y-6">
  <Card>
    <QueryForm onsubmit={handleQuery} {loading} />
  </Card>

  {#if currentResult}
    <QueryResult
      result={currentResult.response}
      domain={currentResult.domain}
      type={currentResult.type}
    />
  {/if}

  <!-- History -->
  <div>
    <div class="mb-3 flex items-center justify-between">
      <h3 class="text-sm font-medium text-text-secondary">Query History</h3>
      {#if history.length > 0}
        <button
          onclick={clearHistory}
          class="text-xs text-text-muted hover:text-text-primary cursor-pointer"
        >
          Clear
        </button>
      {/if}
    </div>

    {#if history.length === 0}
      <EmptyState
        title="No queries yet"
        description="Run a DNS query above to see results"
        icon="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582m15.686 0A11.953 11.953 0 0 1 12 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0 1 21 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0 1 12 16.5a17.92 17.92 0 0 1-8.716-2.247m0 0A8.966 8.966 0 0 1 3 12c0-1.264.26-2.467.732-3.558"
      />
    {:else}
      <div class="space-y-2">
        {#each history as entry}
          <button
            onclick={() => handleQuery(entry.domain, entry.type)}
            class="flex w-full items-center gap-3 rounded-lg border border-surface-border bg-surface-primary px-4 py-3 text-left
							transition-colors hover:bg-surface-hover cursor-pointer"
          >
            <span class="font-mono text-sm text-accent-600 dark:text-accent-400"
              >{entry.domain}</span
            >
            <span
              class="rounded bg-surface-secondary px-2 py-0.5 text-xs text-text-secondary"
              >{entry.type}</span
            >
            <span
              class="rounded px-2 py-0.5 text-xs
							{entry.response.returnCode === 'NOERROR'
                ? 'bg-green-500/10 text-green-600 dark:text-green-400'
                : 'bg-red-500/10 text-red-600 dark:text-red-400'}"
            >
              {entry.response.returnCode}
            </span>
            <span class="ml-auto text-xs text-text-faint"
              >{(void tick, formatRelativeTime(entry.timestamp))}</span
            >
          </button>
        {/each}
      </div>
    {/if}
  </div>
</div>
