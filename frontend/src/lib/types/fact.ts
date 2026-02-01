export type FactType = 'basic' | 'cloze' | 'image_occlusion';
export type CardState = 'new' | 'learning' | 'review' | 'relearning';

export type FactTypeFilter = 'all' | FactType;
export type FactSortField = 'updated' | 'created' | 'cards' | 'due';

export interface FactStats {
	totalFacts: number;
	totalCards: number;
	totalDue: number;
	byType: { basic: number; cloze: number; imageOcclusion: number };
}

export interface Fact {
	id: string;
	notebookId: string;
	factType: FactType;
	content: Record<string, unknown>;
	sourceId: string | null;
	cardCount: number;
	tags: string[];
	dueCount: number;
	createdAt: string;
	updatedAt: string;
}

export interface FactDetail extends Fact {
	cards: Card[];
}

export interface Card {
	id: string;
	factId: string;
	notebookId: string;
	elementId: string;
	state: CardState;
	stability: number | null;
	difficulty: number | null;
	due: string | null;
	reps: number;
	lapses: number;
	suspendedAt: string | null;
	buriedUntil: string | null;
	createdAt: string;
	updatedAt: string;
}
