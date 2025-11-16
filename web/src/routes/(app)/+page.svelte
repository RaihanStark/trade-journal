<script lang="ts">
	import MetricCard from '$lib/components/MetricCard.svelte';
	import TradesTable from '$lib/components/TradesTable.svelte';
	import StatsGrid from '$lib/components/StatsGrid.svelte';
	import AddTradeModal from '$lib/components/AddTradeModal.svelte';
	import DepositModal from '$lib/components/DepositModal.svelte';
	import WithdrawModal from '$lib/components/WithdrawModal.svelte';
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import { apiClient, type Trade } from '$lib/api/client';
	import { authStore } from '$lib/stores/auth.svelte';
	import { accountsStore } from '$lib/stores/accounts.svelte';

	interface Props {
		data: {
			trades: Promise<Trade[]>;
		};
	}

	let { data }: Props = $props();

	let isModalOpen = $state(false);
	let isDepositModalOpen = $state(false);
	let isWithdrawModalOpen = $state(false);
	let isDeleteConfirmOpen = $state(false);
	let tradeToDelete = $state<number | null>(null);

	// Filters
	let selectedAccount = $state('all');
	let startDate = $state('');
	let endDate = $state('');

	async function reloadData() {
		// Reload accounts store and trades
		if (!authStore.token) return;

		await accountsStore.reload();

		const tradesPromise = apiClient.getTrades(authStore.token).then(({ data: tradesData, error }) => {
			if (error) {
				console.error('Failed to load trades:', error);
				return [];
			}
			return tradesData || [];
		});

		// Update the trades promise
		data = {
			trades: tradesPromise
		};
	}

	function openModal() {
		isModalOpen = true;
	}

	function closeModal() {
		isModalOpen = false;
	}

	function openDepositModal() {
		isDepositModalOpen = true;
	}

	function closeDepositModal() {
		isDepositModalOpen = false;
	}

	function openWithdrawModal() {
		isWithdrawModalOpen = true;
	}

	function closeWithdrawModal() {
		isWithdrawModalOpen = false;
	}

	function clearFilters() {
		selectedAccount = 'all';
		startDate = '';
		endDate = '';
	}

	function openDeleteConfirm(id: number) {
		tradeToDelete = id;
		isDeleteConfirmOpen = true;
	}

	function closeDeleteConfirm() {
		tradeToDelete = null;
		isDeleteConfirmOpen = false;
	}

	async function handleDeleteTrade() {
		if (!authStore.token || tradeToDelete === null) return;

		const { error } = await apiClient.deleteTrade(tradeToDelete, authStore.token);
		if (error) {
			console.error('Failed to delete trade:', error);
			return;
		}

		// Reload accounts and trades after delete (to update balance for deposit/withdraw deletions)
		await reloadData();

		// Close the confirmation dialog
		closeDeleteConfirm();
	}

	const metrics = {
		totalPL: 15420.5,
		winRate: 64.5,
		totalTrades: 128,
		avgWin: 240.3,
		avgLoss: -180.5,
		profitFactor: 1.85,
		sharpeRatio: 2.14,
		maxDrawdown: -2340.2
	};

</script>

<div class="grid h-full grid-cols-12 grid-rows-[auto_auto_auto_1fr] bg-slate-950">
	<!-- Header -->
	<div
		class="col-span-12 flex items-center justify-between border-b border-slate-800 bg-slate-900 px-4 py-2"
	>
		<div class="flex items-center gap-4">
			<span class="font-mono text-xs text-slate-500">BAL:</span>
			<span class="font-mono text-sm font-bold text-emerald-400">${accountsStore.accounts.reduce((sum, account) => sum + account.current_balance, 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
			<span class="font-mono text-xs text-slate-500">EQ:</span>
			<span class="font-mono text-sm font-bold text-blue-400">${accountsStore.accounts.reduce((sum, account) => sum + account.current_balance, 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
		</div>
		<div class="flex items-center gap-2">
			<button onclick={openDepositModal} class="bg-emerald-800 px-3 py-1.5 text-xs text-emerald-100 hover:bg-emerald-700">+ DEPOSIT</button>
			<button onclick={openWithdrawModal} class="bg-red-800 px-3 py-1.5 text-xs text-red-100 hover:bg-red-700">- WITHDRAW</button>
			<button onclick={openModal} class="bg-slate-800 px-3 py-1.5 text-xs text-slate-300 hover:bg-slate-700">+ TRADE</button>
		</div>
	</div>

	<!-- Filters -->
	<div class="col-span-12 border-b border-slate-800 bg-slate-900/50 px-4 py-2">
		<div class="flex items-center gap-4">
			<div class="flex items-center gap-2">
				<label for="account-filter" class="text-xs font-medium text-slate-400">Account:</label>
				<select
					id="account-filter"
					bind:value={selectedAccount}
					class="border border-slate-700 bg-slate-800 px-2 py-1 text-xs text-slate-300 focus:border-slate-600 focus:outline-none"
				>
					<option value="all">All Accounts</option>
					{#each accountsStore.accounts as account}
						<option value={account.id}>{account.name} - {account.broker}</option>
					{/each}
				</select>
			</div>

			<div class="flex items-center gap-2">
				<label for="start-date" class="text-xs font-medium text-slate-400">From:</label>
				<input
					id="start-date"
					type="date"
					bind:value={startDate}
					class="border border-slate-700 bg-slate-800 px-2 py-1 text-xs text-slate-300 focus:border-slate-600 focus:outline-none"
				/>
			</div>

			<div class="flex items-center gap-2">
				<label for="end-date" class="text-xs font-medium text-slate-400">To:</label>
				<input
					id="end-date"
					type="date"
					bind:value={endDate}
					class="border border-slate-700 bg-slate-800 px-2 py-1 text-xs text-slate-300 focus:border-slate-600 focus:outline-none"
				/>
			</div>

			{#if selectedAccount !== 'all' || startDate || endDate}
				<button
					onclick={clearFilters}
					class="text-xs text-slate-500 hover:text-slate-300"
				>
					Clear Filters
				</button>
			{/if}
		</div>
	</div>

	<!-- Metrics -->
	<div class="col-span-12 grid grid-cols-8 border-b border-slate-800">
		<MetricCard title="P/L" value="$15,420" change="+12.4%" trend="up" />
		<MetricCard title="WIN%" value="64.5%" change="+2.3%" trend="up" />
		<MetricCard title="TRADES" value="128" change="+8" trend="neutral" />
		<MetricCard title="PF" value="1.85" change="+0.15" trend="up" />
		<MetricCard title="AVG W" value="$240" change="+5.2%" trend="up" />
		<MetricCard title="AVG L" value="$181" change="-3.1%" trend="down" />
		<MetricCard title="SHARPE" value="2.14" change="+0.08" trend="up" />
		<MetricCard title="DD" value="$2,340" change="-2.1%" trend="down" />
	</div>

	<!-- Main Content -->
	<div class="col-span-10 row-span-1 border-r border-slate-800">
		{#await data.trades}
			<div class="flex h-64 items-center justify-center">
				<p class="text-slate-500">Loading trades...</p>
			</div>
		{:then trades}
			<TradesTable {trades} onDelete={openDeleteConfirm} onUpdate={reloadData} />
		{/await}
	</div>

	<div class="col-span-2 row-span-1">
		<StatsGrid {metrics} />
	</div>
</div>

<AddTradeModal isOpen={isModalOpen} onClose={closeModal} onSuccess={reloadData} />
<DepositModal isOpen={isDepositModalOpen} onClose={closeDepositModal} onSuccess={reloadData} />
<WithdrawModal isOpen={isWithdrawModalOpen} onClose={closeWithdrawModal} onSuccess={reloadData} />
<ConfirmDialog
	isOpen={isDeleteConfirmOpen}
	title="Delete Trade"
	message="Are you sure you want to delete this trade? This action cannot be undone."
	confirmText="Delete"
	cancelText="Cancel"
	variant="danger"
	onConfirm={handleDeleteTrade}
	onCancel={closeDeleteConfirm}
/>
