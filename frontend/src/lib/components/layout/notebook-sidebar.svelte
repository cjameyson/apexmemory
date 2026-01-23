<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import SidebarSection from '$lib/components/notebook/sidebar-section.svelte';
	import SourceListItem from '$lib/components/notebook/source-list-item.svelte';
	import {
		FileStackIcon,
		LayersIcon,
		PlusIcon,
		PanelLeftCloseIcon,
		PanelLeftOpenIcon,
		FolderIcon
	} from '@lucide/svelte';
	import type { Notebook, Source, Card } from '$lib/types';

	interface Props {
		notebook: Notebook;
		sources: Source[];
		cards: Card[];
		selectedSource?: Source | null;
		isCollapsed?: boolean;
		onSelectSource?: (source: Source | null) => void;
		onToggleCollapse?: () => void;
		class?: string;
	}

	let {
		notebook,
		sources,
		cards,
		selectedSource = null,
		isCollapsed = $bindable(false),
		onSelectSource,
		onToggleCollapse,
		class: className
	}: Props = $props();

	let sourcesOpen = $state(true);
	let cardsOpen = $state(true);

	let dueCount = $derived(cards.filter((c) => c.due).length);

	function selectSource(source: Source | null) {
		onSelectSource?.(source);
	}

	function toggleCollapse() {
		isCollapsed = !isCollapsed;
		onToggleCollapse?.();
	}
</script>

<aside
	class={cn(
		'bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex flex-col transition-all duration-200',
		isCollapsed ? 'w-12' : 'w-72',
		className
	)}
>
	<!-- Sidebar header -->
	<div class="p-2 flex items-center gap-2 border-b border-slate-200 dark:border-slate-800">
		{#if !isCollapsed}
			<div class="flex items-center gap-2 flex-1 min-w-0 px-2">
				<span class="text-xl">{notebook.emoji}</span>
				<span class="font-semibold text-slate-900 dark:text-white truncate">
					{notebook.name}
				</span>
			</div>
		{/if}

		<Button
			variant="ghost"
			size="icon"
			class="shrink-0 size-8"
			onclick={toggleCollapse}
			aria-label={isCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
		>
			{#if isCollapsed}
				<PanelLeftOpenIcon class="size-4" />
			{:else}
				<PanelLeftCloseIcon class="size-4" />
			{/if}
		</Button>
	</div>

	<!-- Sidebar content -->
	<div class="flex-1 overflow-auto py-2">
		{#if isCollapsed}
			<!-- Collapsed state: icon buttons only -->
			<div class="flex flex-col items-center gap-1 px-1.5">
				<Button
					variant={!selectedSource ? 'secondary' : 'ghost'}
					size="icon"
					class="size-8"
					onclick={() => selectSource(null)}
					aria-label="All sources"
				>
					<FolderIcon class="size-4" />
				</Button>
				<Button
					variant="ghost"
					size="icon"
					class="size-8"
					aria-label="Sources"
				>
					<FileStackIcon class="size-4" />
				</Button>
				<Button
					variant="ghost"
					size="icon"
					class="size-8"
					aria-label="Cards"
				>
					<LayersIcon class="size-4" />
				</Button>
			</div>
		{:else}
			<!-- Expanded state: full sidebar -->

			<!-- All sources button -->
			<div class="px-2 mb-2">
				<button
					type="button"
					onclick={() => selectSource(null)}
					class={cn(
						'w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors',
						!selectedSource
							? 'bg-sky-100 dark:bg-sky-900/30 text-sky-900 dark:text-sky-100'
							: 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300'
					)}
				>
					<FolderIcon class="size-4 text-slate-500" />
					<span class="text-sm font-medium">All sources</span>
				</button>
			</div>

			<!-- Sources section -->
			<SidebarSection
				title="Sources"
				icon={FileStackIcon}
				count={sources.length}
				bind:isOpen={sourcesOpen}
			>
				{#snippet actions()}
					<Button variant="ghost" size="icon" class="size-6" aria-label="Add source">
						<PlusIcon class="size-3.5" />
					</Button>
				{/snippet}

				<div class="space-y-0.5 pr-2">
					{#each sources as source (source.id)}
						<SourceListItem
							{source}
							isSelected={selectedSource?.id === source.id}
							onclick={() => selectSource(source)}
						/>
					{/each}

					{#if sources.length === 0}
						<div class="px-3 py-2 text-sm text-slate-400 dark:text-slate-500">
							No sources yet
						</div>
					{/if}
				</div>
			</SidebarSection>

			<!-- Cards section -->
			<SidebarSection
				title="Cards"
				icon={LayersIcon}
				count={cards.length}
				bind:isOpen={cardsOpen}
			>
				<div class="px-3 py-2 space-y-1">
					<div class="flex items-center justify-between text-sm">
						<span class="text-slate-600 dark:text-slate-400">Total</span>
						<span class="font-medium text-slate-900 dark:text-white">{cards.length}</span>
					</div>
					{#if dueCount > 0}
						<div class="flex items-center justify-between text-sm">
							<span class="text-slate-600 dark:text-slate-400">Due</span>
							<span class="font-medium text-sky-600 dark:text-sky-400">{dueCount}</span>
						</div>
					{/if}
				</div>
			</SidebarSection>
		{/if}
	</div>
</aside>
