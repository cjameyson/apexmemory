<script lang="ts">
	import { goto } from '$app/navigation';
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { ChevronDownIcon, BookOpenIcon, PlusIcon, SearchIcon, CheckIcon } from '@lucide/svelte';
	import type { Notebook } from '$lib/types';

	interface Props {
		notebooks: Notebook[];
		current?: Notebook;
		isInNotebook: boolean;
		onCreateNotebook?: () => void;
	}

	let { notebooks, current, isInNotebook, onCreateNotebook }: Props = $props();

	// Search state (only used when > 8 notebooks)
	let searchQuery = $state('');
	let showSearch = $derived(notebooks.length > 8);

	let filteredNotebooks = $derived(
		searchQuery
			? notebooks.filter((n) => n.name.toLowerCase().includes(searchQuery.toLowerCase()))
			: notebooks
	);

	// Sort alphabetically
	let sortedNotebooks = $derived(
		[...filteredNotebooks].sort((a, b) => a.name.localeCompare(b.name))
	);

	function handleSelect(notebook: Notebook) {
		goto(`/notebooks/${notebook.id}`);
	}

	function handleCreateClick() {
		onCreateNotebook?.();
	}

	// Clear search when dropdown closes
	function handleOpenChange(open: boolean) {
		if (!open) {
			searchQuery = '';
		}
	}
</script>

<DropdownMenu.Root onOpenChange={handleOpenChange}>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button
				{...props}
				variant="ghost"
				size="sm"
				class={cn('gap-2', isInNotebook && 'bg-slate-100 dark:bg-slate-800')}
			>
				<BookOpenIcon class="size-4" />
				<span class="hidden sm:inline">Notebooks</span>
				<ChevronDownIcon class="size-4 text-slate-400" />
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>

	<DropdownMenu.Content class="w-64" align="start">
		<DropdownMenu.Label>Your Notebooks</DropdownMenu.Label>

		{#if showSearch}
			<div class="px-2 py-1.5">
				<div class="relative">
					<SearchIcon
						class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-slate-400"
					/>
					<Input
						type="text"
						placeholder="Search notebooks..."
						class="h-8 pl-8 text-sm"
						bind:value={searchQuery}
						onclick={(e) => e.stopPropagation()}
						onkeydown={(e) => e.stopPropagation()}
					/>
				</div>
			</div>
		{/if}

		<DropdownMenu.Separator />

		<div class="max-h-64 overflow-y-auto">
			{#each sortedNotebooks as notebook (notebook.id)}
				<DropdownMenu.Item
					class={cn(
						'cursor-pointer gap-3',
						current?.id === notebook.id && 'bg-primary/10 dark:bg-primary/20'
					)}
					onclick={() => handleSelect(notebook)}
				>
					<span class="text-lg">{notebook.emoji}</span>
					<div class="min-w-0 flex-1">
						<div class="truncate font-medium">{notebook.name}</div>
						<div class="text-xs text-slate-500 dark:text-slate-400">
							{notebook.totalCards} cards{#if notebook.retention > 0}&nbsp;- {Math.round(
									notebook.retention * 100
								)}%{/if}
						</div>
					</div>
					{#if notebook.dueCount > 0}
						<span class="text-primary text-xs font-medium">
							{notebook.dueCount} due
						</span>
					{:else if notebook.dueCount === 0}
						<CheckIcon class="size-4" />
					{/if}
				</DropdownMenu.Item>
			{/each}

			{#if sortedNotebooks.length === 0}
				<div class="px-4 py-6 text-center">
					<BookOpenIcon class="mx-auto mb-2 size-8 text-slate-300 dark:text-slate-600" />
					{#if searchQuery}
						<p class="text-sm text-slate-500">No notebooks match "{searchQuery}"</p>
					{:else}
						<p class="mb-2 text-sm text-slate-500">No notebooks yet</p>
						<Button size="sm" variant="outline" onclick={handleCreateClick}>
							Create your first notebook
						</Button>
					{/if}
				</div>
			{/if}
		</div>

		<DropdownMenu.Separator />
		<DropdownMenu.Item class="cursor-pointer gap-2" onclick={handleCreateClick}>
			<PlusIcon class="size-4" />
			<span>Create notebook</span>
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
