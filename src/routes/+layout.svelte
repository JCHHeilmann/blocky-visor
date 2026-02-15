<script lang="ts">
	import '../app.css';
	import AppShell from '$lib/components/layout/AppShell.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<AppShell>
	{@render children()}
</AppShell>
