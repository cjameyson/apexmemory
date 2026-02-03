import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiNotebook, CreateNotebookRequest } from '$lib/api/types';

export const POST: RequestHandler = async ({ request, locals }) => {
	const token = locals.sessionToken;
	if (!token) {
		error(401, { message: 'Unauthorized' });
	}

	let body: CreateNotebookRequest;
	try {
		body = await request.json();
	} catch {
		error(400, { message: 'Invalid JSON in request body' });
	}

	// Server-side validation
	const fieldErrors: Record<string, string> = {};
	const name = body.name?.trim() ?? '';

	if (!name) {
		fieldErrors.name = 'Name is required';
	} else if (name.length > 255) {
		fieldErrors.name = 'Name must be 255 characters or less';
	}

	if (body.description && body.description.length > 10000) {
		fieldErrors.description = 'Description must be 10,000 characters or less';
	}

	if (Object.keys(fieldErrors).length > 0) {
		return json({ message: 'Validation failed', fieldErrors }, { status: 400 });
	}

	const result = await apiRequest<ApiNotebook>('/v1/notebooks', {
		method: 'POST',
		token,
		body: {
			name,
			description: body.description?.trim() || undefined,
			emoji: body.emoji || undefined,
			color: body.color || undefined
		}
	});

	if (!result.ok) {
		error(result.status, { message: result.error.error, ...result.error });
	}

	return json(result.data);
};
