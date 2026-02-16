<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import ConfigEditor from "$lib/components/config/ConfigEditor.svelte";
  import ServiceControl from "$lib/components/config/ServiceControl.svelte";
  import { sidecarStore } from "$lib/stores/sidecar.svelte";
</script>

<div class="space-y-6">
  {#if !sidecarStore.configured}
    <Card>
      <p class="text-sm text-text-muted">
        Configure a sidecar connection in Settings to manage Blocky
        configuration.
      </p>
    </Card>
  {:else}
    <div class="grid gap-6 xl:grid-cols-[1fr_320px]">
      <!-- Config editor (main area) -->
      <Card>
        <h2 class="mb-4 text-lg font-semibold text-text-primary">
          Configuration File
        </h2>
        <ConfigEditor />
      </Card>

      <!-- Sidebar: service control -->
      <div class="space-y-6">
        <Card>
          <h2 class="mb-4 text-lg font-semibold text-text-primary">
            Service Status
          </h2>
          <ServiceControl />
        </Card>

        <Card>
          <h2 class="mb-3 text-sm font-medium text-text-secondary">Tips</h2>
          <ul class="space-y-2 text-xs text-text-muted list-disc list-inside">
            <li>Changes are saved with a timestamped backup</li>
            <li>Restart Blocky after saving for changes to take effect</li>
            <li>YAML syntax is validated before saving</li>
            <li>Check Blocky logs after restart to verify config</li>
          </ul>
        </Card>
      </div>
    </div>
  {/if}
</div>
