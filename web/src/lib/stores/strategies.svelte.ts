import { apiClient, type Strategy } from '$lib/api/client';
import { authStore } from './auth.svelte';

interface StrategiesState {
	strategies: Strategy[];
	isLoading: boolean;
	error: string | null;
}

class StrategiesStore {
	private state = $state<StrategiesState>({
		strategies: [],
		isLoading: false,
		error: null
	});

	get strategies() {
		return this.state.strategies;
	}

	get isLoading() {
		return this.state.isLoading;
	}

	get error() {
		return this.state.error;
	}

	async load(force = false) {
		// Don't reload if we already have data unless forced
		if (this.state.strategies.length > 0 && !force) {
			return;
		}

		if (!authStore.token) {
			this.state.error = 'Not authenticated';
			return;
		}

		this.state.isLoading = true;
		this.state.error = null;

		const { data, error } = await apiClient.getStrategies(authStore.token);

		if (error) {
			this.state.error = error;
			this.state.strategies = [];
		} else if (data) {
			this.state.strategies = data;
		}

		this.state.isLoading = false;
	}

	async reload() {
		return this.load(true);
	}

	async add(strategy: Strategy) {
		this.state.strategies = [...this.state.strategies, strategy];
	}

	async update(updatedStrategy: Strategy) {
		this.state.strategies = this.state.strategies.map((s) =>
			s.id === updatedStrategy.id ? updatedStrategy : s
		);
	}

	async remove(id: number) {
		this.state.strategies = this.state.strategies.filter((s) => s.id !== id);
	}

	clear() {
		this.state.strategies = [];
		this.state.error = null;
		this.state.isLoading = false;
	}
}

export const strategiesStore = new StrategiesStore();
