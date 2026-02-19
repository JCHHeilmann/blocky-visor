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

  <h1 class="text-sm font-medium text-text-secondary">{title}</h1>

  {#if !blockingStore.initialLoad}
    {@const active = blockingStore.enabled}
    <span
      class="ml-auto inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium
        {active
        ? 'bg-green-500/10 text-green-600 dark:text-green-400'
        : 'bg-red-500/10 text-red-600 dark:text-red-400'}"
    >
      <span
        class="h-1.5 w-1.5 rounded-full {active
          ? 'bg-green-500 dark:bg-green-400'
          : 'bg-red-500 dark:bg-red-400'}"
      ></span>
      {active ? "Blocking Active" : "Blocking Disabled"}
    </span>
  {/if}
</header>
