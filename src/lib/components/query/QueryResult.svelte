<script lang="ts">
	import type { DnsQueryResponse } from '$lib/types/api';

	interface Props {
		result: DnsQueryResponse;
		domain: string;
		type: string;
	}

	let { result, domain, type }: Props = $props();
</script>

<div class="rounded-xl border border-surface-border bg-surface-primary p-5">
	<div class="mb-3 flex items-center gap-2">
		<span class="font-mono text-sm text-accent-600 dark:text-accent-400">{domain}</span>
		<span class="rounded bg-surface-secondary px-2 py-0.5 text-xs text-text-secondary">{type}</span>
		<span class="rounded px-2 py-0.5 text-xs
			{result.returnCode === 'NOERROR' ? 'bg-green-500/10 text-green-600 dark:text-green-400' : 'bg-red-500/10 text-red-600 dark:text-red-400'}">
			{result.returnCode}
		</span>
	</div>
	<div class="space-y-2 text-sm">
		<div class="flex gap-2">
			<span class="text-text-muted w-24 shrink-0">Response:</span>
			<pre class="flex-1 whitespace-pre-wrap break-all font-mono text-text-primary">{result.response}</pre>
		</div>
		<div class="flex gap-2">
			<span class="text-text-muted w-24 shrink-0">Reason:</span>
			<span class="text-text-secondary">{result.reason}</span>
		</div>
		<div class="flex gap-2">
			<span class="text-text-muted w-24 shrink-0">Type:</span>
			<span class="text-text-secondary">{result.responseType}</span>
		</div>
	</div>
</div>
