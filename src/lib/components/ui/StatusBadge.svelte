<script lang="ts">
  interface Props {
    status: "online" | "offline" | "warning" | "loading";
    label?: string;
  }

  let { status, label }: Props = $props();

  const colors = {
    online: "bg-green-500",
    offline: "bg-red-500",
    warning: "bg-yellow-500",
    loading: "bg-gray-400 animate-pulse",
  };

  const textColors = {
    online: "text-green-400",
    offline: "text-red-400",
    warning: "text-yellow-400",
    loading: "text-gray-400",
  };
</script>

<span class="inline-flex items-center gap-2">
  <span class="relative flex h-2.5 w-2.5">
    {#if status === "online"}
      <span
        class="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-400 opacity-75"
      ></span>
    {/if}
    <span class="relative inline-flex h-2.5 w-2.5 rounded-full {colors[status]}"
    ></span>
  </span>
  {#if label}
    <span class="text-sm {textColors[status]}">{label}</span>
  {/if}
</span>
