<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		open: boolean;
		onclose: () => void;
		title: string;
		children: Snippet;
		actions?: Snippet;
	}

	let { open, onclose, title, children, actions }: Props = $props();

	function handleBackdrop(e: MouseEvent) {
		if (e.target === e.currentTarget) onclose();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onclose();
	}
</script>

{#if open}
	<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
	<!-- svelte-ignore a11y_interactive_supports_focus -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4"
		onclick={handleBackdrop}
		onkeydown={handleKeydown}
		role="dialog"
		aria-modal="true"
		aria-label={title}
	>
		<div class="w-full max-w-md rounded-xl border border-surface-border bg-surface-primary shadow-2xl">
			<div class="flex items-center justify-between border-b border-surface-border px-5 py-4">
				<h2 class="text-lg font-semibold text-text-primary">{title}</h2>
				<button
					onclick={onclose}
					aria-label="Close"
					class="rounded-lg p-1 text-text-secondary hover:bg-surface-hover hover:text-text-primary cursor-pointer"
				>
					<svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
						<path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
					</svg>
				</button>
			</div>
			<div class="px-5 py-4">
				{@render children()}
			</div>
			{#if actions}
				<div class="flex justify-end gap-3 border-t border-surface-border px-5 py-4">
					{@render actions()}
				</div>
			{/if}
		</div>
	</div>
{/if}
