<script lang="ts">
  interface Props {
    label: string;
    value: string;
    icon: string;
    trend?: string;
    color?: "accent" | "rose" | "amber" | "emerald";
  }

  let { label, value, icon, trend, color = "accent" }: Props = $props();

  const colorMap = {
    accent: {
      bg: "bg-accent-600/15",
      text: "text-accent-600 dark:text-accent-400",
    },
    rose: { bg: "bg-rose-500/15", text: "text-rose-600 dark:text-rose-400" },
    amber: {
      bg: "bg-amber-500/15",
      text: "text-amber-600 dark:text-amber-400",
    },
    emerald: {
      bg: "bg-emerald-500/15",
      text: "text-emerald-600 dark:text-emerald-400",
    },
  };

  let colors = $derived(colorMap[color]);
</script>

<div
  class="rounded-xl border border-surface-border bg-surface-primary p-5 shadow-[var(--shadow-card)]"
>
  <div class="flex items-center gap-3">
    <div
      class="flex h-10 w-10 items-center justify-center rounded-lg {colors.bg}"
    >
      <svg
        class="h-5 w-5 {colors.text}"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
      >
        <path stroke-linecap="round" stroke-linejoin="round" d={icon} />
      </svg>
    </div>
    <div class="min-w-0 flex-1">
      <p class="text-sm text-text-muted">{label}</p>
      <p class="text-2xl font-bold tracking-tight text-text-primary truncate">
        {value}
      </p>
    </div>
  </div>
  {#if trend}
    <p class="mt-2 text-xs text-text-muted">{trend}</p>
  {/if}
</div>
