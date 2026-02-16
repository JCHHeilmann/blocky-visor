<script lang="ts">
  interface Props {
    size?: "sm" | "md" | "lg";
  }

  let { size = "md" }: Props = $props();

  const sizes = { sm: "h-4 w-4", md: "h-6 w-6", lg: "h-8 w-8" };
</script>

<svg
  class="{sizes[size]} animate-spin text-accent-600 dark:text-accent-400"
  viewBox="0 0 24 24"
  fill="none"
>
  <circle
    cx="12"
    cy="12"
    r="10"
    stroke="currentColor"
    stroke-width="4"
    class="opacity-25"
  />
  <path
    fill="currentColor"
    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
    class="opacity-75"
  />
</svg>
