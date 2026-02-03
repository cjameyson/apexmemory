import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiNotebook, ApiStudyCountsResponse } from '$lib/api/types';

export const load: LayoutServerLoad = async ({ locals }) => {
	// Redirect unauthenticated users to login
	if (!locals.user) {
		redirect(302, '/login');
	}

	// Fetch notebooks and study counts in parallel
	const [notebooksResult, countsResult] = await Promise.all([
		apiRequest<ApiNotebook[]>('/v1/notebooks?ui=true', {
			token: locals.sessionToken!
		}),
		apiRequest<ApiStudyCountsResponse>('/v1/reviews/study-counts', {
			token: locals.sessionToken!
		})
	]);

	return {
		user: locals.user,
		notebooks: notebooksResult.ok ? notebooksResult.data : [],
		studyCounts: countsResult.ok
			? countsResult.data
			: { counts: {}, total_due: 0, total_new: 0 }
	};
};
