// Mock notebook data for development

import type { Notebook } from '$lib/types';

export const mockNotebooks: Notebook[] = [
	{
		id: 'nb-biology-101',
		name: 'Biology 101',
		emoji: '🧬',
		color: 'emerald',
		dueCount: 23,
		streak: 12,
		totalCards: 156,
		retention: 87
	},
	{
		id: 'nb-spanish-b2',
		name: 'Spanish B2',
		emoji: '🇪🇸',
		color: 'amber',
		dueCount: 45,
		streak: 8,
		totalCards: 412,
		retention: 82
	},
	{
		id: 'nb-calculus',
		name: 'Calculus',
		emoji: '📐',
		color: 'blue',
		dueCount: 15,
		streak: 5,
		totalCards: 89,
		retention: 91
	},
	{
		id: 'nb-us-history',
		name: 'US History',
		emoji: '🏛️',
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
