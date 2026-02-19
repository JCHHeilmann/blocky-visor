<script lang="ts">
  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import Spinner from "$lib/components/ui/Spinner.svelte";
  import YamlEditor from "./YamlEditor.svelte";
  import { fetchConfig, saveConfig } from "$lib/api/sidecar-config";
  import { toastStore } from "$lib/stores/toasts.svelte";

  let content = $state("");
  let original = $state("");
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
      error = err instanceof Error ? err.message : "Failed to load config";
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
      const msg = err instanceof Error ? err.message : "Failed to save config";
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

<div class="flex flex-col h-full min-h-0">
  {#if error}
    <div
      class="shrink-0 rounded-lg border border-red-600/30 bg-red-600/10 px-4 py-3 text-sm text-red-400 mb-3"
    >
      {error}
    </div>
  {/if}

  <div class="shrink-0 flex items-center gap-3 pb-3">
    <Button
      onclick={() => (showConfirm = true)}
      disabled={!hasChanges}
      loading={saving}
    >
      Save
    </Button>
    <Button variant="secondary" onclick={handleRevert} disabled={!hasChanges}>
      Revert
    </Button>
    <Button variant="ghost" onclick={load}>Reload</Button>
    {#if hasChanges}
      <span class="text-xs text-warning">Unsaved changes</span>
    {/if}
  </div>

  {#if loading}
    <div class="flex-1 flex items-center justify-center">
      <Spinner />
    </div>
  {:else}
    <div class="flex-1 min-h-0">
      <YamlEditor value={content} onchange={(v) => (content = v)} />
    </div>
  {/if}
</div>

<Modal
  open={showConfirm}
  onclose={() => (showConfirm = false)}
  title="Save Configuration"
>
  {#snippet children()}
    <p class="text-sm text-text-secondary">
      This will save the configuration to disk and create a timestamped backup.
      You may need to restart Blocky for changes to take effect.
    </p>
  {/snippet}
  {#snippet actions()}
    <Button variant="secondary" onclick={() => (showConfirm = false)}
      >Cancel</Button
    >
    <Button onclick={handleSave}>Save</Button>
  {/snippet}
</Modal>
