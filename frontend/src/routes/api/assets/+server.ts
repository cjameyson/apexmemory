import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { API_BASE_URL } from '$env/static/private';

export const POST: RequestHandler = async ({ request, locals }) => {
	const token = locals.sessionToken;
	if (!token) {
		error(401, { message: 'Unauthorized' });
	}

	const formData = await request.formData();

	const response = await fetch(`${API_BASE_URL}/v1/assets`, {
		method: 'POST',
		headers: {
			Authorization: `Bearer ${token}`
		},
		body: formData
	});

	if (!response.ok) {
		const err = await response.json().catch(() => ({ error: 'Upload failed' }));
		error(response.status, { message: err.error });
	}

	const data = await response.json();
	return json(data, { status: 201 });
};
