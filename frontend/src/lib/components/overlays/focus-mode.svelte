<script lang="ts">
	import { onMount } from 'svelte';
	import { cn } from '$lib/utils';
	import RatingButtons from '$lib/components/cards/rating-buttons.svelte';
	import { XIcon, EyeIcon } from '@lucide/svelte';
	import type { Card, ReviewScope } from '$lib/types';
	import { getAllDueCards, getDueCardsForNotebook, getCardsForSource } from '$lib/mocks';

	interface Props {
		scope: ReviewScope;
		initialIndex?: number;
		onProgressChange?: (index: number) => void;
		onClose?: () => void;
		class?: string;
	}

	let { scope, initialIndex = 0, onProgressChange, onClose, class: className }: Props = $props();

	// Get cards based on scope
	let allCards = $derived.by(() => {
		if (scope.type === 'all') {
			return getAllDueCards();
		} else if (scope.type === 'notebook') {
			return getDueCardsForNotebook(scope.notebook.id);
		} else if (scope.type === 'source') {
			return getCardsForSource(scope.notebook.id, scope.source.id).filter(c => c.due);
		}
		return [];
	});

	// Session state - initialize from prop for session restore
	let currentIndex = $state(initialIndex);
	let isRevealed = $state(false);
	let sessionComplete = $state(initialIndex >= allCards.length && allCards.length > 0);

	let currentCard = $derived(allCards[currentIndex] as Card | undefined);
	let progress = $derived(allCards.length > 0 ? ((currentIndex) / allCards.length) * 100 : 0);

	// Title based on scope
	let scopeTitle = $derived.by(() => {
		if (scope.type === 'all') return 'All Due Cards';
		if (scope.type === 'notebook') return scope.notebook.name;
		if (scope.type === 'source') return scope.source.name;
		return 'Review';
	});

	function reveal() {
		isRevealed = true;
	}

	function handleRate(_rating: 1 | 2 | 3 | 4) {
		// TODO: Send rating to backend via API

		// Move to next card
		if (currentIndex < allCards.length - 1) {
			currentIndex++;
			isRevealed = false;
			// Update URL state for session restore (without creating history entry)
			onProgressChange?.(currentIndex);
		} else {
			sessionComplete = true;
		}
	}

	function handleClose() {
		onClose?.();
	}

	// Keyboard shortcuts
	onMount(() => {
		function handleKeydown(e: KeyboardEvent) {
			if (sessionComplete) return;

			// Space to reveal
			if (e.key === ' ' && !isRevealed) {
				e.preventDefault();
				reveal();
			}

			// 1-4 to rate (only when revealed)
			if (isRevealed && ['1', '2', '3', '4'].includes(e.key)) {
				e.preventDefault();
				handleRate(parseInt(e.key) as 1 | 2 | 3 | 4);
			}
		}

		window.addEventListener('keydown', handleKeydown);
		return () => window.removeEventListener('keydown', handleKeydown);
	});
</script>

<div
	class={cn(
		'fixed inset-0 z-50 bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex flex-col',
		className
	)}
	role="dialog"
	aria-modal="true"
	aria-label="Focus mode review session"
>
	<!-- Header -->
	<div class="flex items-center justify-between px-6 py-4">
		<button
			type="button"
			onclick={handleClose}
			class="p-2 rounded-lg text-white/70 hover:text-white hover:bg-white/10 transition-colors"
			aria-label="Exit focus mode"
		>
			<XIcon class="size-6" />
		</button>

		<div class="flex items-center gap-3 text-white">
			{#if scope.type === 'notebook' || scope.type === 'source'}
				<span class="text-xl">{scope.notebook.emoji}</span>
			{/if}
			<span class="font-medium">{scopeTitle}</span>
		</div>

		<div class="flex items-center gap-3 text-white/70">
			<span class="text-sm">{currentIndex + 1}/{allCards.length}</span>
			<div class="w-24 h-2 bg-white/20 rounded-full overflow-hidden">
				<div
					class="h-full bg-white transition-all duration-300"
					style="width: {progress}%"
				></div>
			</div>
		</div>
	</div>

	<!-- Main content -->
	<div class="flex-1 flex items-center justify-center px-6 pb-6">
		{#if sessionComplete}
			<!-- Completion state -->
			<div class="text-center text-white">
				<div class="text-6xl mb-6">🎉</div>
				<h2 class="text-3xl font-bold mb-2">All done!</h2>
				<p class="text-white/70 mb-8">
					You reviewed {allCards.length} cards.
				</p>
				<button
					type="button"
					onclick={handleClose}
					class="px-6 py-3 bg-white/10 hover:bg-white/20 rounded-xl font-medium transition-colors"
				>
					Close
				</button>
			</div>
		{:else if currentCard}
			<!-- Card display -->
			<div class="w-full max-w-2xl">
				<div
					class="bg-white/10 backdrop-blur-sm rounded-3xl p-8 mb-8 min-h-64 flex flex-col justify-center"
				>
					<!-- Front (question) -->
					<div class={cn('text-center', isRevealed && 'opacity-60')}>
						<p class="text-2xl font-medium text-white">
							{currentCard.front}
						</p>
					</div>

					{#if isRevealed}
						<!-- Separator -->
						<div class="my-6 border-t border-white/20"></div>

						<!-- Back (answer) -->
						<div class="text-center">
							<p class="text-xl text-white/90">
								{currentCard.back}
							</p>
						</div>
					{:else}
						<!-- Reveal prompt -->
						<button
							type="button"
							onclick={reveal}
							class="mt-8 flex items-center justify-center gap-2 text-white/50 hover:text-white/70 transition-colors"
						>
							<EyeIcon class="size-5" />
							<span>Tap to reveal</span>
						</button>
					{/if}
				</div>

				<!-- Rating buttons (shown after reveal) -->
				{#if isRevealed}
					<RatingButtons onRate={handleRate} />
				{/if}
			</div>
		{:else}
			<!-- No cards state -->
			<div class="text-center text-white">
				<div class="text-6xl mb-6">📚</div>
				<h2 class="text-2xl font-bold mb-2">No cards due</h2>
				<p class="text-white/70 mb-8">
					You're all caught up! Check back later.
				</p>
				<button
					type="button"
					onclick={handleClose}
					class="px-6 py-3 bg-white/10 hover:bg-white/20 rounded-xl font-medium transition-colors"
				>
					Close
				</button>
			</div>
		{/if}
	</div>

	<!-- Footer hint -->
	{#if !sessionComplete && currentCard}
		<div class="text-center py-4 text-white/40 text-sm">
			Space to flip • 1-4 to rate • Esc to exit
		</div>
	{/if}
</div>
