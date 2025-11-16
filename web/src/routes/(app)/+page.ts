import { apiClient } from '$lib/api/client';
import { authStore } from '$lib/stores/auth.svelte';
import { accountsStore } from '$lib/stores/accounts.svelte';
import { strategiesStore } from '$lib/stores/strategies.svelte';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	const token = authStore.token;

	// Load accounts and strategies into stores
	await Promise.all([
		accountsStore.load(),
		strategiesStore.load()
	]);

	// Return trades promise
	const tradesPromise = token
		? apiClient.getTrades(token).then(({ data, error }) => {
				if (error) {
					console.error('Failed to load trades:', error);
					return [];
				}
				return data || [];
		  })
		: Promise.resolve([]);

	return {
		trades: tradesPromise
	};
};
