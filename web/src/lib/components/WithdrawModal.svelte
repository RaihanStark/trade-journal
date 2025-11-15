<script lang="ts">
	import Modal from './Modal.svelte';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
	}

	let { isOpen, onClose }: Props = $props();

	let amount = $state('');
	let accountId = $state('');
	let date = $state(new Date().toISOString().split('T')[0]);
	let notes = $state('');

	function handleSubmit() {
		console.log('Withdraw:', { amount, accountId, date, notes });
		// TODO: Implement withdrawal submission
		resetForm();
		onClose();
	}

	function resetForm() {
		amount = '';
		accountId = '';
		date = new Date().toISOString().split('T')[0];
		notes = '';
	}
</script>

<Modal {isOpen} title="Withdraw Funds" size="md" {onClose}>
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
					class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:border-red-500 focus:outline-none"
				>
					<option value="">Select account</option>
					<!-- TODO: Load accounts from API -->
					<option value="1">Demo Account - XM</option>
					<option value="2">Live Account - IC Markets</option>
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
					class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-red-500 focus:outline-none"
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
					class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:border-red-500 focus:outline-none"
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
					placeholder="Add notes about this withdrawal..."
					class="mt-1 w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-red-500 focus:outline-none"
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
				class="bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700"
			>
				Withdraw
			</button>
		</div>
	{/snippet}
</Modal>
