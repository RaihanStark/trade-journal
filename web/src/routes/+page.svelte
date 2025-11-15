<script lang="ts">
	import MetricCard from '$lib/components/MetricCard.svelte';
	import TradesTable from '$lib/components/TradesTable.svelte';
	import StatsGrid from '$lib/components/StatsGrid.svelte';
	import AddTradeModal from '$lib/components/AddTradeModal.svelte';

	let isModalOpen = $state(false);

	function openModal() {
		isModalOpen = true;
	}

	function closeModal() {
		isModalOpen = false;
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

	const recentTrades: Array<{
		id: number;
		date: string;
		time: string;
		pair: string;
		type: 'BUY' | 'SELL';
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
	}> = [
		{
			id: 1,
			date: '2025-01-13',
			time: '08:45',
			pair: 'EUR/USD',
			type: 'BUY' as const,
			entry: 1.0915,
			exit: null,
			lots: 0.5,
			pips: null,
			pl: null,
			rr: null,
			status: 'open' as const,
			stopLoss: 1.0895,
			takeProfit: 1.0955,
			notes: 'Bullish trend continuation setup',
			strategy: 'Trend Following'
		},
		{
			id: 2,
			date: '2025-01-13',
			time: '07:20',
			pair: 'GBP/USD',
			type: 'SELL' as const,
			entry: 1.265,
			exit: null,
			lots: 0.3,
			pips: null,
			pl: null,
			rr: null,
			status: 'open' as const,
			stopLoss: 1.267,
			takeProfit: 1.26,
			notes: 'Resistance level rejection',
			strategy: 'Support/Resistance'
		},
		{
			id: 3,
			date: '2025-01-12',
			time: '14:30',
			pair: 'EUR/USD',
			type: 'BUY' as const,
			entry: 1.0925,
			exit: 1.0965,
			lots: 0.5,
			pips: 40,
			pl: 200,
			rr: 2.5,
			status: 'closed' as const,
			stopLoss: 1.0905,
			takeProfit: 1.0965,
			notes: 'Perfect entry at support level. Hit TP exactly.',
			strategy: 'Support/Resistance'
		},
		{
			id: 4,
			date: '2025-01-12',
			time: '09:15',
			pair: 'GBP/JPY',
			type: 'SELL' as const,
			entry: 188.45,
			exit: 188.05,
			lots: 0.3,
			pips: 40,
			pl: 120,
			rr: 2.0,
			status: 'closed' as const
		},
		{
			id: 5,
			date: '2025-01-11',
			time: '16:45',
			pair: 'USD/JPY',
			type: 'BUY' as const,
			entry: 145.2,
			exit: 145.05,
			lots: 0.4,
			pips: -15,
			pl: -60,
			rr: -0.5,
			status: 'closed' as const
		},
		{
			id: 6,
			date: '2025-01-11',
			time: '11:20',
			pair: 'EUR/GBP',
			type: 'SELL' as const,
			entry: 0.858,
			exit: 0.8545,
			lots: 0.6,
			pips: 35,
			pl: 210,
			rr: 2.8,
			status: 'closed' as const
		},
		{
			id: 7,
			date: '2025-01-10',
			time: '13:00',
			pair: 'AUD/USD',
			type: 'BUY' as const,
			entry: 0.672,
			exit: 0.6705,
			lots: 0.5,
			pips: -15,
			pl: -75,
			rr: -0.6,
			status: 'closed' as const
		},
		{
			id: 8,
			date: '2025-01-10',
			time: '10:30',
			pair: 'EUR/USD',
			type: 'BUY' as const,
			entry: 1.088,
			exit: 1.092,
			lots: 0.7,
			pips: 40,
			pl: 280,
			rr: 3.2,
			status: 'closed' as const
		},
		{
			id: 9,
			date: '2025-01-09',
			time: '15:00',
			pair: 'GBP/USD',
			type: 'SELL' as const,
			entry: 1.265,
			exit: 1.26,
			lots: 0.4,
			pips: 50,
			pl: 200,
			rr: 2.5,
			status: 'closed' as const
		},
		{
			id: 10,
			date: '2025-01-09',
			time: '11:45',
			pair: 'USD/CAD',
			type: 'BUY' as const,
			entry: 1.352,
			exit: 1.348,
			lots: 0.5,
			pips: -40,
			pl: -200,
			rr: -1.5,
			status: 'closed' as const
		}
	];
</script>

<div class="grid h-full grid-cols-12 grid-rows-[auto_auto_1fr] bg-slate-950">
	<!-- Header -->
	<div
		class="col-span-12 flex items-center justify-between border-b border-slate-800 bg-slate-900 px-4 py-2"
	>
		<div class="flex items-center gap-4">
			<span class="font-mono text-xs text-slate-500">BAL:</span>
			<span class="font-mono text-sm font-bold text-emerald-400">$25,420.50</span>
			<span class="font-mono text-xs text-slate-500">EQ:</span>
			<span class="font-mono text-sm font-bold text-blue-400">$25,890.30</span>
		</div>
		<button onclick={openModal} class="bg-slate-800 px-3 py-1.5 text-xs text-slate-300 hover:bg-slate-700">+ TRADE</button>
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
		<TradesTable trades={recentTrades} />
	</div>

	<div class="col-span-2 row-span-1">
		<StatsGrid {metrics} />
	</div>
</div>

<AddTradeModal isOpen={isModalOpen} onClose={closeModal} />
