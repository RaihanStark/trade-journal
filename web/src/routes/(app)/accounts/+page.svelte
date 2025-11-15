<script lang="ts">
	import AddAccountModal from '$lib/components/AddAccountModal.svelte';
	import EditAccountModal from '$lib/components/EditAccountModal.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';

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

	let accounts = $state<Account[]>([
		{
			id: 1,
			name: 'Main Demo Account',
			broker: 'OANDA',
			accountNumber: 'DEMO-12345',
			accountType: 'demo',
			currency: 'USD',
			isActive: true
		},
		{
			id: 2,
			name: 'Live Trading Account',
			broker: 'IC Markets',
			accountNumber: 'LIVE-67890',
			accountType: 'live',
			currency: 'USD',
			isActive: true
		},
		{
			id: 3,
			name: 'Secondary Demo',
			broker: 'XM',
			accountNumber: 'DEMO-54321',
			accountType: 'demo',
			currency: 'EUR',
			isActive: false
		}
	]);

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

	function handleAddAccount(account: Omit<Account, 'id'>) {
		const newAccount: Account = {
			...account,
			id: accounts.length + 1
		};
		accounts = [...accounts, newAccount];
		closeAddModal();
	}

	function handleEditAccount(updatedAccount: Account) {
		accounts = accounts.map((acc) =>
			acc.id === updatedAccount.id ? updatedAccount : acc
		);
		closeEditModal();
	}

	function toggleAccountStatus(id: number) {
		accounts = accounts.map((acc) =>
			acc.id === id ? { ...acc, isActive: !acc.isActive } : acc
		);
	}

	function openDeleteConfirm(id: number) {
		accountToDelete = id;
		isConfirmOpen = true;
	}

	function handleDeleteConfirm() {
		if (accountToDelete !== null) {
			accounts = accounts.filter((acc) => acc.id !== accountToDelete);
			accountToDelete = null;
		}
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
		{#if accounts.length === 0}
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
