<script lang="ts">
	interface Props {
		label: string;
		value: string;
		icon: string;
		trend?: string;
	}

	let { label, value, icon, trend }: Props = $props();
</script>

<div class="rounded-xl border border-surface-border bg-surface-primary p-5">
	<div class="flex items-center gap-3">
		<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-accent-600/15">
			<svg class="h-5 w-5 text-accent-600 dark:text-accent-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" d={icon} />
			</svg>
		</div>
		<div class="min-w-0 flex-1">
			<p class="text-sm text-text-muted">{label}</p>
			<p class="text-2xl font-bold text-text-primary truncate">{value}</p>
		</div>
	</div>
	{#if trend}
		<p class="mt-2 text-xs text-text-muted">{trend}</p>
	{/if}
</div>
