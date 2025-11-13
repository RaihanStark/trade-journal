<script lang="ts">
	interface Trade {
		id: number;
		date: string;
		time: string;
		pair: string;
		type: 'BUY' | 'SELL';
		entry: number;
		exit: number;
		lots: number;
		pips: number;
		pl: number;
		rr: number;
		status: string;
	}

	interface Props {
		trades: Trade[];
	}

	let { trades }: Props = $props();

	function getPLColor(pl: number): string {
		if (pl > 0) return 'text-emerald-400';
		if (pl < 0) return 'text-red-400';
		return 'text-slate-400';
	}
</script>

<div class="flex h-full flex-col bg-slate-900">
	<div class="flex-1 overflow-auto">
		<table class="w-full">
			<thead class="sticky top-0 border-b border-slate-800 bg-slate-950">
				<tr>
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
				</tr>
			</thead>
			<tbody>
				{#each trades as trade (trade.id)}
					<tr class="border-b border-slate-800/50 hover:bg-slate-800/30">
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
							{trade.exit.toFixed(trade.pair.includes('JPY') ? 2 : 4)}
						</td>
						<td class="px-3 py-2 text-right font-mono text-sm text-slate-300">
							{trade.lots.toFixed(2)}
						</td>
						<td class="px-3 py-2 text-right">
							<span class={getPLColor(trade.pips) + ' font-mono text-sm font-bold'}>
								{trade.pips > 0 ? '+' : ''}{trade.pips}
							</span>
						</td>
						<td class="px-3 py-2 text-right">
							<span class={getPLColor(trade.pl) + ' font-mono text-base font-bold'}>
								{trade.pl > 0 ? '+' : ''}${trade.pl.toFixed(0)}
							</span>
						</td>
						<td class="px-3 py-2 text-right">
							<span class={getPLColor(trade.rr) + ' font-mono text-sm'}>
								{trade.rr > 0 ? '+' : ''}{trade.rr.toFixed(1)}
							</span>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
