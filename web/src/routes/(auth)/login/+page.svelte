<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let isLoading = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		isLoading = true;

		try {
			const err = await authStore.login(email, password);
			if (err) {
				error = err;
			}
		} catch (err) {
			error = 'An unexpected error occurred';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="flex h-full items-center justify-center bg-slate-950">
	<div class="w-full max-w-md">
		<!-- Header -->
		<div class="mb-8 text-center">
			<div class="mb-4 flex justify-center">
				<div class="flex h-12 w-12 items-center justify-center bg-slate-800 font-mono text-sm font-bold text-emerald-400">
					FX
				</div>
			</div>
			<h1 class="mb-2 text-2xl font-bold text-slate-100">FOREX JOURNAL</h1>
			<p class="text-sm text-slate-500">Sign in to your account</p>
		</div>

		<!-- Login Form -->
		<div class="border border-slate-800 bg-slate-900 p-8">
			<form onsubmit={handleSubmit} class="space-y-6">
				{#if error}
					<div class="border border-red-800 bg-red-900/20 px-4 py-3 text-sm text-red-400">
						{error}
					</div>
				{/if}

				<!-- Email Field -->
				<div>
					<label for="email" class="mb-2 block text-xs font-bold uppercase text-slate-400">
						Email
					</label>
					<input
						type="email"
						id="email"
						bind:value={email}
						required
						class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:border-emerald-500 focus:outline-none"
						placeholder="your@email.com"
					/>
				</div>

				<!-- Password Field -->
				<div>
					<label for="password" class="mb-2 block text-xs font-bold uppercase text-slate-400">
						Password
					</label>
					<input
						type="password"
						id="password"
						bind:value={password}
						required
						class="w-full border border-slate-700 bg-slate-950 px-4 py-3 text-sm text-slate-100 transition-colors placeholder:text-slate-600 focus:border-emerald-500 focus:outline-none"
						placeholder="Enter your password"
					/>
				</div>

				<!-- Submit Button -->
				<button
					type="submit"
					disabled={isLoading}
					class="w-full bg-emerald-600 px-4 py-3 text-sm font-bold uppercase text-white transition-colors hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
				>
					{isLoading ? 'Signing in...' : 'Sign In'}
				</button>
			</form>

			<!-- Register Link -->
			<div class="mt-6 border-t border-slate-800 pt-6 text-center">
				<p class="text-sm text-slate-500">
					Don't have an account?
					<a href="/register" class="font-medium text-emerald-400 hover:text-emerald-300">
						Create one
					</a>
				</p>
			</div>
		</div>
	</div>
</div>
