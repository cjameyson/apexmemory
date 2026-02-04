<script lang="ts">
	import type { DisplayCard, Source } from '$lib/types';
	import { SearchIcon } from '@lucide/svelte';

	interface Props {
		cards: DisplayCard[];
		sources: Source[];
		searchQuery?: string;
		selectedSourceId?: string | null;
		onSearchChange?: (query: string) => void;
		onSourceFilterChange?: (sourceId: string | null) => void;
	}

	let {
		cards,
		sources,
		searchQuery = '',
		selectedSourceId = null,
		onSearchChange,
		onSourceFilterChange
	}: Props = $props();

	// Local input state for debouncing (initial capture intentional, $effect syncs)
	// svelte-ignore state_referenced_locally
	let localSearchQuery = $state(searchQuery);

	// Sync local state when prop changes
	$effect(() => {
		localSearchQuery = searchQuery;
	});

	// Debounce search input
	let debounceTimeout: ReturnType<typeof setTimeout>;
	function handleSearchInput(e: Event) {
		const value = (e.target as HTMLInputElement).value;
		localSearchQuery = value;
		clearTimeout(debounceTimeout);
		debounceTimeout = setTimeout(() => {
			onSearchChange?.(value);
		}, 300);
	}

	function handleSourceChange(e: Event) {
		const value = (e.target as HTMLSelectElement).value;
		onSourceFilterChange?.(value === '' ? null : value);
	}

	let filteredCards = $derived.by(() => {
		let result = cards;

		if (searchQuery) {
			const query = searchQuery.toLowerCase();
			result = result.filter(
				(c) =>
					c.front.toLowerCase().includes(query) || c.back.toLowerCase().includes(query)
			);
		}

		if (selectedSourceId) {
			result = result.filter((c) => c.sourceId === selectedSourceId);
		}

		return result;
	});

	function getSourceName(sourceId: string | null): string {
		if (!sourceId) return 'No source';
		return sources.find((s) => s.id === sourceId)?.name ?? 'Unknown';
	}
</script>

<div class="space-y-4">
	<!-- Search and filters -->
	<div class="flex items-center gap-3">
		<div class="relative flex-1 max-w-md">
			<SearchIcon class="absolute left-2.5 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
			<input
				type="text"
				placeholder="Search cards..."
				value={localSearchQuery}
				oninput={handleSearchInput}
				class="w-full pl-9 pr-4 py-2 text-sm bg-card border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent"
			/>
		</div>

		<select
			value={selectedSourceId ?? ''}
			onchange={handleSourceChange}
			class="px-3 py-2 text-sm bg-card border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-ring"
		>
			<option value="">All sources</option>
			{#each sources as source (source.id)}
				<option value={source.id}>{source.name}</option>
			{/each}
		</select>
	</div>

	<!-- Results count -->
	<p class="text-sm text-muted-foreground">
		{filteredCards.length} of {cards.length} cards
	</p>

	<!-- Cards table -->
	<div class="border border-border rounded-lg overflow-hidden">
		<table class="w-full text-sm">
			<thead class="bg-muted">
				<tr class="border-b border-border">
					<th class="text-left px-4 py-3 font-medium text-muted-foreground">Front</th>
					<th class="text-left px-4 py-3 font-medium text-muted-foreground">Back</th>
					<th class="text-left px-4 py-3 font-medium text-muted-foreground">Source</th>
					<th class="text-left px-4 py-3 font-medium text-muted-foreground">Status</th>
					<th class="text-left px-4 py-3 font-medium text-muted-foreground">Interval</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-border">
				{#each filteredCards as card (card.id)}
					<tr class="hover:bg-accent cursor-pointer">
						<td class="px-4 py-3 text-foreground max-w-xs truncate">{card.front}</td>
						<td class="px-4 py-3 text-muted-foreground max-w-xs truncate">{card.back}</td>
						<td class="px-4 py-3 text-muted-foreground">{getSourceName(card.sourceId)}</td>
						<td class="px-4 py-3">
							{#if card.due}
								<span class="px-2 py-0.5 text-xs font-medium bg-primary/10 text-primary rounded">Due</span>
							{:else}
								<span class="text-muted-foreground">-</span>
							{/if}
						</td>
						<td class="px-4 py-3 text-muted-foreground">{card.interval}</td>
					</tr>
				{:else}
					<tr>
						<td colspan="5" class="px-4 py-12 text-center text-muted-foreground">
							{#if searchQuery || selectedSourceId}
								No cards match your filters
							{:else}
								No cards yet
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
