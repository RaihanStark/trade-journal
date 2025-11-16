<script lang="ts">
	import { fade, scale } from 'svelte/transition';

	interface Props {
		isOpen: boolean;
		title: string;
		message: string;
		confirmText?: string;
		cancelText?: string;
		variant?: 'danger' | 'warning' | 'info';
		onConfirm: () => void;
		onCancel: () => void;
	}

	let {
		isOpen,
		title,
		message,
		confirmText = 'Confirm',
		cancelText = 'Cancel',
		variant = 'danger',
		onConfirm,
		onCancel
	}: Props = $props();

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			onCancel();
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			onCancel();
		}
	}
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/80"
		onclick={handleBackdropClick}
		onkeydown={handleKeydown}
		role="dialog"
		aria-modal="true"
		tabindex="-1"
		transition:fade={{ duration: 200 }}
	>
		<div
			class="w-full max-w-md border border-slate-800 bg-slate-900 p-6"
			transition:scale={{ duration: 200, start: 0.95 }}
		>
			<!-- Icon & Title -->
			<div class="mb-4 flex items-start gap-4">
				{#if variant === 'danger'}
					<div class="flex h-10 w-10 items-center justify-center rounded-full bg-red-900/30">
						<svg class="h-6 w-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
							></path>
						</svg>
					</div>
				{:else if variant === 'warning'}
					<div class="flex h-10 w-10 items-center justify-center rounded-full bg-yellow-900/30">
						<svg
							class="h-6 w-6 text-yellow-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
							></path>
						</svg>
					</div>
				{:else}
					<div class="flex h-10 w-10 items-center justify-center rounded-full bg-blue-900/30">
						<svg class="h-6 w-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
							></path>
						</svg>
					</div>
				{/if}
				<div class="flex-1">
					<h3 class="text-lg font-bold text-slate-100">{title}</h3>
					<p class="mt-2 text-sm text-slate-400">{message}</p>
				</div>
			</div>

			<!-- Actions -->
			<div class="mt-6 flex justify-end gap-3">
				<button
					onclick={onCancel}
					class="border border-slate-700 px-4 py-2 text-sm font-bold uppercase text-slate-400 transition-colors hover:bg-slate-800 hover:text-slate-300"
				>
					{cancelText}
				</button>
				<button
					onclick={onConfirm}
					class="px-4 py-2 text-sm font-bold uppercase text-white transition-colors {variant ===
					'danger'
						? 'bg-red-600 hover:bg-red-700'
						: variant === 'warning'
							? 'bg-yellow-600 hover:bg-yellow-700'
							: 'bg-emerald-600 hover:bg-emerald-700'}"
				>
					{confirmText}
				</button>
			</div>
		</div>
	</div>
{/if}
