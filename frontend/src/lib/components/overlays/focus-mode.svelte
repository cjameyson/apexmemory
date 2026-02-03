<script lang="ts">
	import { onMount } from 'svelte';
	import { cn } from '$lib/utils';
	import RatingButtons from '$lib/components/cards/rating-buttons.svelte';
	import { XIcon, EyeIcon } from '@lucide/svelte';
	import type { ReviewScope, StudyCard } from '$lib/types';
	import type { ReviewMode } from '$lib/types/review';
	import { extractCardDisplay, ratingToString } from '$lib/services/reviews';
	import type { ApiReviewResponse } from '$lib/api/types';
	import { toast } from 'svelte-sonner';

	interface Props {
		cards: StudyCard[];
		mode: ReviewMode;
		scope: ReviewScope;
		initialIndex?: number;
		onProgressChange?: (index: number) => void;
		onClose?: () => void;
		class?: string;
	}

	let { cards: initialCards, mode, scope, initialIndex = 0, onProgressChange, onClose, class: className }: Props = $props();

	// Mutable queue -- learning cards may be re-inserted
	let cardQueue = $state<StudyCard[]>([...initialCards]);
	let currentIndex = $state(initialIndex);
	let isRevealed = $state(false);
	let sessionComplete = $state(initialCards.length === 0);
	let reviewedCount = $state(0);
	let totalCards = $state(initialCards.length);
	let reviewStartTime = $state(Date.now());
	let isSubmitting = $state(false);

	let currentCard = $derived(cardQueue[currentIndex] as StudyCard | undefined);
	let display = $derived(currentCard ? extractCardDisplay(currentCard) : null);
	let progress = $derived(totalCards > 0 ? (reviewedCount / totalCards) * 100 : 0);

	let scopeTitle = $derived.by(() => {
		if (scope.type === 'all') return 'All Due Cards';
		if (scope.type === 'notebook') return scope.notebook.name;
		if (scope.type === 'source') return scope.source.name;
		return 'Review';
	});

	function reveal() {
		isRevealed = true;
	}

	async function submitReview(card: StudyCard, rating: 1 | 2 | 3 | 4, durationMs: number): Promise<ApiReviewResponse | null> {
		const body = {
			id: crypto.randomUUID(),
			card_id: card.id,
			rating: ratingToString(rating),
			duration_ms: durationMs,
			mode
		};

		for (let attempt = 0; attempt < 2; attempt++) {
			try {
				const res = await fetch('/api/reviews', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(body)
				});
				if (res.ok) {
					return await res.json();
				}
			} catch {
				// retry after delay
			}
			if (attempt === 0) await new Promise((r) => setTimeout(r, 2000));
		}
		return null;
	}

	async function handleRate(rating: 1 | 2 | 3 | 4) {
		if (!currentCard || isSubmitting) return;

		const card = currentCard;
		const durationMs = Date.now() - reviewStartTime;
		isSubmitting = true;

		// Fire-and-forget with retry
		const responsePromise = submitReview(card, rating, durationMs);

		// Toast on failure for non-scheduled (practice) mode where we don't .then() for re-queue
		if (mode !== 'scheduled') {
			responsePromise.then((resp) => {
				if (!resp) toast.error('Review failed to save. Please try again.');
			});
		}

		reviewedCount++;

		// In scheduled mode, check if learning card should be re-queued
		if (mode === 'scheduled') {
			responsePromise.then((resp) => {
				if (!resp) {
					toast.error('Review failed to save. Please try again.');
					return;
				}
				const updatedDue = resp.card.due;
				if (updatedDue) {
					const dueTime = new Date(updatedDue).getTime();
					const now = Date.now();
					const tenMinutes = 10 * 60 * 1000;
					if (dueTime - now < tenMinutes && dueTime > now) {
						// Re-insert card at end of queue with updated state
						const requeued: StudyCard = {
							...card,
							state: resp.card.state,
							due: resp.card.due,
							stability: resp.card.stability,
							difficulty: resp.card.difficulty,
							reps: resp.card.reps,
							lapses: resp.card.lapses
						};
						cardQueue = [...cardQueue, requeued];
						totalCards++;
					}
				}
			});
		}

		// Advance to next card
		if (currentIndex < cardQueue.length - 1) {
			currentIndex++;
			isRevealed = false;
			reviewStartTime = Date.now();
			isSubmitting = false;
			onProgressChange?.(currentIndex);
		} else {
			// Wait for response to check re-queue before completing
			await responsePromise;
			isSubmitting = false;
			if (currentIndex < cardQueue.length - 1) {
				// Card was re-queued
				currentIndex++;
				isRevealed = false;
				reviewStartTime = Date.now();
				onProgressChange?.(currentIndex);
			} else {
				sessionComplete = true;
			}
		}
	}

	function handleClose() {
		onClose?.();
	}

	let dialogEl = $state<HTMLDivElement>();

	onMount(() => {
		dialogEl?.focus();
		reviewStartTime = Date.now();

		function handleKeydown(e: KeyboardEvent) {
			if (sessionComplete || isSubmitting) return;

			if (e.key === ' ' && !isRevealed) {
				e.preventDefault();
				reveal();
			}

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
	bind:this={dialogEl}
	tabindex="-1"
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
			{#if mode === 'practice'}
				<span class="px-2 py-0.5 text-xs font-medium bg-amber-500/20 text-amber-300 rounded-full">Practice Mode</span>
			{/if}
			{#if scope.type === 'notebook' || scope.type === 'source'}
				<span class="text-xl">{scope.notebook.emoji}</span>
			{/if}
			<span class="font-medium">{scopeTitle}</span>
		</div>

		<div class="flex items-center gap-3 text-white/70">
			<span class="text-sm">{currentIndex + 1}/{cardQueue.length}</span>
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
			<div class="text-center text-white">
				<div class="text-6xl mb-6">&#127881;</div>
				<h2 class="text-3xl font-bold mb-2">All done!</h2>
				<p class="text-white/70 mb-8">
					You reviewed {reviewedCount} card{reviewedCount === 1 ? '' : 's'}.
				</p>
				<button
					type="button"
					onclick={handleClose}
					class="px-6 py-3 bg-white/10 hover:bg-white/20 rounded-xl font-medium transition-colors"
				>
					Close
				</button>
			</div>
		{:else if currentCard && display}
			<div class="w-full max-w-2xl">
				<div
					class="bg-white/10 backdrop-blur-sm rounded-3xl p-8 mb-8 min-h-64 flex flex-col justify-center"
				>
					<div class={cn('text-center', isRevealed && 'opacity-60')}>
						<p class="text-2xl font-medium text-white">
							{display.front}
						</p>
					</div>

					{#if isRevealed}
						<div class="my-6 border-t border-white/20"></div>
						<div class="text-center">
							<p class="text-xl text-white/90">
								{display.back}
							</p>
						</div>
					{:else}
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

				{#if isRevealed}
					<RatingButtons
						onRate={handleRate}
						intervals={currentCard.intervals}
						disabled={isSubmitting}
					/>
				{/if}
			</div>
		{:else}
			<div class="text-center text-white">
				<div class="text-6xl mb-6">&#128218;</div>
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

	{#if !sessionComplete && currentCard}
		<div class="text-center py-4 text-white/40 text-sm">
			Space to flip &bull; 1-4 to rate &bull; Esc to exit
		</div>
	{/if}
</div>
