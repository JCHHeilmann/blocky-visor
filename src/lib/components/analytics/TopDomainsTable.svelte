<script lang="ts">
	interface DomainEntry {
		domain: string;
		count: number;
		reason?: string;
	}

	interface Props {
		domains: DomainEntry[];
		showReason?: boolean;
	}

	let { domains, showReason = false }: Props = $props();

	let maxCount = $derived(domains.length > 0 ? domains[0].count : 1);
</script>

{#if domains.length === 0}
	<p class="text-sm text-text-faint py-4 text-center">No data available</p>
{:else}
	<div class="space-y-2">
		{#each domains.slice(0, 10) as entry, i}
			<div class="flex items-center gap-3 text-sm">
				<span class="w-5 text-right text-text-faint font-mono text-xs">{i + 1}</span>
				<div class="flex-1 min-w-0">
					<div class="flex items-center justify-between gap-2">
						<span class="truncate font-mono text-text-primary text-xs">{entry.domain}</span>
						<span class="shrink-0 text-text-muted tabular-nums">{entry.count.toLocaleString()}</span>
					</div>
					<div class="mt-0.5 h-1 rounded-full bg-surface-secondary overflow-hidden">
						<div
							class="h-full rounded-full {showReason ? 'bg-red-500/60' : 'bg-accent-600/60'}"
							style="width: {(entry.count / maxCount) * 100}%"
						></div>
					</div>
					{#if showReason && entry.reason}
						<span class="text-[10px] text-text-faint">{entry.reason}</span>
					{/if}
				</div>
			</div>
		{/each}
	</div>
{/if}
