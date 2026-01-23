import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	// Redirect unauthenticated users to login
	if (!locals.user) {
		redirect(302, '/login');
	}

	return {
		user: locals.user
	};
};
