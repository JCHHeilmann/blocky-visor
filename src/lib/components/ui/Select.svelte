<script lang="ts">
  interface Option {
    value: string;
    label: string;
  }

  interface Props {
    value: string;
    options: Option[];
    onchange?: (value: string) => void;
    label?: string;
    disabled?: boolean;
  }

  let {
    value = $bindable(),
    options,
    onchange,
    label,
    disabled = false,
  }: Props = $props();

  let id = $derived(
    label ? `select-${label.replace(/\s+/g, "-").toLowerCase()}` : undefined,
  );

  function handleChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    value = target.value;
    onchange?.(target.value);
  }
</script>

{#if label}
  <label for={id} class="block text-sm font-medium text-text-secondary mb-1.5"
    >{label}</label
  >
{/if}
<select
  {id}
  {value}
  {disabled}
  onchange={handleChange}
  class="w-full rounded-lg border border-surface-border bg-surface-secondary px-3 py-2 text-sm text-text-primary
		focus:border-accent-500 focus:outline-none focus:ring-1 focus:ring-accent-500
		disabled:opacity-50 cursor-pointer"
>
  {#each options as opt}
    <option value={opt.value}>{opt.label}</option>
  {/each}
</select>
