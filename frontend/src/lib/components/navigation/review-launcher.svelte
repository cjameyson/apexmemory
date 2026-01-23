<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { ZapIcon, ChevronDownIcon } from '@lucide/svelte';
	import type { Notebook, ReviewScope } from '$lib/types';

	interface Props {
		notebooks: Notebook[];
		currentNotebook?: Notebook;
		onStartReview?: (scope: ReviewScope) => void;
		class?: string;
	}

	let { notebooks, currentNotebook, onStartReview, class: className }: Props = $props();

	// Calculate total due across all notebooks
	let totalDue = $derived(notebooks.reduce((sum, nb) => sum + nb.dueCount, 0));

	// Notebooks with due cards
	let notebooksWithDue = $derived(notebooks.filter((nb) => nb.dueCount > 0));

	function startReviewAll() {
		onStartReview?.({ type: 'all' });
	}

	function startReviewNotebook(notebook: Notebook) {
		onStartReview?.({ type: 'notebook', notebook });
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button
				{...props}
				variant="default"
				size="sm"
				class={cn('gap-2 bg-sky-500 hover:bg-sky-600 text-white', className)}
			>
				<ZapIcon class="size-4" />
				<span class="hidden sm:inline">Review</span>
				{#if totalDue > 0}
					<span class="hidden sm:inline text-white/80">({totalDue})</span>
				{/if}
				<ChevronDownIcon class="size-4 text-white/70" />
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>

	<DropdownMenu.Content class="w-64" align="end">
		<!-- Review All option -->
		<DropdownMenu.Item
			class="gap-3 cursor-pointer"
			onclick={startReviewAll}
		>
			<div class="w-8 h-8 bg-gradient-to-br from-sky-500 to-cyan-600 rounded-lg flex items-center justify-center">
				<ZapIcon class="size-4 text-white" />
			</div>
			<div class="flex-1">
				<div class="font-medium">Review All</div>
				<div class="text-xs text-slate-500 dark:text-slate-400">
					{totalDue} cards due
				</div>
			</div>
		</DropdownMenu.Item>

		{#if notebooksWithDue.length > 0}
			<DropdownMenu.Separator />
			<DropdownMenu.Label>By Notebook</DropdownMenu.Label>

			{#each notebooksWithDue as notebook (notebook.id)}
				<DropdownMenu.Item
					class={cn(
						'gap-3 cursor-pointer',
						currentNotebook?.id === notebook.id && 'bg-sky-50 dark:bg-sky-900/20'
					)}
					onclick={() => startReviewNotebook(notebook)}
				>
					<span class="text-lg">{notebook.emoji}</span>
					<div class="flex-1 min-w-0">
						<div class="font-medium truncate">{notebook.name}</div>
					</div>
					<span class="text-sm font-medium text-sky-600 dark:text-sky-400">
						{notebook.dueCount}
					</span>
				</DropdownMenu.Item>
			{/each}
		{/if}

		{#if totalDue === 0}
			<div class="px-3 py-4 text-center text-sm text-slate-500 dark:text-slate-400">
				No cards due for review
			</div>
		{/if}
	</DropdownMenu.Content>
</DropdownMenu.Root>
