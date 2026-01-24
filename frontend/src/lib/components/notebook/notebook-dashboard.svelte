<script lang="ts">
	import type { Notebook, Source, Card } from '$lib/types';
	import { Button } from '$lib/components/ui/button';
	import { PlayIcon, TrendingUpIcon, LayersIcon, FileStackIcon, CalendarIcon } from '@lucide/svelte';
	import { appState } from '$lib/stores/app.svelte';

	interface Props {
		notebook: Notebook;
		sources: Source[];
		cards: Card[];
	}

	let { notebook, sources, cards }: Props = $props();

	let dueCount = $derived(cards.filter((c) => c.due).length);

	// Mock data for activity chart
	let reviewsThisWeek = [12, 8, 15, 20, 5, 18, 10];

	function handleStartReview() {
		appState.startFocusMode({ type: 'notebook', notebook });
	}
</script>

<div class="space-y-6">
	<!-- Review CTA -->
	{#if dueCount > 0}
		<div class="flex justify-end">
			<Button onclick={handleStartReview} class="gap-2">
				<PlayIcon class="size-4" />
				Review {dueCount} cards
			</Button>
		</div>
	{/if}

	<!-- Stats Grid -->
	<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
		<div class="bg-white dark:bg-slate-800 rounded-lg p-4 border border-slate-200 dark:border-slate-700">
			<div class="flex items-center gap-2 text-slate-500 dark:text-slate-400 mb-1">
				<LayersIcon class="size-4" />
				<span class="text-sm">Total Cards</span>
			</div>
			<p class="text-2xl font-bold text-slate-900 dark:text-white">{cards.length}</p>
		</div>

		<div class="bg-white dark:bg-slate-800 rounded-lg p-4 border border-slate-200 dark:border-slate-700">
			<div class="flex items-center gap-2 text-sky-500 mb-1">
				<CalendarIcon class="size-4" />
				<span class="text-sm">Due Today</span>
			</div>
			<p class="text-2xl font-bold text-sky-600 dark:text-sky-400">{dueCount}</p>
		</div>

		<div class="bg-white dark:bg-slate-800 rounded-lg p-4 border border-slate-200 dark:border-slate-700">
			<div class="flex items-center gap-2 text-emerald-500 mb-1">
				<TrendingUpIcon class="size-4" />
				<span class="text-sm">Retention</span>
			</div>
			<p class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">{notebook.retention}%</p>
		</div>

		<div class="bg-white dark:bg-slate-800 rounded-lg p-4 border border-slate-200 dark:border-slate-700">
			<div class="flex items-center gap-2 text-amber-500 mb-1">
				<FileStackIcon class="size-4" />
				<span class="text-sm">Sources</span>
			</div>
			<p class="text-2xl font-bold text-slate-900 dark:text-white">{sources.length}</p>
		</div>
	</div>

	<!-- Activity (mock) -->
	<div class="bg-white dark:bg-slate-800 rounded-lg p-4 border border-slate-200 dark:border-slate-700">
		<h3 class="font-semibold text-slate-900 dark:text-white mb-4">This Week</h3>
		<div class="flex items-end gap-2 h-24">
			{#each reviewsThisWeek as count, i}
				<div class="flex-1 flex flex-col items-center gap-1">
					<div
						class="w-full bg-sky-500 rounded-t"
						style="height: {(count / Math.max(...reviewsThisWeek)) * 100}%"
					></div>
					<span class="text-xs text-slate-400">
						{['M', 'T', 'W', 'T', 'F', 'S', 'S'][i]}
					</span>
				</div>
			{/each}
		</div>
	</div>
</div>
