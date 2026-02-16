<script lang="ts">
  import { formatNumber } from "$lib/utils/formatters";
  import { dnsQuery } from "$lib/api/query";
  import { tooltipStore } from "$lib/stores/tooltip.svelte";

  interface ClientData {
    ip: string;
    name: string;
    total: number;
    blocked: number;
  }

  interface Props {
    clients: ClientData[];
  }

  let { clients }: Props = $props();

  let maxTotal = $derived(clients.length > 0 ? clients[0].total : 1);
  let grandTotal = $derived(clients.reduce((sum, c) => sum + c.total, 0));

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
    } catch {
      // ignore
    }
    return null;
  }

  $effect(() => {
    if (!clients.length) return;
    for (const client of clients) {
      if (client.ip in hostnames) continue;
      hostnames[client.ip] = "";
      resolveHostname(client.ip).then((name) => {
        if (name) {
          hostnames = { ...hostnames, [client.ip]: name };
        }
      });
    }
  });

  function displayName(client: ClientData): string {
    if (hostnames[client.ip]) return hostnames[client.ip];
    if (client.name && client.name !== client.ip) return client.name;
    return client.ip;
  }

  function hasHostname(client: ClientData): boolean {
    const resolved = hostnames[client.ip];
    if (resolved && resolved !== client.ip) return true;
    if (!resolved && client.name && client.name !== client.ip) return true;
    return false;
  }

  function showClientTooltip(e: MouseEvent, client: ClientData) {
    const name = displayName(client);
    const pct =
      grandTotal > 0 ? ((client.total / grandTotal) * 100).toFixed(1) : "0";
    const html = `
			<div class="font-medium text-text-primary">${name}</div>
			${hasHostname(client) ? `<div class="text-xs text-text-muted font-mono">${client.ip}</div>` : ""}
			<div class="mt-1 text-xs space-y-0.5">
				<div class="flex justify-between gap-4">
					<span class="text-text-muted">Queries</span>
					<span class="font-mono text-text-primary">${formatNumber(client.total)}</span>
				</div>
				${
          client.blocked > 0
            ? `<div class="flex justify-between gap-4">
					<span class="text-text-muted">Blocked</span>
					<span class="font-mono text-red-400">${formatNumber(client.blocked)}</span>
				</div>`
            : ""
        }
				<div class="flex justify-between gap-4">
					<span class="text-text-muted">Share</span>
					<span class="font-mono text-text-primary">${pct}%</span>
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

{#if clients.length === 0}
  <p class="text-sm text-text-muted">No client data available.</p>
{:else}
  <div class="space-y-2.5">
    {#each clients as client}
      {@const totalPct = (client.total / maxTotal) * 100}
      {@const blockedPct = (client.blocked / maxTotal) * 100}
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class="cursor-default"
        onmouseenter={(e) => showClientTooltip(e, client)}
        onmousemove={moveTooltip}
        onmouseleave={() => tooltipStore.hide()}
      >
        <div class="flex items-center justify-between mb-1">
          <span class="truncate text-sm min-w-0">
            {#if hasHostname(client)}
              <span class="text-text-primary">{displayName(client)}</span>
              <span class="text-text-faint ml-1.5 font-mono text-xs"
                >{client.ip}</span
              >
            {:else}
              <span class="text-text-secondary font-mono">{client.ip}</span>
            {/if}
          </span>
          <span class="shrink-0 text-sm font-mono text-text-primary ml-2"
            >{formatNumber(client.total)}</span
          >
        </div>
        <div
          class="h-1.5 w-full overflow-hidden rounded-full bg-surface-secondary"
        >
          <div class="flex h-full">
            <div
              class="h-full rounded-l-full bg-accent-500/70 transition-all"
              style="width: {totalPct - blockedPct}%"
            ></div>
            {#if blockedPct > 0}
              <div
                class="h-full rounded-r-full bg-red-500/60 transition-all"
                style="width: {blockedPct}%"
              ></div>
            {/if}
          </div>
        </div>
      </div>
    {/each}
  </div>
{/if}
