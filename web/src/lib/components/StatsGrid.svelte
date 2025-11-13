<script lang="ts">
	interface Metrics {
		avgWin: number;
		avgLoss: number;
		sharpeRatio: number;
		maxDrawdown: number;
	}

	interface Props {
		metrics: Metrics;
	}

	let { metrics }: Props = $props();

	const stats = [
		{ label: 'AVG WIN', value: `$${metrics.avgWin.toFixed(0)}`, color: 'text-emerald-400' },
		{ label: 'AVG LOSS', value: `$${metrics.avgLoss.toFixed(0)}`, color: 'text-red-400' },
		{ label: 'SHARPE', value: metrics.sharpeRatio.toFixed(2), color: 'text-blue-400' },
		{ label: 'MAX DD', value: `$${metrics.maxDrawdown.toFixed(0)}`, color: 'text-amber-400' },
		{
			label: 'R/R',
			value: (Math.abs(metrics.avgWin / metrics.avgLoss)).toFixed(2),
			color: 'text-purple-400'
		}
	];
</script>

<div class="flex h-full flex-col bg-slate-900">
	<div class="flex-1 overflow-auto">
		{#each stats as stat}
			<div class="flex flex-col border-b border-slate-800/50 px-3 py-3">
				<span class="mb-1 text-xs font-bold uppercase text-slate-500">{stat.label}</span>
				<span class="font-mono text-xl font-bold {stat.color}">{stat.value}</span>
			</div>
		{/each}
	</div>
</div>
