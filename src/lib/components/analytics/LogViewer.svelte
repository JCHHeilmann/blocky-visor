<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import type { SidecarLogEntry, SidecarLogsResponse } from '$lib/types/api';
	import type { StatsRange } from '$lib/api/sidecar-stats';
	import { fetchLogs } from '$lib/api/sidecar-stats';

	interface Props {
		range: StatsRange;
	}

	let { range }: Props = $props();

	let entries = $state<SidecarLogEntry[]>([]);
	let total = $state(0);
	let loading = $state(false);
	let error = $state<string | null>(null);
	let offset = $state(0);
	const limit = 50;

	let filterDomain = $state('');
	let filterClient = $state('');
	let filterType = $state('');

	async function load(reset = false) {
		if (reset) {
			offset = 0;
			entries = [];
		}
		loading = true;
		error = null;
		try {
			const result = await fetchLogs({
				range,
				limit,
				offset,
				domain: filterDomain || undefined,
				client: filterClient || undefined,
				type: filterType || undefined
			});
			if (reset) {
				entries = result.entries;
			} else {
				entries = [...entries, ...result.entries];
			}
			total = result.total;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load logs';
		} finally {
			loading = false;
		}
	}

	function loadMore() {
		offset = entries.length;
		load();
	}

	// Load on mount and when range changes
	$effect(() => {
		range; // track dependency
		load(true);
	});

	const typeOptions = [
		{ value: '', label: 'All' },
		{ value: 'blocked', label: 'Blocked' },
		{ value: 'cached', label: 'Cached' },
		{ value: 'resolved', label: 'Resolved' }
	];

	function formatTime(ts: string): string {
		const d = new Date(ts);
		return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
	}
</script>

<div class="space-y-3">
	<!-- Filters -->
	<div class="flex flex-wrap items-end gap-3">
		<div class="w-40">
			<Input bind:value={filterDomain} label="Domain" placeholder="Filter..." />
		</div>
		<div class="w-36">
			<Input bind:value={filterClient} label="Client" placeholder="Filter..." />
		</div>
		<div>
			<span class="mb-1 block text-sm font-medium text-text-secondary">Type</span>
			<div class="flex gap-1">
				{#each typeOptions as opt}
					<button
						onclick={() => { filterType = opt.value; load(true); }}
						class="rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors cursor-pointer
							{filterType === opt.value
								? 'bg-accent-600/15 text-accent-400 border border-accent-500/30'
								: 'bg-surface-secondary text-text-secondary border border-surface-border hover:border-text-muted'}"
					>
						{opt.label}
					</button>
				{/each}
			</div>
		</div>
		<Button size="sm" variant="secondary" onclick={() => load(true)}>Apply</Button>
	</div>

	{#if error}
		<div class="rounded-lg border border-red-600/30 bg-red-600/10 px-4 py-3 text-sm text-red-400">
			{error}
		</div>
	{/if}

	<!-- Log table -->
	<div class="overflow-x-auto rounded-lg border border-surface-border">
		<table class="w-full text-xs">
			<thead>
				<tr class="border-b border-surface-border bg-surface-secondary text-left text-text-muted">
					<th class="px-3 py-2 font-medium">Time</th>
					<th class="px-3 py-2 font-medium">Client</th>
					<th class="px-3 py-2 font-medium">Domain</th>
					<th class="px-3 py-2 font-medium">Type</th>
					<th class="px-3 py-2 font-medium">Response</th>
					<th class="px-3 py-2 font-medium">Code</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-surface-border">
				{#each entries as entry}
					{@const isBlocked = entry.response_reason.toUpperCase().startsWith('BLOCKED')}
					<tr class="hover:bg-surface-hover">
						<td class="px-3 py-1.5 text-text-muted whitespace-nowrap font-mono">{formatTime(entry.timestamp)}</td>
						<td class="px-3 py-1.5 text-text-secondary whitespace-nowrap">
							{entry.client_name && entry.client_name !== entry.client_ip ? entry.client_name : entry.client_ip}
						</td>
						<td class="px-3 py-1.5 font-mono text-text-primary max-w-xs truncate">{entry.domain}</td>
						<td class="px-3 py-1.5 text-text-muted">{entry.query_type}</td>
						<td class="px-3 py-1.5 {isBlocked ? 'text-red-400' : 'text-text-secondary'}">{entry.response_reason}</td>
						<td class="px-3 py-1.5 text-text-muted">{entry.return_code}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<!-- Pagination -->
	<div class="flex items-center justify-between text-sm text-text-muted">
		<span>Showing {entries.length} of {total.toLocaleString()} entries</span>
		{#if entries.length < total}
			<Button size="sm" variant="secondary" onclick={loadMore} {loading}>
				Load More
			</Button>
		{/if}
	</div>
</div>
