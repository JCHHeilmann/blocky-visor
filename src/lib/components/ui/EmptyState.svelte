<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		icon?: string;
		title: string;
		description?: string;
		children?: Snippet;
	}

	let { icon, title, description, children }: Props = $props();
</script>

<div class="flex flex-col items-center justify-center py-12 text-center">
	{#if icon}
		<svg class="mb-4 h-12 w-12 text-text-faint" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor">
			<path stroke-linecap="round" stroke-linejoin="round" d={icon} />
		</svg>
	{/if}
	<h3 class="text-lg font-medium text-text-secondary">{title}</h3>
	{#if description}
		<p class="mt-1 text-sm text-text-muted">{description}</p>
	{/if}
	{#if children}
		<div class="mt-4">
			{@render children()}
		</div>
	{/if}
</div>
