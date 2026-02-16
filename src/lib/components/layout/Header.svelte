<script lang="ts">
  import { page } from "$app/state";
  import { blockingStore } from "$lib/stores/blocking.svelte";

  interface Props {
    onmenuclick: () => void;
  }

  let { onmenuclick }: Props = $props();

  const titles: Record<string, string> = {
    "/": "Dashboard",
    "/query": "DNS Query",
    "/lists": "Lists",
    "/cache": "Cache",
    "/logs": "Query Logs",
    "/config": "Configuration",
    "/settings": "Settings",
  };

  let title = $derived(titles[page.url.pathname] || "Blocky Visor");
</script>

<header
  class="flex h-14 shrink-0 items-center gap-4 border-b border-surface-border bg-surface-primary/80 backdrop-blur-sm shadow-[var(--shadow-card)] px-4 lg:px-6"
>
  <button
    onclick={onmenuclick}
    aria-label="Menu"
    class="rounded-lg p-2 text-text-secondary hover:bg-surface-hover hover:text-text-primary lg:hidden cursor-pointer"
  >
    <svg
      class="h-5 w-5"
      fill="none"
      viewBox="0 0 24 24"
      stroke-width="1.5"
      stroke="currentColor"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
      />
    </svg>
  </button>

  <h1 class="text-lg font-semibold text-text-primary">{title}</h1>

  <div class="ml-auto">
    {#if !blockingStore.initialLoad}
      {#if blockingStore.enabled}
        <span
          class="inline-flex items-center gap-1.5 rounded-full bg-green-500/10 px-3 py-1 text-xs font-medium text-green-600 dark:text-green-400"
        >
          <span class="h-1.5 w-1.5 rounded-full bg-green-500 dark:bg-green-400"
          ></span>
          Blocking Active
        </span>
      {:else}
        <span
          class="inline-flex items-center gap-1.5 rounded-full bg-red-500/10 px-3 py-1 text-xs font-medium text-red-600 dark:text-red-400"
        >
          <span class="h-1.5 w-1.5 rounded-full bg-red-500 dark:bg-red-400"
          ></span>
          Blocking Disabled
        </span>
      {/if}
    {/if}
  </div>
</header>
