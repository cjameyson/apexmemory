<script lang="ts">
	import type { PageData } from './$types';
	import { appState } from '$lib/stores/app.svelte';
	import NotebookSidebar from '$lib/components/layout/notebook-sidebar.svelte';
	import CardsGridView from '$lib/components/notebook/cards-grid-view.svelte';
	import SourceDetail from '$lib/components/notebook/source-detail.svelte';
	import type { Source, Card } from '$lib/types';

	let { data }: { data: PageData } = $props();

	// Local state for source selection (page-specific)
	let selectedSource = $state<Source | null>(null);

	// Filter cards based on selected source
	let displayedCards = $derived.by(() => {
		const source = selectedSource;
		if (source) {
			return data.cards.filter((c) => c.sourceId === source.id);
		}
		return data.cards;
	});

	function handleSelectSource(source: Source | null) {
		selectedSource = source;
		// Reset expanded state when changing source
		if (!source) {
			appState.sourceExpanded = false;
		}
	}

	function handleViewModeChange(mode: 'all' | 'due' | 'mastered') {
		appState.setCardsViewMode(mode);
	}

	function handleCardClick(_card: Card) {
		// Future: open card editor/viewer
	}

	function handleSourceClose() {
		selectedSource = null;
		appState.sourceExpanded = false;
	}

	function handleToggleExpand() {
		appState.toggleSourceExpanded();
		// Collapse sidebar when expanding source
		if (appState.sourceExpanded) {
			appState.sidebarCollapsed = true;
		}
	}

	function handleStartReview() {
		if (selectedSource) {
			appState.startFocusMode({
				type: 'source',
				notebook: data.notebook,
				source: selectedSource
			});
		}
	}
</script>

<div class="flex-1 flex overflow-hidden">
	<!-- Sidebar (hidden when source is expanded) -->
	{#if !appState.sourceExpanded}
		<NotebookSidebar
			notebook={data.notebook}
			sources={data.sources}
			cards={data.cards}
			{selectedSource}
			bind:isCollapsed={appState.sidebarCollapsed}
			onSelectSource={handleSelectSource}
		/>
	{/if}

	<!-- Main content -->
	<div class="flex-1 flex flex-col overflow-auto p-6">
		{#if selectedSource}
			<!-- Source detail view -->
			<SourceDetail
				source={selectedSource}
				cards={data.cards}
				isExpanded={appState.sourceExpanded}
				onClose={handleSourceClose}
				onStartReview={handleStartReview}
				onToggleExpand={handleToggleExpand}
			/>
		{:else}
			<!-- Cards grid view -->
			<div class="mb-6">
				<div class="flex items-center gap-3 mb-2">
					<span class="text-2xl">{data.notebook.emoji}</span>
					<h1 class="text-xl font-bold text-slate-900 dark:text-white">
						{data.notebook.name}
					</h1>
				</div>
				<p class="text-sm text-slate-500 dark:text-slate-400">
					{data.cards.length} cards across {data.sources.length} sources
				</p>
			</div>

			<CardsGridView
				cards={displayedCards}
				sources={data.sources}
				viewMode={appState.cardsViewMode}
				onViewModeChange={handleViewModeChange}
				onCardClick={handleCardClick}
			/>
		{/if}
	</div>
</div>
