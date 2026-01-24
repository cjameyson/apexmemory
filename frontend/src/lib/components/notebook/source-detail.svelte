<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import SourceViewer from './source-viewer.svelte';
	import SourceContextSidebar from './source-context-sidebar.svelte';
	import { ZapIcon, Maximize2Icon, Minimize2Icon } from '@lucide/svelte';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { appState } from '$lib/stores/app.svelte';
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

	function handleCardClick(card: Card) {
		// Future: Scroll source viewer to card's linked section
		// appState.setActiveSourceSection(card.sourceSection);
	}
</script>

<div
	id="source-detail"
	class={cn(
		'flex flex-1 flex-col overflow-hidden rounded-none border-none border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900',
		className
	)}
>
	<!-- Header -->
	<div class="flex items-center gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
		<!-- Source info -->
		<div class="flex min-w-0 flex-1 items-center gap-2">
			<div class="min-w-0">
				<h2 class="truncate font-semibold text-slate-900 dark:text-white">
					{source.name}
				</h2>
				<div class="text-xs text-slate-500 dark:text-slate-400">
					{sourceCardsCount} cards
					{#if dueCount > 0}
						<span class="text-sky-600 dark:text-sky-400">({dueCount} due)</span>
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

			<Tooltip.Provider delayDuration={300}>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								variant="ghost"
								size="icon"
								class="size-8"
								{...props}
								onclick={onToggleExpand}
							>
								{#if isExpanded}
									<Minimize2Icon class="size-4" />
								{:else}
									<Maximize2Icon class="size-4" />
								{/if}
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>{isExpanded ? 'Collapse panel' : 'Expand panel'}</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		</div>
	</div>

	<!-- Two-pane content area -->
	<div class="flex flex-1 overflow-hidden">
		<!-- Left pane: Source viewer (always visible) -->
		<SourceViewer {source} onGenerateCards={handleGenerateCards} class="min-w-0 flex-1" />

		<!-- Right pane: Context sidebar (collapsible/resizable) -->
		<SourceContextSidebar
			{source}
			{cards}
			bind:isCollapsed={appState.sourceContextSidebarCollapsed}
			bind:sidebarWidth={appState.sourceContextSidebarWidth}
			highlightedCardIds={appState.highlightedCardIds}
			onCardClick={handleCardClick}
			class="hidden md:flex"
		/>
	</div>
</div>
