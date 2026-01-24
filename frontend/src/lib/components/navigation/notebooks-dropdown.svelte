<script lang="ts">
	import { goto } from '$app/navigation';
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { ChevronDownIcon, BookOpenIcon, PlusIcon } from '@lucide/svelte';
	import type { Notebook } from '$lib/types';

	interface Props {
		notebooks: Notebook[];
		current?: Notebook;
		isInNotebook: boolean;
	}

	let { notebooks, current, isInNotebook }: Props = $props();

	function handleSelect(notebook: Notebook) {
		goto(`/notebooks/${notebook.id}`);
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button
				{...props}
				variant="ghost"
				size="sm"
				class={cn(
					'gap-2',
					isInNotebook && 'bg-slate-100 dark:bg-slate-800'
				)}
			>
				<BookOpenIcon class="size-4" />
				<span class="hidden sm:inline">Notebooks</span>
				<ChevronDownIcon class="size-4 text-slate-400" />
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>

	<DropdownMenu.Content class="w-64" align="start">
		<DropdownMenu.Label>Your Notebooks</DropdownMenu.Label>
		<DropdownMenu.Separator />

		{#each notebooks as notebook (notebook.id)}
			<DropdownMenu.Item
				class={cn(
					'gap-3 cursor-pointer',
					current?.id === notebook.id && 'bg-sky-50 dark:bg-sky-900/20'
				)}
				onclick={() => handleSelect(notebook)}
			>
				<span class="text-lg">{notebook.emoji}</span>
				<div class="flex-1 min-w-0">
					<div class="font-medium truncate">{notebook.name}</div>
					<div class="text-xs text-slate-500 dark:text-slate-400">
						{notebook.totalCards} cards
					</div>
				</div>
				{#if notebook.dueCount > 0}
					<span class="text-xs font-medium text-sky-600 dark:text-sky-400">
						{notebook.dueCount} due
					</span>
				{/if}
			</DropdownMenu.Item>
		{/each}

		{#if notebooks.length === 0}
			<DropdownMenu.Item disabled class="text-slate-500">
				No notebooks yet
			</DropdownMenu.Item>
		{/if}

		<DropdownMenu.Separator />

		<DropdownMenu.Item class="gap-2 cursor-pointer">
			<PlusIcon class="size-4" />
			<span>Create notebook</span>
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
