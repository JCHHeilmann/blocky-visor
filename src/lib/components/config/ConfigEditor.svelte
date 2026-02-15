<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import { fetchConfig, saveConfig } from '$lib/api/sidecar-config';
	import { toastStore } from '$lib/stores/toasts.svelte';

	let content = $state('');
	let original = $state('');
	let loading = $state(false);
	let saving = $state(false);
	let error = $state<string | null>(null);
	let showConfirm = $state(false);

	let hasChanges = $derived(content !== original);

	async function load() {
		loading = true;
		error = null;
		try {
			content = await fetchConfig();
			original = content;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load config';
		} finally {
			loading = false;
		}
	}

	async function handleSave() {
		showConfirm = false;
		saving = true;
		try {
			const result = await saveConfig(content);
			original = content;
			toastStore.success(`Config saved. Backup: ${result.backup}`);
		} catch (err) {
			const msg = err instanceof Error ? err.message : 'Failed to save config';
			toastStore.error(msg);
		} finally {
			saving = false;
		}
	}

	function handleRevert() {
		content = original;
	}

	$effect(() => {
		load();
	});
</script>

<div class="space-y-4">
	{#if error}
		<div class="rounded-lg border border-red-600/30 bg-red-600/10 px-4 py-3 text-sm text-red-400">
			{error}
		</div>
	{/if}

	{#if loading}
		<div class="flex justify-center py-8">
			<svg class="h-6 w-6 animate-spin text-accent-600" viewBox="0 0 24 24" fill="none">
				<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" class="opacity-25" />
				<path fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" class="opacity-75" />
			</svg>
		</div>
	{:else}
		<textarea
			bind:value={content}
			rows={25}
			class="w-full rounded-lg border border-surface-border bg-surface-bg p-4 font-mono text-sm text-text-primary
				focus:border-accent-500 focus:outline-none focus:ring-1 focus:ring-accent-500 resize-y"
			spellcheck="false"
		></textarea>

		<div class="flex items-center gap-3">
			<Button onclick={() => (showConfirm = true)} disabled={!hasChanges} loading={saving}>
				Save Config
			</Button>
			<Button variant="secondary" onclick={handleRevert} disabled={!hasChanges}>
				Revert Changes
			</Button>
			<Button variant="ghost" onclick={load}>
				Reload
			</Button>
			{#if hasChanges}
				<span class="text-sm text-warning">Unsaved changes</span>
			{/if}
		</div>
	{/if}
</div>

<Modal open={showConfirm} onclose={() => (showConfirm = false)} title="Save Configuration">
	{#snippet children()}
		<p class="text-sm text-text-secondary">
			This will save the configuration to disk and create a timestamped backup. You may need to restart Blocky for changes to take effect.
		</p>
	{/snippet}
	{#snippet actions()}
		<Button variant="secondary" onclick={() => (showConfirm = false)}>Cancel</Button>
		<Button onclick={handleSave}>Save</Button>
	{/snippet}
</Modal>
