<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
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

	let addSourceOpen = $state(false);
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
			const newWidth = Math.min(400, Math.max(200, e.clientX));
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

	let sourcesOpen = $state(false);
	let cardsOpen = $state(true);

	function selectSource(source: Source | null) {
		onSelectSource?.(source);
	}

	function toggleCollapse() {
		isCollapsed = !isCollapsed;
		onToggleCollapse?.();
	}
</script>

<svelte:window onmousemove={handleResizeMove} onmouseup={handleResizeEnd} />

<aside
	class={cn(
		'border-border bg-card relative flex flex-col border-r',
		!isResizing && 'transition-all duration-200',
		isCollapsed && 'w-12',
		className
	)}
	style={!isCollapsed ? `width: ${sidebarWidth}px` : undefined}
>
	<!-- Sidebar header -->
	<div class="border-border flex h-12 items-center gap-2 border-b px-2">
		{#if !isCollapsed}
			<a
				href="/notebooks/{notebook.id}"
				class="hover:bg-accent flex min-w-0 flex-1 items-center gap-2 rounded-lg px-2 transition-colors"
			>
				<span class="text-xl">{notebook.emoji}</span>
				<span class="text-foreground truncate font-semibold">
					{notebook.name}
				</span>
			</a>
		{/if}

		<!-- Collapse button -->
		<Button
			variant="ghost"
			size="icon"
			class="size-8 shrink-0"
			title={isCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
			onclick={toggleCollapse}
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
				<Button variant="ghost" size="icon" class="size-8" title="Flashcards">
					<LayersIcon class="size-4" />
				</Button>
				<Button variant="ghost" size="icon" class="size-8" title="Sources">
					<FileStackIcon class="size-4" />
				</Button>
			</div>
		{:else}
			<!-- Expanded state: full sidebar -->

			<!-- Cards section -->
			<SidebarSection title="Flashcards" count={cards.length} bind:isOpen={cardsOpen} class="border-border border-b pb-2">
				<div class="px-2 py-1">
					<a
						href="/notebooks/{notebook.id}/facts"
						class="text-muted-foreground hover:bg-accent flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
					>
						<SearchIcon class="size-4" />
						<span>Browse</span>
					</a>
				</div>
			</SidebarSection>

			<!-- Sources section -->
			<SidebarSection
				title="Sources"
				count={sources.length}
				bind:isOpen={sourcesOpen}
			>
				{#snippet actions()}
					<Button variant="ghost" size="icon" class="size-6" title="Add source" onclick={() => (addSourceOpen = true)}>
						<PlusIcon class="size-3.5" />
					</Button>
				{/snippet}

				<div class="space-y-0.5 pr-2">
					{#each sources as source (source.id)}
						<SourceListItem
							{source}
							isSelected={selectedSource?.id === source.id}
							href="/notebooks/{notebook.id}/sources/{source.id}"
						/>
					{/each}

					{#if sources.length === 0}
						<div class="text-muted-foreground px-3 py-2 text-sm">No sources yet</div>
					{/if}
				</div>
			</SidebarSection>
		{/if}
	</div>

	<!-- Resize handle (only when expanded) -->
	{#if !isCollapsed}
		<div
			class={cn(
				'hover:bg-primary/50 absolute top-0 right-0 h-full w-1 cursor-col-resize transition-colors',
				isResizing && 'bg-primary/50'
			)}
			onmousedown={handleResizeStart}
			role="separator"
			aria-orientation="vertical"
			aria-label="Resize sidebar"
			tabindex="0"
		></div>
	{/if}
</aside>

<Dialog.Root bind:open={addSourceOpen}>
	<Dialog.Content class="sm:max-w-md">
		<Dialog.Header>
			<Dialog.Title>Add Source</Dialog.Title>
			<Dialog.Description>Upload PDFs, slides, audio, links, and other study materials.</Dialog.Description>
		</Dialog.Header>
		<div class="text-muted-foreground flex flex-col items-center gap-2 py-8 text-center text-sm">
			<FileStackIcon class="text-muted-foreground/50 size-10" />
			<p>Source uploads are coming soon.</p>
		</div>
		<Dialog.Footer>
			<Button variant="secondary" size="sm" onclick={() => (addSourceOpen = false)}>Close</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
