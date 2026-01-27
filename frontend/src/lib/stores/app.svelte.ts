// Global application state using Svelte 5 runes
// Note: Focus mode and command palette state are managed via shallow routing (page.state)
// See +layout.svelte for implementation

import { getNotebookPreferences, setNotebookPreferences } from '$lib/services/notebooks';

class AppState {
	// UI state
	sidebarCollapsed = $state(false);
	sidebarWidth = $state(288);
	sourceExpanded = $state(false);

	// Source context sidebar state (right sidebar)
	sourceContextSidebarWidth = $state(320);
	sourceContextSidebarCollapsed = $state(false);

	// Future: bidirectional communication between source viewer and context sidebar
	activeSourceSection = $state<string | null>(null);
	highlightedCardIds = $state<string[]>([]);

	// Methods
	toggleSidebar() {
		this.sidebarCollapsed = !this.sidebarCollapsed;
	}

	toggleSourceExpanded() {
		this.sourceExpanded = !this.sourceExpanded;
	}

	// Source context sidebar methods
	toggleSourceContextSidebar() {
		this.sourceContextSidebarCollapsed = !this.sourceContextSidebarCollapsed;
	}

	setSourceContextSidebarWidth(width: number) {
		this.sourceContextSidebarWidth = Math.min(500, Math.max(280, width));
	}

	// Future: Called when source viewer scrolls to a new section
	setActiveSourceSection(section: string | null) {
		this.activeSourceSection = section;
	}

	// Future: Called when cards in viewport change
	setHighlightedCards(cardIds: string[]) {
		this.highlightedCardIds = cardIds;
	}

	// Notebook layout persistence
	loadNotebookLayout(notebookId: string) {
		const prefs = getNotebookPreferences(notebookId);
		this.sidebarCollapsed = prefs?.sidebarCollapsed ?? false;
		this.sidebarWidth = prefs?.sidebarWidth ?? 288;
		this.sourceExpanded = false;
		this.sourceContextSidebarCollapsed = prefs?.contextSidebarCollapsed ?? false;
		this.sourceContextSidebarWidth = prefs?.contextSidebarWidth ?? 320;
	}

	saveNotebookLayout(notebookId: string) {
		setNotebookPreferences(notebookId, {
			sidebarCollapsed: this.sidebarCollapsed,
			sidebarWidth: this.sidebarWidth,
			contextSidebarCollapsed: this.sourceContextSidebarCollapsed,
			contextSidebarWidth: this.sourceContextSidebarWidth
		});
	}

	// Reset state when navigating away
	resetNotebookState() {
		this.sourceExpanded = false;
		this.sourceContextSidebarCollapsed = false;
		this.activeSourceSection = null;
		this.highlightedCardIds = [];
	}
}

export const appState = new AppState();
