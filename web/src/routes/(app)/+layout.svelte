<script lang="ts">
	import Sidebar from '$lib/components/Sidebar.svelte';

	let { children } = $props();
</script>

<div class="flex h-screen overflow-hidden bg-slate-950">
	<Sidebar />
	<main class="flex-1 overflow-hidden">
		{@render children()}
	</main>
</div>
