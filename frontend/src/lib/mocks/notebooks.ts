// Mock notebook data for development

import type { Notebook } from '$lib/types';

const now = new Date().toISOString();
const dayAgo = new Date(Date.now() - 86400000).toISOString();
const weekAgo = new Date(Date.now() - 604800000).toISOString();

export const mockNotebooks: Notebook[] = [
	{
		id: 'nb-biology-101',
		name: 'Biology 101',
		description: 'Cell biology, genetics, and evolution fundamentals',
		desiredRetention: 0.9,
		position: 0,
		createdAt: weekAgo,
		updatedAt: dayAgo,
		emoji: '\u{1F9EC}',
		color: 'emerald',
		dueCount: 23,
		streak: 12,
		totalCards: 156,
		retention: 87
	},
	{
		id: 'nb-spanish-b2',
		name: 'Spanish B2',
		description: 'Intermediate Spanish vocabulary and grammar',
		desiredRetention: 0.85,
		position: 1,
		createdAt: weekAgo,
		updatedAt: now,
		emoji: '\u{1F1EA}\u{1F1F8}',
		color: 'amber',
		dueCount: 45,
		streak: 8,
		totalCards: 412,
		retention: 82
	},
	{
		id: 'nb-calculus',
		name: 'Calculus',
		description: null,
		desiredRetention: 0.9,
		position: 2,
		createdAt: weekAgo,
		updatedAt: weekAgo,
		emoji: '\u{1F4D0}',
		color: 'blue',
		dueCount: 15,
		streak: 5,
		totalCards: 89,
		retention: 91
	},
	{
		id: 'nb-us-history',
		name: 'US History',
		description: 'American history from colonial era to present',
		desiredRetention: 0.85,
		position: 3,
		createdAt: weekAgo,
		updatedAt: dayAgo,
		emoji: '\u{1F3DB}\u{FE0F}',
		color: 'slate',
		dueCount: 5,
		streak: 3,
		totalCards: 234,
		retention: 78
	}
];

export function getNotebook(id: string): Notebook | undefined {
	return mockNotebooks.find((nb) => nb.id === id);
}

export function getAllNotebooks(): Notebook[] {
	return mockNotebooks;
}

export function getTotalDueCount(): number {
	return mockNotebooks.reduce((sum, nb) => sum + nb.dueCount, 0);
}
