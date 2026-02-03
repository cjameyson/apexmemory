<script lang="ts">
	import type { Notebook, Source, Card } from '$lib/types';
	import { pushState } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import {
		PlayIcon,
		TrendingUpIcon,
		LayersIcon,
		FileStackIcon,
		CalendarIcon
	} from '@lucide/svelte';

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
		pushState('', {
			focusMode: {
				type: 'notebook',
				mode: 'scheduled',
				notebookId: notebook.id,
				notebookName: notebook.name,
				notebookEmoji: notebook.emoji,
				currentIndex: 0
			}
		});
	}
</script>

<div class="space-y-6 p-3">
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
	<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
		<div
			class="rounded-lg border border-border bg-card p-4"
		>
			<div class="mb-1 flex items-center gap-2 text-muted-foreground">
				<LayersIcon class="size-4" />
				<span class="text-sm">Total Cards</span>
			</div>
			<p class="text-2xl font-bold text-foreground">{cards.length}</p>
		</div>

		<div
			class="rounded-lg border border-border bg-card p-4"
		>
			<div class="mb-1 flex items-center gap-2 text-primary">
				<CalendarIcon class="size-4" />
				<span class="text-sm">Due Today</span>
			</div>
			<p class="text-2xl font-bold text-primary">{dueCount}</p>
		</div>

		<div
			class="rounded-lg border border-border bg-card p-4"
		>
			<div class="mb-1 flex items-center gap-2 text-success">
				<TrendingUpIcon class="size-4" />
				<span class="text-sm">Retention</span>
			</div>
			<p class="text-2xl font-bold text-success">{notebook.retention}%</p>
		</div>

		<div
			class="rounded-lg border border-border bg-card p-4"
		>
			<div class="mb-1 flex items-center gap-2 text-warning">
				<FileStackIcon class="size-4" />
				<span class="text-sm">Sources</span>
			</div>
			<p class="text-2xl font-bold text-foreground">{sources.length}</p>
		</div>
	</div>

	<!-- Activity (mock) -->
	<div
		class="rounded-lg border border-border bg-card p-4"
	>
		<h3 class="mb-4 font-semibold text-foreground">This Week</h3>
		<div class="flex h-24 items-end gap-2">
			{#each reviewsThisWeek as count, i}
				<div class="flex flex-1 flex-col items-center gap-1">
					<div
						class="w-full rounded-t bg-primary"
						style="height: {(count / Math.max(...reviewsThisWeek)) * 100}%"
					></div>
					<span class="text-xs text-muted-foreground">
						{['M', 'T', 'W', 'T', 'F', 'S', 'S'][i]}
					</span>
				</div>
			{/each}
		</div>
	</div>
</div>
