import type { ApiNotebook } from '$lib/api/types';
import type { Notebook } from '$lib/types/notebook';

// ============================================================================
// API Response Adapter
// ============================================================================

/**
 * Transform API response to frontend Notebook type.
 * Handles snake_case -> camelCase conversion and adds UI-only fields.
 */
export function toNotebook(api: ApiNotebook): Notebook {
	return {
		// API fields (snake_case -> camelCase)
		id: api.id,
		name: api.name,
		description: api.description,
		desiredRetention: api.desired_retention,
		position: api.position,
		createdAt: api.created_at,
		updatedAt: api.updated_at,

		// UI fields (mock until stats API)
		emoji: '\u{1F4D8}',
		color: 'blue',
		dueCount: 0,
		streak: 0,
		totalCards: 0,
		retention: api.desired_retention
	};
}

/**
 * Transform array of API responses to frontend Notebook types.
 */
export function toNotebooks(apiList: ApiNotebook[]): Notebook[] {
	return apiList.map(toNotebook);
}

// ============================================================================
// localStorage Helpers for UI Preferences
// ============================================================================

const STORAGE_KEY = 'notebook-preferences';

export interface NotebookPreference {
	sidebarCollapsed?: boolean;
	lastViewedSourceId?: string;
}

type NotebookPreferences = Record<string, NotebookPreference>;

/**
 * Get UI preferences for a specific notebook.
 * Returns null if no preferences exist or localStorage is unavailable.
 */
export function getNotebookPreferences(notebookId: string): NotebookPreference | null {
	if (typeof window === 'undefined') return null;

	try {
		const stored = localStorage.getItem(STORAGE_KEY);
		if (!stored) return null;

		const prefs: NotebookPreferences = JSON.parse(stored);
		return prefs[notebookId] ?? null;
	} catch {
		return null;
	}
}

/**
 * Set UI preferences for a specific notebook.
 * Merges with existing preferences.
 */
export function setNotebookPreferences(
	notebookId: string,
	prefs: Partial<NotebookPreference>
): void {
	if (typeof window === 'undefined') return;

	try {
		const stored = localStorage.getItem(STORAGE_KEY);
		const allPrefs: NotebookPreferences = stored ? JSON.parse(stored) : {};

		allPrefs[notebookId] = {
			...allPrefs[notebookId],
			...prefs
		};

		localStorage.setItem(STORAGE_KEY, JSON.stringify(allPrefs));
	} catch {
		// Silently fail - localStorage may be full or unavailable
	}
}

/**
 * Clear UI preferences for a specific notebook.
 * Call this when a notebook is deleted.
 */
export function clearNotebookPreferences(notebookId: string): void {
	if (typeof window === 'undefined') return;

	try {
		const stored = localStorage.getItem(STORAGE_KEY);
		if (!stored) return;

		const allPrefs: NotebookPreferences = JSON.parse(stored);
		delete allPrefs[notebookId];

		localStorage.setItem(STORAGE_KEY, JSON.stringify(allPrefs));
	} catch {
		// Silently fail
	}
}
