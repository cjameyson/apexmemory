import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiNotebook } from '$lib/api/types';
import { toNotebook } from '$lib/services/notebooks';
import type { Source, Card } from '$lib/types';

export const load: PageServerLoad = async ({ params, locals }) => {
	const result = await apiRequest<ApiNotebook>(`/v1/notebooks/${params.id}`, {
		token: locals.sessionToken!
	});

	if (!result.ok) {
		if (result.status === 404) {
			error(404, { message: 'Notebook not found' });
		}
		error(500, { message: 'Failed to load notebook' });
	}

	const notebook = toNotebook(result.data);

	// TODO: Fetch from API when sources/cards endpoints exist
	const sources: Source[] = [];
	const cards: Card[] = [];

	return {
		notebook,
		sources,
		cards
	};
};
