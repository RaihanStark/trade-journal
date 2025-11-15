<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import Sidebar from '$lib/components/Sidebar.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex h-screen overflow-hidden bg-slate-950">
	<Sidebar />
	<main class="flex-1 overflow-hidden">
		{@render children()}
	</main>
</div>
