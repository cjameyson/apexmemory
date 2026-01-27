// Notebook and related domain types

export interface Notebook {
	// === API fields (from ApiNotebook, snake_case -> camelCase) ===
	id: string;
	name: string;
	description: string | null;
	desiredRetention: number;
	position: number;
	createdAt: string;
	updatedAt: string;

	// === UI-only fields (computed/mock until stats API exists) ===
	emoji: string; // Static '📘' for now
	color: string; // Deferred - keep for component compatibility
	dueCount: number; // Mock until stats endpoint
	streak: number; // Mock until stats endpoint
	totalCards: number; // Mock until stats endpoint
	retention: number; // Mock until stats endpoint (uses desiredRetention as fallback)
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
