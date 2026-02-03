import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiStudyCountsResponse } from '$lib/api/types';

export const GET: RequestHandler = async ({ locals }) => {
	const token = locals.sessionToken;
	if (!token) {
		error(401, { message: 'Unauthorized' });
	}

	const result = await apiRequest<ApiStudyCountsResponse>(
		'/v1/reviews/study-counts',
		{ token }
	);

	if (!result.ok) {
		error(result.status, { message: result.error.error, ...result.error });
	}

	return json(result.data);
};
