// Global application state using Svelte 5 runes

import type { ReviewScope } from '$lib/types';

class AppState {
	// UI state
	sidebarCollapsed = $state(false);
	sourceExpanded = $state(false);
	cardsViewMode = $state<'all' | 'due' | 'mastered'>('all');

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

	// Reset state when navigating away
	resetNotebookState() {
		this.sourceExpanded = false;
		this.cardsViewMode = 'all';
	}
}

export const appState = new AppState();
