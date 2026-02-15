<script lang="ts">
	import { dnsQuery } from '$lib/api/query';

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

	let hostnames = $state<Record<string, string>>({});

	function ipToArpa(ip: string): string | null {
		const parts = ip.split('.');
		if (parts.length !== 4) return null;
		return parts.reverse().join('.') + '.in-addr.arpa';
	}

	async function resolveHostname(ip: string): Promise<string | null> {
		const arpa = ipToArpa(ip);
		if (!arpa) return null;
		try {
			const result = await dnsQuery(arpa, 'PTR');
			if (result.returnCode === 'NOERROR' && result.response && result.response !== '') {
				const match = result.response.match(/PTR\s*\(([^)]+)\)/);
				if (match) return match[1].replace(/\.$/, '');
				return result.response.replace(/\.$/, '').trim();
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
			hostnames[client.ip] = '';
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

	function showIp(client: ClientData): boolean {
		const resolved = hostnames[client.ip];
		if (resolved && resolved !== client.ip) return true;
		if (!resolved && client.name && client.name !== client.ip) return true;
		return false;
	}
</script>

{#if clients.length === 0}
	<p class="text-sm text-text-faint py-4 text-center">No data available</p>
{:else}
	<div class="space-y-2">
		{#each clients as client}
			{@const totalPct = (client.total / maxTotal) * 100}
			{@const blockedPct = (client.blocked / maxTotal) * 100}
			<div class="text-sm">
				<div class="flex items-center justify-between gap-2">
					<span class="truncate min-w-0">
						<span class="text-text-primary">{displayName(client)}</span>
					</span>
					<span class="shrink-0 text-text-muted tabular-nums">
						{client.total.toLocaleString()}
						{#if client.blocked > 0}
							<span class="text-red-400 text-xs">({client.blocked} blocked)</span>
						{/if}
					</span>
				</div>
				<div class="mt-0.5 h-1.5 rounded-full bg-surface-secondary overflow-hidden">
					<div class="flex h-full">
						<div class="h-full bg-accent-600/60" style="width: {totalPct - blockedPct}%"></div>
						{#if blockedPct > 0}
							<div class="h-full bg-red-500/60" style="width: {blockedPct}%"></div>
						{/if}
					</div>
				</div>
				{#if showIp(client)}
					<span class="text-[10px] text-text-faint font-mono">{client.ip}</span>
				{/if}
			</div>
		{/each}
	</div>
{/if}
