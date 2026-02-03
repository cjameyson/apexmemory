import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiStudyCard, PaginatedResponse } from '$lib/api/types';

export const GET: RequestHandler = async ({ url, locals }) => {
	const token = locals.sessionToken;
	if (!token) {
		error(401, { message: 'Unauthorized' });
	}

	const notebookId = url.searchParams.get('notebook_id');
	const limit = url.searchParams.get('limit') ?? '50';
	const offset = url.searchParams.get('offset') ?? '0';

	const params = new URLSearchParams({ limit, offset });
	if (notebookId) params.set('notebook_id', notebookId);

	const result = await apiRequest<PaginatedResponse<ApiStudyCard>>(
		`/v1/reviews/practice?${params}`,
		{ token }
	);

	if (!result.ok) {
		error(result.status, { message: result.error.error, ...result.error });
	}

	return json(result.data);
};
