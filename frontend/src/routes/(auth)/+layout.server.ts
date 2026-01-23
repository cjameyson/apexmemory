import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	// Redirect authenticated users away from auth pages
	if (locals.user) {
		redirect(302, '/home');
	}

	return {
		user: null,
	};
};
