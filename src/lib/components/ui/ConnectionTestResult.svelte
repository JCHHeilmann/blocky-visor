<script lang="ts">
  interface Props {
    result: boolean | null;
  }

  let { result }: Props = $props();
</script>

{#if result === true}
  <span
    class="flex items-center gap-1.5 text-sm text-green-600 dark:text-green-400"
  >
    <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
      <path
        fill-rule="evenodd"
        d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z"
        clip-rule="evenodd"
      />
    </svg>
    Connected
  </span>
{:else if result === false}
  <span
    class="flex items-center gap-1.5 text-sm text-red-600 dark:text-red-400"
  >
    <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
      <path
        fill-rule="evenodd"
        d="M18 10a8 8 0 1 1-16 0 8 8 0 0 1 16 0Zm-8-5a.75.75 0 0 1 .75.75v4.5a.75.75 0 0 1-1.5 0v-4.5A.75.75 0 0 1 10 5Zm0 10a1 1 0 1 0 0-2 1 1 0 0 0 0 2Z"
        clip-rule="evenodd"
      />
    </svg>
    Connection failed
  </span>
{/if}
