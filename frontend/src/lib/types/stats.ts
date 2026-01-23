// Stats and review scope types

import type { Notebook, Source } from './notebook';

export interface GlobalStats {
	totalCards: number;
	cardsReviewedToday: number;
	currentStreak: number;
	longestStreak: number;
	averageRetention: number;
	reviewsThisWeek: number[];
	upcomingDue: {
		today: number;
		tomorrow: number;
		thisWeek: number;
	};
}

export type ReviewScope =
	| { type: 'all' }
	| { type: 'notebook'; notebook: Notebook }
	| { type: 'source'; notebook: Notebook; source: Source };
