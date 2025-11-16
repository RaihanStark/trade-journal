<script lang="ts">
	import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
	import { apiClient } from '$lib/api/client';
	import { authStore } from '$lib/stores/auth.svelte';
	import { strategiesStore } from '$lib/stores/strategies.svelte';
	import type { Strategy } from '$lib/api/client';
	import { onMount } from 'svelte';

	let isAddMode = $state(false);
	let editingId = $state<number | null>(null);
	let isConfirmOpen = $state(false);
	let strategyToDelete = $state<number | null>(null);

	onMount(async () => {
		await strategiesStore.load();
	});

	let formData = $state({
		name: '',
		description: ''
	});

	function startAdd() {
		isAddMode = true;
		formData = { name: '', description: '' };
	}

	function cancelAdd() {
		isAddMode = false;
		formData = { name: '', description: '' };
	}

	async function saveAdd() {
		if (!formData.name.trim() || !authStore.token) return;

		const { data, error } = await apiClient.createStrategy(
			{
				name: formData.name,
				description: formData.description
			},
			authStore.token
		);

		if (error) {
			console.error('Failed to create strategy:', error);
			return;
		}

		if (data) {
			await strategiesStore.add(data);
		}

		cancelAdd();
	}

	function startEdit(strategy: Strategy) {
		editingId = strategy.id;
		formData = { name: strategy.name, description: strategy.description };
	}

	function cancelEdit() {
		editingId = null;
		formData = { name: '', description: '' };
	}

	async function saveEdit() {
		if (!formData.name.trim() || editingId === null || !authStore.token) return;

		const { data, error } = await apiClient.updateStrategy(
			editingId,
			{
				name: formData.name,
				description: formData.description
			},
			authStore.token
		);

		if (error) {
			console.error('Failed to update strategy:', error);
			return;
		}

		if (data) {
			await strategiesStore.update(data);
		}

		cancelEdit();
	}

	function openDeleteConfirm(id: number) {
		strategyToDelete = id;
		isConfirmOpen = true;
	}

	async function handleDeleteConfirm() {
		if (strategyToDelete === null || !authStore.token) return;

		const { error } = await apiClient.deleteStrategy(strategyToDelete, authStore.token);

		if (error) {
			console.error('Failed to delete strategy:', error);
			return;
		}

		await strategiesStore.remove(strategyToDelete);
		strategyToDelete = null;
		isConfirmOpen = false;
	}

	function handleDeleteCancel() {
		strategyToDelete = null;
		isConfirmOpen = false;
	}
</script>

<div class="flex h-full flex-col bg-slate-950">
	<!-- Header -->
	<div class="border-b border-slate-800 bg-slate-900 px-6 py-4">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-100">Trading Strategies</h1>
				<p class="mt-1 text-xs text-slate-500">Manage your trading strategies and setups</p>
			</div>
			{#if !isAddMode}
				<button
					onclick={startAdd}
					class="bg-emerald-600 px-4 py-2 text-sm font-bold uppercase text-white transition-colors hover:bg-emerald-700"
				>
					+ New Strategy
				</button>
			{/if}
		</div>
	</div>

	<!-- Strategies Table -->
	<div class="flex-1 overflow-auto">
		{#if strategiesStore.isLoading}
			<div class="flex h-64 items-center justify-center">
				<div class="text-center">
					<p class="text-slate-500">Loading strategies...</p>
				</div>
			</div>
		{:else if strategiesStore.strategies.length === 0 && !isAddMode}
			<div class="flex h-64 items-center justify-center">
				<div class="text-center">
					<svg
						class="mx-auto h-12 w-12 text-slate-700"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
						></path>
					</svg>
					<p class="mt-4 text-slate-500">No strategies yet</p>
					<button
						onclick={startAdd}
						class="mt-4 text-sm font-medium text-emerald-400 hover:text-emerald-300"
					>
						Create your first strategy
					</button>
				</div>
			</div>
		{:else}
			<table class="w-full">
				<thead class="border-b border-slate-800 bg-slate-900">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-bold uppercase text-slate-500">Name</th>
						<th class="px-6 py-3 text-left text-xs font-bold uppercase text-slate-500">
							Description
						</th>
						<th class="px-6 py-3 text-right text-xs font-bold uppercase text-slate-500">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-800/50 bg-slate-950">
				<!-- Add New Row -->
				{#if isAddMode}
					<tr class="bg-slate-900/50">
						<td class="px-6 py-4">
							<input
								type="text"
								bind:value={formData.name}
								placeholder="Strategy name"
								class="w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-emerald-500 focus:outline-none"
								autofocus
							/>
						</td>
						<td class="px-6 py-4">
							<input
								type="text"
								bind:value={formData.description}
								placeholder="Description (optional)"
								class="w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-emerald-500 focus:outline-none"
							/>
						</td>
						<td class="px-6 py-4">
							<div class="flex justify-end gap-2">
								<button
									onclick={cancelAdd}
									class="px-3 py-1 text-sm text-slate-400 transition-colors hover:text-slate-200"
								>
									Cancel
								</button>
								<button
									onclick={saveAdd}
									class="bg-emerald-600 px-3 py-1 text-sm text-white transition-colors hover:bg-emerald-700"
								>
									Save
								</button>
							</div>
						</td>
					</tr>
				{/if}

				<!-- Strategy Rows -->
				{#each strategiesStore.strategies as strategy}
					<tr class="transition-colors hover:bg-slate-900/50">
						<td class="px-6 py-4">
							{#if editingId === strategy.id}
								<input
									type="text"
									bind:value={formData.name}
									class="w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:border-emerald-500 focus:outline-none"
								/>
							{:else}
								<span class="font-medium text-slate-100">{strategy.name}</span>
							{/if}
						</td>
						<td class="px-6 py-4">
							{#if editingId === strategy.id}
								<input
									type="text"
									bind:value={formData.description}
									class="w-full border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:border-emerald-500 focus:outline-none"
								/>
							{:else}
								<span class="text-sm text-slate-400">{strategy.description || '-'}</span>
							{/if}
						</td>
						<td class="px-6 py-4">
							<div class="flex justify-end gap-2">
								{#if editingId === strategy.id}
									<button
										onclick={cancelEdit}
										class="px-3 py-1 text-sm text-slate-400 transition-colors hover:text-slate-200"
									>
										Cancel
									</button>
									<button
										onclick={saveEdit}
										class="bg-emerald-600 px-3 py-1 text-sm text-white transition-colors hover:bg-emerald-700"
									>
										Save
									</button>
								{:else}
									<button
										onclick={() => startEdit(strategy)}
										class="text-slate-400 transition-colors hover:text-slate-200"
										aria-label="Edit strategy"
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
										onclick={() => openDeleteConfirm(strategy.id)}
										class="text-red-400 transition-colors hover:text-red-300"
										aria-label="Delete strategy"
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
								{/if}
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
		{/if}
	</div>
</div>

<ConfirmDialog
	isOpen={isConfirmOpen}
	title="Delete Strategy"
	message="Are you sure you want to delete this strategy? This action cannot be undone."
	confirmText="Delete"
	cancelText="Cancel"
	variant="danger"
	onConfirm={handleDeleteConfirm}
	onCancel={handleDeleteCancel}
/>
