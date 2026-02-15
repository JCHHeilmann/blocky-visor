<script lang="ts">
	import Toggle from '$lib/components/ui/Toggle.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import { blockingStore } from '$lib/stores/blocking.svelte';
	import { enableBlocking, disableBlocking } from '$lib/api/blocking';
	import { toastStore } from '$lib/stores/toasts.svelte';
	import { formatDuration } from '$lib/utils/formatters';

	let showDisableModal = $state(false);
	let selectedDuration = $state('5m');
	let customMinutes = $state('');
	let loading = $state(false);

	const presets = [
		{ value: '5m', label: '5 minutes' },
		{ value: '15m', label: '15 minutes' },
		{ value: '30m', label: '30 minutes' },
		{ value: '1h', label: '1 hour' },
		{ value: '2h', label: '2 hours' },
		{ value: '', label: 'Indefinitely' },
		{ value: 'custom', label: 'Custom' }
	];

	async function handleToggle(checked: boolean) {
		if (checked) {
			await doEnable();
		} else {
			showDisableModal = true;
		}
	}

	async function doEnable() {
		loading = true;
		try {
			await enableBlocking();
			await blockingStore.refresh();
			toastStore.success('Blocking enabled');
		} catch {
			toastStore.error('Failed to enable blocking');
		} finally {
			loading = false;
		}
	}

	async function doDisable() {
		loading = true;
		try {
			let duration: string | undefined;
			if (selectedDuration === 'custom' && customMinutes) {
				duration = `${customMinutes}m`;
			} else if (selectedDuration && selectedDuration !== 'custom') {
				duration = selectedDuration;
			}
			await disableBlocking(duration);
			await blockingStore.refresh();
			showDisableModal = false;
			const label = duration ? `for ${duration}` : 'indefinitely';
			toastStore.warning(`Blocking disabled ${label}`);
		} catch {
			toastStore.error('Failed to disable blocking');
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex items-center gap-4">
	<Toggle
		checked={blockingStore.enabled}
		onchange={handleToggle}
		disabled={loading || !blockingStore.connected}
		size="lg"
	/>
	<div>
		<p class="font-medium {blockingStore.enabled ? 'text-green-600 dark:text-green-400' : 'text-red-500 dark:text-red-400'}">
			{blockingStore.enabled ? 'Blocking Active' : 'Blocking Disabled'}
		</p>
		{#if !blockingStore.enabled && blockingStore.autoEnableInSec}
			<p class="text-sm text-text-muted">
				Re-enables in {formatDuration(blockingStore.autoEnableInSec)}
			</p>
		{/if}
	</div>
</div>

<Modal open={showDisableModal} onclose={() => (showDisableModal = false)} title="Disable Blocking">
	<p class="mb-4 text-sm text-text-secondary">How long should blocking be disabled?</p>
	<div class="grid grid-cols-2 gap-2">
		{#each presets as preset}
			<button
				onclick={() => (selectedDuration = preset.value)}
				class="rounded-lg border px-3 py-2 text-sm transition-colors cursor-pointer
					{selectedDuration === preset.value
						? 'border-accent-500 bg-accent-600/15 text-accent-600 dark:text-accent-400'
						: 'border-surface-border bg-surface-secondary text-text-secondary hover:border-text-muted'}"
			>
				{preset.label}
			</button>
		{/each}
	</div>
	{#if selectedDuration === 'custom'}
		<div class="mt-3">
			<label for="custom-minutes" class="block text-sm text-text-secondary mb-1">Minutes</label>
			<input
				id="custom-minutes"
				type="number"
				bind:value={customMinutes}
				placeholder="e.g. 45"
				min="1"
				class="w-full rounded-lg border border-surface-border bg-surface-secondary px-3 py-2 text-sm text-text-primary
					focus:border-accent-500 focus:outline-none focus:ring-1 focus:ring-accent-500"
			/>
		</div>
	{/if}

	{#snippet actions()}
		<Button variant="ghost" onclick={() => (showDisableModal = false)}>Cancel</Button>
		<Button variant="danger" onclick={doDisable} {loading}>Disable</Button>
	{/snippet}
</Modal>
