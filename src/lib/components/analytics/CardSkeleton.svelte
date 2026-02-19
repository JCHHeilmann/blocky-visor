<script lang="ts">
  interface Props {
    lines?: number;
  }

  let { lines = 4 }: Props = $props();

  const widths = [75, 60, 85, 50, 70, 55, 80, 65];
</script>

<div class="animate-pulse space-y-3">
  {#each Array(lines) as _, i}
    <div class="flex items-center justify-between gap-4">
      <div
        class="h-3 rounded-full bg-surface-hover flex-1"
        style="width: {widths[i % widths.length]}%"
      ></div>
      <div class="h-3 w-12 rounded-full bg-surface-hover shrink-0"></div>
    </div>
  {/each}
</div>
