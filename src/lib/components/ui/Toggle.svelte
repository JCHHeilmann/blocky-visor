<script lang="ts">
	interface Props {
		checked: boolean;
		onchange: (checked: boolean) => void;
		disabled?: boolean;
		size?: 'sm' | 'md' | 'lg';
	}

	let { checked, onchange, disabled = false, size = 'md' }: Props = $props();

	const sizes = {
		sm: { track: 'w-8 h-4', thumb: 'h-3 w-3', translate: 'translate-x-[18px]' },
		md: { track: 'w-11 h-6', thumb: 'h-5 w-5', translate: 'translate-x-[22px]' },
		lg: { track: 'w-14 h-7', thumb: 'h-6 w-6', translate: 'translate-x-[30px]' }
	};

	function handleClick() {
		if (!disabled) onchange(!checked);
	}
</script>

<button
	type="button"
	role="switch"
	aria-checked={checked}
	aria-label="Toggle"
	{disabled}
	onclick={handleClick}
	class="relative inline-flex shrink-0 cursor-pointer rounded-full transition-colors duration-200
		disabled:opacity-50 disabled:cursor-not-allowed
		{sizes[size].track}
		{checked ? 'bg-accent-600' : 'bg-surface-border dark:bg-gray-700'}"
>
	<span
		class="pointer-events-none inline-block transform rounded-full bg-white shadow-sm
			ring-0 transition-transform duration-200
			{sizes[size].thumb}
			{checked ? sizes[size].translate : 'translate-x-0.5'}
			mt-0.5"
	></span>
</button>
