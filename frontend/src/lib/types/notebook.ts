// Notebook and related domain types

export interface Notebook {
	id: string;
	name: string;
	emoji: string;
	color: string;
	dueCount: number;
	streak: number;
	totalCards: number;
	retention: number;
}

export type SourceType = 'pdf' | 'youtube' | 'url' | 'audio' | 'notes';

export interface Source {
	id: string;
	notebookId: string;
	name: string;
	type: SourceType;
	cards: number;
	excerpt: string;
	pages?: number;
	duration?: string;
	addedAt: string;
}

export interface Card {
	id: string;
	notebookId: string;
	sourceId: string;
	front: string;
	back: string;
	due: boolean;
	interval: string;
	tags: string[];
}
