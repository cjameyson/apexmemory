import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiUndoReviewResponse } from '$lib/api/types';

export const DELETE: RequestHandler = async ({ params, locals }) => {
	const token = locals.sessionToken;
	if (!token) {
		error(401, { message: 'Unauthorized' });
	}

	const { reviewId } = params;

	const result = await apiRequest<ApiUndoReviewResponse>(
		`/v1/reviews/${reviewId}`,
		{ method: 'DELETE', token }
	);

	if (!result.ok) {
		error(result.status, { message: result.error.error, ...result.error });
	}

	return json(result.data);
};
