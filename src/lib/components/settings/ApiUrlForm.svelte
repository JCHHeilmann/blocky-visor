<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import { settingsStore } from '$lib/stores/settings.svelte';
	import { testConnection } from '$lib/api/client';

	let url = $state(settingsStore.apiUrl);
	let testing = $state(false);
	let testResult = $state<boolean | null>(null);

	async function handleTest() {
		testing = true;
		testResult = null;
		testResult = await testConnection(url);
		testing = false;
	}

	function handleSave() {
		settingsStore.apiUrl = url;
		testResult = null;
	}
</script>

<div class="space-y-4">
	<Input bind:value={url} label="Blocky API URL" placeholder="http://localhost:4000" />

	<div class="flex items-center gap-3">
		<Button onclick={handleSave} disabled={url === settingsStore.apiUrl}>Save</Button>
		<Button variant="secondary" onclick={handleTest} loading={testing}>Test Connection</Button>

		{#if testResult === true}
			<span class="flex items-center gap-1.5 text-sm text-green-600 dark:text-green-400">
				<svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z" clip-rule="evenodd" />
				</svg>
				Connection successful
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
</div>
