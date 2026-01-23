import { API_BASE_URL } from '$env/static/private';
import type { ApiError, ApiResult } from '$lib/api/types';

interface RequestOptions {
	method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
	body?: unknown;
	token?: string;
}

/**
 * Server-side API request function.
 * Use this in +page.server.ts, +layout.server.ts, +server.ts, and hooks.server.ts files.
 *
 * This module is in $lib/server/ to enforce server-only usage at compile time.
 */
export async function apiRequest<T>(
	endpoint: string,
	options: RequestOptions = {}
): Promise<ApiResult<T>> {
	const { method = 'GET', body, token } = options;

	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
	};

	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	try {
		const response = await fetch(`${API_BASE_URL}${endpoint}`, {
			method,
			headers,
			body: body ? JSON.stringify(body) : undefined,
		});

		if (!response.ok) {
			let error: ApiError;
			try {
				error = await response.json();
			} catch {
				error = { error: response.statusText || 'Request failed' };
			}
			return { ok: false, error, status: response.status };
		}

		// Handle 204 No Content
		if (response.status === 204) {
			return { ok: true, data: undefined as T };
		}

		const data = await response.json();
		return { ok: true, data };
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Network error';
		return {
			ok: false,
			error: { error: message },
			status: 0,
		};
	}
}
