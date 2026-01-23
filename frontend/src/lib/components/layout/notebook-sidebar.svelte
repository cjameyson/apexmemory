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
		SettingsIcon,
		SearchIcon
	} from '@lucide/svelte';
	import type { Notebook, Source, Card } from '$lib/types';

	interface Props {
		notebook: Notebook;
		sources: Source[];
		cards: Card[];
		selectedSource?: Source | null;
		isCollapsed?: boolean;
		sidebarWidth?: number;
		onSelectSource?: (source: Source | null) => void;
		onToggleCollapse?: () => void;
		onOpenSettings?: () => void;
		class?: string;
	}

	let {
		notebook,
		sources,
		cards,
		selectedSource = null,
		isCollapsed = $bindable(false),
		sidebarWidth = $bindable(288),
		onSelectSource,
		onToggleCollapse,
		onOpenSettings,
		class: className
	}: Props = $props();

	let isResizing = $state(false);

	function handleResizeStart(e: MouseEvent) {
		isResizing = true;
		e.preventDefault();
	}

	function handleResizeMove(e: MouseEvent) {
		if (!isResizing) return;
		const newWidth = Math.min(400, Math.max(200, e.clientX));
		sidebarWidth = newWidth;
	}

	function handleResizeEnd() {
		isResizing = false;
	}

	let sourcesOpen = $state(true);
	let cardsOpen = $state(true);

	function selectSource(source: Source | null) {
		onSelectSource?.(source);
	}

	function toggleCollapse() {
		isCollapsed = !isCollapsed;
		onToggleCollapse?.();
	}
</script>

<svelte:window
	onmousemove={handleResizeMove}
	onmouseup={handleResizeEnd}
/>

<aside
	class={cn(
		'bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex flex-col transition-all duration-200 relative',
		isCollapsed && 'w-12',
		className
	)}
	style={!isCollapsed ? `width: ${sidebarWidth}px` : undefined}
>
	<!-- Sidebar header -->
	<div class="p-2 flex items-center gap-2 border-b border-slate-200 dark:border-slate-800">
		{#if !isCollapsed}
			<a
				href="/notebooks/{notebook.id}"
				class="flex items-center gap-2 flex-1 min-w-0 px-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
			>
				<span class="text-xl">{notebook.emoji}</span>
				<span class="font-semibold text-slate-900 dark:text-white truncate">
					{notebook.name}
				</span>
			</a>
		{/if}

		{#if !isCollapsed}
			<!-- Settings button (only when expanded) -->
			<Button
				variant="ghost"
				size="icon"
				class="shrink-0 size-8 text-slate-400 hover:text-slate-600 dark:text-slate-500 dark:hover:text-slate-300"
				onclick={() => onOpenSettings?.()}
				aria-label="Notebook settings"
			>
				<SettingsIcon class="size-4" />
			</Button>
		{/if}

		<!-- Collapse button -->
		<Button
			variant="ghost"
			size="icon"
			class="shrink-0 size-8 text-slate-400 hover:text-slate-600 dark:text-slate-500 dark:hover:text-slate-300"
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

			<!-- Sources section -->
			<SidebarSection
				title="Sources"
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
				count={cards.length}
				bind:isOpen={cardsOpen}
			>
				<div class="px-2 py-1">
					<a
						href="/notebooks/{notebook.id}/cards"
						class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
					>
						<SearchIcon class="size-4" />
						<span>Browse</span>
					</a>
				</div>
			</SidebarSection>
		{/if}
	</div>

	<!-- Resize handle (only when expanded) -->
	{#if !isCollapsed}
		<div
			class={cn(
				'absolute top-0 right-0 w-1 h-full cursor-col-resize hover:bg-sky-500/50 transition-colors',
				isResizing && 'bg-sky-500/50'
			)}
			onmousedown={handleResizeStart}
			role="separator"
			aria-orientation="vertical"
			aria-label="Resize sidebar"
			tabindex="0"
		></div>
	{/if}
</aside>
