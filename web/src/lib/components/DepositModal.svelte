<script lang="ts">
	import Modal from './Modal.svelte';
	import { apiClient } from '$lib/api/client';
	import { authStore } from '$lib/stores/auth.svelte';
	import { accountsStore } from '$lib/stores/accounts.svelte';
	import { z } from 'zod';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
		onSuccess?: () => void;
	}

	let { isOpen, onClose, onSuccess }: Props = $props();

	let amount = $state(0);
	let accountId = $state('');
	let date = $state(new Date().toISOString().split('T')[0]);
	let notes = $state('');
	let isSubmitting = $state(false);
	let errors = $state<Record<string, string>>({});

	// Zod validation schema
	const depositSchema = z.object({
		accountId: z.string()
			.min(1, 'Please select an account')
			.refine((val) => !isNaN(parseInt(val)), 'Please select an account'),
		amount: z.number()
			.min(1, 'Amount is required')
			.refine((val) => val > 0, 'Amount must be greater than 0'),
		date: z.string().min(1, 'Date is required'),
		notes: z.string().optional()
	});

	async function handleSubmit() {
		if (!authStore.token) return;

		// Clear previous errors
		errors = {};

		// Validate form data
		const result = depositSchema.safeParse({
			accountId,
			amount,
			date,
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

		const { data, error } = await apiClient.createTrade({
			account_id: parseInt(accountId),
			date,
			time: new Date().toTimeString().split(' ')[0].substring(0, 5),
			pair: '',
			type: 'DEPOSIT',
			entry: 0,
			exit: null,
			lots: 0,
			stop_loss: null,
			take_profit: null,
			notes,
			mistakes: '',
			amount: amount,
			strategy_ids: []
		}, authStore.token);

		isSubmitting = false;

		if (error) {
			console.error('Failed to create deposit:', error);
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
			const fieldSchema = depositSchema.shape[field as keyof typeof depositSchema.shape];
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
		amount = 0;
		accountId = '';
		date = new Date().toISOString().split('T')[0];
		notes = '';
		errors = {};
	}
</script>

<Modal {isOpen} title="Deposit Funds" size="md" {onClose}>
	{#snippet children()}
		<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
			{#if errors.submit}
				<div class="rounded border border-red-800 bg-red-900/20 px-4 py-3 text-sm text-red-400">
					{errors.submit}
				</div>
			{/if}

			<div>
				<label for="account" class="block text-sm font-medium text-slate-300">
					Trading Account
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

			<div>
				<label for="amount" class="block text-sm font-medium text-slate-300">
					Amount
				</label>
				<input
					id="amount"
					type="number"
					step="0.01"
					bind:value={amount}
					onblur={() => validateField('amount', amount)}
					placeholder="0.00"
					class="mt-1 w-full border px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:outline-none {errors.amount
						? 'border-red-500 bg-red-900/10 focus:border-red-500'
						: 'border-slate-700 bg-slate-800 focus:border-emerald-500'}"
				/>
				{#if errors.amount}
					<p class="mt-1 text-xs text-red-400">{errors.amount}</p>
				{/if}
			</div>

			<div>
				<label for="date" class="block text-sm font-medium text-slate-300">
					Date
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
				<label for="notes" class="block text-sm font-medium text-slate-300">
					Notes (Optional)
				</label>
				<textarea
					id="notes"
					bind:value={notes}
					rows="3"
					placeholder="Add notes about this deposit..."
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
				class="bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
			>
				{isSubmitting ? 'Processing...' : 'Deposit'}
			</button>
		</div>
	{/snippet}
</Modal>
