import { error, redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiNotebook } from '$lib/api/types';
import { toNotebook } from '$lib/services/notebooks';
import { getMockSourcesForNotebook, getMockCardsForNotebook } from '$lib/mocks/notebook-content';

export const load: LayoutServerLoad = async ({ params, locals }) => {
	if (!locals.sessionToken) {
		redirect(302, '/login');
	}

	const result = await apiRequest<ApiNotebook>(`/v1/notebooks/${params.id}`, {
		token: locals.sessionToken
	});

	if (!result.ok) {
		if (result.status === 404) {
			error(404, { message: 'Notebook not found' });
		}
		if (result.status === 401) {
			redirect(302, '/login');
		}
		console.error('Failed to load notebook:', result.error);
		error(500, { message: 'Failed to load notebook' });
	}

	let notebook;
	try {
		notebook = toNotebook(result.data);
	} catch (e) {
		console.error('Failed to transform notebook:', e);
		error(500, { message: 'Failed to process notebook data' });
	}

	// TODO: Replace with API calls when sources/cards endpoints exist
	const sources = getMockSourcesForNotebook(params.id);
	const cards = getMockCardsForNotebook(params.id);

	return {
		notebook,
		sources,
		cards
	};
};
