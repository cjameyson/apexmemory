<script lang="ts">
	import type { LayoutData } from './$types';
	import type { Snippet } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { appState } from '$lib/stores/app.svelte';
	import NotebookSidebar from '$lib/components/layout/notebook-sidebar.svelte';
	import type { Source } from '$lib/types';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();

	// Track which notebook layout is loaded to avoid save-on-load
	let layoutLoadedForId = $state<string | null>(null);

	// Load notebook layout on notebook change
	$effect(() => {
		const id = data.notebook.id;
		layoutLoadedForId = null;
		appState.loadNotebookLayout(id);
		layoutLoadedForId = id;
	});

	// Auto-save layout when any layout value changes
	$effect(() => {
		const collapsed = appState.sidebarCollapsed;
		const width = appState.sidebarWidth;
		const ctxCollapsed = appState.sourceContextSidebarCollapsed;
		const ctxWidth = appState.sourceContextSidebarWidth;
		if (layoutLoadedForId && layoutLoadedForId === data.notebook.id) {
			appState.saveNotebookLayout(data.notebook.id);
		}
	});

	// Derive selected source from URL pathname
	let selectedSource = $derived.by(() => {
		const match = page.url.pathname.match(/\/sources\/([^/]+)/);
		if (match) {
			return data.sources.find((s) => s.id === match[1]) ?? null;
		}
		return null;
	});

	function handleSelectSource(source: Source | null) {
		if (source) {
			goto(`/notebooks/${data.notebook.id}/sources/${source.id}`);
		} else {
			appState.sourceExpanded = false;
			goto(`/notebooks/${data.notebook.id}`);
		}
	}

	function handleOpenSettings() {
		// TODO: Open notebook settings modal
		console.log('Open notebook settings');
	}
</script>

<div id="notebook-page" class="flex flex-1 overflow-hidden">
	{#if !appState.sourceExpanded}
		<NotebookSidebar
			notebook={data.notebook}
			sources={data.sources}
			cards={data.cards}
			{selectedSource}
			bind:isCollapsed={appState.sidebarCollapsed}
			bind:sidebarWidth={appState.sidebarWidth}
			onSelectSource={handleSelectSource}
			onOpenSettings={handleOpenSettings}
		/>
	{/if}

	{@render children()}
</div>
