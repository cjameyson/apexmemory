<script lang="ts">
	import { page } from '$app/state';
	import { goto, pushState } from '$app/navigation';
	import { error } from '@sveltejs/kit';
	import { appState } from '$lib/stores/app.svelte';
	import SourceDetail from '$lib/components/notebook/source-detail.svelte';
	import SourceContextSidebar from '$lib/components/notebook/source-context-sidebar.svelte';
	import { type SidebarTab } from '$lib/components/notebook/source-sidebar-tabs.svelte';
	import type { DisplayCard } from '$lib/types';

	let { data } = $props();

	let source = $derived.by(() => {
		const s = data.sources.find((s) => s.id === page.params.sourceId);
		if (!s) return null;
		return s;
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

	function handleSourceClose() {
		appState.sourceExpanded = false;
		goto(`/notebooks/${page.params.id}`, { keepFocus: true });
	}

	function handleToggleExpand() {
		appState.toggleSourceExpanded();
	}

	function handleStartReview() {
		if (source) {
			pushState('', {
				focusMode: {
					type: 'source',
					mode: 'scheduled',
					notebookId: data.notebook.id,
					notebookName: data.notebook.name,
					notebookEmoji: data.notebook.emoji,
					sourceId: source.id,
					sourceName: source.name,
					currentIndex: 0
				}
			});
		}
	}

	function handleCardClick(card: DisplayCard) {
		// Future: Scroll source viewer to card's linked section
	}

	function handleTabChange(tab: SidebarTab) {
		const url = new URL(page.url);
		if (tab === 'cards') {
			url.searchParams.delete('tab');
		} else {
			url.searchParams.set('tab', tab);
		}
		goto(url, { replaceState: true, keepFocus: true, noScroll: true });
	}
</script>

{#if source}
	<!-- Main content -->
	<div id="notebook-main" class="flex flex-1 flex-col overflow-auto">
		<SourceDetail
			{source}
			cards={data.cards}
			isExpanded={appState.sourceExpanded}
			onClose={handleSourceClose}
			onStartReview={handleStartReview}
			onToggleExpand={handleToggleExpand}
		/>
	</div>

	<!-- Right context sidebar (hidden when expanded) -->
	{#if !appState.sourceExpanded}
		<SourceContextSidebar
			{source}
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
{:else}
	<div class="flex flex-1 items-center justify-center">
		<p class="text-muted-foreground">Source not found</p>
	</div>
{/if}
