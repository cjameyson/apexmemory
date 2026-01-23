import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { apiRequest } from '$lib/server/api';

export const POST: RequestHandler = async ({ cookies, locals }) => {
	const token = locals.sessionToken;

	// Call backend logout endpoint if we have a token
	if (token) {
		await apiRequest('/v1/auth/logout', {
			method: 'POST',
			token,
		});
	}

	// Always clear the cookie
	cookies.delete('session_token', { path: '/' });

	redirect(302, '/login');
};
