<script lang="ts">
	import { page } from '$app/stores';
	import type { PageData } from './$types';
	import FactsHeader from './FactsHeader.svelte';
	import FactsToolbar from './FactsToolbar.svelte';
	import FactsTable from './FactsTable.svelte';
	import BulkActionsBar from './BulkActionsBar.svelte';
	import Pagination from './Pagination.svelte';
	import EmptyState from './EmptyState.svelte';

	let { data }: { data: PageData } = $props();

	let factsHeader: FactsHeader;

	let searchQuery = $derived($page.url.searchParams.get('q') || '');
	let typeFilter = $derived($page.url.searchParams.get('type') || '');

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
	<FactsHeader bind:this={factsHeader} stats={data.stats} notebookId={data.notebookId} />
	<FactsToolbar />

	{#if selectedIds.size > 0}
		<BulkActionsBar selectedCount={selectedIds.size} onClear={() => (selectedIds = new Set())} />
	{/if}

	{#if data.facts.length > 0}
		<FactsTable
			facts={data.facts}
			{selectedIds}
			notebookId={data.notebookId}
			onToggleSelect={toggleSelect}
			onToggleAll={toggleSelectAll}
			onedit={(factId) => factsHeader.openEdit(factId)}
		/>
		<Pagination
			page={data.pagination.page}
			totalPages={data.pagination.totalPages}
			total={data.pagination.total}
			pageSize={data.pagination.pageSize}
		/>
	{:else}
		<EmptyState totalFacts={data.stats.totalFacts} {searchQuery} {typeFilter} />
	{/if}
</div>
