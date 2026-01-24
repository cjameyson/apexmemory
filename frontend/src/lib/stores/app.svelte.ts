// Global application state using Svelte 5 runes
// Note: Focus mode and command palette state are managed via shallow routing (page.state)
// See +layout.svelte for implementation

class AppState {
	// UI state
	sidebarCollapsed = $state(false);
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

	// Reset state when navigating away
	resetNotebookState() {
		this.sourceExpanded = false;
		this.sourceContextSidebarCollapsed = false;
		this.activeSourceSection = null;
		this.highlightedCardIds = [];
	}
}

export const appState = new AppState();
