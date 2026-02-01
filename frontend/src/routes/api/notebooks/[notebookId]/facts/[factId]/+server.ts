import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiFactDetail } from '$lib/api/types';

export const GET: RequestHandler = async ({ params, locals }) => {
	const token = locals.sessionToken;
	if (!token) {
		error(401, { message: 'Unauthorized' });
	}

	const result = await apiRequest<ApiFactDetail>(
		`/v1/notebooks/${params.notebookId}/facts/${params.factId}`,
		{ token }
	);

	if (!result.ok) {
		error(result.status, { message: result.error.error });
	}

	return json(result.data);
};
