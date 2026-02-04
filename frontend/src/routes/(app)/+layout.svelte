<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto, invalidateAll, pushState, replaceState } from '$app/navigation';
	import TopNavBar from '$lib/components/layout/top-nav-bar.svelte';
	import CommandPalette from '$lib/components/overlays/command-palette.svelte';
	import FocusMode from '$lib/components/overlays/focus-mode.svelte';
	import CreateNotebookModal from '$lib/components/notebooks/create-notebook-modal.svelte';
	import type { Notebook, ReviewScope, StudyCard } from '$lib/types';
	import type { ReviewMode } from '$lib/types/review';
	import { toNotebooksWithCounts } from '$lib/services/notebooks';
	import { fetchStudyCards, fetchPracticeCards } from '$lib/services/reviews';
	import { studyCounts } from '$lib/stores/study-counts.svelte';

	let { data, children } = $props();

	// Initialize study counts store from SSR data and start polling
	onMount(() => {
		studyCounts.initialize(data.studyCounts);
		studyCounts.startPolling();

		return () => {
			studyCounts.stopPolling();
		};
	});

	// Transform API data to frontend types with study counts
	let notebooks = $derived(toNotebooksWithCounts(data.notebooks, data.studyCounts));

	// Create lookup map for notebook by ID
	let notebookMap = $derived(new Map(notebooks.map((n) => [n.id, n])));

	function getNotebook(id: string) {
		return notebookMap.get(id);
	}

	let currentNotebook = $derived.by(() => {
		const match = page.url.pathname.match(/\/notebooks\/([^/]+)/);
		if (match) {
			return getNotebook(match[1]);
		}
		return undefined;
	});

	// Focus mode cards stored in component state (too large for URL state)
	// Can also be passed via page state for direct pushState calls
	let focusModeCards = $state<StudyCard[]>([]);
	let focusModeMode = $state<ReviewMode>('scheduled');

	// Sync focusModeCards from page state when it contains cards
	$effect(() => {
		const fm = page.state.focusMode;
		if (fm?.cards && fm.cards.length > 0) {
			focusModeCards = fm.cards;
			focusModeMode = fm.mode ?? 'scheduled';
		}
	});

	// Reconstruct ReviewScope from page state for focus mode
	let focusModeScope = $derived.by((): ReviewScope | null => {
		const fm = page.state.focusMode;
		if (!fm) return null;

		const mode = fm.mode ?? 'scheduled';

		if (fm.type === 'all') {
			return { type: 'all', mode };
		} else if (fm.type === 'notebook' && fm.notebookId) {
			const notebook = getNotebook(fm.notebookId);
			if (notebook) {
				return { type: 'notebook', notebook, mode };
			}
		} else if (fm.type === 'source' && fm.notebookId && fm.sourceId) {
			const notebook = getNotebook(fm.notebookId);
			if (notebook) {
				return {
					type: 'source',
					notebook,
					mode,
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
	function startFocusMode(scope: ReviewScope, cards: StudyCard[]) {
		focusModeCards = cards;
		focusModeMode = scope.mode;

		const focusState: App.PageState['focusMode'] = {
			type: scope.type,
			mode: scope.mode,
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

	function handleProgressChange(index: number) {
		if (page.state.focusMode) {
			replaceState('', {
				...page.state,
				focusMode: { ...page.state.focusMode, currentIndex: index }
			});
		}
	}

	async function exitFocusMode() {
		// Refresh study counts after session ends
		await studyCounts.refresh();
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
			if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
				e.preventDefault();
				if (commandPaletteOpen) {
					closeCommandPalette();
				} else {
					openCommandPalette();
				}
			}
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
		// Command palette triggers global scheduled review -- need to fetch cards
		fetchAndStartReview('scheduled');
	}

	async function fetchAndStartReview(mode: ReviewMode, notebookId?: string) {
		const cards = mode === 'scheduled'
			? await fetchStudyCards(notebookId)
			: await fetchPracticeCards(notebookId);

		const scope: ReviewScope = { type: 'all', mode };
		startFocusMode(scope, cards);
	}

	// Create notebook modal state
	let createNotebookOpen = $state(false);

	function handleCreateNotebook() {
		createNotebookOpen = true;
	}

	async function handleNotebookCreated(notebook: Notebook) {
		// Refresh notebooks list
		await invalidateAll();
		// Navigate to the new notebook
		goto(`/notebooks/${notebook.id}`);
	}
</script>

<div class="flex h-screen flex-col bg-background">
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
{#if focusModeScope && focusModeCards.length > 0}
	<FocusMode
		cards={focusModeCards}
		mode={focusModeMode}
		scope={focusModeScope}
		initialIndex={page.state.focusMode?.currentIndex ?? 0}
		onProgressChange={handleProgressChange}
		onClose={exitFocusMode}
	/>
{:else if focusModeScope && focusModeCards.length === 0}
	<FocusMode
		cards={[]}
		mode={focusModeMode}
		scope={focusModeScope}
		onClose={exitFocusMode}
	/>
{/if}

<!-- Create Notebook Modal -->
<CreateNotebookModal
	bind:open={createNotebookOpen}
	onOpenChange={(open) => (createNotebookOpen = open)}
	onSuccess={handleNotebookCreated}
/>
