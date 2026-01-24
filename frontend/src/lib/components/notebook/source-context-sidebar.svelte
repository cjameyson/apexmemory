<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import SourceSidebarTabs, { type SidebarTab } from './source-sidebar-tabs.svelte';
	import CardItem from './card-item.svelte';
	import { PanelRightCloseIcon, PanelRightOpenIcon } from '@lucide/svelte';
	import type { Source, Card } from '$lib/types';

	interface Props {
		source: Source;
		cards: Card[];
		activeTab?: SidebarTab;
		isCollapsed?: boolean;
		sidebarWidth?: number;
		highlightedCardIds?: string[];
		onTabChange?: (tab: SidebarTab) => void;
		onCardClick?: (card: Card) => void;
		class?: string;
	}

	let {
		source,
		cards,
		activeTab = $bindable<SidebarTab>('cards'),
		isCollapsed = $bindable(false),
		sidebarWidth = $bindable(320),
		highlightedCardIds = [],
		onTabChange,
		onCardClick,
		class: className
	}: Props = $props();

	// Filter cards for this source
	let sourceCards = $derived(cards.filter((c) => c.sourceId === source.id));

	// Resize state
	let isResizing = $state(false);
	let rafId: number | null = null;

	function handleResizeStart(e: MouseEvent) {
		isResizing = true;
		document.body.style.userSelect = 'none';
		document.body.style.cursor = 'col-resize';
		e.preventDefault();
	}

	function handleResizeMove(e: MouseEvent) {
		if (!isResizing) return;
		if (rafId) return;

		rafId = requestAnimationFrame(() => {
			// For right sidebar: calculate from right edge
			const newWidth = Math.min(500, Math.max(280, window.innerWidth - e.clientX));
			sidebarWidth = newWidth;
			rafId = null;
		});
	}

	function handleResizeEnd() {
		isResizing = false;
		document.body.style.userSelect = '';
		document.body.style.cursor = '';
		if (rafId) {
			cancelAnimationFrame(rafId);
			rafId = null;
		}
	}

	function toggleCollapse() {
		isCollapsed = !isCollapsed;
	}

	function handleTabChange(tab: SidebarTab) {
		activeTab = tab;
		onTabChange?.(tab);
	}
</script>

<svelte:window onmousemove={handleResizeMove} onmouseup={handleResizeEnd} />

<Tooltip.Provider delayDuration={300}>
	<aside
		class={cn(
			'relative flex flex-col border-l border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-900',
			!isResizing && 'transition-all duration-200',
			isCollapsed && 'w-10',
			className
		)}
		style={!isCollapsed ? `width: ${sidebarWidth}px` : undefined}
	>
		<!-- Resize handle (left edge for right sidebar) -->
		{#if !isCollapsed}
			<div
				class={cn(
					'absolute left-0 top-0 h-full w-1 cursor-col-resize transition-colors hover:bg-sky-500/50',
					isResizing && 'bg-sky-500/50'
				)}
				onmousedown={handleResizeStart}
				role="separator"
				aria-orientation="vertical"
				aria-label="Resize sidebar"
				tabindex="0"
			></div>
		{/if}

		{#if isCollapsed}
			<!-- Collapsed state: just show expand button -->
			<div class="flex flex-col items-center py-2">
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								variant="ghost"
								size="icon"
								class="size-8"
								{...props}
								onclick={toggleCollapse}
							>
								<PanelRightOpenIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content side="left">Expand sidebar</Tooltip.Content>
				</Tooltip.Root>
			</div>
		{:else}
			<!-- Expanded state -->
			<!-- Header with collapse button -->
			<div
				class="flex items-center justify-between border-b border-slate-200 px-3 py-2 dark:border-slate-700"
			>
				<span class="text-sm font-medium text-slate-700 dark:text-slate-300">Context</span>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								variant="ghost"
								size="icon"
								class="size-7"
								{...props}
								onclick={toggleCollapse}
							>
								<PanelRightCloseIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content side="left">Collapse sidebar</Tooltip.Content>
				</Tooltip.Root>
			</div>

			<!-- Tabs -->
			<SourceSidebarTabs {activeTab} onTabChange={handleTabChange} />

			<!-- Tab content -->
			<div class="flex-1 overflow-auto">
				{#if activeTab === 'cards'}
					<!-- Cards tab -->
					<div class="p-4">
						{#if sourceCards.length > 0}
							<div class="flex flex-col gap-3">
								{#each sourceCards as card (card.id)}
									<CardItem
										{card}
										onclick={() => onCardClick?.(card)}
										class={cn(highlightedCardIds.includes(card.id) && 'ring-2 ring-sky-500')}
									/>
								{/each}
							</div>
						{:else}
							<div class="py-12 text-center text-slate-500 dark:text-slate-400">
								No cards generated from this source yet.
							</div>
						{/if}
					</div>
				{:else if activeTab === 'summary'}
					<!-- Summary tab -->
					<div class="p-4">
						<div class="rounded-lg bg-slate-50 p-4 dark:bg-slate-800">
							<h3 class="mb-3 font-semibold text-slate-900 dark:text-white">AI Summary</h3>
							<p class="text-sm text-slate-600 dark:text-slate-400">
								{source.excerpt}
							</p>
							<p class="mt-4 text-sm text-slate-400 dark:text-slate-500">
								[AI-generated summary would appear here]
							</p>
						</div>
					</div>
				{:else if activeTab === 'chat'}
					<!-- Chat tab -->
					<div class="flex h-full flex-col">
						<div class="flex-1 p-4">
							<div class="py-12 text-center text-slate-500 dark:text-slate-400">
								Ask questions about this source...
							</div>
						</div>
						<div class="border-t border-slate-200 p-3 dark:border-slate-700">
							<input
								type="text"
								placeholder="Ask a question..."
								class="w-full rounded-lg border-0 bg-slate-100 px-3 py-2 text-sm text-slate-900 placeholder:text-slate-400 focus:ring-2 focus:ring-sky-500 dark:bg-slate-800 dark:text-white"
							/>
						</div>
					</div>
				{/if}
			</div>
		{/if}
	</aside>
</Tooltip.Provider>
