<script lang="ts">
	import Modal from './Modal.svelte';

	interface Account {
		id: number;
		name: string;
		broker: string;
		accountNumber: string;
		accountType: 'demo' | 'live';
		currency: string;
		isActive: boolean;
	}

	interface Props {
		isOpen: boolean;
		account: Account | null;
		onClose: () => void;
		onSubmit: (account: Account) => void;
	}

	let { isOpen, account, onClose, onSubmit }: Props = $props();

	let name = $state('');
	let broker = $state('');
	let accountNumber = $state('');
	let accountType = $state<'demo' | 'live'>('demo');
	let currency = $state('USD');

	// Update form when account prop changes
	$effect(() => {
		if (account) {
			name = account.name;
			broker = account.broker;
			accountNumber = account.accountNumber;
			accountType = account.accountType;
			currency = account.currency;
		}
	});

	function handleSubmit(e: Event) {
		e.preventDefault();

		if (account) {
			onSubmit({
				...account,
				name,
				broker,
				accountNumber,
				accountType,
				currency
			});
		}
	}
</script>

<Modal {isOpen} title="Edit Trading Account" size="md" {onClose}>
	{#snippet children()}
		<form onsubmit={handleSubmit} id="edit-account-form" class="space-y-6">
			<div class="grid gap-6 md:grid-cols-2">
				<!-- Account Name -->
				<div>
					<label for="edit-name" class="mb-2 block text-xs font-bold uppercase text-slate-400">
						Account Name
					</label>
					<input
						type="text"
						id="edit-name"
						bind:value={name}
						required
						class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:border-emerald-500 focus:outline-none"
						placeholder="e.g., Main Demo Account"
					/>
				</div>

				<!-- Broker -->
				<div>
					<label for="edit-broker" class="mb-2 block text-xs font-bold uppercase text-slate-400">
						Broker
					</label>
					<input
						type="text"
						id="edit-broker"
						bind:value={broker}
						required
						class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:border-emerald-500 focus:outline-none"
						placeholder="e.g., OANDA, IC Markets"
					/>
				</div>

				<!-- Account Number -->
				<div>
					<label
						for="edit-accountNumber"
						class="mb-2 block text-xs font-bold uppercase text-slate-400"
					>
						Account Number
					</label>
					<input
						type="text"
						id="edit-accountNumber"
						bind:value={accountNumber}
						required
						class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:border-emerald-500 focus:outline-none"
						placeholder="e.g., 12345678"
					/>
				</div>

				<!-- Account Type -->
				<div>
					<label
						for="edit-accountType"
						class="mb-2 block text-xs font-bold uppercase text-slate-400"
					>
						Account Type
					</label>
					<select
						id="edit-accountType"
						bind:value={accountType}
						class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors focus:border-emerald-500 focus:outline-none"
					>
						<option value="demo">Demo</option>
						<option value="live">Live</option>
					</select>
				</div>

				<!-- Currency -->
				<div>
					<label for="edit-currency" class="mb-2 block text-xs font-bold uppercase text-slate-400">
						Currency
					</label>
					<select
						id="edit-currency"
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
				form="edit-account-form"
				class="bg-emerald-600 px-6 py-3 text-sm font-bold uppercase text-white transition-colors hover:bg-emerald-700"
			>
				Save Changes
			</button>
		</div>
	{/snippet}
</Modal>
