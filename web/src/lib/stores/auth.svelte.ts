import { apiClient, type User } from '$lib/api/client';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

interface AuthState {
	user: User | null;
	token: string | null;
	isAuthenticated: boolean;
	isLoading: boolean;
}

class AuthStore {
	private state = $state<AuthState>({
		user: null,
		token: null,
		isAuthenticated: false,
		isLoading: true
	});

	constructor() {
		if (browser) {
			this.initialize();
		}
	}

	private initialize() {
		const token = localStorage.getItem('auth_token');
		const userStr = localStorage.getItem('auth_user');

		if (token && userStr) {
			try {
				const user = JSON.parse(userStr);
				this.state.token = token;
				this.state.user = user;
				this.state.isAuthenticated = true;
			} catch (err) {
				this.clearAuth();
			}
		}
		this.state.isLoading = false;
	}

	get user() {
		return this.state.user;
	}

	get token() {
		return this.state.token;
	}

	get isAuthenticated() {
		return this.state.isAuthenticated;
	}

	get isLoading() {
		return this.state.isLoading;
	}

	async register(email: string, password: string): Promise<string | null> {
		const { data, error } = await apiClient.register({ email, password });

		if (error) {
			return error;
		}

		if (data) {
			this.setAuth(data.token, data.user);
			goto('/');
		}

		return null;
	}

	async login(email: string, password: string): Promise<string | null> {
		const { data, error } = await apiClient.login({ email, password });

		if (error) {
			return error;
		}

		if (data) {
			this.setAuth(data.token, data.user);
			goto('/');
		}

		return null;
	}

	logout() {
		this.clearAuth();
		goto('/login');
	}

	private setAuth(token: string, user: User) {
		this.state.token = token;
		this.state.user = user;
		this.state.isAuthenticated = true;

		if (browser) {
			localStorage.setItem('auth_token', token);
			localStorage.setItem('auth_user', JSON.stringify(user));
		}
	}

	private clearAuth() {
		this.state.token = null;
		this.state.user = null;
		this.state.isAuthenticated = false;

		if (browser) {
			localStorage.removeItem('auth_token');
			localStorage.removeItem('auth_user');
		}
	}
}

export const authStore = new AuthStore();
