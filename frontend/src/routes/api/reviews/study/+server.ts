import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiStudyCard } from '$lib/api/types';

export const GET: RequestHandler = async ({ url, locals }) => {
	const token = locals.sessionToken;
	if (!token) {
		error(401, { message: 'Unauthorized' });
	}

	const notebookId = url.searchParams.get('notebook_id');
	const limit = url.searchParams.get('limit') ?? '50';

	const params = new URLSearchParams({ limit });
	if (notebookId) params.set('notebook_id', notebookId);

	const result = await apiRequest<ApiStudyCard[]>(
		`/v1/reviews/study?${params}`,
		{ token }
	);

	if (!result.ok) {
		error(result.status, { message: result.error.error, ...result.error });
	}

	return json(result.data);
};
