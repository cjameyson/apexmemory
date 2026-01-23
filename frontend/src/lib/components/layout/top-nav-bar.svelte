<script lang="ts">
	import { page } from '$app/stores';
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import NotebooksDropdown from '$lib/components/navigation/notebooks-dropdown.svelte';
	import SearchTrigger from '$lib/components/navigation/search-trigger.svelte';
	import ReviewLauncher from '$lib/components/navigation/review-launcher.svelte';
	import { HomeIcon, BrainIcon } from '@lucide/svelte';
	import type { Notebook, ReviewScope } from '$lib/types';
	import type { User } from '$lib/api/types';
	import { appState } from '$lib/stores/app.svelte';

	interface Props {
		user: User;
		notebooks: Notebook[];
		currentNotebook?: Notebook;
	}

	let { user, notebooks, currentNotebook }: Props = $props();

	let isHome = $derived($page.url.pathname === '/home');
	let isInNotebook = $derived($page.url.pathname.startsWith('/notebooks/'));

	function handleOpenSearch() {
		appState.openCommandPalette();
	}

	function handleStartReview(scope: ReviewScope) {
		appState.startFocusMode(scope);
	}
</script>

<header class="bg-white dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 px-4 py-2 shrink-0">
	<div class="flex items-center justify-between">
		<!-- Left group -->
		<div class="flex items-center gap-1">
			<!-- Logo -->
			<a
				href="/home"
				class="flex items-center gap-2 px-2 py-1.5 rounded-lg text-slate-900 dark:text-white hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
			>
				<BrainIcon class="size-6 text-sky-500" />
				<span class="font-semibold text-sm hidden sm:inline">Apex Memory</span>
			</a>

			<!-- Home button -->
			<Button
				href="/home"
				variant="ghost"
				size="sm"
				class={cn(
					'gap-2',
					isHome && 'bg-slate-100 dark:bg-slate-800'
				)}
			>
				<HomeIcon class="size-4" />
				<span class="hidden sm:inline">Home</span>
			</Button>

			<!-- Notebooks dropdown -->
			<NotebooksDropdown
				{notebooks}
				current={currentNotebook}
				{isInNotebook}
			/>
		</div>

		<!-- Right group -->
		<div class="flex items-center gap-2">
			<!-- Search trigger -->
			<SearchTrigger onclick={handleOpenSearch} />

			<!-- Review launcher -->
			<ReviewLauncher
				{notebooks}
				{currentNotebook}
				onStartReview={handleStartReview}
			/>
		</div>
	</div>
</header>
