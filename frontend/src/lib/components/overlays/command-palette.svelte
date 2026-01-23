<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { cn } from '$lib/utils';
	import {
		SearchIcon,
		PlusIcon,
		UploadIcon,
		ZapIcon,
		BookOpenIcon,
		LayersIcon,
		FileTextIcon,
		ArrowRightIcon
	} from '@lucide/svelte';
	import type { Notebook, Source, Card } from '$lib/types';
	import { getAllNotebooks, getAllSources, getAllDueCards } from '$lib/mocks';

	interface Props {
		currentNotebook?: Notebook;
		onClose?: () => void;
		onStartReview?: () => void;
		class?: string;
	}

	let { currentNotebook, onClose, onStartReview, class: className }: Props = $props();

	let searchQuery = $state('');
	let scope = $state<'notebook' | 'all'>('all');
	let selectedIndex = $state(0);
	let inputRef = $state<HTMLInputElement | null>(null);

	// Initialize scope based on currentNotebook
	$effect(() => {
		if (currentNotebook) {
			scope = 'notebook';
		}
	});

	// Get data for search
	let notebooks = $derived(getAllNotebooks());
	let sources = $derived(getAllSources());
	let cards = $derived(getAllDueCards());

	// Quick actions when no query
	const quickActions = [
		{ id: 'create-card', icon: PlusIcon, label: 'Create new card', shortcut: '⌘ C' },
		{ id: 'add-source', icon: UploadIcon, label: 'Add source', shortcut: '⌘ S' },
		{ id: 'start-review', icon: ZapIcon, label: 'Start review', shortcut: '⌘ R' },
		{ id: 'switch-notebook', icon: BookOpenIcon, label: 'Switch notebook', shortcut: '⌘ N' }
	];

	// Filter results based on query
	let searchResults = $derived.by(() => {
		if (!searchQuery.trim()) return [];

		const query = searchQuery.toLowerCase();
		const results: { type: 'notebook' | 'source' | 'card'; item: Notebook | Source | Card }[] = [];

		// Filter by scope
		const scopedNotebooks = scope === 'notebook' && currentNotebook
			? [currentNotebook]
			: notebooks;

		const scopedSources = scope === 'notebook' && currentNotebook
			? sources.filter(s => s.notebookId === currentNotebook.id)
			: sources;

		// Search notebooks (only in 'all' scope)
		if (scope === 'all') {
			for (const nb of scopedNotebooks) {
				if (nb.name.toLowerCase().includes(query)) {
					results.push({ type: 'notebook', item: nb });
				}
			}
		}

		// Search sources
		for (const src of scopedSources) {
			if (src.name.toLowerCase().includes(query)) {
				results.push({ type: 'source', item: src });
			}
		}

		// Search cards (limited to top 5)
		const scopedCards = scope === 'notebook' && currentNotebook
			? cards.filter(c => c.notebookId === currentNotebook.id)
			: cards;

		for (const card of scopedCards.slice(0, 5)) {
			if (card.front.toLowerCase().includes(query) || card.back.toLowerCase().includes(query)) {
				results.push({ type: 'card', item: card });
			}
		}

		return results.slice(0, 10);
	});

	let displayItems = $derived(searchQuery.trim() ? searchResults : []);

	function handleAction(actionId: string) {
		switch (actionId) {
			case 'create-card':
				console.log('Create card');
				break;
			case 'add-source':
				console.log('Add source');
				break;
			case 'start-review':
				onStartReview?.();
				onClose?.();
				break;
			case 'switch-notebook':
				// Clear query to show notebooks
				searchQuery = '';
				break;
		}
	}

	function handleResultClick(result: typeof searchResults[0]) {
		if (result.type === 'notebook') {
			const nb = result.item as Notebook;
			goto(`/notebooks/${nb.id}`);
			onClose?.();
		} else if (result.type === 'source') {
			const src = result.item as Source;
			goto(`/notebooks/${src.notebookId}`);
			onClose?.();
		} else if (result.type === 'card') {
			// Open card editor (future)
			console.log('Open card:', result.item);
		}
	}

	function toggleScope() {
		scope = scope === 'notebook' ? 'all' : 'notebook';
	}

	// Focus input on mount
	onMount(() => {
		inputRef?.focus();
	});

	// Keyboard navigation
	function handleKeydown(e: KeyboardEvent) {
		const items = searchQuery.trim() ? displayItems : quickActions;
		const maxIndex = items.length - 1;

		if (e.key === 'ArrowDown') {
			e.preventDefault();
			selectedIndex = Math.min(selectedIndex + 1, maxIndex);
		} else if (e.key === 'ArrowUp') {
			e.preventDefault();
			selectedIndex = Math.max(selectedIndex - 1, 0);
		} else if (e.key === 'Enter') {
			e.preventDefault();
			if (searchQuery.trim() && displayItems[selectedIndex]) {
				handleResultClick(displayItems[selectedIndex]);
			} else if (!searchQuery.trim() && quickActions[selectedIndex]) {
				handleAction(quickActions[selectedIndex].id);
			}
		}
	}
</script>

<div
	class={cn(
		'fixed inset-0 z-50 flex items-start justify-center pt-[15vh]',
		className
	)}
	role="dialog"
	aria-modal="true"
	aria-label="Command palette"
>
	<!-- Backdrop -->
	<button
		type="button"
		class="absolute inset-0 bg-black/50"
		onclick={onClose}
		aria-label="Close"
	></button>

	<!-- Modal -->
	<div class="relative bg-white dark:bg-slate-900 rounded-xl shadow-2xl w-full max-w-lg mx-4 overflow-hidden">
		<!-- Search input -->
		<div class="flex items-center gap-3 px-4 py-3 border-b border-slate-200 dark:border-slate-800">
			<SearchIcon class="size-5 text-slate-400" />
			<input
				bind:this={inputRef}
				bind:value={searchQuery}
				onkeydown={handleKeydown}
				type="text"
				placeholder={currentNotebook ? `Search in ${currentNotebook.name}...` : 'Search everywhere...'}
				class="flex-1 bg-transparent border-0 focus:ring-0 text-slate-900 dark:text-white placeholder:text-slate-400 text-sm"
			/>

			<!-- Scope toggle -->
			{#if currentNotebook}
				<button
					type="button"
					onclick={toggleScope}
					class={cn(
						'flex items-center gap-1.5 px-2 py-1 rounded-md text-xs font-medium transition-colors',
						scope === 'notebook'
							? 'bg-sky-100 dark:bg-sky-900/30 text-sky-700 dark:text-sky-300'
							: 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400'
					)}
				>
					{#if scope === 'notebook'}
						<BookOpenIcon class="size-3.5" />
						This notebook
					{:else}
						<LayersIcon class="size-3.5" />
						All
					{/if}
				</button>
			{/if}
		</div>

		<!-- Results / Quick actions -->
		<div class="max-h-80 overflow-auto">
			{#if !searchQuery.trim()}
				<!-- Quick actions -->
				<div class="p-2">
					<div class="px-2 py-1.5 text-xs font-medium text-slate-500 dark:text-slate-400">
						Quick actions
					</div>
					{#each quickActions as action, i (action.id)}
						<button
							type="button"
							onclick={() => handleAction(action.id)}
							class={cn(
								'w-full flex items-center gap-3 px-2 py-2 rounded-lg text-left transition-colors',
								i === selectedIndex
									? 'bg-sky-100 dark:bg-sky-900/30 text-sky-900 dark:text-sky-100'
									: 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300'
							)}
						>
							<action.icon class="size-4 text-slate-500" />
							<span class="flex-1 text-sm">{action.label}</span>
							<kbd class="text-xs text-slate-400 bg-slate-100 dark:bg-slate-800 px-1.5 py-0.5 rounded">
								{action.shortcut}
							</kbd>
						</button>
					{/each}
				</div>
			{:else if displayItems.length > 0}
				<!-- Search results -->
				<div class="p-2">
					{#each displayItems as result, i (i)}
						<button
							type="button"
							onclick={() => handleResultClick(result)}
							class={cn(
								'w-full flex items-center gap-3 px-2 py-2 rounded-lg text-left transition-colors',
								i === selectedIndex
									? 'bg-sky-100 dark:bg-sky-900/30 text-sky-900 dark:text-sky-100'
									: 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300'
							)}
						>
							{#if result.type === 'notebook'}
								{@const nb = result.item as Notebook}
								<span class="text-lg">{nb.emoji}</span>
								<span class="flex-1 text-sm truncate">{nb.name}</span>
								<span class="text-xs text-slate-400">Notebook</span>
							{:else if result.type === 'source'}
								{@const src = result.item as Source}
								<FileTextIcon class="size-4 text-slate-500" />
								<span class="flex-1 text-sm truncate">{src.name}</span>
								<span class="text-xs text-slate-400">Source</span>
							{:else if result.type === 'card'}
								{@const card = result.item as Card}
								<LayersIcon class="size-4 text-slate-500" />
								<span class="flex-1 text-sm truncate">{card.front}</span>
								<span class="text-xs text-slate-400">Card</span>
							{/if}
						</button>
					{/each}
				</div>
			{:else}
				<!-- No results -->
				<div class="p-8 text-center text-slate-500 dark:text-slate-400">
					No results found for "{searchQuery}"
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<div class="flex items-center justify-between px-4 py-2 border-t border-slate-200 dark:border-slate-800 text-xs text-slate-400">
			<div class="flex items-center gap-3">
				<span><kbd class="px-1 bg-slate-100 dark:bg-slate-800 rounded">↑↓</kbd> Navigate</span>
				<span><kbd class="px-1 bg-slate-100 dark:bg-slate-800 rounded">↵</kbd> Select</span>
				<span><kbd class="px-1 bg-slate-100 dark:bg-slate-800 rounded">esc</kbd> Close</span>
			</div>
			<span><kbd class="px-1 bg-slate-100 dark:bg-slate-800 rounded">⌘K</kbd></span>
		</div>
	</div>
</div>
