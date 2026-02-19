<script lang="ts">
  interface Props {
    status: "online" | "offline" | "warning" | "loading";
    label?: string;
  }

  let { status, label }: Props = $props();

  const styles = {
    online: {
      dot: "bg-green-500",
      ping: "bg-green-400",
      text: "text-green-400",
    },
    offline: { dot: "bg-red-500", text: "text-red-400" },
    warning: { dot: "bg-yellow-500", text: "text-yellow-400" },
    loading: { dot: "bg-gray-400 animate-pulse", text: "text-gray-400" },
  };
</script>

<span class="inline-flex items-center gap-2">
  <span class="relative flex h-2.5 w-2.5">
    {#if status === "online"}
      <span
        class="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-400 opacity-75"
      ></span>
    {/if}
    <span
      class="relative inline-flex h-2.5 w-2.5 rounded-full {styles[status].dot}"
    ></span>
  </span>
  {#if label}
    <span class="text-sm {styles[status].text}">{label}</span>
  {/if}
</span>
