import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiNotebook } from '$lib/api/types';

export const load: LayoutServerLoad = async ({ locals }) => {
	// Redirect unauthenticated users to login
	if (!locals.user) {
		redirect(302, '/login');
	}

	// Fetch notebooks - graceful fallback on error
	const result = await apiRequest<ApiNotebook[]>('/v1/notebooks', {
		token: locals.sessionToken!
	});

	return {
		user: locals.user,
		notebooks: result.ok ? result.data : []
	};
};
