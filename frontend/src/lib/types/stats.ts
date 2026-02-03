// Stats and review scope types

import type { Notebook, Source } from './notebook';
import type { ReviewMode } from './review';

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
	| { type: 'all'; mode: ReviewMode }
	| { type: 'notebook'; notebook: Notebook; mode: ReviewMode }
	| { type: 'source'; notebook: Notebook; source: Source; mode: ReviewMode };
