import { error, redirect } from '@sveltejs/kit';
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

	if (!notebooksResult.ok) {
		if (notebooksResult.status === 401 || notebooksResult.status === 403) {
			redirect(302, '/login');
		}
		console.error('Failed to load notebooks:', notebooksResult.error);
		error(notebooksResult.status || 502, { message: 'Failed to load notebooks' });
	}

	if (!countsResult.ok) {
		if (countsResult.status === 401 || countsResult.status === 403) {
			redirect(302, '/login');
		}
		console.error('Failed to load study counts:', countsResult.error);
		error(countsResult.status || 502, { message: 'Failed to load study counts' });
	}

	return {
		user: locals.user,
		notebooks: notebooksResult.data,
		studyCounts: countsResult.data
	};
};
