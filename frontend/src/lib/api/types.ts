// User returned from /v1/auth/me and in auth responses
export interface User {
	id: string;
	email: string;
	username: string;
}

// Response from /v1/auth/login and /v1/auth/register
export interface AuthResponse {
	session_token: string;
	expires_at: string;
	user: User;
}

// API error response structure
export interface ApiError {
	error: string;
	code?: string;
	details?: Record<string, string>;
	fieldErrors?: Record<string, string>;
}

// Discriminated union for API results
export type ApiResult<T> =
	| { ok: true; data: T }
	| { ok: false; error: ApiError; status: number };
