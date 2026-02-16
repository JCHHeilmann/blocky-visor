<script lang="ts">
  import type { Toast } from "$lib/stores/toasts.svelte";

  interface Props {
    toast: Toast;
    onclose: () => void;
  }

  let { toast, onclose }: Props = $props();

  const icons: Record<Toast["type"], string> = {
    success: "M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z",
    error:
      "M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z",
    info: "m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z",
    warning:
      "M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z",
  };

  const colors: Record<Toast["type"], string> = {
    success:
      "text-green-600 dark:text-green-400 border-green-300 dark:border-green-600/30 bg-green-50 dark:bg-green-950/50",
    error:
      "text-red-600 dark:text-red-400 border-red-300 dark:border-red-600/30 bg-red-50 dark:bg-red-950/50",
    info: "text-accent-600 dark:text-accent-400 border-accent-300 dark:border-accent-600/30 bg-accent-50 dark:bg-accent-950/50",
    warning:
      "text-yellow-600 dark:text-yellow-400 border-yellow-300 dark:border-yellow-600/30 bg-yellow-50 dark:bg-yellow-950/50",
  };
</script>

<div
  class="flex items-start gap-3 rounded-lg border px-4 py-3 shadow-lg {colors[
    toast.type
  ]}"
>
  <svg
    class="mt-0.5 h-5 w-5 shrink-0"
    fill="none"
    viewBox="0 0 24 24"
    stroke-width="1.5"
    stroke="currentColor"
  >
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      d={icons[toast.type]}
    />
  </svg>
  <p class="flex-1 text-sm text-text-primary">{toast.message}</p>
  <button
    onclick={onclose}
    aria-label="Dismiss"
    class="shrink-0 text-text-muted hover:text-text-primary cursor-pointer"
  >
    <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
      <path
        d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z"
      />
    </svg>
  </button>
</div>
