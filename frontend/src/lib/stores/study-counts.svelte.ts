/**
 * Study counts store for managing due/new card counts across the app.
 * Uses Svelte 5 runes for reactivity.
 *
 * Refresh triggers:
 * - Layout mount: initialize from SSR data + start polling
 * - Focus mode exits: explicit refresh
 * - Fact created/updated: explicit refresh
 * - Periodic interval: every 2 minutes (when tab visible)
 * - Tab becomes visible: if >30s since last refresh
 */

import type { ApiStudyCountsResponse } from '$lib/api/types';

type NotebookCounts = { due: number; new: number; total: number };

type StudyCountsData = {
	byNotebook: Record<string, NotebookCounts>;
	totalDue: number;
	totalNew: number;
};

// State
let counts = $state<StudyCountsData | null>(null);
let polling = $state(false);
let pollIntervalId: ReturnType<typeof setInterval> | null = null;
let lastRefresh = $state(0);

// Constants
const DEFAULT_POLL_INTERVAL = 2 * 60 * 1000; // 2 minutes
const VISIBILITY_STALE_THRESHOLD = 30 * 1000; // 30 seconds

// Bound visibility handler (created once to enable removeEventListener)
function handleVisibilityChange(): void {
	if (document.visibilityState === 'visible') {
		const elapsed = Date.now() - lastRefresh;
		if (elapsed > VISIBILITY_STALE_THRESHOLD) {
			studyCounts.refresh();
		}
	}
}

export const studyCounts = {
	// Getters
	get value() {
		return counts;
	},
	get isPolling() {
		return polling;
	},

	/**
	 * Initialize store from SSR data. Called once on layout mount.
	 */
	initialize(data: ApiStudyCountsResponse) {
		counts = {
			byNotebook: data.counts,
			totalDue: data.total_due,
			totalNew: data.total_new
		};
		lastRefresh = Date.now();
	},

	/**
	 * Fetch fresh counts from the API.
	 */
	async refresh(): Promise<void> {
		try {
			const response = await fetch('/api/reviews/study-counts');
			if (response.ok) {
				const data: ApiStudyCountsResponse = await response.json();
				counts = {
					byNotebook: data.counts,
					totalDue: data.total_due,
					totalNew: data.total_new
				};
				lastRefresh = Date.now();
			}
		} catch (error) {
			console.error('Failed to refresh study counts:', error);
		}
	},

	/**
	 * Start periodic polling. Only refreshes when document is visible.
	 */
	startPolling(intervalMs = DEFAULT_POLL_INTERVAL): void {
		if (polling) return;
		polling = true;

		pollIntervalId = setInterval(() => {
			if (document.visibilityState === 'visible') {
				this.refresh();
			}
		}, intervalMs);

		// Also refresh when tab becomes visible if stale
		document.addEventListener('visibilitychange', handleVisibilityChange);
	},

	/**
	 * Stop periodic polling.
	 */
	stopPolling(): void {
		if (pollIntervalId) {
			clearInterval(pollIntervalId);
			pollIntervalId = null;
		}
		polling = false;
		document.removeEventListener('visibilitychange', handleVisibilityChange);
	},

	// Helpers for components
	getDueCount(notebookId: string): number {
		const notebook = counts?.byNotebook[notebookId];
		if (!notebook) return 0;
		return notebook.due + notebook.new;
	},

	getTotalDue(): number {
		return (counts?.totalDue ?? 0) + (counts?.totalNew ?? 0);
	},

	getTotalNew(): number {
		return counts?.totalNew ?? 0;
	},

	getNotebookCounts(notebookId: string): NotebookCounts | null {
		return counts?.byNotebook[notebookId] ?? null;
	}
};
