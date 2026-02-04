// Notebook and related domain types

export interface FSRSSettings {
	desiredRetention: number;
	version: number;
	params: number[];
	learningSteps: number[];
	relearningSteps: number[];
	maximumInterval: number;
	enableFuzzing: boolean;
}

export interface Notebook {
	// === API fields (from ApiNotebook, snake_case -> camelCase) ===
	id: string;
	name: string;
	description: string | null;
	emoji: string;
	color: string;
	fsrsSettings: FSRSSettings;
	position: number;
	createdAt: string;
	updatedAt: string;

	// === UI-only fields (computed/mock until stats API exists) ===
	dueCount: number; // Mock until stats endpoint
	streak: number; // Mock until stats endpoint
	totalCards: number; // Mock until stats endpoint
	retention: number; // Mock until stats endpoint (uses fsrsSettings.desiredRetention as fallback)
}

/**
 * DisplayCard - UI-friendly card representation for display in lists, grids, and search.
 * Constructed by combining Card (FSRS scheduling) + Fact (content) data.
 * Used throughout the UI where card content needs to be displayed.
 * For raw FSRS scheduling data, use Card from '$lib/types/fact'.
 */
export interface DisplayCard {
	id: string;
	notebookId: string;
	sourceId: string | null;
	front: string;
	back: string;
	due: boolean;
	interval: string;
	tags: string[];
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

