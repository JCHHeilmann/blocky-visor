<script lang="ts">
  import { formatNumber } from "$lib/utils/formatters";
  import { dnsQuery } from "$lib/api/query";
  import { tooltipStore } from "$lib/stores/tooltip.svelte";

  interface Props {
    queriesByClient?: Record<string, number>;
    total?: number;
  }

  let { queriesByClient, total }: Props = $props();

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
        // Response format: "PTR (hostname.local.)" — extract hostname from parens
        const match = result.response.match(/PTR\s*\(([^)]+)\)/);
        if (match) return match[1].replace(/\.$/, "");
        return result.response.replace(/\.$/, "").trim();
      }
    } catch {
      // Ignore resolution failures
    }
    return null;
  }

  $effect(() => {
    if (!queriesByClient) return;
    for (const ip of Object.keys(queriesByClient)) {
      if (ip in hostnames) continue;
      // Mark as resolving to avoid duplicate requests
      hostnames[ip] = "";
      resolveHostname(ip).then((name) => {
        if (name) {
          hostnames = { ...hostnames, [ip]: name };
        }
      });
    }
  });

  let clients = $derived.by(() => {
    if (!queriesByClient) return [];
    return Object.entries(queriesByClient)
      .sort(([, a], [, b]) => b - a)
      .map(([ip, count]) => ({
        ip,
        hostname: hostnames[ip] || null,
        count,
        pct: total && total > 0 ? ((count / total) * 100).toFixed(1) : "0",
      }));
  });

  function showClientTooltip(e: MouseEvent, client: (typeof clients)[0]) {
    const html = `
			<div class="font-medium text-text-primary">${client.hostname || client.ip}</div>
			${client.hostname ? `<div class="text-xs text-text-muted font-mono">${client.ip}</div>` : ""}
			<div class="mt-1 text-xs space-y-0.5">
				<div class="flex justify-between gap-4">
					<span class="text-text-muted">Queries</span>
					<span class="font-mono text-text-primary">${formatNumber(client.count)}</span>
				</div>
				<div class="flex justify-between gap-4">
					<span class="text-text-muted">Share</span>
					<span class="font-mono text-text-primary">${client.pct}%</span>
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
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class="flex items-center gap-3 cursor-default"
        onmouseenter={(e) => showClientTooltip(e, client)}
        onmousemove={moveTooltip}
        onmouseleave={() => tooltipStore.hide()}
      >
        <div class="min-w-0 flex-1">
          <div class="flex items-center justify-between mb-1">
            <span class="truncate text-sm min-w-0">
              {#if client.hostname}
                <span class="text-text-primary">{client.hostname}</span>
                <span class="text-text-faint ml-1.5 font-mono text-xs"
                  >{client.ip}</span
                >
              {:else}
                <span class="text-text-secondary font-mono">{client.ip}</span>
              {/if}
            </span>
            <span class="shrink-0 text-sm font-mono text-text-primary ml-2"
              >{formatNumber(client.count)}</span
            >
          </div>
          <div
            class="h-1.5 w-full overflow-hidden rounded-full bg-surface-secondary"
          >
            <div
              class="h-full rounded-full bg-accent-500/70 transition-all"
              style="width: {client.pct}%"
            ></div>
          </div>
        </div>
      </div>
    {/each}
  </div>
{/if}
