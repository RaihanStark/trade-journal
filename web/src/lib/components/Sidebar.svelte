<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let currentPath = $derived($page.url.pathname);

	const menuItems = [
		{ id: 'dashboard', label: 'Dashboard', path: '/', icon: '▪' },
		{ id: 'accounts', label: 'Accounts', path: '/accounts', icon: '▪' },
		{ id: 'strategies', label: 'Strategies', path: '/strategies', icon: '▪' },
	];

	function handleNavigation(path: string) {
		goto(path);
	}

	function handleLogout() {
		authStore.logout();
	}
</script>

<aside class="flex h-full w-48 flex-col border-r border-slate-800 bg-slate-900">
	<!-- Logo/Brand -->
	<div class="border-b border-slate-800 px-4 py-3">
		<div class="flex items-center gap-2">
			<div class="flex h-8 w-8 items-center justify-center bg-slate-800 font-mono text-xs font-bold text-emerald-400">
				FX
			</div>
			<span class="text-xs font-bold text-slate-100">JOURNAL</span>
		</div>
	</div>

	<!-- Navigation Menu -->
	<nav class="flex-1 overflow-y-auto py-2">
		{#each menuItems as item}
			<button
				onclick={() => handleNavigation(item.path)}
				class="group flex w-full items-center gap-3 px-4 py-2.5 text-left transition-colors {currentPath === item.path
					? 'bg-slate-800 text-slate-100'
					: 'text-slate-400 hover:bg-slate-800/50 hover:text-slate-200'}"
			>
				<span class="text-base {currentPath === item.path ? 'text-emerald-400' : 'text-slate-600 group-hover:text-slate-500'}">{item.icon}</span>
				<span class="text-xs font-medium">{item.label}</span>
			</button>
		{/each}
	</nav>

	<!-- Footer -->
	<div class="border-t border-slate-800 px-4 py-3">
		<div class="mb-3 flex items-center gap-2">
			<div class="flex h-6 w-6 items-center justify-center rounded-full bg-slate-800 text-[10px] font-bold text-slate-400">
				{authStore.user?.email?.[0].toUpperCase() || 'U'}
			</div>
			<div class="flex-1 overflow-hidden">
				<div class="truncate text-[10px] font-medium text-slate-400">
					{authStore.user?.email || 'User'}
				</div>
				<div class="truncate text-[9px] text-slate-600">Trader</div>
			</div>
		</div>
		<button
			onclick={handleLogout}
			class="w-full border border-slate-700 px-2 py-1.5 text-[10px] font-medium uppercase text-slate-400 transition-colors hover:bg-slate-800 hover:text-slate-300"
		>
			Logout
		</button>
	</div>
</aside>
