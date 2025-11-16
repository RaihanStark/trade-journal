import { apiClient, type Account } from '$lib/api/client';
import { authStore } from './auth.svelte';

interface AccountsState {
	accounts: Account[];
	isLoading: boolean;
	error: string | null;
}

class AccountsStore {
	private state = $state<AccountsState>({
		accounts: [],
		isLoading: false,
		error: null
	});

	get accounts() {
		return this.state.accounts;
	}

	get isLoading() {
		return this.state.isLoading;
	}

	get error() {
		return this.state.error;
	}

	async load(force = false) {
		// Don't reload if we already have data unless forced
		if (this.state.accounts.length > 0 && !force) {
			return;
		}

		if (!authStore.token) {
			this.state.error = 'Not authenticated';
			return;
		}

		this.state.isLoading = true;
		this.state.error = null;

		const { data, error } = await apiClient.getAccounts(authStore.token);

		if (error) {
			this.state.error = error;
			this.state.accounts = [];
		} else if (data) {
			this.state.accounts = data;
		}

		this.state.isLoading = false;
	}

	async reload() {
		return this.load(true);
	}

	async add(account: Account) {
		this.state.accounts = [...this.state.accounts, account];
	}

	async update(updatedAccount: Account) {
		this.state.accounts = this.state.accounts.map((a) =>
			a.id === updatedAccount.id ? updatedAccount : a
		);
	}

	async remove(id: number) {
		this.state.accounts = this.state.accounts.filter((a) => a.id !== id);
	}

	clear() {
		this.state.accounts = [];
		this.state.error = null;
		this.state.isLoading = false;
	}
}

export const accountsStore = new AccountsStore();
