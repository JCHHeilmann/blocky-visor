<script lang="ts">
  import Card from "$lib/components/ui/Card.svelte";
  import ConfigEditor from "$lib/components/config/ConfigEditor.svelte";
  import ServiceControl from "$lib/components/config/ServiceControl.svelte";
  import { sidecarStore } from "$lib/stores/sidecar.svelte";
</script>

{#if !sidecarStore.configured}
  <Card>
    <p class="text-sm text-text-muted">
      Configure a sidecar connection in Settings to manage Blocky
      configuration.
    </p>
  </Card>
{:else}
  <div class="flex flex-col h-full gap-4 lg:flex-row">
    <!-- Config editor (main area, fills height) -->
    <Card class="flex-1 min-h-0 flex flex-col overflow-hidden !p-4">
      <ConfigEditor />
    </Card>

    <!-- Sidebar: service control -->
    <div class="shrink-0 lg:w-72">
      <Card>
        <h2 class="mb-4 text-sm font-semibold text-text-primary">
          Service Status
        </h2>
        <ServiceControl />
      </Card>
    </div>
  </div>
{/if}
