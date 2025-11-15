<script lang="ts">
	import { slide } from 'svelte/transition';

	interface Trade {
		id: number;
		date: string;
		time: string;
		pair: string;
		type: 'BUY' | 'SELL' | 'DEPOSIT' | 'WITHDRAW';
		entry: number;
		exit: number | null;
		lots: number;
		pips: number | null;
		pl: number | null;
		rr: number | null;
		status: 'open' | 'closed';
		stopLoss?: number;
		takeProfit?: number;
		notes?: string;
		strategy?: string;
		mistakes?: string;
		amount?: number;
	}

	function isTransaction(trade: Trade): boolean {
		return trade.type === 'DEPOSIT' || trade.type === 'WITHDRAW';
	}

	interface Props {
		trades: Trade[];
	}

	let { trades }: Props = $props();

	let expandedRows = $state<Set<number>>(new Set());
	let editingField = $state<{ tradeId: number; field: string } | null>(null);
	let editValue = $state<string>('');

	function toggleRow(tradeId: number) {
		if (expandedRows.has(tradeId)) {
			expandedRows.delete(tradeId);
		} else {
			expandedRows.add(tradeId);
		}
		expandedRows = new Set(expandedRows);
	}

	function startEdit(tradeId: number, field: string, currentValue: any) {
		editingField = { tradeId, field };
		editValue = currentValue?.toString() || '';
	}

	function saveEdit() {
		// TODO: Implement save logic
		console.log('Saving:', editingField, 'Value:', editValue);
		editingField = null;
	}

	function cancelEdit() {
		editingField = null;
		editValue = '';
	}

	function getPLColor(pl: number): string {
		if (pl > 0) return 'text-emerald-400';
		if (pl < 0) return 'text-red-400';
		return 'text-slate-400';
	}

	function getTradeStatus(trade: Trade): 'ongoing' | 'win' | 'loss' {
		if (trade.exit === null) {
			return 'ongoing';
		}
		if (trade.pl !== null && trade.pl > 0) {
			return 'win';
		}
		return 'loss';
	}
</script>

<div class="flex h-full flex-col bg-slate-900">
	<div class="flex-1 overflow-auto">
		<table class="w-full">
			<thead class="sticky top-0 border-b border-slate-800 bg-slate-950">
				<tr>
					<th class="w-8 px-3 py-2"></th>
					<th class="px-3 py-2 text-left text-xs font-bold uppercase text-slate-500">DATE</th>
					<th class="px-3 py-2 text-left text-xs font-bold uppercase text-slate-500">TIME</th>
					<th class="px-3 py-2 text-left text-xs font-bold uppercase text-slate-500">PAIR</th>
					<th class="px-3 py-2 text-left text-xs font-bold uppercase text-slate-500">TYPE</th>
					<th class="px-3 py-2 text-right text-xs font-bold uppercase text-slate-500">ENTRY</th>
					<th class="px-3 py-2 text-right text-xs font-bold uppercase text-slate-500">EXIT</th>
					<th class="px-3 py-2 text-right text-xs font-bold uppercase text-slate-500">LOTS</th>
					<th class="px-3 py-2 text-right text-xs font-bold uppercase text-slate-500">PIPS</th>
					<th class="px-3 py-2 text-right text-xs font-bold uppercase text-slate-500">P/L</th>
					<th class="px-3 py-2 text-right text-xs font-bold uppercase text-slate-500">R:R</th>
					<th class="px-3 py-2 text-center text-xs font-bold uppercase text-slate-500">STATUS</th>
				</tr>
			</thead>
			<tbody>
				{#each trades as trade (trade.id)}
					{#if isTransaction(trade)}
						<!-- Transaction Row (Deposit/Withdraw) -->
						<tr
							class="border-b border-slate-800/50 hover:bg-slate-800/30 {trade.type === 'DEPOSIT'
								? 'bg-emerald-950/10'
								: 'bg-red-950/10'}"
						>
							<td class="px-3 py-2"></td>
							<td class="px-3 py-2 font-mono text-sm text-slate-400">{trade.date}</td>
							<td class="px-3 py-2 font-mono text-sm text-slate-500">{trade.time}</td>
							<td class="px-3 py-2 font-mono text-sm font-bold text-slate-200" colspan="6">
								<div class="flex items-center gap-3">
									<span
										class="px-2 py-0.5 font-mono text-xs font-bold {trade.type === 'DEPOSIT'
											? 'bg-blue-500/20 text-blue-400'
											: 'bg-orange-500/20 text-orange-400'}"
									>
										{trade.type}
									</span>
									{#if trade.notes}
										<span class="text-sm text-slate-400">{trade.notes}</span>
									{/if}
								</div>
							</td>
							<td class="px-3 py-2 text-right">
								{#if trade.pl !== null}
									<span class={getPLColor(trade.pl) + ' font-mono text-base font-bold'}>
										{trade.pl > 0 ? '+' : ''}${trade.pl.toFixed(0)}
									</span>
								{:else}
									<span class="font-mono text-base text-slate-600">-</span>
								{/if}
							</td>
							<td class="px-3 py-2"></td>
							<td class="px-3 py-2 text-center">
								<span
									class="inline-block {trade.type === 'DEPOSIT'
										? 'bg-blue-600'
										: 'bg-orange-600'} px-2 py-0.5 font-mono text-xs font-bold uppercase text-white"
								>
									{trade.type}
								</span>
							</td>
						</tr>
					{:else}
						<!-- Trade Row -->
						<tr
							class="cursor-pointer border-b border-slate-800/50 hover:bg-slate-800/30 {getTradeStatus(
								trade
							) === 'ongoing'
								? 'bg-blue-950/20'
								: ''}"
							onclick={() => toggleRow(trade.id)}
						>
							<td class="px-3 py-2">
								<svg
									class="h-4 w-4 text-slate-500 transition-transform {expandedRows.has(trade.id)
										? 'rotate-90'
										: ''}"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"
									></path>
								</svg>
							</td>
							<td class="px-3 py-2 font-mono text-sm text-slate-400">{trade.date}</td>
							<td class="px-3 py-2 font-mono text-sm text-slate-500">{trade.time}</td>
							<td class="px-3 py-2 font-mono text-sm font-bold text-slate-200">{trade.pair}</td>
							<td class="px-3 py-2">
								<span
									class="px-2 py-0.5 font-mono text-xs font-bold {trade.type === 'BUY'
										? 'bg-emerald-500/20 text-emerald-400'
										: 'bg-red-500/20 text-red-400'}"
								>
									{trade.type}
								</span>
							</td>
							<td class="px-3 py-2 text-right font-mono text-sm text-slate-300">
								{trade.entry.toFixed(trade.pair.includes('JPY') ? 2 : 4)}
							</td>
							<td class="px-3 py-2 text-right font-mono text-sm text-slate-300">
								{#if trade.exit !== null}
									{trade.exit.toFixed(trade.pair.includes('JPY') ? 2 : 4)}
								{:else}
									<span class="text-slate-600">-</span>
								{/if}
							</td>
							<td class="px-3 py-2 text-right font-mono text-sm text-slate-300">
								{trade.lots.toFixed(2)}
							</td>
							<td class="px-3 py-2 text-right">
								{#if trade.pips !== null}
									<span class={getPLColor(trade.pips) + ' font-mono text-sm font-bold'}>
										{trade.pips > 0 ? '+' : ''}{trade.pips}
									</span>
								{:else}
									<span class="font-mono text-sm text-slate-600">-</span>
								{/if}
							</td>
							<td class="px-3 py-2 text-right">
								{#if trade.pl !== null}
									<span class={getPLColor(trade.pl) + ' font-mono text-base font-bold'}>
										{trade.pl > 0 ? '+' : ''}${trade.pl.toFixed(0)}
									</span>
								{:else}
									<span class="font-mono text-base text-slate-600">-</span>
								{/if}
							</td>
							<td class="px-3 py-2 text-right">
								{#if trade.rr !== null}
									<span class={getPLColor(trade.rr) + ' font-mono text-sm'}>
										{trade.rr > 0 ? '+' : ''}{trade.rr.toFixed(1)}
									</span>
								{:else}
									<span class="font-mono text-sm text-slate-600">-</span>
								{/if}
							</td>
							<td class="px-3 py-2 text-center">
								{#if getTradeStatus(trade) === 'ongoing'}
									<span
										class="inline-block bg-blue-600 px-2 py-0.5 font-mono text-xs font-bold uppercase text-white"
									>
										OPEN
									</span>
								{:else if getTradeStatus(trade) === 'win'}
									<span
										class="inline-block bg-emerald-600 px-2 py-0.5 font-mono text-xs font-bold uppercase text-white"
									>
										WIN
									</span>
								{:else}
									<span
										class="inline-block bg-red-600 px-2 py-0.5 font-mono text-xs font-bold uppercase text-white"
									>
										LOSS
									</span>
								{/if}
							</td>
						</tr>
					{/if}

					<!-- Expanded Details Row -->
					{#if expandedRows.has(trade.id)}
						<tr class="border-b border-slate-800 bg-slate-950">
							<td colspan="12" class="overflow-hidden p-0">
								<div transition:slide={{ duration: 200 }} class="px-3 py-4">
									<div class="grid grid-cols-2 gap-4 text-sm" onclick={(e) => e.stopPropagation()}>
										<!-- Exit Price -->
										<div class="flex items-center justify-between border-b border-slate-800/50 pb-2">
										<span class="text-xs font-bold uppercase text-slate-500">EXIT PRICE</span>
										<div class="flex items-center gap-2">
											{#if editingField?.tradeId === trade.id && editingField?.field === 'exit'}
												<input
													type="number"
													step="0.00001"
													bind:value={editValue}
													class="w-32 border border-slate-700 bg-slate-900 px-2 py-1 font-mono text-xs text-slate-100 focus:border-blue-600 focus:outline-none"
													autofocus
												/>
												<button
													onclick={saveEdit}
													class="text-emerald-400 hover:text-emerald-300"
													aria-label="Save"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M5 13l4 4L19 7"
														></path>
													</svg>
												</button>
												<button
													onclick={cancelEdit}
													class="text-red-400 hover:text-red-300"
													aria-label="Cancel"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M6 18L18 6M6 6l12 12"
														></path>
													</svg>
												</button>
											{:else}
												<span class="font-mono text-sm text-slate-300">
													{#if trade.exit !== null}
														{trade.exit.toFixed(trade.pair.includes('JPY') ? 2 : 4)}
													{:else}
														-
													{/if}
												</span>
												<button
													onclick={() => startEdit(trade.id, 'exit', trade.exit)}
													class="text-slate-500 hover:text-slate-300"
													aria-label="Edit exit price"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
														></path>
													</svg>
												</button>
											{/if}
										</div>
									</div>

									<!-- Stop Loss -->
									<div class="flex items-center justify-between border-b border-slate-800/50 pb-2">
										<span class="text-xs font-bold uppercase text-slate-500">STOP LOSS</span>
										<div class="flex items-center gap-2">
											{#if editingField?.tradeId === trade.id && editingField?.field === 'stopLoss'}
												<input
													type="number"
													step="0.00001"
													bind:value={editValue}
													class="w-32 border border-slate-700 bg-slate-900 px-2 py-1 font-mono text-xs text-slate-100 focus:border-blue-600 focus:outline-none"
													autofocus
												/>
												<button
													onclick={saveEdit}
													class="text-emerald-400 hover:text-emerald-300"
													aria-label="Save"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M5 13l4 4L19 7"
														></path>
													</svg>
												</button>
												<button
													onclick={cancelEdit}
													class="text-red-400 hover:text-red-300"
													aria-label="Cancel"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M6 18L18 6M6 6l12 12"
														></path>
													</svg>
												</button>
											{:else}
												<span class="font-mono text-sm text-slate-300">
													{trade.stopLoss?.toFixed(trade.pair.includes('JPY') ? 2 : 4) || '-'}
												</span>
												<button
													onclick={() => startEdit(trade.id, 'stopLoss', trade.stopLoss)}
													class="text-slate-500 hover:text-slate-300"
													aria-label="Edit stop loss"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
														></path>
													</svg>
												</button>
											{/if}
										</div>
									</div>

									<!-- Take Profit -->
									<div class="flex items-center justify-between border-b border-slate-800/50 pb-2">
										<span class="text-xs font-bold uppercase text-slate-500">TAKE PROFIT</span>
										<div class="flex items-center gap-2">
											{#if editingField?.tradeId === trade.id && editingField?.field === 'takeProfit'}
												<input
													type="number"
													step="0.00001"
													bind:value={editValue}
													class="w-32 border border-slate-700 bg-slate-900 px-2 py-1 font-mono text-xs text-slate-100 focus:border-blue-600 focus:outline-none"
													autofocus
												/>
												<button
													onclick={saveEdit}
													class="text-emerald-400 hover:text-emerald-300"
													aria-label="Save"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M5 13l4 4L19 7"
														></path>
													</svg>
												</button>
												<button
													onclick={cancelEdit}
													class="text-red-400 hover:text-red-300"
													aria-label="Cancel"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M6 18L18 6M6 6l12 12"
														></path>
													</svg>
												</button>
											{:else}
												<span class="font-mono text-sm text-slate-300">
													{trade.takeProfit?.toFixed(trade.pair.includes('JPY') ? 2 : 4) || '-'}
												</span>
												<button
													onclick={() => startEdit(trade.id, 'takeProfit', trade.takeProfit)}
													class="text-slate-500 hover:text-slate-300"
													aria-label="Edit take profit"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
														></path>
													</svg>
												</button>
											{/if}
										</div>
									</div>

									<!-- Strategy -->
									<div class="flex items-center justify-between border-b border-slate-800/50 pb-2">
										<span class="text-xs font-bold uppercase text-slate-500">STRATEGY</span>
										<div class="flex items-center gap-2">
											{#if editingField?.tradeId === trade.id && editingField?.field === 'strategy'}
												<input
													type="text"
													bind:value={editValue}
													class="w-48 border border-slate-700 bg-slate-900 px-2 py-1 font-mono text-xs text-slate-100 focus:border-blue-600 focus:outline-none"
													autofocus
												/>
												<button
													onclick={saveEdit}
													class="text-emerald-400 hover:text-emerald-300"
													aria-label="Save"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M5 13l4 4L19 7"
														></path>
													</svg>
												</button>
												<button
													onclick={cancelEdit}
													class="text-red-400 hover:text-red-300"
													aria-label="Cancel"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M6 18L18 6M6 6l12 12"
														></path>
													</svg>
												</button>
											{:else}
												<span class="text-sm text-slate-300">{trade.strategy || '-'}</span>
												<button
													onclick={() => startEdit(trade.id, 'strategy', trade.strategy)}
													class="text-slate-500 hover:text-slate-300"
													aria-label="Edit strategy"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
														></path>
													</svg>
												</button>
											{/if}
										</div>
									</div>

									<!-- Notes -->
									<div class="col-span-2 flex flex-col gap-2 border-b border-slate-800/50 pb-2">
										<div class="flex items-center justify-between">
											<span class="text-xs font-bold uppercase text-slate-500">NOTES</span>
											{#if editingField?.tradeId === trade.id && editingField?.field === 'notes'}
												<div class="flex items-center gap-2">
													<button
														onclick={saveEdit}
														class="text-emerald-400 hover:text-emerald-300"
														aria-label="Save"
													>
														<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																stroke-width="2"
																d="M5 13l4 4L19 7"
															></path>
														</svg>
													</button>
													<button
														onclick={cancelEdit}
														class="text-red-400 hover:text-red-300"
														aria-label="Cancel"
													>
														<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																stroke-width="2"
																d="M6 18L18 6M6 6l12 12"
															></path>
														</svg>
													</button>
												</div>
											{:else}
												<button
													onclick={() => startEdit(trade.id, 'notes', trade.notes)}
													class="text-slate-500 hover:text-slate-300"
													aria-label="Edit notes"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
														></path>
													</svg>
												</button>
											{/if}
										</div>
										{#if editingField?.tradeId === trade.id && editingField?.field === 'notes'}
											<textarea
												bind:value={editValue}
												rows="3"
												class="w-full border border-slate-700 bg-slate-900 px-2 py-1 text-xs text-slate-100 focus:border-blue-600 focus:outline-none"
												autofocus
											></textarea>
										{:else}
											<p class="text-sm text-slate-300">{trade.notes || 'No notes'}</p>
										{/if}
									</div>

									<!-- Mistakes -->
									<div class="col-span-2 flex flex-col gap-2 border-b border-slate-800/50 pb-2">
										<div class="flex items-center justify-between">
											<span class="text-xs font-bold uppercase text-slate-500">MISTAKES</span>
											{#if editingField?.tradeId === trade.id && editingField?.field === 'mistakes'}
												<div class="flex items-center gap-2">
													<button
														onclick={saveEdit}
														class="text-emerald-400 hover:text-emerald-300"
														aria-label="Save"
													>
														<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																stroke-width="2"
																d="M5 13l4 4L19 7"
															></path>
														</svg>
													</button>
													<button
														onclick={cancelEdit}
														class="text-red-400 hover:text-red-300"
														aria-label="Cancel"
													>
														<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																stroke-width="2"
																d="M6 18L18 6M6 6l12 12"
															></path>
														</svg>
													</button>
												</div>
											{:else}
												<button
													onclick={() => startEdit(trade.id, 'mistakes', trade.mistakes)}
													class="text-slate-500 hover:text-slate-300"
													aria-label="Edit mistakes"
												>
													<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
														></path>
													</svg>
												</button>
											{/if}
										</div>
										{#if editingField?.tradeId === trade.id && editingField?.field === 'mistakes'}
											<textarea
												bind:value={editValue}
												rows="3"
												class="w-full border border-slate-700 bg-slate-900 px-2 py-1 text-xs text-slate-100 focus:border-blue-600 focus:outline-none"
												autofocus
											></textarea>
										{:else}
											<p class="text-sm text-slate-300">{trade.mistakes || 'No mistakes noted'}</p>
										{/if}
									</div>
								</div>
							</div>
						</td>
						</tr>
					{/if}
				{/each}
			</tbody>
		</table>
	</div>
</div>
