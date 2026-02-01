<script lang="ts">
	import type { PageData } from './$types';
	import FactsHeader from './FactsHeader.svelte';
	import FactsToolbar from './FactsToolbar.svelte';

	let { data }: { data: PageData } = $props();
</script>

<div id="notebook-main" class="flex flex-1 flex-col overflow-auto">
	<FactsHeader stats={data.stats} notebookId={data.notebookId} />
	<FactsToolbar />

	<div class="space-y-2 p-6">
		{#each data.facts as fact}
			<div class="rounded-lg border border-border bg-card p-3">
				<div class="flex items-center gap-2">
					<span class="text-xs font-medium uppercase text-muted-foreground">{fact.factType}</span>
					<span class="text-xs text-muted-foreground">Cards: {fact.cardCount}</span>
					<span class="text-xs text-muted-foreground">Due: {fact.dueCount}</span>
				</div>
				<pre class="mt-1 overflow-hidden text-ellipsis text-xs text-foreground">{JSON.stringify(fact.content, null, 2).slice(0, 200)}</pre>
			</div>
		{/each}

		{#if data.facts.length === 0}
			<p class="text-muted-foreground">No facts found.</p>
		{/if}
	</div>
</div>
