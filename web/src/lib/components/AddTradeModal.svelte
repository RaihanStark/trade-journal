<script lang="ts">
	import Modal from './Modal.svelte';
	import TagInput from './TagInput.svelte';
	import { apiClient } from '$lib/api/client';
	import { authStore } from '$lib/stores/auth.svelte';
	import { accountsStore } from '$lib/stores/accounts.svelte';
	import { strategiesStore } from '$lib/stores/strategies.svelte';
	import { z } from 'zod';

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
	let entry = $state(0);
	let exit = $state(0);
	let lots = $state(0);
	let stopLoss = $state(0);
	let takeProfit = $state(0);
	let selectedStrategyNames = $state<string[]>([]);
	let notes = $state('');
	let isSubmitting = $state(false);
	let errors = $state<Record<string, string>>({});

	// Zod validation schema
	const tradeSchema = z.object({
		accountId: z.string()
			.min(1, 'Please select an account')
			.refine((val) => !isNaN(parseInt(val)), 'Please select an account'),
		date: z.string().min(1, 'Date is required'),
		time: z.string().min(1, 'Time is required'),
		pair: z.string().min(1, 'Currency pair is required'),
		type: z.enum(['BUY', 'SELL'], { errorMap: () => ({ message: 'Type must be BUY or SELL' }) }),
		entry: z.number()
			.refine((val) => val > 0, 'Entry price must be greater than 0'),
		exit: z.number().optional(),
		lots: z.number()
			.refine((val) => val > 0, 'Lot size must be greater than 0'),
		stopLoss: z.number().optional(),
		takeProfit: z.number().optional(),
		notes: z.string().optional()
	});

	async function handleSubmit() {
		if (!authStore.token || isSubmitting) return;

		// Clear previous errors
		errors = {};

		// Validate form data
		const result = tradeSchema.safeParse({
			accountId,
			date,
			time,
			pair,
			type,
			entry,
			exit: exit || undefined,
			lots,
			stopLoss: stopLoss || undefined,
			takeProfit: takeProfit || undefined,
			notes
		});

		if (!result.success) {
			// Map Zod errors to field errors
			result.error.issues.forEach((err) => {
				if (err.path[0]) {
					errors[err.path[0] as string] = err.message;
				}
			});
			return;
		}

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
			entry,
			exit: exit || null,
			lots,
			stop_loss: stopLoss || null,
			take_profit: takeProfit || null,
			notes,
			mistakes: '',
			amount: null,
			strategy_ids: strategyIds
		}, authStore.token);

		isSubmitting = false;

		if (error) {
			console.error('Failed to create trade:', error);
			errors.submit = error;
			return;
		}

		resetForm();
		onClose();
		if (onSuccess) {
			onSuccess();
		}
	}

	function validateField(field: string, value: any) {
		// Validate single field
		try {
			const fieldSchema = tradeSchema.shape[field as keyof typeof tradeSchema.shape];
			if (fieldSchema) {
				fieldSchema.parse(value);
				// Clear error if validation passes
				const { [field]: _, ...rest } = errors;
				errors = rest;
			}
		} catch (err) {
			if (err instanceof z.ZodError) {
				errors[field] = err.issues[0]?.message || 'Invalid value';
			}
		}
	}

	function resetForm() {
		accountId = '';
		date = new Date().toISOString().split('T')[0];
		time = new Date().toTimeString().split(' ')[0].substring(0, 5);
		pair = 'EUR/USD';
		type = 'BUY';
		entry = 0;
		exit = 0;
		lots = 0;
		stopLoss = 0;
		takeProfit = 0;
		selectedStrategyNames = [];
		notes = '';
		errors = {};
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
			{#if errors.submit}
				<div class="rounded border border-red-800 bg-red-900/20 px-4 py-3 text-sm text-red-400">
					{errors.submit}
				</div>
			{/if}

			<!-- Account Selection -->
			<div>
				<label for="account" class="block text-sm font-medium text-slate-300">
					Trading Account <span class="text-red-400">*</span>
				</label>
				<select
					id="account"
					bind:value={accountId}
					onblur={() => validateField('accountId', accountId)}
					class="mt-1 w-full border px-3 py-2 text-sm text-slate-100 focus:outline-none {errors.accountId
						? 'border-red-500 bg-red-900/10 focus:border-red-500'
						: 'border-slate-700 bg-slate-800 focus:border-emerald-500'}"
				>
					<option value="">Select account</option>
					{#each accountsStore.accounts as account}
						<option value={account.id.toString()}>{account.name} - {account.broker}</option>
					{/each}
				</select>
				{#if errors.accountId}
					<p class="mt-1 text-xs text-red-400">{errors.accountId}</p>
				{/if}
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
						onblur={() => validateField('date', date)}
						class="mt-1 w-full border px-3 py-2 text-sm text-slate-100 focus:outline-none {errors.date
							? 'border-red-500 bg-red-900/10 focus:border-red-500'
							: 'border-slate-700 bg-slate-800 focus:border-emerald-500'}"
					/>
					{#if errors.date}
						<p class="mt-1 text-xs text-red-400">{errors.date}</p>
					{/if}
				</div>

				<div>
					<label for="time" class="block text-sm font-medium text-slate-300">
						Time <span class="text-red-400">*</span>
					</label>
					<input
						id="time"
						type="time"
						bind:value={time}
						onblur={() => validateField('time', time)}
						class="mt-1 w-full border px-3 py-2 text-sm text-slate-100 focus:outline-none {errors.time
							? 'border-red-500 bg-red-900/10 focus:border-red-500'
							: 'border-slate-700 bg-slate-800 focus:border-emerald-500'}"
					/>
					{#if errors.time}
						<p class="mt-1 text-xs text-red-400">{errors.time}</p>
					{/if}
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
						onblur={() => validateField('entry', entry)}
						placeholder="0.00000"
						class="mt-1 w-full border px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:outline-none {errors.entry
							? 'border-red-500 bg-red-900/10 focus:border-red-500'
							: 'border-slate-700 bg-slate-800 focus:border-emerald-500'}"
					/>
					{#if errors.entry}
						<p class="mt-1 text-xs text-red-400">{errors.entry}</p>
					{/if}
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
						onblur={() => validateField('lots', lots)}
						placeholder="0.00"
						class="mt-1 w-full border px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:outline-none {errors.lots
							? 'border-red-500 bg-red-900/10 focus:border-red-500'
							: 'border-slate-700 bg-slate-800 focus:border-emerald-500'}"
					/>
					{#if errors.lots}
						<p class="mt-1 text-xs text-red-400">{errors.lots}</p>
					{/if}
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
