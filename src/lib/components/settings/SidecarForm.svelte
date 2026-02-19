<script lang="ts">
  import Button from "$lib/components/ui/Button.svelte";
  import ConnectionTestResult from "$lib/components/ui/ConnectionTestResult.svelte";
  import Input from "$lib/components/ui/Input.svelte";
  import { testSidecarConnection } from "$lib/api/sidecar";
  import { settingsStore } from "$lib/stores/settings.svelte";
  import { sidecarStore } from "$lib/stores/sidecar.svelte";

  let url = $state(settingsStore.sidecarUrl);
  let apiKey = $state(settingsStore.sidecarApiKey);
  let testing = $state(false);
  let testResult = $state<boolean | null>(null);

  let hasUnsavedChanges = $derived(
    url !== settingsStore.sidecarUrl || apiKey !== settingsStore.sidecarApiKey,
  );

  async function handleTest() {
    testing = true;
    testResult = null;
    testResult = await testSidecarConnection(url);
    testing = false;
  }

  function handleSave() {
    settingsStore.sidecarUrl = url;
    settingsStore.sidecarApiKey = apiKey;
    testResult = null;
    sidecarStore.checkConnection();
  }

  function handleClear() {
    url = "";
    apiKey = "";
    settingsStore.sidecarUrl = "";
    settingsStore.sidecarApiKey = "";
    testResult = null;
  }
</script>

<div class="space-y-4">
  <Input
    bind:value={url}
    label="Sidecar URL"
    placeholder="http://localhost:8550"
  />
  <Input
    bind:value={apiKey}
    label="API Key"
    type="password"
    placeholder="your-secret-api-key"
  />

  <div class="flex items-center gap-3">
    <Button onclick={handleSave} disabled={!hasUnsavedChanges}>Save</Button>
    <Button variant="secondary" onclick={handleTest} loading={testing}
      >Test Connection</Button
    >
    {#if settingsStore.sidecarConfigured}
      <Button variant="ghost" onclick={handleClear}>Clear</Button>
    {/if}
    <ConnectionTestResult result={testResult} />
  </div>

  <p class="text-xs text-text-faint">
    Optional. Connect to a Blocky Visor sidecar for query log analytics, config
    management, and service control.
  </p>
</div>
