<script lang="ts">
	import type { PageData } from './$types';
	import { appState } from '$lib/stores/app.svelte';
	import NotebookSidebar from '$lib/components/layout/notebook-sidebar.svelte';
	import NotebookDashboard from '$lib/components/notebook/notebook-dashboard.svelte';
	import SourceDetail from '$lib/components/notebook/source-detail.svelte';
	import type { Source } from '$lib/types';

	let { data }: { data: PageData } = $props();

	// Local state for source selection (page-specific)
	let selectedSource = $state<Source | null>(null);

	function handleSelectSource(source: Source | null) {
		selectedSource = source;
		// Reset expanded state when changing source
		if (!source) {
			appState.sourceExpanded = false;
		}
	}

	function handleOpenSettings() {
		// TODO: Open notebook settings modal
		console.log('Open notebook settings');
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

<div id="notebook-page" class="flex flex-1 overflow-hidden">
	<!-- Sidebar (hidden when source is expanded) -->
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
</div>
