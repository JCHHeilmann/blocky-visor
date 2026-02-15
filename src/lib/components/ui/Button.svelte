<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		loading?: boolean;
		onclick?: () => void;
		type?: 'button' | 'submit';
		children: Snippet;
	}

	let {
		variant = 'primary',
		size = 'md',
		disabled = false,
		loading = false,
		onclick,
		type = 'button',
		children
	}: Props = $props();

	const variants = {
		primary: 'bg-accent-600 hover:bg-accent-500 text-white',
		secondary: 'bg-surface-secondary hover:bg-surface-hover text-text-primary border border-surface-border',
		danger: 'bg-red-600/20 hover:bg-red-600/30 text-red-500 dark:text-red-400 border border-red-600/30',
		ghost: 'hover:bg-surface-hover text-text-secondary'
	};

	const sizes = {
		sm: 'px-3 py-1.5 text-sm',
		md: 'px-4 py-2 text-sm',
		lg: 'px-6 py-3 text-base'
	};
</script>

<button
	{type}
	{onclick}
	disabled={disabled || loading}
	class="inline-flex items-center justify-center gap-2 rounded-lg font-medium transition-colors
		disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer
		{variants[variant]} {sizes[size]}"
>
	{#if loading}
		<svg class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
			<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" class="opacity-25" />
			<path fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" class="opacity-75" />
		</svg>
	{/if}
	{@render children()}
</button>
