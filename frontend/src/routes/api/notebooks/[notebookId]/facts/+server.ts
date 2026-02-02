import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiFactDetail } from '$lib/api/types';

export const POST: RequestHandler = async ({ params, request, locals }) => {
	const token = locals.sessionToken;
	if (!token) {
		error(401, { message: 'Unauthorized' });
	}

	let body: unknown;
	try {
		body = await request.json();
	} catch {
		error(400, { message: 'Invalid JSON in request body' });
	}

	const result = await apiRequest<ApiFactDetail>(
		`/v1/notebooks/${params.notebookId}/facts`,
		{
			method: 'POST',
			token,
			body
		}
	);

	if (!result.ok) {
		error(result.status, { message: result.error.error, ...result.error });
	}

	return json(result.data);
};
