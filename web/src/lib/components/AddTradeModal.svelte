<script lang="ts">
	import { fade, scale } from 'svelte/transition';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
	}

	let { isOpen, onClose }: Props = $props();

	let formData = $state({
		date: new Date().toISOString().split('T')[0],
		time: new Date().toTimeString().split(' ')[0].substring(0, 5),
		pair: 'EUR/USD',
		type: 'BUY' as 'BUY' | 'SELL',
		entry: '',
		exit: '',
		lots: '',
		stopLoss: '',
		takeProfit: '',
		notes: '',
		mistakes: ''
	});

	function handleSubmit(e: Event) {
		e.preventDefault();
		console.log('Trade submitted:', formData);
		// TODO: Add trade to database/state
		onClose();
		resetForm();
	}

	function resetForm() {
		formData = {
			date: new Date().toISOString().split('T')[0],
			time: new Date().toTimeString().split(' ')[0].substring(0, 5),
			pair: 'EUR/USD',
			type: 'BUY',
			entry: '',
			exit: '',
			lots: '',
			stopLoss: '',
			takeProfit: '',
			notes: '',
			mistakes: ''
		};
	}

	const pairs = [
		'EUR/USD',
		'GBP/USD',
		'USD/JPY',
		'USD/CHF',
		'AUD/USD',
		'USD/CAD',
		'NZD/USD',
		'EUR/GBP',
		'EUR/JPY',
		'GBP/JPY'
	];
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/80 backdrop-blur-sm"
		onclick={onClose}
		onkeydown={(e) => e.key === 'Escape' && onClose()}
		role="button"
		tabindex="-1"
		transition:fade={{ duration: 200 }}
	>
		<div
			class="w-full max-w-2xl border border-slate-800 bg-slate-900 shadow-2xl"
			onclick={(e) => e.stopPropagation()}
			transition:scale={{ duration: 200, start: 0.95 }}
			onkeydown={(e) => e.stopPropagation()}
			role="dialog"
			tabindex="-1"
		>
			<!-- Modal Header -->
			<div class="flex items-center justify-between border-b border-slate-800 bg-slate-950 px-4 py-3">
				<h2 class="text-sm font-bold uppercase tracking-wider text-slate-100">ADD NEW TRADE</h2>
				<button
					onclick={onClose}
					class="text-slate-400 transition hover:text-slate-200"
					type="button"
					aria-label="Close modal"
				>
					<svg
						class="h-5 w-5"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						></path>
					</svg>
				</button>
			</div>

			<!-- Modal Body -->
			<form onsubmit={handleSubmit} class="p-4">
				<div class="grid grid-cols-2 gap-4">
					<!-- Date -->
					<div>
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="date"
							>DATE</label
						>
						<input
							type="date"
							id="date"
							bind:value={formData.date}
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							required
						/>
					</div>

					<!-- Time -->
					<div>
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="time"
							>TIME</label
						>
						<input
							type="time"
							id="time"
							bind:value={formData.time}
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							required
						/>
					</div>

					<!-- Pair -->
					<div>
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="pair"
							>PAIR</label
						>
						<select
							id="pair"
							bind:value={formData.pair}
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							required
						>
							{#each pairs as pair}
								<option value={pair}>{pair}</option>
							{/each}
						</select>
					</div>

					<!-- Type -->
					<div>
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="type"
							>TYPE</label
						>
						<select
							id="type"
							bind:value={formData.type}
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							required
						>
							<option value="BUY">BUY</option>
							<option value="SELL">SELL</option>
						</select>
					</div>

					<!-- Entry Price -->
					<div>
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="entry"
							>ENTRY PRICE</label
						>
						<input
							type="number"
							id="entry"
							step="0.00001"
							bind:value={formData.entry}
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							placeholder="0.00000"
							required
						/>
					</div>

					<!-- Exit Price -->
					<div>
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="exit"
							>EXIT PRICE</label
						>
						<input
							type="number"
							id="exit"
							step="0.00001"
							bind:value={formData.exit}
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							placeholder="0.00000"
						/>
					</div>

					<!-- Lots -->
					<div>
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="lots"
							>LOTS</label
						>
						<input
							type="number"
							id="lots"
							step="0.01"
							bind:value={formData.lots}
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							placeholder="0.00"
							required
						/>
					</div>

					<!-- Stop Loss -->
					<div>
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="stopLoss"
							>STOP LOSS</label
						>
						<input
							type="number"
							id="stopLoss"
							step="0.00001"
							bind:value={formData.stopLoss}
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							placeholder="0.00000"
						/>
					</div>

					<!-- Take Profit -->
					<div>
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="takeProfit"
							>TAKE PROFIT</label
						>
						<input
							type="number"
							id="takeProfit"
							step="0.00001"
							bind:value={formData.takeProfit}
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							placeholder="0.00000"
						/>
					</div>

					<!-- Notes -->
					<div class="col-span-2">
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="notes"
							>NOTES</label
						>
						<textarea
							id="notes"
							bind:value={formData.notes}
							rows="3"
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							placeholder="Trade notes..."
						></textarea>
					</div>

					<!-- Mistakes -->
					<div class="col-span-2">
						<label class="mb-1 block text-xs font-bold uppercase text-slate-500" for="mistakes"
							>MISTAKES</label
						>
						<textarea
							id="mistakes"
							bind:value={formData.mistakes}
							rows="3"
							class="w-full border border-slate-800 bg-slate-950 px-3 py-2 font-mono text-sm text-slate-100 focus:border-slate-700 focus:outline-none"
							placeholder="What went wrong..."
						></textarea>
					</div>
				</div>

				<!-- Modal Footer -->
				<div class="mt-6 flex justify-end gap-2">
					<button
						type="button"
						onclick={onClose}
						class="border border-slate-700 bg-slate-800 px-4 py-2 text-sm font-medium text-slate-300 transition hover:bg-slate-700"
					>
						CANCEL
					</button>
					<button
						type="submit"
						class="bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-emerald-500"
					>
						ADD TRADE
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
