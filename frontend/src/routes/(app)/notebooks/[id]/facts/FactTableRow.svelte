<script lang="ts">
	import { ChevronRightIcon, SquareIcon, SquareCheckBigIcon, PencilIcon, Trash2Icon, MoreHorizontalIcon } from '@lucide/svelte';
	import type { Fact } from '$lib/types/fact';
	import { getFactDisplayText } from '$lib/utils/fact-display';
	import FactTypeBadge from '$lib/components/ui/FactTypeBadge.svelte';

	let {
		fact,
		selected,
		onToggleSelect
	}: {
		fact: Fact;
		selected: boolean;
		onToggleSelect: (id: string) => void;
	} = $props();

	let expanded = $state(false);

	const display = $derived(getFactDisplayText(fact.content, fact.factType));
</script>

<tr
	class="group border-b border-border transition-colors {selected ? 'bg-primary/5' : 'hover:bg-accent'}"
>
	<!-- Checkbox -->
	<td class="w-12 px-3 py-2.5 text-center">
		<button onclick={() => onToggleSelect(fact.id)} class="text-muted-foreground hover:text-foreground" aria-label={selected ? 'Deselect fact' : 'Select fact'}>
			{#if selected}
				<SquareCheckBigIcon class="h-4 w-4 text-primary" />
			{:else}
				<SquareIcon class="h-4 w-4" />
			{/if}
		</button>
	</td>

	<!-- Content -->
	<td class="px-3 py-2.5">
		<button
			onclick={() => (expanded = !expanded)}
			class="flex w-full items-start gap-2 text-left"
			aria-expanded={expanded}
			aria-label="Toggle fact details"
		>
			<ChevronRightIcon
				class="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground transition-transform {expanded ? 'rotate-90' : ''}"
			/>
			<div class="min-w-0 flex-1">
				<p class="text-sm font-medium text-foreground">{display.primary}</p>
				{#if display.secondary}
					<p class="mt-0.5 text-xs text-muted-foreground">{display.secondary}</p>
				{/if}
				{#if fact.tags.length > 0}
					<div class="mt-1 flex flex-wrap gap-1">
						{#each fact.tags as tag}
							<span class="rounded bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">{tag}</span>
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
	<td class="w-20 px-3 py-2.5 text-center text-sm text-foreground">
		{fact.cardCount}
	</td>

	<!-- Due -->
	<td class="w-20 px-3 py-2.5 text-center">
		{#if fact.dueCount > 0}
			<span class="inline-flex h-6 w-6 items-center justify-center rounded-full bg-info/15 text-xs font-medium text-info">
				{fact.dueCount}
			</span>
		{:else}
			<span class="text-sm text-muted-foreground">&mdash;</span>
		{/if}
	</td>

	<!-- Actions -->
	<td class="w-28 px-3 py-2.5">
		<div class="flex items-center justify-end gap-1 opacity-0 transition-opacity group-hover:opacity-100">
			<button class="rounded p-1 text-muted-foreground hover:bg-accent hover:text-foreground" aria-label="Edit">
				<PencilIcon class="h-4 w-4" />
			</button>
			<button class="rounded p-1 text-muted-foreground hover:bg-destructive/10 hover:text-destructive" aria-label="Delete">
				<Trash2Icon class="h-4 w-4" />
			</button>
			<button class="rounded p-1 text-muted-foreground hover:bg-accent hover:text-foreground" aria-label="More options">
				<MoreHorizontalIcon class="h-4 w-4" />
			</button>
		</div>
	</td>
</tr>

{#if expanded}
	<tr class="border-b border-border bg-muted/30">
		<td colspan="6" class="px-12 py-4">
			<p class="text-sm text-muted-foreground">Cards detail coming in Phase 5</p>
		</td>
	</tr>
{/if}
