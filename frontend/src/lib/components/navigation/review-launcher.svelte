<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { ZapIcon, ChevronDownIcon, RepeatIcon, LoaderCircleIcon } from '@lucide/svelte';
	import type { Notebook, ReviewScope, StudyCard } from '$lib/types';
	import type { ReviewMode } from '$lib/types/review';
	import { fetchStudyCards, fetchPracticeCards } from '$lib/services/reviews';

	interface Props {
		notebooks: Notebook[];
		currentNotebook?: Notebook;
		onStartReview?: (scope: ReviewScope, cards: StudyCard[]) => void;
		class?: string;
	}

	let { notebooks, currentNotebook, onStartReview, class: className }: Props = $props();

	let open = $state(false);
	let loading = $state(false);

	let totalDue = $derived(notebooks.reduce((sum, nb) => sum + nb.dueCount, 0));
	let notebooksWithDue = $derived(notebooks.filter((nb) => nb.dueCount > 0));

	async function startReview(mode: ReviewMode, notebook?: Notebook) {
		open = false;
		loading = true;

		try {
			const cards = mode === 'scheduled'
				? await fetchStudyCards(notebook?.id)
				: await fetchPracticeCards(notebook?.id);

			const scope: ReviewScope = notebook
				? { type: 'notebook', notebook, mode }
				: { type: 'all', mode };

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
				onclick={() => startReview('scheduled')}
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
						onclick={() => startReview('scheduled', notebook)}
					>
						<span class="text-lg">{notebook.emoji}</span>
						<div class="flex-1 min-w-0">
							<div class="font-medium truncate">{notebook.name}</div>
						</div>
						<span class="text-sm font-medium text-primary">
							{notebook.dueCount}
						</span>
					</DropdownMenu.Item>
				{/each}
			{/if}

			<!-- Practice Section -->
			<DropdownMenu.Separator />
			<DropdownMenu.Label>Practice</DropdownMenu.Label>

			<DropdownMenu.Item
				class="gap-3 cursor-pointer"
				onclick={() => startReview('practice')}
			>
				<div class="w-8 h-8 bg-amber-500/20 rounded-lg flex items-center justify-center">
					<RepeatIcon class="size-4 text-amber-500" />
				</div>
				<div class="flex-1">
					<div class="font-medium">Practice All</div>
					<div class="text-xs text-muted-foreground">
						Review without affecting schedule
					</div>
				</div>
			</DropdownMenu.Item>

			{#each notebooks as notebook (notebook.id)}
				{#if notebook.totalCards > 0}
					<DropdownMenu.Item
						class="gap-3 cursor-pointer"
						onclick={() => startReview('practice', notebook)}
					>
						<span class="text-lg">{notebook.emoji}</span>
						<div class="flex-1 min-w-0">
							<div class="font-medium truncate">{notebook.name}</div>
						</div>
						<span class="text-xs text-muted-foreground">Practice</span>
					</DropdownMenu.Item>
				{/if}
			{/each}

			{#if totalDue === 0 && notebooks.every((nb) => nb.totalCards === 0)}
				<div class="px-3 py-4 text-center text-sm text-muted-foreground">
					No cards available
				</div>
			{/if}
		</DropdownMenu.Content>
	</DropdownMenu.Root>
{/if}
