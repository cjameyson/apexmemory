<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import TopNavBar from '$lib/components/layout/top-nav-bar.svelte';
	import CommandPalette from '$lib/components/overlays/command-palette.svelte';
	import FocusMode from '$lib/components/overlays/focus-mode.svelte';
	import { appState } from '$lib/stores/app.svelte';
	import type { Notebook } from '$lib/types';
	import { getAllNotebooks, getNotebook } from '$lib/mocks';

	let { data, children } = $props();

	// Get notebooks from mock data
	let notebooks: Notebook[] = $state(getAllNotebooks());

	// Derive current notebook from URL params
	let currentNotebook = $derived.by(() => {
		const match = $page.url.pathname.match(/\/notebooks\/([^/]+)/);
		if (match) {
			return getNotebook(match[1]);
		}
		return undefined;
	});

	// Global keyboard shortcuts
	onMount(() => {
		function handleKeydown(e: KeyboardEvent) {
			// Cmd+K or Ctrl+K to open command palette
			if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
				e.preventDefault();
				appState.toggleCommandPalette();
			}
			// Escape to close overlays (handled by individual overlays)
		}

		window.addEventListener('keydown', handleKeydown);
		return () => window.removeEventListener('keydown', handleKeydown);
	});

	function handleStartReview() {
		appState.closeCommandPalette();
		appState.startFocusMode({ type: 'all' });
	}
</script>

<div class="h-screen flex flex-col bg-slate-100 dark:bg-slate-950">
	<TopNavBar
		user={data.user}
		{notebooks}
		{currentNotebook}
	/>

	<main class="flex-1 flex overflow-hidden">
		{@render children()}
	</main>
</div>

<!-- Command Palette -->
{#if appState.commandPaletteOpen}
	<CommandPalette
		{currentNotebook}
		onClose={() => appState.closeCommandPalette()}
		onStartReview={handleStartReview}
	/>
{/if}

<!-- Focus Mode -->
{#if appState.focusMode.active && appState.focusMode.scope}
	<FocusMode
		scope={appState.focusMode.scope}
		onClose={() => appState.exitFocusMode()}
	/>
{/if}
