<script lang="ts">
  import type { Snippet } from "svelte";
  import Sidebar from "./Sidebar.svelte";
  import Header from "./Header.svelte";
  import ToastContainer from "$lib/components/ui/ToastContainer.svelte";
  import TooltipOverlay from "$lib/components/ui/TooltipOverlay.svelte";
  import Spinner from "$lib/components/ui/Spinner.svelte";
  import { blockingStore } from "$lib/stores/blocking.svelte";
  import { settingsStore } from "$lib/stores/settings.svelte";

  interface Props {
    children: Snippet;
  }

  let { children }: Props = $props();
  let mobileOpen = $state(false);

  // Polling for blocking status
  $effect(() => {
    blockingStore.refresh();

    const interval = setInterval(() => {
      if (!document.hidden) {
        blockingStore.refresh();
      }
    }, settingsStore.refreshInterval * 1000);

    return () => clearInterval(interval);
  });
</script>

<div class="flex h-dvh overflow-hidden bg-surface-bg">
  <!-- Desktop sidebar -->
  <div class="hidden lg:flex lg:w-60 lg:shrink-0">
    <Sidebar />
  </div>

  <!-- Tablet sidebar (icons only) -->
  <div class="hidden md:flex md:w-16 md:shrink-0 lg:hidden">
    <Sidebar collapsed />
  </div>

  <!-- Mobile drawer overlay -->
  {#if mobileOpen}
    <div class="fixed inset-0 z-40 md:hidden">
      <button
        class="absolute inset-0 bg-black/60 backdrop-blur-[2px] animate-fade-in cursor-pointer"
        onclick={() => (mobileOpen = false)}
        aria-label="Close menu"
      ></button>
      <div class="relative z-50 h-full w-64 animate-slide-in-left">
        <Sidebar onnavigate={() => (mobileOpen = false)} />
      </div>
    </div>
  {/if}

  <!-- Main content -->
  <div class="flex flex-1 flex-col overflow-hidden">
    <Header onmenuclick={() => (mobileOpen = !mobileOpen)} />

    {#if blockingStore.initialLoad}
      <!-- Loading state on initial page load -->
      <main class="flex flex-1 items-center justify-center">
        <div class="flex flex-col items-center gap-3">
          <Spinner size="lg" />
          <p class="text-sm text-text-muted">Connecting to Blocky...</p>
        </div>
      </main>
    {:else}
      <!-- Connection-lost banner -->
      {#if !blockingStore.connected && blockingStore.error}
        <div
          class="border-b border-red-600/30 bg-red-950/30 dark:bg-red-950/50 px-4 py-2 text-center text-sm text-red-500 dark:text-red-400"
        >
          Unable to connect to Blocky — check your API URL in Settings
        </div>
      {/if}

      <main class="flex-1 overflow-y-auto p-4 lg:p-6">
        {@render children()}
      </main>
    {/if}
  </div>
</div>

<ToastContainer />
<TooltipOverlay />
