<script lang="ts">
	import type { PageData } from './$types';
	import FactsHeader from './FactsHeader.svelte';
	import FactsToolbar from './FactsToolbar.svelte';
	import FactsTable from './FactsTable.svelte';

	let { data }: { data: PageData } = $props();

	let selectedIds = $state<Set<string>>(new Set());

	// Clear selection when facts list changes (filter/search/page navigation)
	$effect(() => {
		data.facts;
		selectedIds = new Set();
	});

	function toggleSelect(id: string) {
		const next = new Set(selectedIds);
		if (next.has(id)) {
			next.delete(id);
		} else {
			next.add(id);
		}
		selectedIds = next;
	}

	function toggleSelectAll() {
		if (selectedIds.size === data.facts.length) {
			selectedIds = new Set();
		} else {
			selectedIds = new Set(data.facts.map((f) => f.id));
		}
	}
</script>

<div id="notebook-main" class="flex flex-1 flex-col overflow-auto">
	<FactsHeader stats={data.stats} notebookId={data.notebookId} />
	<FactsToolbar />

	{#if data.facts.length > 0}
		<FactsTable
			facts={data.facts}
			{selectedIds}
			onToggleSelect={toggleSelect}
			onToggleAll={toggleSelectAll}
		/>
	{:else}
		<div class="p-6">
			<p class="text-muted-foreground">No facts found.</p>
		</div>
	{/if}
</div>
