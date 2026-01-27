<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { pushState, replaceState } from '$app/navigation';
	import TopNavBar from '$lib/components/layout/top-nav-bar.svelte';
	import CommandPalette from '$lib/components/overlays/command-palette.svelte';
	import FocusMode from '$lib/components/overlays/focus-mode.svelte';
	import type { Notebook, ReviewScope } from '$lib/types';
	import { toNotebooks } from '$lib/services/notebooks';

	let { data, children } = $props();

	// Transform API data to frontend types
	let notebooks = $derived(toNotebooks(data.notebooks));

	// Create lookup map for notebook by ID
	let notebookMap = $derived(new Map(notebooks.map((n) => [n.id, n])));

	// Helper to get notebook by ID (replaces mock getNotebook)
	function getNotebook(id: string) {
		return notebookMap.get(id);
	}

	// Derive current notebook from URL params
	let currentNotebook = $derived.by(() => {
		const match = page.url.pathname.match(/\/notebooks\/([^/]+)/);
		if (match) {
			return getNotebook(match[1]);
		}
		return undefined;
	});

	// Reconstruct ReviewScope from page state for focus mode
	let focusModeScope = $derived.by((): ReviewScope | null => {
		const fm = page.state.focusMode;
		if (!fm) return null;

		if (fm.type === 'all') {
			return { type: 'all' };
		} else if (fm.type === 'notebook' && fm.notebookId) {
			const notebook = getNotebook(fm.notebookId);
			if (notebook) {
				return { type: 'notebook', notebook };
			}
		} else if (fm.type === 'source' && fm.notebookId && fm.sourceId) {
			const notebook = getNotebook(fm.notebookId);
			if (notebook) {
				return {
					type: 'source',
					notebook,
					source: {
						id: fm.sourceId,
						notebookId: fm.notebookId,
						name: fm.sourceName || 'Source',
						type: 'pdf',
						cards: 0,
						excerpt: '',
						addedAt: ''
					}
				};
			}
		}
		return null;
	});

	// Start focus mode using shallow routing (pushState)
	function startFocusMode(scope: ReviewScope) {
		const focusState: App.PageState['focusMode'] = {
			type: scope.type,
			currentIndex: 0
		};

		if (scope.type === 'notebook') {
			focusState.notebookId = scope.notebook.id;
			focusState.notebookName = scope.notebook.name;
			focusState.notebookEmoji = scope.notebook.emoji;
		} else if (scope.type === 'source') {
			focusState.notebookId = scope.notebook.id;
			focusState.notebookName = scope.notebook.name;
			focusState.notebookEmoji = scope.notebook.emoji;
			focusState.sourceId = scope.source.id;
			focusState.sourceName = scope.source.name;
		}

		pushState('', { focusMode: focusState });
	}

	// Update progress without creating new history entry
	function handleProgressChange(index: number) {
		if (page.state.focusMode) {
			replaceState('', {
				...page.state,
				focusMode: { ...page.state.focusMode, currentIndex: index }
			});
		}
	}

	// Exit focus mode by going back in history
	function exitFocusMode() {
		history.back();
	}

	// Command palette state derived from shallow routing
	let commandPaletteOpen = $derived(!!page.state.commandPalette);

	function openCommandPalette() {
		if (!commandPaletteOpen) {
			pushState('', { ...page.state, commandPalette: true });
		}
	}

	function closeCommandPalette() {
		if (commandPaletteOpen) {
			history.back();
		}
	}

	// Global keyboard shortcuts
	onMount(() => {
		function handleKeydown(e: KeyboardEvent) {
			// Cmd+K or Ctrl+K to open command palette
			if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
				e.preventDefault();
				if (commandPaletteOpen) {
					closeCommandPalette();
				} else {
					openCommandPalette();
				}
			}
			// Escape to close overlays
			if (e.key === 'Escape') {
				if (commandPaletteOpen) {
					e.preventDefault();
					closeCommandPalette();
				} else if (page.state.focusMode) {
					e.preventDefault();
					exitFocusMode();
				}
			}
		}

		window.addEventListener('keydown', handleKeydown);
		return () => window.removeEventListener('keydown', handleKeydown);
	});

	function handleStartReview() {
		closeCommandPalette();
		startFocusMode({ type: 'all' });
	}

	function handleCreateNotebook() {
		// TODO: Wire up create notebook modal
		console.log('TODO: Create notebook');
	}
</script>

<div class="flex h-screen flex-col bg-white dark:bg-slate-950">
	<TopNavBar
		user={data.user}
		{notebooks}
		{currentNotebook}
		onStartFocusMode={startFocusMode}
		onOpenSearch={openCommandPalette}
		onCreateNotebook={handleCreateNotebook}
	/>

	<main class="flex flex-1 overflow-hidden">
		{@render children()}
	</main>
</div>

<!-- Command Palette (shallow routing - Back button closes) -->
{#if commandPaletteOpen}
	<CommandPalette
		{currentNotebook}
		onClose={closeCommandPalette}
		onStartReview={handleStartReview}
	/>
{/if}

<!-- Focus Mode (rendered based on shallow routing state) -->
{#if focusModeScope}
	<FocusMode
		scope={focusModeScope}
		initialIndex={page.state.focusMode?.currentIndex ?? 0}
		onProgressChange={handleProgressChange}
		onClose={exitFocusMode}
	/>
{/if}
