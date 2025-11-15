<script lang="ts">
	import { onMount } from 'svelte';
	import AddAccountModal from '$lib/components/AddAccountModal.svelte';
	import EditAccountModal from '$lib/components/EditAccountModal.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import { apiClient } from '$lib/api/client';
	import { authStore } from '$lib/stores/auth.svelte';

	interface Account {
		id: number;
		name: string;
		broker: string;
		accountNumber: string;
		accountType: 'demo' | 'live';
		currency: string;
		isActive: boolean;
	}

	let isAddModalOpen = $state(false);
	let isEditModalOpen = $state(false);
	let isConfirmOpen = $state(false);
	let accountToDelete = $state<number | null>(null);
	let accountToEdit = $state<Account | null>(null);
	let isLoading = $state(true);

	let accounts = $state<Account[]>([]);

	onMount(async () => {
		await loadAccounts();
	});

	async function loadAccounts() {
		if (!authStore.token) return;

		isLoading = true;
		const { data, error } = await apiClient.getAccounts(authStore.token);
		isLoading = false;

		if (error) {
			console.error('Failed to load accounts:', error);
			return;
		}

		if (data) {
			accounts = data.map((acc) => ({
				id: acc.id,
				name: acc.name,
				broker: acc.broker,
				accountNumber: acc.account_number,
				accountType: acc.account_type,
				currency: acc.currency,
				isActive: acc.is_active
			}));
		}
	}

	function openAddModal() {
		isAddModalOpen = true;
	}

	function closeAddModal() {
		isAddModalOpen = false;
	}

	function openEditModal(account: Account) {
		accountToEdit = account;
		isEditModalOpen = true;
	}

	function closeEditModal() {
		accountToEdit = null;
		isEditModalOpen = false;
	}

	async function handleAddAccount(account: Omit<Account, 'id'>) {
		if (!authStore.token) return;

		const { data, error } = await apiClient.createAccount(
			{
				name: account.name,
				broker: account.broker,
				account_number: account.accountNumber,
				account_type: account.accountType,
				currency: account.currency,
				is_active: account.isActive
			},
			authStore.token
		);

		if (error) {
			console.error('Failed to create account:', error);
			return;
		}

		if (data) {
			accounts = [
				...accounts,
				{
					id: data.id,
					name: data.name,
					broker: data.broker,
					accountNumber: data.account_number,
					accountType: data.account_type,
					currency: data.currency,
					isActive: data.is_active
				}
			];
		}

		closeAddModal();
	}

	async function handleEditAccount(updatedAccount: Account) {
		if (!authStore.token) return;

		const { data, error } = await apiClient.updateAccount(
			updatedAccount.id,
			{
				name: updatedAccount.name,
				broker: updatedAccount.broker,
				account_number: updatedAccount.accountNumber,
				account_type: updatedAccount.accountType,
				currency: updatedAccount.currency,
				is_active: updatedAccount.isActive
			},
			authStore.token
		);

		if (error) {
			console.error('Failed to update account:', error);
			return;
		}

		if (data) {
			accounts = accounts.map((acc) =>
				acc.id === data.id
					? {
							id: data.id,
							name: data.name,
							broker: data.broker,
							accountNumber: data.account_number,
							accountType: data.account_type,
							currency: data.currency,
							isActive: data.is_active
						}
					: acc
			);
		}

		closeEditModal();
	}

	async function toggleAccountStatus(id: number) {
		const account = accounts.find((acc) => acc.id === id);
		if (!account || !authStore.token) return;

		const updatedAccount = { ...account, isActive: !account.isActive };

		const { data, error } = await apiClient.updateAccount(
			id,
			{
				name: updatedAccount.name,
				broker: updatedAccount.broker,
				account_number: updatedAccount.accountNumber,
				account_type: updatedAccount.accountType,
				currency: updatedAccount.currency,
				is_active: updatedAccount.isActive
			},
			authStore.token
		);

		if (error) {
			console.error('Failed to toggle account status:', error);
			return;
		}

		if (data) {
			accounts = accounts.map((acc) =>
				acc.id === id ? { ...acc, isActive: data.is_active } : acc
			);
		}
	}

	function openDeleteConfirm(id: number) {
		accountToDelete = id;
		isConfirmOpen = true;
	}

	async function handleDeleteConfirm() {
		if (accountToDelete === null || !authStore.token) return;

		const { error } = await apiClient.deleteAccount(accountToDelete, authStore.token);

		if (error) {
			console.error('Failed to delete account:', error);
			return;
		}

		accounts = accounts.filter((acc) => acc.id !== accountToDelete);
		accountToDelete = null;
		isConfirmOpen = false;
	}

	function handleDeleteCancel() {
		accountToDelete = null;
		isConfirmOpen = false;
	}
</script>

<div class="flex h-full flex-col bg-slate-950">
	<!-- Header -->
	<div class="border-b border-slate-800 bg-slate-900 px-6 py-4">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-100">Trading Accounts</h1>
				<p class="mt-1 text-xs text-slate-500">Manage your demo and live trading accounts</p>
			</div>
			<button
				onclick={openAddModal}
				class="bg-emerald-600 px-4 py-2 text-sm font-bold uppercase text-white transition-colors hover:bg-emerald-700"
			>
				+ New Account
			</button>
		</div>
	</div>

	<!-- Accounts Table -->
	<div class="flex-1 overflow-auto">
		{#if isLoading}
			<div class="flex h-64 items-center justify-center">
				<div class="text-center">
					<p class="text-slate-500">Loading accounts...</p>
				</div>
			</div>
		{:else if accounts.length === 0}
			<div class="flex h-64 items-center justify-center">
				<div class="text-center">
					<p class="text-slate-500">No trading accounts yet</p>
					<button
						onclick={openAddModal}
						class="mt-4 text-sm font-medium text-emerald-400 hover:text-emerald-300"
					>
						Create your first account
					</button>
				</div>
			</div>
		{:else}
			<table class="w-full">
				<thead class="border-b border-slate-800 bg-slate-900">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-bold uppercase text-slate-500">Name</th>
						<th class="px-6 py-3 text-left text-xs font-bold uppercase text-slate-500">Broker</th>
						<th class="px-6 py-3 text-left text-xs font-bold uppercase text-slate-500">
							Account Number
						</th>
						<th class="px-6 py-3 text-left text-xs font-bold uppercase text-slate-500">Type</th>
						<th class="px-6 py-3 text-left text-xs font-bold uppercase text-slate-500">
							Currency
						</th>
						<th class="px-6 py-3 text-left text-xs font-bold uppercase text-slate-500">Status</th>
						<th class="px-6 py-3 text-right text-xs font-bold uppercase text-slate-500">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-800/50 bg-slate-950">
					{#each accounts as account}
						<tr class="transition-colors hover:bg-slate-900/50 {!account.isActive ? 'opacity-50' : ''}">
							<td class="px-6 py-4">
								<span class="font-medium text-slate-100">{account.name}</span>
							</td>
							<td class="px-6 py-4">
								<span class="text-sm text-slate-400">{account.broker}</span>
							</td>
							<td class="px-6 py-4">
								<span class="font-mono text-sm text-slate-400">{account.accountNumber}</span>
							</td>
							<td class="px-6 py-4">
								<span
									class="inline-block px-2 py-1 text-[10px] font-bold uppercase {account.accountType ===
									'live'
										? 'bg-emerald-900/30 text-emerald-400'
										: 'bg-blue-900/30 text-blue-400'}"
								>
									{account.accountType}
								</span>
							</td>
							<td class="px-6 py-4">
								<span class="text-sm text-slate-400">{account.currency}</span>
							</td>
							<td class="px-6 py-4">
								<button
									onclick={() => toggleAccountStatus(account.id)}
									class="text-xs font-medium uppercase transition-colors cursor-pointer {account.isActive
										? 'text-emerald-400 hover:text-emerald-300'
										: 'text-slate-600 hover:text-slate-500'}"
								>
									{account.isActive ? 'Active' : 'Inactive'}
								</button>
							</td>
							<td class="px-6 py-4">
								<div class="flex justify-end gap-2">
									<button
										onclick={() => openEditModal(account)}
										class="text-slate-400 transition-colors hover:text-slate-200"
										aria-label="Edit account"
									>
										<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
											></path>
										</svg>
									</button>
									<button
										onclick={() => openDeleteConfirm(account.id)}
										class="text-red-400 transition-colors hover:text-red-300"
										aria-label="Delete account"
									>
										<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
											></path>
										</svg>
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{/if}
	</div>
</div>

<AddAccountModal isOpen={isAddModalOpen} onClose={closeAddModal} onSubmit={handleAddAccount} />

<EditAccountModal
	isOpen={isEditModalOpen}
	account={accountToEdit}
	onClose={closeEditModal}
	onSubmit={handleEditAccount}
/>

<ConfirmDialog
	isOpen={isConfirmOpen}
	title="Delete Account"
	message="Are you sure you want to delete this account? This action cannot be undone."
	confirmText="Delete"
	cancelText="Cancel"
	variant="danger"
	onConfirm={handleDeleteConfirm}
	onCancel={handleDeleteCancel}
/>
