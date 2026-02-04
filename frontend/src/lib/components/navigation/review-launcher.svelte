<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { ZapIcon, ChevronDownIcon, LoaderCircleIcon } from '@lucide/svelte';
	import type { Notebook, ReviewScope, StudyCard } from '$lib/types';
	import { fetchStudyCards } from '$lib/services/reviews';
	import { studyCounts } from '$lib/stores/study-counts.svelte';

	interface Props {
		notebooks: Notebook[];
		currentNotebook?: Notebook;
		onStartReview?: (scope: ReviewScope, cards: StudyCard[]) => void;
		class?: string;
	}

	let { notebooks, currentNotebook, onStartReview, class: className }: Props = $props();

	let open = $state(false);
	let loading = $state(false);

	// Use store for totalDue (stays in sync after reviews)
	let totalDue = $derived(studyCounts.getTotalDue());
	// Filter notebooks with due cards using store
	let notebooksWithDue = $derived(
		notebooks.filter((nb) => studyCounts.getDueCount(nb.id) > 0)
	);

	async function startReview(notebook?: Notebook) {
		open = false;
		loading = true;

		try {
			const cards = await fetchStudyCards(notebook?.id);

			const scope: ReviewScope = notebook
				? { type: 'notebook', notebook, mode: 'scheduled' }
				: { type: 'all', mode: 'scheduled' };

			onStartReview?.(scope, cards);
		} finally {
			loading = false;
		}
	}
</script>

{#if loading}
	<Button variant="default" size="sm" class={cn('gap-2', className)} disabled>
		<LoaderCircleIcon class="size-4 animate-spin" />
		<span class="hidden sm:inline">Loading...</span>
	</Button>
{:else}
	<DropdownMenu.Root bind:open>
		<DropdownMenu.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					variant="default"
					size="sm"
					class={cn('gap-2', className)}
				>
					<ZapIcon class="size-4" />
					<span class="hidden sm:inline">Review</span>
					{#if totalDue > 0}
						<span class="hidden sm:inline text-primary-foreground/80">({totalDue})</span>
					{/if}
					<ChevronDownIcon class="size-4 text-primary-foreground/70" />
				</Button>
			{/snippet}
		</DropdownMenu.Trigger>

		<DropdownMenu.Content class="w-64" align="end">
			<!-- Scheduled Review -->
			<DropdownMenu.Item
				class="gap-3 cursor-pointer"
				onclick={() => startReview()}
			>
				<div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
					<ZapIcon class="size-4 text-primary-foreground" />
				</div>
				<div class="flex-1">
					<div class="font-medium">Review All</div>
					<div class="text-xs text-muted-foreground">
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
							currentNotebook?.id === notebook.id && 'bg-primary/10 dark:bg-primary/20'
						)}
						onclick={() => startReview(notebook)}
					>
						<span class="text-lg">{notebook.emoji}</span>
						<div class="flex-1 min-w-0">
							<div class="font-medium truncate">{notebook.name}</div>
						</div>
						<span class="text-sm font-medium text-primary">
							{studyCounts.getDueCount(notebook.id)}
						</span>
					</DropdownMenu.Item>
				{/each}
			{/if}

			{#if totalDue === 0}
				<div class="px-3 py-4 text-center text-sm text-muted-foreground">
					No cards available
				</div>
			{/if}
		</DropdownMenu.Content>
	</DropdownMenu.Root>
{/if}
