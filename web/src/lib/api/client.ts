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
	current_balance: number;
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

export interface Strategy {
	id: number;
	name: string;
	description: string;
	created_at: string;
	updated_at: string;
}

export interface CreateStrategyRequest {
	name: string;
	description: string;
}

export interface UpdateStrategyRequest {
	name: string;
	description: string;
}

export interface Trade {
	id: number;
	account_id: number | null;
	date: string;
	time: string;
	pair: string;
	type: 'BUY' | 'SELL' | 'DEPOSIT' | 'WITHDRAW';
	entry: number;
	exit: number | null;
	lots: number;
	pips: number | null;
	pl: number | null;
	rr: string;
	status: 'open' | 'closed';
	stop_loss: number | null;
	take_profit: number | null;
	notes: string;
	mistakes: string;
	amount: number | null;
	strategies: Array<{ id: number; name: string }>;
	created_at: string;
	updated_at: string;
}

export interface CreateTradeRequest {
	account_id: number | null;
	date: string;
	time: string;
	pair: string;
	type: string;
	entry: number;
	exit: number | null;
	lots: number;
	stop_loss: number | null;
	take_profit: number | null;
	notes: string;
	mistakes: string;
	amount: number | null;
	strategy_ids: number[];
}

export interface UpdateTradeRequest {
	account_id: number | null;
	date: string;
	time: string;
	pair: string;
	type: string;
	entry: number;
	exit: number | null;
	lots: number;
	stop_loss: number | null;
	take_profit: number | null;
	notes: string;
	mistakes: string;
	amount: number | null;
	strategy_ids: number[];
}

export interface Analytics {
	total_pl: number;
	win_rate: number;
	total_trades: number;
	winning_trades: number;
	losing_trades: number;
	avg_win: number;
	avg_loss: number;
	profit_factor: number;
	sharpe_ratio: number;
	max_drawdown: number;
	largest_win: number;
	largest_loss: number;
	avg_rr: number;
	consecutive_wins: number;
	consecutive_losses: number;
	best_streak: number;
	worst_streak: number;
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

			// Handle 204 No Content responses (e.g., DELETE operations)
			if (response.status === 204) {
				return { data: undefined as T };
			}

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

	// Strategy APIs
	async getStrategies(token: string): Promise<{ data?: Strategy[]; error?: string }> {
		return this.request<Strategy[]>('/api/strategies', {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}

	async getStrategy(id: number, token: string): Promise<{ data?: Strategy; error?: string }> {
		return this.request<Strategy>(`/api/strategies/${id}`, {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}

	async createStrategy(
		req: CreateStrategyRequest,
		token: string
	): Promise<{ data?: Strategy; error?: string }> {
		return this.request<Strategy>('/api/strategies', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(req)
		});
	}

	async updateStrategy(
		id: number,
		req: UpdateStrategyRequest,
		token: string
	): Promise<{ data?: Strategy; error?: string }> {
		return this.request<Strategy>(`/api/strategies/${id}`, {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(req)
		});
	}

	async deleteStrategy(id: number, token: string): Promise<{ data?: any; error?: string }> {
		return this.request<any>(`/api/strategies/${id}`, {
			method: 'DELETE',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}

	// Trade APIs
	async getTrades(token: string): Promise<{ data?: Trade[]; error?: string }> {
		return this.request<Trade[]>('/api/trades', {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}

	async getTrade(id: number, token: string): Promise<{ data?: Trade; error?: string }> {
		return this.request<Trade>(`/api/trades/${id}`, {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}

	async createTrade(
		req: CreateTradeRequest,
		token: string
	): Promise<{ data?: Trade; error?: string }> {
		return this.request<Trade>('/api/trades', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(req)
		});
	}

	async updateTrade(
		id: number,
		req: UpdateTradeRequest,
		token: string
	): Promise<{ data?: Trade; error?: string }> {
		return this.request<Trade>(`/api/trades/${id}`, {
			method: 'PUT',
			headers: {
				Authorization: `Bearer ${token}`
			},
			body: JSON.stringify(req)
		});
	}

	async deleteTrade(id: number, token: string): Promise<{ data?: any; error?: string }> {
		return this.request<any>(`/api/trades/${id}`, {
			method: 'DELETE',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}

	// Analytics APIs
	async getAnalytics(token: string): Promise<{ data?: Analytics; error?: string }> {
		return this.request<Analytics>('/api/analytics', {
			method: 'GET',
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
	}
}

export const apiClient = new ApiClient(API_BASE_URL);
