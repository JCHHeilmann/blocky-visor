<script lang="ts">
  import type { StatsRange } from "$lib/api/sidecar-stats";

  interface Props {
    value: StatsRange;
    onchange: (range: StatsRange) => void;
  }

  let { value, onchange }: Props = $props();

  const options: { value: StatsRange; label: string }[] = [
    { value: "today", label: "Today" },
    { value: "yesterday", label: "Yesterday" },
    { value: "7d", label: "7 Days" },
    { value: "30d", label: "30 Days" },
  ];
</script>

<div class="flex gap-1.5">
  {#each options as option}
    <button
      onclick={() => onchange(option.value)}
      class="rounded-lg px-3 py-1.5 text-sm font-medium transition-colors cursor-pointer
				{value === option.value
        ? 'bg-accent-600/15 text-accent-600 dark:text-accent-400 border border-accent-500/30'
        : 'bg-surface-secondary text-text-secondary border border-surface-border hover:border-text-muted'}"
    >
      {option.label}
    </button>
  {/each}
</div>
