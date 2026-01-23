import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	// If user is authenticated, redirect to app home
	if (locals.user) {
		redirect(302, '/home');
	}

	return {};
};
