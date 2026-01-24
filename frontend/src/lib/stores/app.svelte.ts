// Global application state using Svelte 5 runes

import type { ReviewScope } from '$lib/types';

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

	// Overlay state
	commandPaletteOpen = $state(false);
	focusMode = $state<{
		active: boolean;
		scope: ReviewScope | null;
	}>({ active: false, scope: null });

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

	startFocusMode(scope: ReviewScope) {
		this.focusMode = { active: true, scope };
	}

	exitFocusMode() {
		this.focusMode = { active: false, scope: null };
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
