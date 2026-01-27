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

// ============================================================================
// Notebook API Types
// ============================================================================

// Mirror backend NotebookResponse exactly
export interface ApiNotebook {
	id: string;
	name: string;
	description: string | null;
	desired_retention: number;
	position: number;
	created_at: string;
	updated_at: string;
}

// POST /v1/notebooks request body
export interface CreateNotebookRequest {
	name: string;
	description?: string;
}

// PATCH /v1/notebooks/{id} request body
export interface UpdateNotebookRequest {
	name?: string;
	description?: string | null; // null = clear, undefined = no change
	desired_retention?: number;
	position?: number;
}

// Discriminated union for API results
export type ApiResult<T> =
	| { ok: true; data: T }
	| { ok: false; error: ApiError; status: number };
