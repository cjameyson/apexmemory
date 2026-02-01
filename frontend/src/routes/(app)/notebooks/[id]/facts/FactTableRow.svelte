<script lang="ts">
	import {
		ChevronRightIcon,
		SquareIcon,
		SquareCheckBigIcon,
		PencilIcon,
		Trash2Icon,
		EyeIcon
	} from '@lucide/svelte';
	import type { Fact } from '$lib/types/fact';
	import { getFactDisplayText } from '$lib/utils/fact-display';
	import FactTypeBadge from '$lib/components/ui/FactTypeBadge.svelte';
	import FactExpandedContent from './FactExpandedContent.svelte';

	let {
		fact,
		selected,
		notebookId,
		onToggleSelect
	}: {
		fact: Fact;
		selected: boolean;
		notebookId: string;
		onToggleSelect: (id: string) => void;
	} = $props();

	let expanded = $state(false);

	const display = $derived(getFactDisplayText(fact.content, fact.factType));
</script>

<tr
	class="group border-border border-b transition-colors {selected
		? 'bg-primary/5'
		: 'hover:bg-accent'}"
>
	<!-- Checkbox -->
	<td class="w-12 px-3 py-2.5 text-center">
		<button
			onclick={() => onToggleSelect(fact.id)}
			class="text-muted-foreground hover:text-foreground"
			aria-label={selected ? 'Deselect fact' : 'Select fact'}
		>
			{#if selected}
				<SquareCheckBigIcon class="text-primary h-4 w-4" />
			{:else}
				<SquareIcon class="h-4 w-4" />
			{/if}
		</button>
	</td>

	<!-- Content -->
	<td class="px-3 py-2.5">
		<button
			onclick={() => (expanded = !expanded)}
			class="flex w-full items-center gap-2 text-left"
			aria-expanded={expanded}
			aria-label="Toggle fact details"
		>
			<ChevronRightIcon
				class="text-muted-foreground h-4 w-4 shrink-0 transition-transform {expanded
					? 'rotate-90'
					: ''}"
			/>
			<div class="min-w-0 flex-1">
				<p class="text-foreground line-clamp-2 text-sm font-medium">{display.primary}</p>
				{#if display.secondary}
					<p class="text-muted-foreground mt-0.5 line-clamp-1 text-xs">{display.secondary}</p>
				{/if}
				{#if fact.tags.length > 0}
					<div class="mt-1 flex flex-wrap gap-1">
						{#each fact.tags as tag}
							<span class="bg-muted text-muted-foreground rounded px-1.5 py-0.5 text-xs">{tag}</span
							>
						{/each}
					</div>
				{/if}
			</div>
		</button>
	</td>

	<!-- Type -->
	<td class="w-24 px-3 py-2.5">
		<FactTypeBadge factType={fact.factType} />
	</td>

	<!-- Cards -->
	<td class="text-foreground w-20 px-3 py-2.5 text-center text-sm">
		{fact.cardCount}
	</td>

	<!-- Due -->
	<td class="w-20 px-3 py-2.5 text-center">
		{#if fact.dueCount > 0}
			<span
				class="bg-info/15 text-info inline-flex h-6 w-6 items-center justify-center rounded-full text-xs font-medium"
			>
				{fact.dueCount}
			</span>
		{:else}
			<span class="text-muted-foreground text-sm">&mdash;</span>
		{/if}
	</td>

	<!-- Actions -->
	<td class="w-28 px-3 py-2.5">
		<div
			class="flex items-center justify-end gap-1 opacity-0 transition-opacity group-hover:opacity-100"
		>
			<button
				class="text-muted-foreground hover:bg-accent hover:text-foreground rounded p-1"
				title="Preview fact"
				aria-label="Preview fact"
			>
				<EyeIcon class="h-4 w-4" />
			</button>
			<button
				class="text-muted-foreground hover:bg-accent hover:text-foreground rounded p-1"
				title="Edit fact"
				aria-label="Edit fact"
			>
				<PencilIcon class="h-4 w-4" />
			</button>
			<button
				class="text-muted-foreground hover:bg-destructive/10 hover:text-destructive rounded p-1"
				title="Delete fact"
				aria-label="Delete fact"
			>
				<Trash2Icon class="h-4 w-4" />
			</button>
		</div>
	</td>
</tr>

{#if expanded}
	<tr class="border-border bg-muted/30 border-b">
		<td colspan="6" class="px-12 py-4">
			<FactExpandedContent {fact} {notebookId} />
		</td>
	</tr>
{/if}
