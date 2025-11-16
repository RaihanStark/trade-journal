<script lang="ts">
	import Modal from './Modal.svelte';
	import TagInput from './TagInput.svelte';
	import { apiClient } from '$lib/api/client';
	import { authStore } from '$lib/stores/auth.svelte';
	import { accountsStore } from '$lib/stores/accounts.svelte';
	import { strategiesStore } from '$lib/stores/strategies.svelte';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
		onSuccess?: () => void;
	}

	let { isOpen, onClose, onSuccess }: Props = $props();

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
	let selectedStrategyNames = $state<string[]>([]);
	let notes = $state('');
	let isSubmitting = $state(false);

	async function handleSubmit() {
		if (!authStore.token || isSubmitting) return;

		isSubmitting = true;

		// Find strategy IDs from selected names
		const strategyIds: number[] = [];
		for (const name of selectedStrategyNames) {
			let strategy = strategiesStore.strategies.find(s => s.name === name);
			if (!strategy) {
				// Create new strategy if it doesn't exist
				const { data } = await apiClient.createStrategy({ name, description: '' }, authStore.token);
				if (data) {
					strategy = data;
					await strategiesStore.add(data);
				}
			}
			if (strategy) {
				strategyIds.push(strategy.id);
			}
		}

		const { data, error } = await apiClient.createTrade({
			account_id: accountId ? parseInt(accountId) : null,
			date,
			time,
			pair,
			type,
			entry: parseFloat(entry),
			exit: exit ? parseFloat(exit) : null,
			lots: parseFloat(lots),
			stop_loss: stopLoss ? parseFloat(stopLoss) : null,
			take_profit: takeProfit ? parseFloat(takeProfit) : null,
			notes,
			mistakes: '',
			amount: null,
			strategy_ids: strategyIds
		}, authStore.token);

		isSubmitting = false;

		if (error) {
			console.error('Failed to create trade:', error);
			return;
		}

		resetForm();
		onClose();
		if (onSuccess) {
			onSuccess();
		}
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
		selectedStrategyNames = [];
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

	// Get strategy names for suggestions
	let strategyNames = $derived(strategiesStore.strategies.map(s => s.name));
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
					{#each accountsStore.accounts as account}
						<option value={account.id}>{account.name} - {account.broker}</option>
					{/each}
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
					value={selectedStrategyNames}
					suggestions={strategyNames}
					label="Strategies"
					placeholder="Type to add or create strategies..."
					onUpdate={(tags) => { selectedStrategyNames = tags; }}
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
				disabled={isSubmitting}
				class="bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed"
			>
				{isSubmitting ? 'Adding...' : 'Add Trade'}
			</button>
		</div>
	{/snippet}
</Modal>
