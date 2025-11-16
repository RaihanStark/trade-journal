<script lang="ts">
	import Modal from './Modal.svelte';
	import { apiClient } from '$lib/api/client';
	import { authStore } from '$lib/stores/auth.svelte';
	import { accountsStore } from '$lib/stores/accounts.svelte';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
		onSuccess?: () => void;
	}

	let { isOpen, onClose, onSuccess }: Props = $props();

	let amount = $state('');
	let accountId = $state('');
	let date = $state(new Date().toISOString().split('T')[0]);
	let notes = $state('');
	let isSubmitting = $state(false);

	async function handleSubmit() {
		if (!authStore.token || !accountId || !amount) return;

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
			amount: parseFloat(amount),
			strategy_ids: []
		}, authStore.token);

		isSubmitting = false;

		if (error) {
			console.error('Failed to create deposit:', error);
			return;
		}

		resetForm();
		onClose();
		if (onSuccess) {
			onSuccess();
		}
	}

	function resetForm() {
		amount = '';
		accountId = '';
		date = new Date().toISOString().split('T')[0];
		notes = '';
	}
</script>

<Modal {isOpen} title="Deposit Funds" size="md" {onClose}>
	{#snippet children()}
		<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
			<div>
				<label for="account" class="block text-sm font-medium text-slate-300">
					Trading Account
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

			<div>
				<label for="amount" class="block text-sm font-medium text-slate-300">
					Amount
				</label>
				<input
					id="amount"
					type="number"
					step="0.01"
					bind:value={amount}
					required
					placeholder="0.00"
					class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-emerald-500 focus:outline-none"
				/>
			</div>

			<div>
				<label for="date" class="block text-sm font-medium text-slate-300">
					Date
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
				class="bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-emerald-700"
			>
				Deposit
			</button>
		</div>
	{/snippet}
</Modal>
