<script lang="ts">
  import Button from "$lib/components/ui/Button.svelte";
  import ConnectionTestResult from "$lib/components/ui/ConnectionTestResult.svelte";
  import Input from "$lib/components/ui/Input.svelte";
  import { testConnection } from "$lib/api/client";
  import { settingsStore } from "$lib/stores/settings.svelte";

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
  <Input
    bind:value={url}
    label="Blocky API URL"
    placeholder="http://localhost:4000"
  />

  <div class="flex items-center gap-3">
    <Button onclick={handleSave} disabled={url === settingsStore.apiUrl}
      >Save</Button
    >
    <Button variant="secondary" onclick={handleTest} loading={testing}
      >Test Connection</Button
    >
    <ConnectionTestResult result={testResult} />
  </div>
</div>
