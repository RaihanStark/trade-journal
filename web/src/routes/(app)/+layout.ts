import { authStore } from '$lib/stores/auth.svelte';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

export const ssr = false;

export async function load() {
	if (browser) {
		// Wait for auth store to initialize
		await new Promise((resolve) => setTimeout(resolve, 0));

		if (!authStore.isAuthenticated) {
			goto('/login');
			return;
		}
	}

	return {
		user: authStore.user
	};
}
