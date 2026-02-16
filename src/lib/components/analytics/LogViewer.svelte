<script lang="ts">
  import { untrack } from "svelte";
  import type { SidecarLogEntry } from "$lib/types/api";
  import type { StatsRange } from "$lib/api/sidecar-stats";
  import { fetchLogs, buildLogStreamUrl } from "$lib/api/sidecar-stats";
  import { dnsQuery } from "$lib/api/query";

  interface Props {
    range: StatsRange;
    live: boolean;
  }

  let { range, live }: Props = $props();

  type KeyedEntry = SidecarLogEntry & { _id: number };

  let entries = $state<KeyedEntry[]>([]);
  let total = $state(0);
  let loading = $state(false);
  let error = $state<string | null>(null);
  let offset = $state(0);
  const limit = 50;
  let nextId = 0;

  function keyEntries(raw: SidecarLogEntry[]): KeyedEntry[] {
    return raw.map((e) => ({ ...e, _id: nextId++ }));
  }

  let filterDomain = $state("");
  let filterClient = $state("");
  let filterType = $state("");

  // Infinite scroll
  let sentinel: HTMLDivElement | undefined = $state();
  let scrollContainer: HTMLDivElement | undefined = $state();
  let observer: IntersectionObserver | undefined;

  // Live streaming
  let eventSource: EventSource | null = $state(null);
  let sseConnected = $state(false);
  const LIVE_CAP = 200;

  // Client hostname resolution
  let hostnames = $state<Record<string, string>>({});

  function ipToArpa(ip: string): string | null {
    const parts = ip.split(".");
    if (parts.length !== 4) return null;
    return parts.reverse().join(".") + ".in-addr.arpa";
  }

  async function resolveHostname(ip: string): Promise<string | null> {
    const arpa = ipToArpa(ip);
    if (!arpa) return null;
    try {
      const result = await dnsQuery(arpa, "PTR");
      if (
        result.returnCode === "NOERROR" &&
        result.response &&
        result.response !== ""
      ) {
        const match = result.response.match(/PTR\s*\(([^)]+)\)/);
        if (match) return match[1].replace(/\.$/, "");
        return result.response.replace(/\.$/, "").trim();
      }
    } catch {}
    return null;
  }

  function resolveNewClients(ips: string[]) {
    for (const ip of ips) {
      if (ip in hostnames) continue;
      hostnames[ip] = "";
      resolveHostname(ip).then((name) => {
        if (name) {
          hostnames = { ...hostnames, [ip]: name };
        }
      });
    }
  }

  function clientDisplay(entry: SidecarLogEntry): string {
    const resolved = hostnames[entry.client_ip];
    if (resolved) return resolved;
    if (entry.client_name && entry.client_name !== entry.client_ip)
      return entry.client_name;
    return entry.client_ip;
  }

  function stripDot(domain: string): string {
    return domain.endsWith(".") ? domain.slice(0, -1) : domain;
  }

  // Resolve a client filter to a value the server can match.
  // If the filter matches a PTR-resolved hostname, return the IP instead.
  function resolveClientFilter(): string | undefined {
    const f = filterClient.trim();
    if (!f) return undefined;
    // Already looks like an IP — send directly
    if (/^\d/.test(f)) return f;
    // Check resolved hostnames for a match
    const lf = f.toLowerCase();
    for (const [ip, hostname] of Object.entries(hostnames)) {
      if (hostname && hostname.toLowerCase().includes(lf)) {
        return ip;
      }
    }
    // Return as-is — might match ClientName on server
    return f;
  }

  let loadGen = 0;

  async function load(reset = false) {
    if (reset) {
      offset = 0;
      entries = [];
    }
    loading = true;
    error = null;
    const gen = ++loadGen;
    try {
      const currentOffset = offset;
      const result = await fetchLogs({
        range,
        limit,
        offset: currentOffset,
        domain: filterDomain || undefined,
        client: resolveClientFilter(),
        type: filterType || undefined,
      });
      if (gen !== loadGen) return; // Stale response, discard
      const keyed = keyEntries(result.entries);
      resolveNewClients(keyed.map((e) => e.client_ip));
      if (currentOffset === 0) {
        entries = keyed;
      } else {
        entries = [...entries, ...keyed];
      }
      total = result.total;
    } catch (err) {
      if (gen !== loadGen) return;
      error = err instanceof Error ? err.message : "Failed to load logs";
    } finally {
      if (gen === loadGen) loading = false;
    }
  }

  function loadNextPage() {
    if (loading || entries.length >= total) return;
    offset = entries.length;
    load();
  }

  // Reload on range changes only (untrack load internals)
  $effect(() => {
    range;
    if (!live) {
      untrack(() => load(true));
    }
  });

  // Set up IntersectionObserver for infinite scroll
  $effect(() => {
    if (live || !sentinel) return;

    observer = new IntersectionObserver(
      (observerEntries) => {
        if (observerEntries[0]?.isIntersecting) {
          untrack(() => loadNextPage());
        }
      },
      { root: scrollContainer, threshold: 0.1 },
    );
    observer.observe(sentinel);

    return () => {
      observer?.disconnect();
    };
  });

  function connectSse() {
    closeSse();
    entries = [];
    total = 0;
    error = null;

    const url = buildLogStreamUrl({
      client: resolveClientFilter(),
      domain: filterDomain || undefined,
      type: filterType || undefined,
    });

    const es = new EventSource(url);
    eventSource = es;

    es.onopen = () => {
      sseConnected = true;
    };

    es.addEventListener("backfill", (event: MessageEvent) => {
      try {
        const backfillEntries = JSON.parse(event.data) as SidecarLogEntry[];
        const keyed = keyEntries([...backfillEntries].reverse().slice(0, LIVE_CAP));
        resolveNewClients(keyed.map((e) => e.client_ip));
        entries = keyed;
        total = entries.length;
      } catch {}
    });

    es.onmessage = (event) => {
      try {
        const entry = JSON.parse(event.data) as SidecarLogEntry;
        const [keyed] = keyEntries([entry]);
        resolveNewClients([keyed.client_ip]);
        entries = [keyed, ...entries.slice(0, LIVE_CAP - 1)];
        total = entries.length;
      } catch {}
    };

    es.onerror = () => {
      sseConnected = false;
    };
  }

  // Live mode SSE
  $effect(() => {
    if (!live) {
      closeSse();
      return;
    }

    untrack(() => connectSse());

    return () => {
      closeSse();
    };
  });

  function closeSse() {
    if (eventSource) {
      eventSource.close();
      eventSource = null;
      sseConnected = false;
    }
  }

  function applyFilters() {
    if (live) {
      connectSse();
    } else {
      load(true);
    }
  }

  let hasActiveFilters = $derived(
    filterDomain !== "" || filterClient !== "" || filterType !== "",
  );

  function clearFilters() {
    filterDomain = "";
    filterClient = "";
    filterType = "";
    applyFilters();
  }

  function handleFilterKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      e.preventDefault();
      applyFilters();
    }
  }

  // Auto-apply text filters after typing stops
  let debounceTimer: ReturnType<typeof setTimeout> | undefined;
  let filtersInitialized = false;
  $effect(() => {
    // Track the text filter values
    filterDomain;
    filterClient;
    if (!filtersInitialized) {
      filtersInitialized = true;
      return;
    }
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      untrack(() => applyFilters());
    }, 400);
    return () => clearTimeout(debounceTimer);
  });

  const typeOptions = [
    { value: "", label: "All" },
    { value: "blocked", label: "Blocked" },
    { value: "cached", label: "Cached" },
    { value: "resolved", label: "Resolved" },
  ];

  function formatTime(ts: string): string {
    const d = new Date(ts);
    return d.toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  }
</script>

<div class="flex flex-col h-full min-h-0">
  <!-- Filters row (pinned) -->
  <div class="shrink-0 flex items-center gap-2 pb-3">
    <input
      bind:value={filterDomain}
      onkeydown={handleFilterKeydown}
      placeholder="Domain"
      class="w-36 rounded-md border border-surface-border bg-surface-secondary px-2.5 py-1.5 text-xs text-text-primary
        placeholder:text-text-muted focus:border-accent-500 focus:outline-none focus:ring-1 focus:ring-accent-500"
    />
    <input
      bind:value={filterClient}
      onkeydown={handleFilterKeydown}
      placeholder="Client"
      class="w-32 rounded-md border border-surface-border bg-surface-secondary px-2.5 py-1.5 text-xs text-text-primary
        placeholder:text-text-muted focus:border-accent-500 focus:outline-none focus:ring-1 focus:ring-accent-500"
    />

    <span class="mx-0.5 h-5 w-px bg-surface-border"></span>

    {#each typeOptions as opt}
      <button
        onclick={() => {
          filterType = opt.value;
          applyFilters();
        }}
        class="rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors cursor-pointer
          {filterType === opt.value
          ? 'bg-accent-600/15 text-accent-400 border border-accent-500/30'
          : 'bg-surface-secondary text-text-secondary border border-surface-border hover:border-text-muted'}"
      >
        {opt.label}
      </button>
    {/each}

    {#if hasActiveFilters}
      <button
        onclick={clearFilters}
        class="ml-1 rounded-md px-2 py-1.5 text-xs text-text-muted transition-colors hover:text-text-secondary cursor-pointer"
        title="Clear all filters"
      >
        &times; Clear
      </button>
    {/if}
  </div>

  {#if error}
    <div
      class="shrink-0 rounded-lg border border-red-600/30 bg-red-600/10 px-4 py-3 text-sm text-red-400 mb-3"
    >
      {error}
    </div>
  {/if}

  <!-- Scrollable table area -->
  <div
    bind:this={scrollContainer}
    class="flex-1 min-h-0 overflow-y-auto overflow-x-auto rounded-lg border border-surface-border"
  >
    <table class="w-full text-xs">
      <thead class="sticky top-0 z-10">
        <tr
          class="border-b border-surface-border bg-surface-secondary text-left text-text-muted"
        >
          <th class="px-3 py-2 font-medium">Time</th>
          <th class="px-3 py-2 font-medium">Client</th>
          <th class="px-3 py-2 font-medium">Domain</th>
          <th class="px-3 py-2 font-medium">Type</th>
          <th class="px-3 py-2 font-medium">Response</th>
          <th class="px-3 py-2 font-medium">Code</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-surface-border">
        {#each entries as entry (entry._id)}
          {@const isBlocked = entry.response_reason
            .toUpperCase()
            .startsWith("BLOCKED")}
          <tr class="hover:bg-surface-hover">
            <td class="px-3 py-1.5 text-text-muted whitespace-nowrap font-mono"
              >{formatTime(entry.timestamp)}</td
            >
            <td class="px-3 py-1.5 text-text-secondary whitespace-nowrap">
              {clientDisplay(entry)}
            </td>
            <td
              class="px-3 py-1.5 font-mono text-text-primary max-w-xs truncate"
              >{stripDot(entry.domain)}</td
            >
            <td class="px-3 py-1.5 text-text-muted">{entry.query_type}</td>
            <td
              class="px-3 py-1.5 {isBlocked
                ? 'text-red-400'
                : 'text-text-secondary'}">{entry.response_reason}</td
            >
            <td class="px-3 py-1.5 text-text-muted">{entry.return_code}</td>
          </tr>
        {/each}
      </tbody>
    </table>

    <!-- Infinite scroll sentinel -->
    {#if !live && entries.length < total}
      <div bind:this={sentinel} class="h-8 flex items-center justify-center">
        {#if loading}
          <span class="text-xs text-text-muted">Loading...</span>
        {/if}
      </div>
    {/if}
  </div>

  <!-- Status bar (pinned bottom) -->
  <div
    class="shrink-0 flex items-center justify-between text-xs text-text-muted pt-2"
  >
    {#if live}
      <span class="flex items-center gap-1.5">
        {#if sseConnected}
          <span class="inline-flex h-1.5 w-1.5 rounded-full bg-emerald-400"
          ></span>
          Streaming — {entries.length} entries
        {:else}
          <span class="inline-flex h-1.5 w-1.5 rounded-full bg-amber-400"
          ></span>
          Reconnecting...
        {/if}
      </span>
    {:else}
      <span>Showing {entries.length} of {total.toLocaleString()} entries</span>
      {#if loading}
        <span class="text-text-muted">Loading...</span>
      {/if}
    {/if}
  </div>
</div>
