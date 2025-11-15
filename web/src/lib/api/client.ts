const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

interface ApiError {
	error: string;
}

export interface User {
	id: number;
	email: string;
}

export interface AuthResponse {
	token: string;
	user: User;
}

export interface RegisterRequest {
	email: string;
	password: string;
}

export interface LoginRequest {
	email: string;
	password: string;
}

export interface Account {
	id: number;
	name: string;
	broker: string;
	account_number: string;
	account_type: 'demo' | 'live';
	currency: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateAccountRequest {
	name: string;
	broker: string;
	account_number: string;
	account_type: 'demo' | 'live';
	currency: string;
	is_active: boolean;
}

export interface UpdateAccountRequest {
	name: string;
	broker: string;
	account_number: string;
	account_type: 'demo' | 'live';
	currency: string;
	is_active: boolean;
}

class ApiClient {
	private baseUrl: string;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;
	}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<{ data?: T; error?: string }> {
		try {
			const response = await fetch(`${this.baseUrl}${endpoint}`, {
				...options,
				headers: {
					'Content-Type': 'application/json',
					...options.headers
				}
			});

			const data = await response.json();

			if (!response.ok) {
				return { error: (data as ApiError).error || 'An error occurred' };
			}

			return { data: data as T };
		} catch (err) {
			return { error: 'Network error. Please check your connection.' };
		}
	}

	async register(req: RegisterRequest): Promise<{ data?: AuthResponse; error?: string }> {
		return this.request<AuthResponse>('/api/auth/register', {
			method: 'POST',
			body: JSON.stringify(req)
		});
	}

	async login(req: LoginRequest): Promise<{ data?: AuthResponse; error?: string }> {
		return this.request<AuthResponse>('/api/auth/login', {
			method: 'POST',
			body: JSON.stringify(req)
		});
	}

	async getCurrentUser(token: string): Promise<{ data?: User; error?: string }> {
		return this.request<User>('/api/me', {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}

	// Account APIs
	async getAccounts(token: string): Promise<{ data?: Account[]; error?: string }> {
		return this.request<Account[]>('/api/accounts', {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}

	async getAccount(id: number, token: string): Promise<{ data?: Account; error?: string }> {
		return this.request<Account>(`/api/accounts/${id}`, {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}

	async createAccount(
		req: CreateAccountRequest,
		token: string
	): Promise<{ data?: Account; error?: string }> {
		return this.request<Account>('/api/accounts', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(req)
		});
	}

	async updateAccount(
		id: number,
		req: UpdateAccountRequest,
		token: string
	): Promise<{ data?: Account; error?: string }> {
		return this.request<Account>(`/api/accounts/${id}`, {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(req)
		});
	}

	async deleteAccount(id: number, token: string): Promise<{ data?: any; error?: string }> {
		return this.request<any>(`/api/accounts/${id}`, {
			method: 'DELETE',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}
}

export const apiClient = new ApiClient(API_BASE_URL);
