<script lang="ts">
	const mockData = [
		{ date: 'Jan 1', value: 10000 },
		{ date: 'Jan 5', value: 10500 },
		{ date: 'Jan 9', value: 10200 },
		{ date: 'Jan 13', value: 11000 },
		{ date: 'Jan 17', value: 10800 },
		{ date: 'Jan 21', value: 11500 },
		{ date: 'Jan 25', value: 12200 },
		{ date: 'Jan 29', value: 11900 },
		{ date: 'Feb 2', value: 12800 },
		{ date: 'Feb 6', value: 13200 },
		{ date: 'Feb 10', value: 13000 },
		{ date: 'Feb 14', value: 13800 }
	];

	const maxValue = Math.max(...mockData.map((d) => d.value));
	const minValue = Math.min(...mockData.map((d) => d.value));
	const range = maxValue - minValue;
</script>

<div class="flex h-full flex-col border-b border-slate-800">
	<!-- Chart Header -->
	<div class="grid grid-cols-4 border-b border-slate-800 bg-slate-950">
		<div class="border-r border-slate-800 px-3 py-2">
			<p class="mb-0.5 text-xs font-bold uppercase text-slate-500">BALANCE</p>
			<p class="font-mono text-lg font-bold text-slate-100">$13,800</p>
		</div>
		<div class="border-r border-slate-800 px-3 py-2">
			<p class="mb-0.5 text-xs font-bold uppercase text-slate-500">GROWTH</p>
			<p class="font-mono text-lg font-bold text-emerald-400">+38%</p>
		</div>
		<div class="border-r border-slate-800 px-3 py-2">
			<p class="mb-0.5 text-xs font-bold uppercase text-slate-500">MONTH</p>
			<p class="font-mono text-lg font-bold text-emerald-400">+$3,800</p>
		</div>
		<div class="px-3 py-2">
			<p class="mb-0.5 text-xs font-bold uppercase text-slate-500">EQUITY</p>
			<p class="font-mono text-lg font-bold text-blue-400">$14,120</p>
		</div>
	</div>

	<!-- Chart Area -->
	<div class="relative flex-1 bg-slate-900 p-4">
		<svg class="h-full w-full" viewBox="0 0 1000 200" preserveAspectRatio="none">
			<!-- Grid lines -->
			{#each [0, 1, 2, 3, 4] as i}
				<line
					x1="0"
					y1={i * 50}
					x2="1000"
					y2={i * 50}
					stroke="#1e293b"
					stroke-width="0.5"
					stroke-dasharray="3 3"
				/>
			{/each}

			<!-- Area gradient -->
			<defs>
				<linearGradient id="chartGradient" x1="0" x2="0" y1="0" y2="1">
					<stop offset="0%" stop-color="#10b981" stop-opacity="0.15" />
					<stop offset="100%" stop-color="#10b981" stop-opacity="0.0" />
				</linearGradient>
			</defs>

			<!-- Area fill -->
			<path
				d={mockData
					.map((point, i) => {
						const x = (i / (mockData.length - 1)) * 1000;
						const y = 200 - ((point.value - minValue) / range) * 190 - 5;
						return `${i === 0 ? 'M' : 'L'} ${x} ${y}`;
					})
					.join(' ') + ' L 1000 200 L 0 200 Z'}
				fill="url(#chartGradient)"
			/>

			<!-- Chart line -->
			<path
				d={mockData
					.map((point, i) => {
						const x = (i / (mockData.length - 1)) * 1000;
						const y = 200 - ((point.value - minValue) / range) * 190 - 5;
						return `${i === 0 ? 'M' : 'L'} ${x} ${y}`;
					})
					.join(' ')}
				stroke="#10b981"
				stroke-width="2"
				fill="none"
			/>

			<!-- Data points -->
			{#each mockData as point, i}
				{@const x = (i / (mockData.length - 1)) * 1000}
				{@const y = 200 - ((point.value - minValue) / range) * 190 - 5}
				<circle cx={x} cy={y} r="3" fill="#10b981">
					<title>{point.date}: ${point.value}</title>
				</circle>
			{/each}
		</svg>
	</div>
</div>
