<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import SourceIcon from '$lib/components/ui/source-icon.svelte';
	import SourceTabs from './source-tabs.svelte';
	import SourceToolbar from './source-toolbar.svelte';
	import SourceContent from './source-content.svelte';
	import { ArrowLeftIcon, ZapIcon, MaximizeIcon, MinimizeIcon } from '@lucide/svelte';
	import type { Source, Card } from '$lib/types';

	type Tab = 'source' | 'cards' | 'summary' | 'chat';

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

	let activeTab = $state<Tab>('source');

	// Count cards for this source
	let sourceCardsCount = $derived(cards.filter((c) => c.sourceId === source.id).length);
	let dueCount = $derived(cards.filter((c) => c.sourceId === source.id && c.due).length);

	function handleTabChange(tab: Tab) {
		activeTab = tab;
	}

	function handleGenerateCards() {
		// TODO: Implement AI card generation
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
	<div class="flex items-center gap-3 border-none border-slate-200 px-4 py-3 dark:border-slate-800">
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
			<!-- disable for now -->
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
				onclick={onToggleExpand}
				aria-label={isExpanded ? 'Minimize' : 'Maximize'}
			>
				{#if isExpanded}
					<MinimizeIcon class="size-4" />
				{:else}
					<MaximizeIcon class="size-4" />
				{/if}
			</Button>
		</div>
	</div>

	<!-- Tabs -->
	<SourceTabs {activeTab} onTabChange={handleTabChange} />

	<!-- Toolbar (only on source tab) -->
	{#if activeTab === 'source'}
		<SourceToolbar type={source.type} onGenerateCards={handleGenerateCards} />
	{/if}

	<!-- Content -->
	<SourceContent {source} {cards} {activeTab} />
</div>
