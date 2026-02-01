<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { LayersIcon, PlusIcon } from '@lucide/svelte';

	interface Props {
		totalFacts: number;
		searchQuery: string;
		typeFilter: string;
	}

	let { totalFacts, searchQuery, typeFilter }: Props = $props();

	let hasFilters = $derived(!!searchQuery || !!typeFilter);
	let showNoFacts = $derived(totalFacts === 0 && !hasFilters);

	function clearFilters() {
		const url = new URL($page.url);
		url.searchParams.delete('q');
		url.searchParams.delete('type');
		url.searchParams.delete('page');
		goto(url.toString(), { replaceState: true });
	}
</script>

<div class="flex flex-1 flex-col items-center justify-center gap-4 p-12">
	{#if showNoFacts}
		<div class="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
			<LayersIcon class="h-8 w-8 text-muted-foreground" />
		</div>
		<div class="text-center">
			<h3 class="text-lg font-medium text-foreground">No facts yet</h3>
			<p class="mt-1 text-sm text-muted-foreground">
				Create your first fact to start building your knowledge base.
			</p>
		</div>
		<button
			disabled
			class="mt-2 inline-flex items-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground opacity-50"
		>
			<PlusIcon class="h-4 w-4" />
			Create your first fact
		</button>
	{:else}
		<div class="text-center">
			<h3 class="text-lg font-medium text-foreground">No facts found</h3>
			<p class="mt-1 text-sm text-muted-foreground">
				{#if searchQuery}
					No facts match "{searchQuery}". Try a different search term.
				{:else if typeFilter}
					No {typeFilter} facts found. Try a different filter.
				{:else}
					No facts match the current filters.
				{/if}
			</p>
		</div>
		<button
			onclick={clearFilters}
			class="mt-2 text-sm font-medium text-primary hover:underline"
		>
			Clear filters
		</button>
	{/if}
</div>
