<script lang="ts">
	import type { Region } from './types';
	import { ScrollArea } from '$lib/components/ui/scroll-area';
	import Input from '$lib/components/ui/input/input.svelte';
	import LabelPanelItem from './LabelPanelItem.svelte';
	import { Search, SquareDashed } from '@lucide/svelte';

	interface Props {
		regions: Region[];
		selectedRegionId: string | null;
		focusLabelRegionId?: string | null;
		title?: string;
		titleError?: boolean;
		regionLabelErrors?: Set<string>;
		onSelectRegion?: (id: string) => void;
		onUpdateRegion?: (
			id: string,
			updates: Partial<Pick<Region, 'label' | 'hint' | 'backExtra'>>
		) => void;
		onDeleteRegion?: (id: string) => void;
		onFilterChange?: (visibleIds: Set<string> | null) => void;
		onTitleChange?: (value: string) => void;
	}

	let {
		regions,
		selectedRegionId,
		focusLabelRegionId = null,
		title = '',
		titleError = false,
		regionLabelErrors,
		onSelectRegion,
		onUpdateRegion,
		onDeleteRegion,
		onFilterChange,
		onTitleChange
	}: Props = $props();

	let filterText = $state('');

	// Filter regions based on search text
	let filteredRegions = $derived(
		filterText.trim()
			? regions.filter(
					(r) =>
						r.label.toLowerCase().includes(filterText.toLowerCase()) ||
						r.hint?.toLowerCase().includes(filterText.toLowerCase()) ||
						r.backExtra?.toLowerCase().includes(filterText.toLowerCase())
				)
			: regions
	);

	// Notify parent when visible region set changes
	$effect(() => {
		if (!filterText.trim()) {
			onFilterChange?.(null); // null = show all
		} else {
			onFilterChange?.(new Set(filteredRegions.map((r) => r.id)));
		}
	});

	function handleLabelChange(id: string, value: string) {
		onUpdateRegion?.(id, { label: value });
	}

	function handleHintChange(id: string, value: string) {
		onUpdateRegion?.(id, { hint: value || undefined });
	}

	function handleBackExtraChange(id: string, value: string) {
		onUpdateRegion?.(id, { backExtra: value || undefined });
	}
</script>

<div class="flex h-full min-h-0 w-full flex-col bg-card">
	<!-- Title field -->
	<div class="shrink-0 border-b border-border px-3 py-2">
		<span class="mb-1 block text-xs font-medium text-muted-foreground"
			>Title<span class="text-destructive"> *</span></span
		>
		<Input
			value={title}
			placeholder="e.g. Ant anatomy"
			class="h-8 text-sm"
			aria-label="Card title"
			aria-invalid={titleError || undefined}
			oninput={(e: Event) => onTitleChange?.((e.target as HTMLInputElement).value)}
		/>
	</div>

	<!-- Filter -->
	<div class="flex items-center gap-2 border-b border-border px-3 py-2">
		<div class="relative flex-1">
			<Search class="absolute left-2 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
			<Input
				bind:value={filterText}
				placeholder="Filter regions..."
				class="h-7 pl-7 text-sm"
			/>
		</div>
		<span class="shrink-0 rounded-full bg-muted px-2 py-0.5 text-xs font-medium text-muted-foreground">
			{regions.length}
		</span>
	</div>

	<!-- Region list -->
	<ScrollArea class="min-h-0 flex-1">
		<div class="space-y-2 p-3">
			{#if filteredRegions.length === 0}
				{#if regions.length === 0}
					<!-- No regions at all -->
					<div class="flex flex-col items-center justify-center py-8 text-center">
						<div class="mb-2 rounded-full bg-muted p-3">
							<SquareDashed class="h-6 w-6 text-muted-foreground" />
						</div>
						<p class="text-sm font-medium text-foreground">No regions yet</p>
						<p class="mt-1 text-xs text-muted-foreground">
							Draw regions on the image to create flashcards
						</p>
					</div>
				{:else}
					<!-- No results matching filter -->
					<div class="py-8 text-center">
						<p class="text-sm text-muted-foreground">No regions match "{filterText}"</p>
					</div>
				{/if}
			{:else}
				{#each filteredRegions as region, i (region.id)}
					<LabelPanelItem
						{region}
						index={i + 1}
						isSelected={region.id === selectedRegionId}
						focusLabel={region.id === focusLabelRegionId}
						labelError={regionLabelErrors?.has(region.id) ?? false}
						onSelect={() => onSelectRegion?.(region.id)}
						onLabelChange={(value) => handleLabelChange(region.id, value)}
						onHintChange={(value) => handleHintChange(region.id, value)}
						onBackExtraChange={(value) => handleBackExtraChange(region.id, value)}
						onDelete={() => onDeleteRegion?.(region.id)}
					/>
				{/each}
			{/if}
		</div>
	</ScrollArea>
</div>
