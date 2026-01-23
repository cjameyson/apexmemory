// Mock stats data for development

import type { GlobalStats } from '$lib/types';

export const mockGlobalStats: GlobalStats = {
	totalCards: 891,
	cardsReviewedToday: 47,
	currentStreak: 12,
	longestStreak: 28,
	averageRetention: 85,
	reviewsThisWeek: [32, 45, 28, 51, 38, 47, 0], // Mon-Sun, today is Sat
	upcomingDue: {
		today: 88,
		tomorrow: 42,
		thisWeek: 156
	}
};

export function getGlobalStats(): GlobalStats {
	return mockGlobalStats;
}
