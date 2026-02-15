<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		class?: string;
		children: Snippet;
	}

	let { class: className = '', children }: Props = $props();
</script>

<div class="rounded-xl border border-surface-border bg-surface-primary p-5 {className}">
	{@render children()}
</div>
