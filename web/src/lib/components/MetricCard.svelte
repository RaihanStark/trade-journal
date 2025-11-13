<script lang="ts">
	interface Props {
		title: string;
		value: string | number;
		change?: string;
		trend?: 'up' | 'down' | 'neutral';
	}

	let { title, value, change, trend = 'neutral' }: Props = $props();

	function getTrendColor(trend: string): string {
		switch (trend) {
			case 'up':
				return 'text-emerald-400';
			case 'down':
				return 'text-red-400';
			default:
				return 'text-slate-400';
		}
	}
</script>

<div class="flex h-full flex-col justify-center border-r border-slate-800 bg-slate-900 px-3 py-2">
	<div class="mb-1">
		<h3 class="text-xs font-bold uppercase tracking-wide text-slate-500">{title}</h3>
	</div>
	<div class="mb-1">
		<p class="font-mono text-2xl font-bold text-slate-100">{value}</p>
	</div>
	{#if change}
		<div>
			<span class={getTrendColor(trend) + ' font-mono text-sm font-medium'}>{change}</span>
		</div>
	{/if}
</div>
