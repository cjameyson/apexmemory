<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import SourceViewer from './source-viewer.svelte';
	import { ZapIcon, Maximize2Icon, Minimize2Icon } from '@lucide/svelte';
	import type { Source, Card } from '$lib/types';

	interface Props {
		source: Source;
		cards: Card[];
		isExpanded?: boolean;
		onClose?: () => void;
		onStartReview?: () => void;
		onToggleExpand?: () => void;
		class?: string;
	}

	let {
		source,
		cards,
		isExpanded = false,
		onClose,
		onStartReview,
		onToggleExpand,
		class: className
	}: Props = $props();

	// Count cards for this source
	let sourceCardsCount = $derived(cards.filter((c) => c.sourceId === source.id).length);
	let dueCount = $derived(cards.filter((c) => c.sourceId === source.id && c.due).length);

	function handleGenerateCards() {
		// TODO: Implement AI card generation
	}
</script>

<div
	id="source-detail"
	class={cn(
		'flex flex-1 flex-col overflow-hidden rounded-none border-none border-border bg-card',
		className
	)}
>
	<!-- Header -->
	<div class="flex h-12 items-center gap-3 border-b border-border px-4">
		<!-- Source info -->
		<div class="flex min-w-0 flex-1 items-center gap-2">
			<div class="min-w-0">
				<h2 class="truncate font-semibold text-foreground">
					{source.name}
				</h2>
				<div class="text-xs text-muted-foreground">
					{sourceCardsCount} cards
					{#if dueCount > 0}
						<span class="text-primary">({dueCount} due)</span>
					{/if}
				</div>
			</div>
		</div>

		<!-- Actions -->
		<div class="flex items-center gap-2">
			{#if dueCount > 0}
				<Button variant="outline" size="sm" class="hidden sm:inline-flex" onclick={onStartReview}>
					<ZapIcon class="size-4" />
					<span class="hidden sm:inline">Review ({dueCount})</span>
				</Button>
			{/if}

			<Button
				variant="ghost"
				size="icon"
				class="size-8"
				title={isExpanded ? 'Collapse panel' : 'Expand panel'}
				onclick={onToggleExpand}
			>
				{#if isExpanded}
					<Minimize2Icon class="size-4" />
				{:else}
					<Maximize2Icon class="size-4" />
				{/if}
			</Button>
		</div>
	</div>

	<!-- Source viewer (full width now, sidebar is at page level) -->
	<SourceViewer {source} onGenerateCards={handleGenerateCards} class="flex-1" />
</div>
