<script lang="ts">
	import type { PageData } from './$types';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { appState } from '$lib/stores/app.svelte';
	import NotebookSidebar from '$lib/components/layout/notebook-sidebar.svelte';
	import CardBrowser from '$lib/components/notebook/card-browser.svelte';
	import type { Source } from '$lib/types';

	let { data }: { data: PageData } = $props();

	// Derive filters from URL query parameters
	let searchQuery = $derived(page.url.searchParams.get('search') ?? '');
	let selectedSourceId = $derived(page.url.searchParams.get('filter') ?? null);

	function handleSearchChange(query: string) {
		const url = new URL(page.url);
		if (query) {
			url.searchParams.set('search', query);
		} else {
			url.searchParams.delete('search');
		}
		goto(url, { replaceState: true, keepFocus: true, noScroll: true });
	}

	function handleSourceFilterChange(sourceId: string | null) {
		const url = new URL(page.url);
		if (sourceId) {
			url.searchParams.set('filter', sourceId);
		} else {
			url.searchParams.delete('filter');
		}
		goto(url, { replaceState: true, keepFocus: true, noScroll: true });
	}

	function handleSelectSource(source: Source | null) {
		if (source) {
			goto(`/notebooks/${data.notebook.id}?source=${source.id}`);
		}
	}

	function handleOpenSettings() {
		console.log('Open notebook settings');
	}
</script>

<div class="flex-1 flex overflow-hidden">
	<NotebookSidebar
		notebook={data.notebook}
		sources={data.sources}
		cards={data.cards}
		bind:isCollapsed={appState.sidebarCollapsed}
		onSelectSource={handleSelectSource}
		onOpenSettings={handleOpenSettings}
	/>

	<div class="flex-1 flex flex-col overflow-auto p-6">
		<div class="mb-6">
			<h1 class="text-xl font-bold text-slate-900 dark:text-white mb-1">
				Card Browser
			</h1>
			<p class="text-sm text-slate-500 dark:text-slate-400">
				Browse and search all cards in this notebook
			</p>
		</div>

		<CardBrowser
			cards={data.cards}
			sources={data.sources}
			{searchQuery}
			{selectedSourceId}
			onSearchChange={handleSearchChange}
			onSourceFilterChange={handleSourceFilterChange}
		/>
	</div>
</div>
