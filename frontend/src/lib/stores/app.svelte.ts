// Global application state using Svelte 5 runes
// Note: Focus mode state is managed via shallow routing (page.state)
// See +layout.svelte for focus mode implementation

class AppState {
	// UI state
	sidebarCollapsed = $state(false);
	sourceExpanded = $state(false);
	cardsViewMode = $state<'all' | 'due' | 'mastered'>('all');

	// Source context sidebar state (right sidebar)
	sourceContextSidebarWidth = $state(320);
	sourceContextSidebarCollapsed = $state(false);

	// Future: bidirectional communication between source viewer and context sidebar
	activeSourceSection = $state<string | null>(null);
	highlightedCardIds = $state<string[]>([]);

	// Overlay state (command palette only - focus mode uses shallow routing)
	commandPaletteOpen = $state(false);

	// Methods
	toggleSidebar() {
		this.sidebarCollapsed = !this.sidebarCollapsed;
	}

	toggleSourceExpanded() {
		this.sourceExpanded = !this.sourceExpanded;
	}

	setCardsViewMode(mode: 'all' | 'due' | 'mastered') {
		this.cardsViewMode = mode;
	}

	openCommandPalette() {
		this.commandPaletteOpen = true;
	}

	closeCommandPalette() {
		this.commandPaletteOpen = false;
	}

	toggleCommandPalette() {
		this.commandPaletteOpen = !this.commandPaletteOpen;
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
		this.cardsViewMode = 'all';
		this.sourceContextSidebarCollapsed = false;
		this.activeSourceSection = null;
		this.highlightedCardIds = [];
	}
}

export const appState = new AppState();
