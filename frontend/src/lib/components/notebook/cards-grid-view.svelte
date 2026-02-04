<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import CardItem from './card-item.svelte';
	import { PlusIcon } from '@lucide/svelte';
	import type { DisplayCard, Source } from '$lib/types';

	interface Props {
		cards: DisplayCard[];
		sources: Source[];
		viewMode?: 'all' | 'due' | 'mastered';
		onViewModeChange?: (mode: 'all' | 'due' | 'mastered') => void;
		onCardClick?: (card: DisplayCard) => void;
		class?: string;
	}

	let {
		cards,
		sources,
		viewMode = 'all',
		onViewModeChange,
		onCardClick,
		class: className
	}: Props = $props();

	// Filter cards based on view mode
	let filteredCards = $derived.by(() => {
		switch (viewMode) {
			case 'due':
				return cards.filter((c) => c.due);
			case 'mastered':
				return cards.filter((c) => !c.due);
			default:
				return cards;
		}
	});

	let dueCount = $derived(cards.filter((c) => c.due).length);
	let masteredCount = $derived(cards.filter((c) => !c.due).length);

	// Get source for a card
	function getSource(card: DisplayCard): Source | undefined {
		return sources.find((s) => s.id === card.sourceId);
	}

	function setViewMode(mode: 'all' | 'due' | 'mastered') {
		onViewModeChange?.(mode);
	}
</script>

<div class={cn('flex-1 bg-card rounded-2xl border border-border p-6', className)}>
	<!-- Header -->
	<div class="flex items-center justify-between mb-6">
		<div>
			<h2 class="text-lg font-semibold text-foreground">Cards</h2>
			<p class="text-sm text-muted-foreground">
				{filteredCards.length} {viewMode === 'all' ? 'total' : viewMode}
			</p>
		</div>

		<div class="flex items-center gap-2">
			<!-- View mode toggle -->
			<div class="flex items-center bg-muted rounded-lg p-1">
				<button
					type="button"
					onclick={() => setViewMode('all')}
					class={cn(
						'px-3 py-1.5 text-sm font-medium rounded-md transition-colors',
						viewMode === 'all'
							? 'bg-background text-foreground shadow-sm'
							: 'text-muted-foreground hover:text-foreground'
					)}
				>
					All
				</button>
				<button
					type="button"
					onclick={() => setViewMode('due')}
					class={cn(
						'px-3 py-1.5 text-sm font-medium rounded-md transition-colors',
						viewMode === 'due'
							? 'bg-background text-foreground shadow-sm'
							: 'text-muted-foreground hover:text-foreground'
					)}
				>
					Due ({dueCount})
				</button>
				<button
					type="button"
					onclick={() => setViewMode('mastered')}
					class={cn(
						'px-3 py-1.5 text-sm font-medium rounded-md transition-colors',
						viewMode === 'mastered'
							? 'bg-background text-foreground shadow-sm'
							: 'text-muted-foreground hover:text-foreground'
					)}
				>
					Mastered
				</button>
			</div>

			<!-- Add card button -->
			<Button variant="outline" size="sm" class="gap-1.5">
				<PlusIcon class="size-4" />
				<span class="hidden sm:inline">Add card</span>
			</Button>
		</div>
	</div>

	<!-- Cards grid -->
	{#if filteredCards.length > 0}
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each filteredCards as card (card.id)}
				<CardItem
					{card}
					source={getSource(card)}
					onclick={() => onCardClick?.(card)}
				/>
			{/each}
		</div>
	{:else}
		<div class="flex flex-col items-center justify-center py-12 text-center">
			<div class="text-muted-foreground mb-2">
				{#if viewMode === 'due'}
					No cards due for review
				{:else if viewMode === 'mastered'}
					No mastered cards yet
				{:else}
					No cards in this notebook
				{/if}
			</div>
			<Button variant="outline" size="sm" class="gap-1.5">
				<PlusIcon class="size-4" />
				Create your first card
			</Button>
		</div>
	{/if}
</div>
