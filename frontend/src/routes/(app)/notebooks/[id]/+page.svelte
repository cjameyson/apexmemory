<script lang="ts">
	import type { PageData } from './$types';
	import { page } from '$app/state';
	import { goto, pushState } from '$app/navigation';
	import { appState } from '$lib/stores/app.svelte';
	import NotebookSidebar from '$lib/components/layout/notebook-sidebar.svelte';
	import NotebookDashboard from '$lib/components/notebook/notebook-dashboard.svelte';
	import SourceDetail from '$lib/components/notebook/source-detail.svelte';
	import SourceContextSidebar from '$lib/components/notebook/source-context-sidebar.svelte';
	import { type SidebarTab } from '$lib/components/notebook/source-sidebar-tabs.svelte';
	import type { Source } from '$lib/types';

	let { data }: { data: PageData } = $props();

	// Derive selected source from URL query parameter (source of truth)
	let selectedSource = $derived.by(() => {
		const sourceId = page.url.searchParams.get('source');
		if (sourceId) {
			return data.sources.find((s) => s.id === sourceId) ?? null;
		}
		return null;
	});

	// Derive active sidebar tab from URL query parameter
	const validTabs: SidebarTab[] = ['cards', 'summary', 'chat'];
	let activeTab = $derived.by((): SidebarTab => {
		const tab = page.url.searchParams.get('tab');
		if (tab && validTabs.includes(tab as SidebarTab)) {
			return tab as SidebarTab;
		}
		return 'cards';
	});

	function handleSelectSource(source: Source | null) {
		const url = new URL(page.url);
		if (source) {
			url.searchParams.set('source', source.id);
		} else {
			url.searchParams.delete('source');
			// Reset expanded state when deselecting source
			appState.sourceExpanded = false;
		}
		goto(url, { replaceState: true, keepFocus: true, noScroll: true });
	}

	function handleOpenSettings() {
		// TODO: Open notebook settings modal
		console.log('Open notebook settings');
	}

	function handleSourceClose() {
		handleSelectSource(null);
	}

	function handleToggleExpand() {
		appState.toggleSourceExpanded();
	}

	function handleStartReview() {
		if (selectedSource) {
			// Start focus mode using shallow routing
			pushState('', {
				focusMode: {
					type: 'source',
					notebookId: data.notebook.id,
					notebookName: data.notebook.name,
					notebookEmoji: data.notebook.emoji,
					sourceId: selectedSource.id,
					sourceName: selectedSource.name,
					currentIndex: 0
				}
			});
		}
	}

	function handleCardClick(card: import('$lib/types').Card) {
		// Future: Scroll source viewer to card's linked section
		// appState.setActiveSourceSection(card.sourceSection);
	}

	function handleTabChange(tab: SidebarTab) {
		const url = new URL(page.url);
		if (tab === 'cards') {
			// Remove tab param when it's the default
			url.searchParams.delete('tab');
		} else {
			url.searchParams.set('tab', tab);
		}
		goto(url, { replaceState: true, keepFocus: true, noScroll: true });
	}
</script>

<div id="notebook-page" class="flex flex-1 overflow-hidden">
	<!-- Left sidebar (hidden when source is expanded) -->
	{#if !appState.sourceExpanded}
		<NotebookSidebar
			notebook={data.notebook}
			sources={data.sources}
			cards={data.cards}
			{selectedSource}
			bind:isCollapsed={appState.sidebarCollapsed}
			onSelectSource={handleSelectSource}
			onOpenSettings={handleOpenSettings}
		/>
	{/if}

	<!-- Main content -->
	<div id="notebook-main" class="flex flex-1 flex-col overflow-auto">
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
			<!-- Notebook dashboard -->
			<NotebookDashboard notebook={data.notebook} sources={data.sources} cards={data.cards} />
		{/if}
	</div>

	<!-- Right context sidebar (only when source is selected, hidden when expanded) -->
	{#if selectedSource && !appState.sourceExpanded}
		<SourceContextSidebar
			source={selectedSource}
			cards={data.cards}
			{activeTab}
			bind:isCollapsed={appState.sourceContextSidebarCollapsed}
			bind:sidebarWidth={appState.sourceContextSidebarWidth}
			highlightedCardIds={appState.highlightedCardIds}
			onTabChange={handleTabChange}
			onCardClick={handleCardClick}
			class="hidden md:flex"
		/>
	{/if}
</div>
