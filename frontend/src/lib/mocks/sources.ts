// Mock source data for development

import type { Source } from '$lib/types';

const sourcesMap: Record<string, Source[]> = {
	'nb-biology-101': [
		{
			id: 'src-cell-biology-pdf',
			notebookId: 'nb-biology-101',
			name: 'Cell Biology Textbook Ch.1-3',
			type: 'pdf',
			cards: 45,
			excerpt: 'The cell is the basic structural and functional unit of all living organisms...',
			pages: 87,
			addedAt: '2024-01-15T10:30:00Z'
		},
		{
			id: 'src-mitosis-video',
			notebookId: 'nb-biology-101',
			name: 'Mitosis Explained - Khan Academy',
			type: 'youtube',
			cards: 28,
			excerpt: 'Learn about the stages of cell division and how DNA is replicated...',
			duration: '14:32',
			addedAt: '2024-01-18T14:00:00Z'
		},
		{
			id: 'src-genetics-notes',
			notebookId: 'nb-biology-101',
			name: 'Genetics Lecture Notes',
			type: 'notes',
			cards: 33,
			excerpt: 'Key concepts: dominant/recessive alleles, Punnett squares, inheritance patterns',
			addedAt: '2024-01-20T09:15:00Z'
		},
		{
			id: 'src-evolution-article',
			notebookId: 'nb-biology-101',
			name: 'Natural Selection - Nature.com',
			type: 'url',
			cards: 12,
			excerpt: 'Darwin theory of evolution through natural selection remains foundational...',
			addedAt: '2024-01-22T16:45:00Z'
		}
	],
	'nb-spanish-b2': [
		{
			id: 'src-spanish-podcast',
			notebookId: 'nb-spanish-b2',
			name: 'Spanish Pod 101 - Episode 42',
			type: 'audio',
			cards: 65,
			excerpt: 'Conversational Spanish with native speakers discussing everyday topics...',
			duration: '28:15',
			addedAt: '2024-01-10T08:00:00Z'
		},
		{
			id: 'src-grammar-pdf',
			notebookId: 'nb-spanish-b2',
			name: 'Advanced Spanish Grammar',
			type: 'pdf',
			cards: 120,
			excerpt: 'Subjunctive mood, conditional tenses, and complex sentence structures...',
			pages: 156,
			addedAt: '2024-01-12T11:30:00Z'
		},
		{
			id: 'src-vocab-notes',
			notebookId: 'nb-spanish-b2',
			name: 'B2 Vocabulary List',
			type: 'notes',
			cards: 227,
			excerpt: 'Essential vocabulary for B2 level proficiency including idioms and expressions',
			addedAt: '2024-01-08T15:00:00Z'
		}
	],
	'nb-calculus': [
		{
			id: 'src-calc-textbook',
			notebookId: 'nb-calculus',
			name: 'Calculus: Early Transcendentals',
			type: 'pdf',
			cards: 56,
			excerpt: 'Limits, derivatives, and integrals with practical applications...',
			pages: 342,
			addedAt: '2024-01-05T10:00:00Z'
		},
		{
			id: 'src-derivatives-video',
			notebookId: 'nb-calculus',
			name: '3Blue1Brown - Essence of Calculus',
			type: 'youtube',
			cards: 33,
			excerpt: 'Visual intuition for derivatives and integrals...',
			duration: '17:04',
			addedAt: '2024-01-14T20:00:00Z'
		}
	],
	'nb-us-history': [
		{
			id: 'src-revolution-pdf',
			notebookId: 'nb-us-history',
			name: 'American Revolution 1765-1783',
			type: 'pdf',
			cards: 89,
			excerpt: 'From the Stamp Act to the Treaty of Paris...',
			pages: 234,
			addedAt: '2024-01-02T09:00:00Z'
		},
		{
			id: 'src-civil-war-article',
			notebookId: 'nb-us-history',
			name: 'Civil War Overview - History.com',
			type: 'url',
			cards: 67,
			excerpt: 'Causes, key battles, and lasting impact of the American Civil War...',
			addedAt: '2024-01-06T14:30:00Z'
		},
		{
			id: 'src-constitution-notes',
			notebookId: 'nb-us-history',
			name: 'Constitution Study Notes',
			type: 'notes',
			cards: 78,
			excerpt: 'Articles, amendments, and key Supreme Court interpretations',
			addedAt: '2024-01-09T11:00:00Z'
		}
	]
};

export function getSourcesForNotebook(notebookId: string): Source[] {
	return sourcesMap[notebookId] || [];
}

export function getSource(notebookId: string, sourceId: string): Source | undefined {
	const sources = sourcesMap[notebookId] || [];
	return sources.find((s) => s.id === sourceId);
}

export function getAllSources(): Source[] {
	return Object.values(sourcesMap).flat();
}
