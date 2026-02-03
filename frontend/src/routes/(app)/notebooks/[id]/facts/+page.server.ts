import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { apiRequest } from '$lib/server/api';
import type { ApiFactsListWithStats } from '$lib/api/types';
import { toFact } from '$lib/services/facts';

export const load: PageServerLoad = async ({ params, locals, url }) => {
	if (!locals.sessionToken) {
		redirect(302, '/login');
	}

	// Read filter params from URL
	const typeFilter = url.searchParams.get('type') || '';
	const search = url.searchParams.get('q') || '';
	const sort = url.searchParams.get('sort') || '';
	const page = Math.max(1, parseInt(url.searchParams.get('page') || '1', 10));
	const pageSize = 20;
	const offset = (page - 1) * pageSize;

	// Build API query string
	const apiParams = new URLSearchParams({
		limit: String(pageSize),
		offset: String(offset),
		stats: 'true'
	});
	if (typeFilter) apiParams.set('type', typeFilter);
	if (search) apiParams.set('q', search);
	if (sort) apiParams.set('sort', sort);

	const result = await apiRequest<ApiFactsListWithStats>(
		`/v1/notebooks/${params.id}/facts?${apiParams}`,
		{ token: locals.sessionToken }
	);

	if (!result.ok) {
		if (result.status === 401) redirect(302, '/login');
		if (result.status === 404) error(404, { message: 'Notebook not found' });
		console.error('Failed to load facts:', result.error);
		error(500, { message: 'Failed to load facts' });
	}

	const { data: apiFacts, total, stats: apiStats } = result.data;

	return {
		notebookId: params.id,
		facts: apiFacts.map(toFact),
		pagination: {
			page,
			pageSize,
			total,
			totalPages: Math.ceil(total / pageSize)
		},
		stats: {
			totalFacts: apiStats.total_facts,
			totalCards: apiStats.total_cards,
			totalDue: apiStats.total_due,
			byType: {
				basic: apiStats.by_type.basic,
				cloze: apiStats.by_type.cloze,
				imageOcclusion: apiStats.by_type.image_occlusion
			}
		}
	};
};
