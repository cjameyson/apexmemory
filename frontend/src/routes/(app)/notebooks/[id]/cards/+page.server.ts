import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { getNotebook, getSourcesForNotebook, getCardsForNotebook } from '$lib/mocks';

export const load: PageServerLoad = async ({ params }) => {
	const notebook = getNotebook(params.id);

	if (!notebook) {
		error(404, { message: 'Notebook not found' });
	}

	const sources = getSourcesForNotebook(params.id);
	const cards = getCardsForNotebook(params.id);

	return {
		notebook,
		sources,
		cards
	};
};
