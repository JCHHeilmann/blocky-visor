<script lang="ts">
  import { toastStore } from "$lib/stores/toasts.svelte";
  import Toast from "./Toast.svelte";
</script>

<div class="fixed bottom-4 right-4 z-[100] flex flex-col gap-2 w-80">
  {#each toastStore.items as toast (toast.id)}
    <Toast {toast} onclose={() => toastStore.remove(toast.id)} />
  {/each}
</div>
