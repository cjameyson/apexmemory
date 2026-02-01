<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import {
		SearchIcon,
		LayersIcon,
		TypeIcon,
		ImageIcon,
		TableIcon,
		LayoutGridIcon
	} from '@lucide/svelte';
	import type { FactTypeFilter } from '$lib/types/fact';

	const VALID_TYPES = new Set<FactTypeFilter>(['all', 'basic', 'cloze', 'image_occlusion']);
	const VALID_SORT_VALUES = new Set([
		'-updated', 'updated', '-created', 'created'
	]);

	// All state derived from URL — single source of truth
	let rawType = $derived($page.url.searchParams.get('type') ?? '');
	let currentType = $derived<FactTypeFilter>(
		VALID_TYPES.has(rawType as FactTypeFilter) ? (rawType as FactTypeFilter) : 'all'
	);
	let sortValue = $derived($page.url.searchParams.get('sort') || '-updated');
	let currentView = $derived<'table' | 'grid'>(
		$page.url.searchParams.get('view') === 'grid' ? 'grid' : 'table'
	);
	let currentSearch = $derived($page.url.searchParams.get('q') || '');

	// Local search input tracks the text field during typing,
	// synced back to URL state after debounce
	let searchInput = $state(currentSearch);
	let debounceTimer: ReturnType<typeof setTimeout> | undefined;

	// Sync input when URL changes externally (back/forward, link click)
	$effect(() => {
		searchInput = currentSearch;
	});

	// Cleanup debounce timer on destroy
	$effect(() => {
		return () => clearTimeout(debounceTimer);
	});

	function updateUrl(params: Record<string, string>) {
		const url = new URL($page.url);
		for (const [key, value] of Object.entries(params)) {
			if (value) {
				url.searchParams.set(key, value);
			} else {
				url.searchParams.delete(key);
			}
		}
		// Reset to page 1 on filter changes (except view toggle)
		if (!('view' in params)) {
			url.searchParams.delete('page');
		}
		goto(url.toString(), { replaceState: true, keepFocus: true });
	}

	function onSearchInput(e: Event) {
		const value = (e.target as HTMLInputElement).value;
		searchInput = value;
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => {
			updateUrl({ q: value });
		}, 300);
	}

	function setType(type: FactTypeFilter) {
		updateUrl({ type: type === 'all' ? '' : type });
	}

	function onSortChange(e: Event) {
		const value = (e.target as HTMLSelectElement).value;
		if (!VALID_SORT_VALUES.has(value)) return;
		updateUrl({ sort: value === '-updated' ? '' : value });
	}

	function setView(view: 'table' | 'grid') {
		updateUrl({ view: view === 'table' ? '' : view });
	}

	const typeOptions: { value: FactTypeFilter; label: string; icon: typeof LayersIcon }[] = [
		{ value: 'all', label: 'All', icon: LayersIcon },
		{ value: 'basic', label: 'Basic', icon: LayersIcon },
		{ value: 'cloze', label: 'Cloze', icon: TypeIcon },
		{ value: 'image_occlusion', label: 'Image', icon: ImageIcon }
	];

	const sortOptions = [
		{ value: '-updated', label: 'Recently Updated' },
		{ value: 'updated', label: 'Oldest Updated' },
		{ value: '-created', label: 'Recently Created' },
		{ value: 'created', label: 'Oldest Created' }
	];
</script>

<div class="flex flex-wrap items-center gap-3 border-b border-border px-6 py-3">
	<!-- Search Input -->
	<div class="relative flex-1" style="max-width: 20rem;">
		<SearchIcon class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
		<input
			type="text"
			placeholder="Search facts..."
			value={searchInput}
			oninput={onSearchInput}
			class="w-full rounded-lg border border-border bg-background py-1.5 pl-9 pr-3 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
		/>
	</div>

	<!-- Type Filter Tabs -->
	<div class="flex rounded-lg bg-muted p-0.5">
		{#each typeOptions as opt}
			<button
				onclick={() => setType(opt.value)}
				class="flex items-center gap-1.5 rounded-md px-3 py-1 text-xs font-medium transition-colors {currentType === opt.value
					? 'bg-background text-foreground shadow-sm'
					: 'text-muted-foreground hover:text-foreground'}"
			>
				<opt.icon class="h-3.5 w-3.5" />
				{opt.label}
			</button>
		{/each}
	</div>

	<!-- Spacer pushes sort + view toggle right -->
	<div class="ml-auto flex items-center gap-2">
		<!-- Sort Dropdown -->
		<select
			value={sortValue}
			onchange={onSortChange}
			class="rounded-lg border border-border bg-background px-3 py-1.5 text-xs text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
		>
			{#each sortOptions as opt}
				<option value={opt.value}>{opt.label}</option>
			{/each}
		</select>

		<!-- View Mode Toggle (hidden until grid view is implemented) -->
		<!-- <div class="flex rounded-lg bg-muted p-0.5">
			<button
				onclick={() => setView('table')}
				class="rounded-md p-1.5 transition-colors {currentView === 'table'
					? 'bg-background text-foreground shadow-sm'
					: 'text-muted-foreground hover:text-foreground'}"
				aria-label="Table view"
			>
				<TableIcon class="h-4 w-4" />
			</button>
			<button
				onclick={() => setView('grid')}
				class="rounded-md p-1.5 transition-colors {currentView === 'grid'
					? 'bg-background text-foreground shadow-sm'
					: 'text-muted-foreground hover:text-foreground'}"
				aria-label="Grid view"
			>
				<LayoutGridIcon class="h-4 w-4" />
			</button>
		</div> -->
	</div>
</div>
