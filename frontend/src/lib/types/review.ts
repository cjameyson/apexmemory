import type { CardState, FactType } from './fact';

export type ReviewMode = 'scheduled' | 'practice';

export interface StudyCard {
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
	factType: FactType;
	factContent: Record<string, unknown>;
	intervals: { again: string; hard: string; good: string; easy: string };
}

export interface CardDisplay {
	front: string;
	back: string;
}
