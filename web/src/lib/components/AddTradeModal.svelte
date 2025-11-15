<script lang="ts">
	import Modal from './Modal.svelte';
	import TagInput from './TagInput.svelte';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
	}

	let { isOpen, onClose }: Props = $props();

	let accountId = $state('');
	let date = $state(new Date().toISOString().split('T')[0]);
	let time = $state(new Date().toTimeString().split(' ')[0].substring(0, 5));
	let pair = $state('EUR/USD');
	let type = $state<'BUY' | 'SELL'>('BUY');
	let entry = $state('');
	let exit = $state('');
	let lots = $state('');
	let stopLoss = $state('');
	let takeProfit = $state('');
	let strategies = $state<string[]>([]);
	let notes = $state('');

	function handleSubmit() {
		console.log('Trade submitted:', {
			accountId,
			date,
			time,
			pair,
			type,
			entry,
			exit,
			lots,
			stopLoss,
			takeProfit,
			strategies,
			notes
		});
		// TODO: Add trade to database/state
		resetForm();
		onClose();
	}

	function resetForm() {
		accountId = '';
		date = new Date().toISOString().split('T')[0];
		time = new Date().toTimeString().split(' ')[0].substring(0, 5);
		pair = 'EUR/USD';
		type = 'BUY';
		entry = '';
		exit = '';
		lots = '';
		stopLoss = '';
		takeProfit = '';
		strategies = [];
		notes = '';
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
		'GBP/JPY',
		'AUD/JPY',
		'NZD/JPY',
		'EUR/CHF',
		'GBP/CHF'
	];

	const suggestedStrategies = [
		'Trend Following',
		'Support/Resistance',
		'Breakout',
		'Reversal',
		'Range Trading',
		'News Trading',
		'Scalping',
		'Swing Trading',
		'Price Action',
		'ICT Concepts',
		'Smart Money Concepts',
		'Supply & Demand',
		'Fibonacci',
		'Moving Average',
		'Candlestick Patterns'
	];
</script>

<Modal {isOpen} title="Add New Trade" size="lg" {onClose}>
	{#snippet children()}
		<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
			<!-- Account Selection -->
			<div>
				<label for="account" class="block text-sm font-medium text-slate-300">
					Trading Account <span class="text-red-400">*</span>
				</label>
				<select
					id="account"
					bind:value={accountId}
					required
					class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:border-emerald-500 focus:outline-none"
				>
					<option value="">Select account</option>
					<!-- TODO: Load accounts from API -->
					<option value="1">Demo Account - XM</option>
					<option value="2">Live Account - IC Markets</option>
				</select>
			</div>

			<div class="grid grid-cols-2 gap-4">
				<!-- Date & Time -->
				<div>
					<label for="date" class="block text-sm font-medium text-slate-300">
						Date <span class="text-red-400">*</span>
					</label>
					<input
						id="date"
						type="date"
						bind:value={date}
						required
						class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:border-emerald-500 focus:outline-none"
					/>
				</div>

				<div>
					<label for="time" class="block text-sm font-medium text-slate-300">
						Time <span class="text-red-400">*</span>
					</label>
					<input
						id="time"
						type="time"
						bind:value={time}
						required
						class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:border-emerald-500 focus:outline-none"
					/>
				</div>

				<!-- Pair & Type -->
				<div>
					<label for="pair" class="block text-sm font-medium text-slate-300">
						Currency Pair <span class="text-red-400">*</span>
					</label>
					<select
						id="pair"
						bind:value={pair}
						required
						class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:border-emerald-500 focus:outline-none"
					>
						{#each pairs as pairOption}
							<option value={pairOption}>{pairOption}</option>
						{/each}
					</select>
				</div>

				<div>
					<label for="type" class="block text-sm font-medium text-slate-300">
						Type <span class="text-red-400">*</span>
					</label>
					<select
						id="type"
						bind:value={type}
						required
						class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:border-emerald-500 focus:outline-none"
					>
						<option value="BUY">BUY</option>
						<option value="SELL">SELL</option>
					</select>
				</div>

				<!-- Entry & Exit Price -->
				<div>
					<label for="entry" class="block text-sm font-medium text-slate-300">
						Entry Price <span class="text-red-400">*</span>
					</label>
					<input
						id="entry"
						type="number"
						step="0.00001"
						bind:value={entry}
						required
						placeholder="0.00000"
						class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-emerald-500 focus:outline-none"
					/>
				</div>

				<div>
					<label for="exit" class="block text-sm font-medium text-slate-300">
						Exit Price
					</label>
					<input
						id="exit"
						type="number"
						step="0.00001"
						bind:value={exit}
						placeholder="0.00000"
						class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-emerald-500 focus:outline-none"
					/>
				</div>

				<!-- Lots -->
				<div class="col-span-2">
					<label for="lots" class="block text-sm font-medium text-slate-300">
						Lot Size <span class="text-red-400">*</span>
					</label>
					<input
						id="lots"
						type="number"
						step="0.01"
						bind:value={lots}
						required
						placeholder="0.00"
						class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-emerald-500 focus:outline-none"
					/>
				</div>
			</div>

			<!-- Strategy Tags -->
			<div>
				<TagInput
					value={strategies}
					suggestions={suggestedStrategies}
					label="Strategies"
					placeholder="Type to add or create strategies..."
					onUpdate={(tags) => { strategies = tags; }}
				/>
			</div>

			<div class="grid grid-cols-2 gap-4">
				<!-- Stop Loss & Take Profit -->
				<div>
					<label for="stopLoss" class="block text-sm font-medium text-slate-300">
						Stop Loss
					</label>
					<input
						id="stopLoss"
						type="number"
						step="0.00001"
						bind:value={stopLoss}
						placeholder="0.00000"
						class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-emerald-500 focus:outline-none"
					/>
				</div>

				<div>
					<label for="takeProfit" class="block text-sm font-medium text-slate-300">
						Take Profit
					</label>
					<input
						id="takeProfit"
						type="number"
						step="0.00001"
						bind:value={takeProfit}
						placeholder="0.00000"
						class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-emerald-500 focus:outline-none"
					/>
				</div>
			</div>

			<!-- Notes -->
			<div>
				<label for="notes" class="block text-sm font-medium text-slate-300">
					Notes
				</label>
				<textarea
					id="notes"
					bind:value={notes}
					rows="3"
					placeholder="Trade analysis, setup details, or any other notes..."
					class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-emerald-500 focus:outline-none"
				></textarea>
			</div>
		</form>
	{/snippet}

	{#snippet footer()}
		<div class="flex justify-end gap-2">
			<button
				onclick={onClose}
				type="button"
				class="px-4 py-2 text-sm font-medium text-slate-400 transition-colors hover:text-slate-200"
			>
				Cancel
			</button>
			<button
				onclick={handleSubmit}
				type="submit"
				class="bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-emerald-700"
			>
				Add Trade
			</button>
		</div>
	{/snippet}
</Modal>
