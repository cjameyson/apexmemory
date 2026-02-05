import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { API_BASE_URL } from '$env/static/private';

export const GET: RequestHandler = async ({ params, request, locals }) => {
	const token = locals.sessionToken;
	if (!token) {
		error(401, { message: 'Unauthorized' });
	}

	const headers: Record<string, string> = {
		Authorization: `Bearer ${token}`
	};

	const ifNoneMatch = request.headers.get('If-None-Match');
	if (ifNoneMatch) {
		headers['If-None-Match'] = ifNoneMatch;
	}

	const response = await fetch(`${API_BASE_URL}/v1/assets/${params.id}/file`, { headers });

	if (response.status === 304) {
		return new Response(null, { status: 304 });
	}

	if (!response.ok) {
		error(response.status, { message: 'Asset not found' });
	}

	return new Response(response.body, {
		status: 200,
		headers: {
			'Content-Type': response.headers.get('Content-Type') ?? 'application/octet-stream',
			'Cache-Control': response.headers.get('Cache-Control') ?? '',
			ETag: response.headers.get('ETag') ?? ''
		}
	});
};
