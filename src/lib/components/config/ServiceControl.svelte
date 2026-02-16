<script lang="ts">
  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import { fetchServiceStatus, restartService } from "$lib/api/sidecar-service";
  import { toastStore } from "$lib/stores/toasts.svelte";
  import type { SidecarServiceStatus } from "$lib/types/api";

  let status = $state<SidecarServiceStatus | null>(null);
  let loading = $state(false);
  let restarting = $state(false);
  let error = $state<string | null>(null);
  let showConfirm = $state(false);

  async function loadStatus() {
    loading = true;
    error = null;
    try {
      status = await fetchServiceStatus();
    } catch (err) {
      error =
        err instanceof Error ? err.message : "Failed to get service status";
    } finally {
      loading = false;
    }
  }

  async function handleRestart() {
    showConfirm = false;
    restarting = true;
    try {
      await restartService();
      toastStore.success("Blocky service restarted");
      // Wait a moment then refresh status
      setTimeout(loadStatus, 2000);
    } catch (err) {
      const msg =
        err instanceof Error ? err.message : "Failed to restart service";
      toastStore.error(msg);
    } finally {
      restarting = false;
    }
  }

  $effect(() => {
    loadStatus();
    const interval = setInterval(() => {
      if (!document.hidden) loadStatus();
    }, 10000);
    return () => clearInterval(interval);
  });

  let isActive = $derived(status?.active === "active");
</script>

<div class="space-y-4">
  {#if error}
    <div
      class="rounded-lg border border-red-600/30 bg-red-600/10 px-4 py-3 text-sm text-red-400"
    >
      {error}
    </div>
  {/if}

  {#if status}
    <div class="space-y-3 text-sm">
      <div class="flex items-center justify-between">
        <span class="text-text-muted">Status</span>
        <span class="flex items-center gap-2">
          <span
            class="h-2 w-2 rounded-full {isActive
              ? 'bg-green-500'
              : 'bg-red-500'}"
          ></span>
          <span
            class="font-medium {isActive ? 'text-green-400' : 'text-red-400'}"
          >
            {status.active} ({status.sub_state})
          </span>
        </span>
      </div>
      {#if status.pid}
        <div class="flex justify-between">
          <span class="text-text-muted">PID</span>
          <span class="font-mono text-text-primary">{status.pid}</span>
        </div>
      {/if}
      {#if status.memory}
        <div class="flex justify-between">
          <span class="text-text-muted">Memory</span>
          <span class="font-mono text-text-primary">{status.memory}</span>
        </div>
      {/if}
    </div>
  {/if}

  <div class="flex gap-3">
    <Button
      variant="danger"
      onclick={() => (showConfirm = true)}
      loading={restarting}
    >
      Restart Blocky
    </Button>
    <Button variant="ghost" onclick={loadStatus} {loading}>
      Refresh Status
    </Button>
  </div>
</div>

<Modal
  open={showConfirm}
  onclose={() => (showConfirm = false)}
  title="Restart Blocky"
>
  {#snippet children()}
    <p class="text-sm text-text-secondary">
      This will restart the Blocky DNS service. DNS resolution will be briefly
      interrupted during the restart.
    </p>
  {/snippet}
  {#snippet actions()}
    <Button variant="secondary" onclick={() => (showConfirm = false)}
      >Cancel</Button
    >
    <Button variant="danger" onclick={handleRestart}>Restart</Button>
  {/snippet}
</Modal>
