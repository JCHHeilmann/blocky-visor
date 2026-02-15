<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import { settingsStore } from '$lib/stores/settings.svelte';
	import { sidecarStore } from '$lib/stores/sidecar.svelte';
	import { testSidecarConnection } from '$lib/api/sidecar';

	let url = $state(settingsStore.sidecarUrl);
	let apiKey = $state(settingsStore.sidecarApiKey);
	let testing = $state(false);
	let testResult = $state<boolean | null>(null);

	async function handleTest() {
		testing = true;
		testResult = null;
		testResult = await testSidecarConnection(url, apiKey);
		testing = false;
	}

	function handleSave() {
		settingsStore.sidecarUrl = url;
		settingsStore.sidecarApiKey = apiKey;
		testResult = null;
		sidecarStore.checkConnection();
	}

	function handleClear() {
		url = '';
		apiKey = '';
		settingsStore.sidecarUrl = '';
		settingsStore.sidecarApiKey = '';
		testResult = null;
	}
</script>

<div class="space-y-4">
	<Input bind:value={url} label="Sidecar URL" placeholder="http://localhost:8550" />
	<Input bind:value={apiKey} label="API Key" type="password" placeholder="your-secret-api-key" />

	<div class="flex items-center gap-3">
		<Button
			onclick={handleSave}
			disabled={url === settingsStore.sidecarUrl && apiKey === settingsStore.sidecarApiKey}
		>Save</Button>
		<Button variant="secondary" onclick={handleTest} loading={testing}>Test Connection</Button>
		{#if settingsStore.sidecarConfigured}
			<Button variant="ghost" onclick={handleClear}>Clear</Button>
		{/if}

		{#if testResult === true}
			<span class="flex items-center gap-1.5 text-sm text-green-600 dark:text-green-400">
				<svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z" clip-rule="evenodd" />
				</svg>
				Connected
			</span>
		{:else if testResult === false}
			<span class="flex items-center gap-1.5 text-sm text-red-600 dark:text-red-400">
				<svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M18 10a8 8 0 1 1-16 0 8 8 0 0 1 16 0Zm-8-5a.75.75 0 0 1 .75.75v4.5a.75.75 0 0 1-1.5 0v-4.5A.75.75 0 0 1 10 5Zm0 10a1 1 0 1 0 0-2 1 1 0 0 0 0 2Z" clip-rule="evenodd" />
				</svg>
				Connection failed
			</span>
		{/if}
	</div>

	<p class="text-xs text-text-faint">
		Optional. Connect to a Blocky Visor sidecar for query log analytics, config management, and service control.
	</p>
</div>
