<script lang="ts">
	import Modal from './Modal.svelte';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
		onSubmit: (account: {
			name: string;
			broker: string;
			accountNumber: string;
			accountType: 'demo' | 'live';
			currency: string;
			isActive: boolean;
		}) => void;
	}

	let { isOpen, onClose, onSubmit }: Props = $props();

	let name = $state('');
	let broker = $state('');
	let accountNumber = $state('');
	let accountType = $state<'demo' | 'live'>('demo');
	let currency = $state('USD');

	function handleSubmit(e: Event) {
		e.preventDefault();

		onSubmit({
			name,
			broker,
			accountNumber,
			accountType,
			currency,
			isActive: true
		});

		// Reset form
		name = '';
		broker = '';
		accountNumber = '';
		accountType = 'demo';
		currency = 'USD';
	}
</script>

<Modal {isOpen} title="Add Trading Account" size="md" {onClose}>
	{#snippet children()}
		<form onsubmit={handleSubmit} id="add-account-form" class="space-y-6">
				<div class="grid gap-6 md:grid-cols-2">
					<!-- Account Name -->
					<div>
						<label for="name" class="mb-2 block text-xs font-bold uppercase text-slate-400">
							Account Name
						</label>
						<input
							type="text"
							id="name"
							bind:value={name}
							required
							class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:border-emerald-500 focus:outline-none"
							placeholder="e.g., Main Demo Account"
						/>
					</div>

					<!-- Broker -->
					<div>
						<label for="broker" class="mb-2 block text-xs font-bold uppercase text-slate-400">
							Broker
						</label>
						<input
							type="text"
							id="broker"
							bind:value={broker}
							required
							class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:border-emerald-500 focus:outline-none"
							placeholder="e.g., OANDA, IC Markets"
						/>
					</div>

					<!-- Account Number -->
					<div>
						<label
							for="accountNumber"
							class="mb-2 block text-xs font-bold uppercase text-slate-400"
						>
							Account Number
						</label>
						<input
							type="text"
							id="accountNumber"
							bind:value={accountNumber}
							required
							class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:border-emerald-500 focus:outline-none"
							placeholder="e.g., 12345678"
						/>
					</div>

					<!-- Account Type -->
					<div>
						<label for="accountType" class="mb-2 block text-xs font-bold uppercase text-slate-400">
							Account Type
						</label>
						<select
							id="accountType"
							bind:value={accountType}
							class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors focus:border-emerald-500 focus:outline-none"
						>
							<option value="demo">Demo</option>
							<option value="live">Live</option>
						</select>
					</div>

					<!-- Currency -->
					<div>
						<label for="currency" class="mb-2 block text-xs font-bold uppercase text-slate-400">
							Currency
						</label>
						<select
							id="currency"
							bind:value={currency}
							class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors focus:border-emerald-500 focus:outline-none"
						>
							<option value="USD">USD</option>
							<option value="EUR">EUR</option>
							<option value="GBP">GBP</option>
							<option value="JPY">JPY</option>
							<option value="AUD">AUD</option>
							<option value="CAD">CAD</option>
						</select>
					</div>
				</div>
			</form>
	{/snippet}

	{#snippet footer()}
		<div class="flex justify-end gap-3">
			<button
				type="button"
				onclick={onClose}
				class="border border-slate-700 px-6 py-3 text-sm font-bold uppercase text-slate-400 transition-colors hover:bg-slate-800 hover:text-slate-300"
			>
				Cancel
			</button>
			<button
				type="submit"
				form="add-account-form"
				class="bg-emerald-600 px-6 py-3 text-sm font-bold uppercase text-white transition-colors hover:bg-emerald-700"
			>
				Add Account
			</button>
		</div>
	{/snippet}
</Modal>
