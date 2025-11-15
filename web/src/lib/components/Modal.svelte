<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		isOpen: boolean;
		title: string;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		onClose: () => void;
		children: Snippet;
		footer?: Snippet;
	}

	let { isOpen, title, size = 'md', onClose, children, footer }: Props = $props();

	const sizeClasses = {
		sm: 'max-w-md',
		md: 'max-w-2xl',
		lg: 'max-w-4xl',
		xl: 'max-w-6xl'
	};

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			onClose();
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			onClose();
		}
	}
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4"
		onclick={handleBackdropClick}
		onkeydown={handleKeydown}
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<div class="w-full {sizeClasses[size]} border border-slate-800 bg-slate-900">
			<!-- Header -->
			<div class="flex items-center justify-between border-b border-slate-800 px-8 py-4">
				<h2 class="text-xl font-bold text-slate-100">{title}</h2>
				<button
					onclick={onClose}
					class="text-slate-500 transition-colors hover:text-slate-300"
					aria-label="Close"
				>
					<svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						></path>
					</svg>
				</button>
			</div>

			<!-- Content -->
			<div class="max-h-[calc(100vh-12rem)] overflow-y-auto px-8 py-6">
				{@render children()}
			</div>

			<!-- Footer (optional) -->
			{#if footer}
				<div class="border-t border-slate-800 px-8 py-4">
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}
