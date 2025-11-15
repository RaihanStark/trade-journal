<script lang="ts">
	let { children } = $props();
</script>

<div class="h-screen bg-slate-950">
	{@render children()}
</div>
