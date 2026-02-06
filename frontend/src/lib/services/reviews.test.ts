import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { ApiStudyCard } from '$lib/api/types';
import { fetchPracticeCards, fetchStudyCards } from './reviews';

const baseApiCard: ApiStudyCard = {
	id: 'card-1',
	fact_id: 'fact-1',
	notebook_id: 'notebook-1',
	element_id: 'c1',
	state: 'new',
	stability: null,
	difficulty: null,
	due: null,
	reps: 0,
	lapses: 0,
	suspended_at: null,
	buried_until: null,
	created_at: '2026-02-01T00:00:00Z',
	updated_at: '2026-02-01T00:00:00Z',
	fact_type: 'cloze',
	fact_content: { version: 1, fields: [] },
	intervals: { again: '1m', hard: '5m', good: '10m', easy: '1d' }
};

describe('review services', () => {
	beforeEach(() => {
		vi.restoreAllMocks();
	});

	it('maps scheduled study cards from array responses', async () => {
		const fetchMock = vi.fn().mockResolvedValue(
			new Response(JSON.stringify([baseApiCard]), {
				status: 200,
				headers: { 'Content-Type': 'application/json' }
			})
		);
		vi.stubGlobal('fetch', fetchMock);

		const cards = await fetchStudyCards('notebook-1');

		expect(fetchMock).toHaveBeenCalledWith('/api/reviews/study?limit=50&notebook_id=notebook-1');
		expect(cards).toHaveLength(1);
		expect(cards[0]).toMatchObject({
			id: 'card-1',
			factId: 'fact-1',
			notebookId: 'notebook-1',
			factType: 'cloze'
		});
	});

	it('maps practice cards from paginated responses', async () => {
		const fetchMock = vi.fn().mockResolvedValue(
			new Response(
				JSON.stringify({
					data: [baseApiCard],
					total: 1,
					has_more: false
				}),
				{
					status: 200,
					headers: { 'Content-Type': 'application/json' }
				}
			)
		);
		vi.stubGlobal('fetch', fetchMock);

		const cards = await fetchPracticeCards('notebook-1');

		expect(fetchMock).toHaveBeenCalledWith('/api/reviews/practice?limit=50&notebook_id=notebook-1');
		expect(cards).toHaveLength(1);
		expect(cards[0].id).toBe('card-1');
	});

	it('throws API message instead of silently returning empty cards', async () => {
		const fetchMock = vi.fn().mockResolvedValue(
			new Response(JSON.stringify({ message: 'Service unavailable' }), {
				status: 503,
				headers: { 'Content-Type': 'application/json' }
			})
		);
		vi.stubGlobal('fetch', fetchMock);

		await expect(fetchStudyCards()).rejects.toThrow('Service unavailable');
	});
});
