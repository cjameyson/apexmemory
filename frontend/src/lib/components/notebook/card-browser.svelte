<script lang="ts">
	import type { Card, Source } from '$lib/types';
	import { SearchIcon } from '@lucide/svelte';

	interface Props {
		cards: Card[];
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

	// Local input state for debouncing
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

	function getSourceName(sourceId: string): string {
		return sources.find((s) => s.id === sourceId)?.name ?? 'Unknown';
	}
</script>

<div class="space-y-4">
	<!-- Search and filters -->
	<div class="flex items-center gap-3">
		<div class="relative flex-1 max-w-md">
			<SearchIcon class="absolute left-2.5 top-1/2 -translate-y-1/2 size-4 text-slate-400" />
			<input
				type="text"
				placeholder="Search cards..."
				value={localSearchQuery}
				oninput={handleSearchInput}
				class="w-full pl-9 pr-4 py-2 text-sm bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent"
			/>
		</div>

		<select
			value={selectedSourceId ?? ''}
			onchange={handleSourceChange}
			class="px-3 py-2 text-sm bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500"
		>
			<option value="">All sources</option>
			{#each sources as source (source.id)}
				<option value={source.id}>{source.name}</option>
			{/each}
		</select>
	</div>

	<!-- Results count -->
	<p class="text-sm text-slate-500 dark:text-slate-400">
		{filteredCards.length} of {cards.length} cards
	</p>

	<!-- Cards table -->
	<div class="border border-slate-200 dark:border-slate-700 rounded-lg overflow-hidden">
		<table class="w-full text-sm">
			<thead class="bg-slate-50 dark:bg-slate-800/50">
				<tr class="border-b border-slate-200 dark:border-slate-700">
					<th class="text-left px-4 py-3 font-medium text-slate-600 dark:text-slate-400">Front</th>
					<th class="text-left px-4 py-3 font-medium text-slate-600 dark:text-slate-400">Back</th>
					<th class="text-left px-4 py-3 font-medium text-slate-600 dark:text-slate-400">Source</th>
					<th class="text-left px-4 py-3 font-medium text-slate-600 dark:text-slate-400">Status</th>
					<th class="text-left px-4 py-3 font-medium text-slate-600 dark:text-slate-400">Interval</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-slate-200 dark:divide-slate-700">
				{#each filteredCards as card (card.id)}
					<tr class="hover:bg-slate-50 dark:hover:bg-slate-800/50 cursor-pointer">
						<td class="px-4 py-3 text-slate-900 dark:text-white max-w-xs truncate">{card.front}</td>
						<td class="px-4 py-3 text-slate-500 dark:text-slate-400 max-w-xs truncate">{card.back}</td>
						<td class="px-4 py-3 text-slate-500 dark:text-slate-400">{getSourceName(card.sourceId)}</td>
						<td class="px-4 py-3">
							{#if card.due}
								<span class="px-2 py-0.5 text-xs font-medium bg-sky-100 dark:bg-sky-900/30 text-sky-700 dark:text-sky-300 rounded">Due</span>
							{:else}
								<span class="text-slate-400">-</span>
							{/if}
						</td>
						<td class="px-4 py-3 text-slate-500 dark:text-slate-400">{card.interval}</td>
					</tr>
				{:else}
					<tr>
						<td colspan="5" class="px-4 py-12 text-center text-slate-500 dark:text-slate-400">
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
