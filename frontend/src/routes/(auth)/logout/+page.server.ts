import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { apiRequest } from '$lib/server/api';

export const load: PageServerLoad = async ({ cookies, locals }) => {
	const token = locals.sessionToken;

	// Call backend logout endpoint if we have a token
	if (token) {
		await apiRequest('/v1/auth/logout', {
			method: 'POST',
			token
		});
	}

	// Clear the session cookie
	cookies.delete('session_token', { path: '/' });

	// Redirect to login page
	redirect(302, '/login');
};
