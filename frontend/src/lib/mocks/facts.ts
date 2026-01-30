import type { FactDetail, FactType, CardState } from '$lib/types/fact';

function makeCard(
	id: string,
	factId: string,
	notebookId: string,
	state: CardState,
	opts: { due?: string; stability?: number; difficulty?: number; reps?: number; lapses?: number } = {}
) {
	const now = new Date();
	return {
		id,
		factId,
		notebookId,
		elementId: `el-${id}`,
		state,
		stability: opts.stability ?? (state === 'new' ? null : 5.0),
		difficulty: opts.difficulty ?? (state === 'new' ? null : 5.5),
		due: opts.due ?? (state === 'new' ? null : new Date(now.getTime() + Math.random() * 7 * 86400000).toISOString()),
		reps: opts.reps ?? (state === 'new' ? 0 : 3),
		lapses: opts.lapses ?? 0,
		suspendedAt: null,
		buriedUntil: null,
		createdAt: '2024-12-01T10:00:00Z',
		updatedAt: '2025-01-15T08:00:00Z'
	};
}

const factsMap: Record<string, FactDetail[]> = {
	'nb-biology-101': [
		{
			id: 'fact-bio-1',
			notebookId: 'nb-biology-101',
			factType: 'basic',
			content: { front: 'What is the powerhouse of the cell?', back: 'The mitochondria. It produces ATP through cellular respiration.' },
			sourceId: 'src-cell-biology-pdf',
			cardCount: 1,
			tags: ['cell-biology', 'organelles'],
			dueCount: 1,
			createdAt: '2024-12-01T10:00:00Z',
			updatedAt: '2025-01-20T08:00:00Z',
			cards: [makeCard('card-bio-1a', 'fact-bio-1', 'nb-biology-101', 'review', { due: '2025-01-21T00:00:00Z', reps: 5, stability: 12.3 })]
		},
		{
			id: 'fact-bio-2',
			notebookId: 'nb-biology-101',
			factType: 'cloze',
			content: { text: 'The {{c1::cell membrane}} controls what enters and exits the cell. It is a {{c2::semi-permeable}} barrier.' },
			sourceId: 'src-cell-biology-pdf',
			cardCount: 2,
			tags: ['cell-biology', 'membrane'],
			dueCount: 1,
			createdAt: '2024-12-02T10:00:00Z',
			updatedAt: '2025-01-19T08:00:00Z',
			cards: [
				makeCard('card-bio-2a', 'fact-bio-2', 'nb-biology-101', 'learning', { due: '2025-01-20T12:00:00Z', reps: 1 }),
				makeCard('card-bio-2b', 'fact-bio-2', 'nb-biology-101', 'new')
			]
		},
		{
			id: 'fact-bio-3',
			notebookId: 'nb-biology-101',
			factType: 'basic',
			content: { front: 'Name the four phases of mitosis in order.', back: 'Prophase, Metaphase, Anaphase, Telophase (PMAT)' },
			sourceId: 'src-mitosis-video',
			cardCount: 1,
			tags: ['mitosis', 'cell-division'],
			dueCount: 0,
			createdAt: '2024-12-05T10:00:00Z',
			updatedAt: '2025-01-18T08:00:00Z',
			cards: [makeCard('card-bio-3a', 'fact-bio-3', 'nb-biology-101', 'review', { due: '2025-02-01T00:00:00Z', reps: 8, stability: 30 })]
		},
		{
			id: 'fact-bio-4',
			notebookId: 'nb-biology-101',
			factType: 'cloze',
			content: { text: 'A {{c1::dominant}} allele expresses its phenotype even when paired with a {{c2::recessive}} allele.' },
			sourceId: 'src-genetics-notes',
			cardCount: 2,
			tags: ['genetics', 'alleles', 'heredity'],
			dueCount: 2,
			createdAt: '2024-12-08T10:00:00Z',
			updatedAt: '2025-01-22T08:00:00Z',
			cards: [
				makeCard('card-bio-4a', 'fact-bio-4', 'nb-biology-101', 'relearning', { due: '2025-01-20T14:00:00Z', reps: 4, lapses: 2 }),
				makeCard('card-bio-4b', 'fact-bio-4', 'nb-biology-101', 'learning', { due: '2025-01-20T16:00:00Z', reps: 1 })
			]
		},
		{
			id: 'fact-bio-5',
			notebookId: 'nb-biology-101',
			factType: 'image_occlusion',
			content: { title: 'Cell Organelle Diagram', imageUrl: '/mock/cell-diagram.png', regions: 5 },
			sourceId: 'src-cell-biology-pdf',
			cardCount: 5,
			tags: ['cell-biology', 'organelles', 'diagram', 'visual'],
			dueCount: 3,
			createdAt: '2024-12-10T10:00:00Z',
			updatedAt: '2025-01-21T08:00:00Z',
			cards: [
				makeCard('card-bio-5a', 'fact-bio-5', 'nb-biology-101', 'new'),
				makeCard('card-bio-5b', 'fact-bio-5', 'nb-biology-101', 'review', { due: '2025-01-20T00:00:00Z', reps: 3 }),
				makeCard('card-bio-5c', 'fact-bio-5', 'nb-biology-101', 'review', { due: '2025-01-21T00:00:00Z', reps: 2 }),
				makeCard('card-bio-5d', 'fact-bio-5', 'nb-biology-101', 'learning', { due: '2025-01-20T18:00:00Z' }),
				makeCard('card-bio-5e', 'fact-bio-5', 'nb-biology-101', 'new')
			]
		},
		{
			id: 'fact-bio-6',
			notebookId: 'nb-biology-101',
			factType: 'basic',
			content: { front: 'What does the Punnett square show?', back: 'Predicts the probability of offspring genotypes from a genetic cross.' },
			sourceId: 'src-genetics-notes',
			cardCount: 1,
			tags: ['genetics'],
			dueCount: 0,
			createdAt: '2024-12-12T10:00:00Z',
			updatedAt: '2025-01-17T08:00:00Z',
			cards: [makeCard('card-bio-6a', 'fact-bio-6', 'nb-biology-101', 'review', { due: '2025-02-05T00:00:00Z', reps: 10, stability: 45 })]
		},
		{
			id: 'fact-bio-7',
			notebookId: 'nb-biology-101',
			factType: 'basic',
			content: { front: 'What is DNA replication?', back: 'The process by which DNA makes a copy of itself during cell division.' },
			sourceId: 'src-cell-biology-pdf',
			cardCount: 1,
			tags: ['dna', 'cell-biology'],
			dueCount: 1,
			createdAt: '2024-12-15T10:00:00Z',
			updatedAt: '2025-01-23T08:00:00Z',
			cards: [makeCard('card-bio-7a', 'fact-bio-7', 'nb-biology-101', 'learning', { due: '2025-01-20T10:00:00Z', reps: 2 })]
		}
	]
};

// Generate generic facts for notebooks that don't have specific data
function generateGenericFacts(notebookId: string): FactDetail[] {
	const types: FactType[] = ['basic', 'basic', 'cloze', 'basic', 'cloze', 'image_occlusion'];
	return types.map((factType, i) => {
		const id = `fact-gen-${notebookId}-${i}`;
		const states: CardState[] = ['new', 'learning', 'review', 'relearning'];
		const cards = factType === 'image_occlusion'
			? [1, 2, 3].map((j) => makeCard(`card-gen-${id}-${j}`, id, notebookId, states[j % 4]))
			: factType === 'cloze'
				? [1, 2].map((j) => makeCard(`card-gen-${id}-${j}`, id, notebookId, states[j % 4]))
				: [makeCard(`card-gen-${id}-1`, id, notebookId, states[i % 4])];

		const contentMap: Record<FactType, Record<string, unknown>> = {
			basic: { front: `Sample question ${i + 1}`, back: `Sample answer ${i + 1}` },
			cloze: { text: `The {{c1::answer}} is found in the {{c2::source material}}.` },
			image_occlusion: { title: `Diagram ${i + 1}`, imageUrl: '/mock/diagram.png', regions: 3 }
		};

		return {
			id,
			notebookId,
			factType,
			content: contentMap[factType],
			sourceId: null,
			cardCount: cards.length,
			tags: ['auto-generated'],
			dueCount: cards.filter((c) => c.state !== 'new' && c.due && new Date(c.due) <= new Date()).length,
			createdAt: new Date(2024, 11, 1 + i).toISOString(),
			updatedAt: new Date(2025, 0, 15 + i).toISOString(),
			cards
		};
	});
}

export function getFactsForNotebook(notebookId: string): FactDetail[] {
	return factsMap[notebookId] ?? generateGenericFacts(notebookId);
}

export function getFactDetailById(notebookId: string, factId: string): FactDetail | undefined {
	return getFactsForNotebook(notebookId).find((f) => f.id === factId);
}

export function countFactsByType(notebookId: string): Record<FactType, number> {
	const facts = getFactsForNotebook(notebookId);
	return {
		basic: facts.filter((f) => f.factType === 'basic').length,
		cloze: facts.filter((f) => f.factType === 'cloze').length,
		image_occlusion: facts.filter((f) => f.factType === 'image_occlusion').length
	};
}
