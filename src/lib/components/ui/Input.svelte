<script lang="ts">
	interface Props {
		value: string;
		oninput?: (value: string) => void;
		placeholder?: string;
		type?: string;
		label?: string;
		disabled?: boolean;
		class?: string;
	}

	let {
		value = $bindable(),
		oninput,
		placeholder = '',
		type = 'text',
		label,
		disabled = false,
		class: className = ''
	}: Props = $props();

	let id = $derived(label ? `input-${label.replace(/\s+/g, '-').toLowerCase()}` : undefined);

	function handleInput(e: Event) {
		const target = e.target as HTMLInputElement;
		value = target.value;
		oninput?.(target.value);
	}
</script>

{#if label}
	<label for={id} class="block text-sm font-medium text-text-secondary mb-1.5">{label}</label>
{/if}
<input
	{id}
	{type}
	{value}
	{placeholder}
	{disabled}
	oninput={handleInput}
	class="w-full rounded-lg border border-surface-border bg-surface-secondary px-3 py-2 text-sm text-text-primary
		placeholder:text-text-muted focus:border-accent-500 focus:outline-none focus:ring-1
		focus:ring-accent-500 disabled:opacity-50 {className}"
/>
