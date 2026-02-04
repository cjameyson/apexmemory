// Mock card data for development

import type { DisplayCard } from '$lib/types';

const cardsMap: Record<string, DisplayCard[]> = {
	'nb-biology-101': [
		{
			id: 'card-bio-1',
			notebookId: 'nb-biology-101',
			sourceId: 'src-cell-biology-pdf',
			front: 'What is the powerhouse of the cell?',
			back: 'The mitochondria. It produces ATP through cellular respiration.',
			due: true,
			interval: '1d',
			tags: ['cell-biology', 'organelles']
		},
		{
			id: 'card-bio-2',
			notebookId: 'nb-biology-101',
			sourceId: 'src-cell-biology-pdf',
			front: 'What is the function of the cell membrane?',
			back: 'Controls what enters and exits the cell; semi-permeable barrier.',
			due: true,
			interval: '4h',
			tags: ['cell-biology', 'membrane']
		},
		{
			id: 'card-bio-3',
			notebookId: 'nb-biology-101',
			sourceId: 'src-mitosis-video',
			front: 'Name the four phases of mitosis in order.',
			back: 'Prophase, Metaphase, Anaphase, Telophase (PMAT)',
			due: false,
			interval: '7d',
			tags: ['mitosis', 'cell-division']
		},
		{
			id: 'card-bio-4',
			notebookId: 'nb-biology-101',
			sourceId: 'src-genetics-notes',
			front: 'What is a dominant allele?',
			back: 'An allele that expresses its phenotype even when paired with a recessive allele.',
			due: true,
			interval: '2d',
			tags: ['genetics', 'alleles']
		},
		{
			id: 'card-bio-5',
			notebookId: 'nb-biology-101',
			sourceId: 'src-genetics-notes',
			front: 'What does the Punnett square show?',
			back: 'Predicts the probability of offspring genotypes from a genetic cross.',
			due: false,
			interval: '14d',
			tags: ['genetics', 'heredity']
		}
	],
	'nb-spanish-b2': [
		{
			id: 'card-esp-1',
			notebookId: 'nb-spanish-b2',
			sourceId: 'src-vocab-notes',
			front: 'Translate: "I would have gone if I had known"',
			back: 'Habria ido si hubiera sabido (conditional perfect + pluperfect subjunctive)',
			due: true,
			interval: '1d',
			tags: ['conditional', 'subjunctive']
		},
		{
			id: 'card-esp-2',
			notebookId: 'nb-spanish-b2',
			sourceId: 'src-grammar-pdf',
			front: 'When do you use the subjunctive mood?',
			back: 'For wishes, doubts, emotions, recommendations, and hypothetical situations.',
			due: true,
			interval: '6h',
			tags: ['grammar', 'subjunctive']
		},
		{
			id: 'card-esp-3',
			notebookId: 'nb-spanish-b2',
			sourceId: 'src-spanish-podcast',
			front: 'What does "echar de menos" mean?',
			back: 'To miss (someone or something)',
			due: false,
			interval: '21d',
			tags: ['idioms', 'vocabulary']
		}
	],
	'nb-calculus': [
		{
			id: 'card-calc-1',
			notebookId: 'nb-calculus',
			sourceId: 'src-calc-textbook',
			front: 'What is the derivative of e^x?',
			back: 'e^x (the function is its own derivative)',
			due: true,
			interval: '3d',
			tags: ['derivatives', 'exponentials']
		},
		{
			id: 'card-calc-2',
			notebookId: 'nb-calculus',
			sourceId: 'src-derivatives-video',
			front: 'What does the derivative represent geometrically?',
			back: 'The slope of the tangent line to the curve at a given point.',
			due: false,
			interval: '10d',
			tags: ['derivatives', 'concepts']
		}
	],
	'nb-us-history': [
		{
			id: 'card-hist-1',
			notebookId: 'nb-us-history',
			sourceId: 'src-revolution-pdf',
			front: 'What event sparked the American Revolution?',
			back: 'The Boston Tea Party (1773) was a key catalyst, protesting taxation without representation.',
			due: true,
			interval: '2d',
			tags: ['revolution', 'events']
		},
		{
			id: 'card-hist-2',
			notebookId: 'nb-us-history',
			sourceId: 'src-constitution-notes',
			front: 'How many amendments are in the Bill of Rights?',
			back: '10 amendments, ratified in 1791.',
			due: false,
			interval: '30d',
			tags: ['constitution', 'amendments']
		}
	]
};

export function getCardsForNotebook(notebookId: string): DisplayCard[] {
	return cardsMap[notebookId] || [];
}

export function getCardsForSource(notebookId: string, sourceId: string): DisplayCard[] {
	const cards = cardsMap[notebookId] || [];
	return cards.filter((c) => c.sourceId === sourceId);
}

export function getDueCardsForNotebook(notebookId: string): DisplayCard[] {
	return getCardsForNotebook(notebookId).filter((c) => c.due);
}

export function getAllDueCards(): DisplayCard[] {
	return Object.values(cardsMap)
		.flat()
		.filter((c) => c.due);
}

export function getCard(cardId: string): DisplayCard | undefined {
	return Object.values(cardsMap)
		.flat()
		.find((c) => c.id === cardId);
}
