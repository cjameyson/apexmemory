<script lang="ts">
	import { SquareIcon, SquareCheckBigIcon, MinusSquareIcon } from '@lucide/svelte';
	import type { Fact } from '$lib/types/fact';
	import FactTableRow from './FactTableRow.svelte';

	let {
		facts,
		selectedIds,
		notebookId,
		onToggleSelect,
		onToggleAll
	}: {
		facts: Fact[];
		selectedIds: Set<string>;
		notebookId: string;
		onToggleSelect: (id: string) => void;
		onToggleAll: () => void;
	} = $props();

	const allSelected = $derived(facts.length > 0 && selectedIds.size === facts.length);
	const someSelected = $derived(selectedIds.size > 0 && !allSelected);
</script>

<div class="overflow-x-auto pb-6">
	<table class="facts-table w-full border-collapse">
		<thead>
			<tr class="border-border border-b">
				<th class="w-12 px-3 py-2 text-center">
					<button onclick={onToggleAll} class="text-muted-foreground hover:text-foreground" aria-label="Select all facts">
						{#if allSelected}
							<SquareCheckBigIcon class="text-primary h-4 w-4" />
						{:else if someSelected}
							<MinusSquareIcon class="text-primary h-4 w-4" />
						{:else}
							<SquareIcon class="h-4 w-4" />
						{/if}
					</button>
				</th>
				<th
					class="text-muted-foreground px-3 py-2 text-left text-xs font-medium tracking-wider uppercase"
					>Content</th
				>
				<th
					class="text-muted-foreground w-24 px-3 py-2 text-left text-xs font-medium tracking-wider uppercase"
					>Type</th
				>
				<th
					class="text-muted-foreground w-20 px-3 py-2 text-center text-xs font-medium tracking-wider uppercase"
					>Cards</th
				>
				<th
					class="text-muted-foreground w-20 px-3 py-2 text-center text-xs font-medium tracking-wider uppercase"
					>Due</th
				>
				<th
					class="text-muted-foreground w-28 px-3 py-2 text-right text-xs font-medium tracking-wider uppercase"
					>Actions</th
				>
			</tr>
		</thead>
		<tbody>
			{#each facts as fact (fact.id)}
				<FactTableRow {fact} selected={selectedIds.has(fact.id)} {notebookId} {onToggleSelect} />
			{/each}
		</tbody>
	</table>
</div>
