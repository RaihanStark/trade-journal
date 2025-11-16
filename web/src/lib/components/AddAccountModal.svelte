<script lang="ts">
	import Modal from './Modal.svelte';
	import { z } from 'zod';

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
	let errors = $state<Record<string, string>>({});

	// Zod validation schema
	const accountSchema = z.object({
		name: z.string().min(1, 'Account name is required'),
		broker: z.string().min(1, 'Broker is required'),
		accountNumber: z.string().min(1, 'Account number is required'),
		accountType: z.enum(['demo', 'live'], { errorMap: () => ({ message: 'Account type must be demo or live' }) }),
		currency: z.string().min(1, 'Currency is required')
	});

	function handleSubmit(e: Event) {
		e.preventDefault();

		// Clear previous errors
		errors = {};

		// Validate form data
		const result = accountSchema.safeParse({
			name,
			broker,
			accountNumber,
			accountType,
			currency
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

		onSubmit({
			name,
			broker,
			accountNumber,
			accountType,
			currency,
			isActive: true
		});

		// Reset form
		resetForm();
	}

	function validateField(field: string, value: any) {
		// Validate single field
		try {
			const fieldSchema = accountSchema.shape[field as keyof typeof accountSchema.shape];
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
		name = '';
		broker = '';
		accountNumber = '';
		accountType = 'demo';
		currency = 'USD';
		errors = {};
	}
</script>

<Modal {isOpen} title="Add Trading Account" size="md" {onClose}>
	{#snippet children()}
		<form onsubmit={handleSubmit} id="add-account-form" class="space-y-6">
				{#if errors.submit}
					<div class="rounded border border-red-800 bg-red-900/20 px-4 py-3 text-sm text-red-400">
						{errors.submit}
					</div>
				{/if}

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
							onblur={() => validateField('name', name)}
							class="w-full border bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:outline-none {errors.name
								? 'border-red-500 bg-red-900/10 focus:border-red-500'
								: 'border-slate-700 focus:border-emerald-500'}"
							placeholder="e.g., Main Demo Account"
						/>
						{#if errors.name}
							<p class="mt-1 text-xs text-red-400">{errors.name}</p>
						{/if}
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
							onblur={() => validateField('broker', broker)}
							class="w-full border bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:outline-none {errors.broker
								? 'border-red-500 bg-red-900/10 focus:border-red-500'
								: 'border-slate-700 focus:border-emerald-500'}"
							placeholder="e.g., OANDA, IC Markets"
						/>
						{#if errors.broker}
							<p class="mt-1 text-xs text-red-400">{errors.broker}</p>
						{/if}
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
							onblur={() => validateField('accountNumber', accountNumber)}
							class="w-full border bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:outline-none {errors.accountNumber
								? 'border-red-500 bg-red-900/10 focus:border-red-500'
								: 'border-slate-700 focus:border-emerald-500'}"
							placeholder="e.g., 12345678"
						/>
						{#if errors.accountNumber}
							<p class="mt-1 text-xs text-red-400">{errors.accountNumber}</p>
						{/if}
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
