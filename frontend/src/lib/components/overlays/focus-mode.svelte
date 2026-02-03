<script lang="ts">
	import { onMount } from 'svelte';
	import { cn } from '$lib/utils';
	import { generateUUID } from '$lib/utils/uuid';
	import RatingButtons from '$lib/components/cards/rating-buttons.svelte';
	import { XIcon, EyeIcon } from '@lucide/svelte';
	import type { ReviewScope, StudyCard } from '$lib/types';
	import type { ReviewMode } from '$lib/types/review';
	import { extractCardDisplay, ratingToString } from '$lib/services/reviews';
	import type { ApiReviewResponse, ApiUndoReviewResponse } from '$lib/api/types';
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

	// Undo state
	interface UndoState {
		reviewId: string;
		cardId: string;
		cardBefore: StudyCard;
		insertPosition: number;
		wasRequeued: boolean;
	}
	let lastReview = $state<UndoState | null>(null);
	let undoTimeoutId = $state<ReturnType<typeof setTimeout> | null>(null);
	let isUndoing = $state(false);

	const UNDO_WINDOW_MS = 8000;

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

	function clearUndoState() {
		if (undoTimeoutId) {
			clearTimeout(undoTimeoutId);
			undoTimeoutId = null;
		}
		lastReview = null;
	}

	async function submitReview(card: StudyCard, rating: 1 | 2 | 3 | 4, durationMs: number): Promise<{ response: ApiReviewResponse; reviewId: string } | null> {
		const reviewId = generateUUID();
		const body = {
			id: reviewId,
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
					const response = await res.json();
					return { response, reviewId };
				}
			} catch {
				// retry after delay
			}
			if (attempt === 0) await new Promise((r) => setTimeout(r, 2000));
		}
		return null;
	}

	async function handleUndo() {
		if (!lastReview || isUndoing) return;

		isUndoing = true;
		const reviewToUndo = lastReview;
		clearUndoState();

		try {
			const res = await fetch(`/api/reviews/${reviewToUndo.reviewId}`, {
				method: 'DELETE'
			});

			if (res.ok) {
				const undoResp: ApiUndoReviewResponse = await res.json();

				// Restore queue state
				if (reviewToUndo.wasRequeued) {
					// Remove the requeued card from end of queue
					cardQueue = cardQueue.slice(0, -1);
					totalCards--;
				}

				// Decrement reviewed count
				reviewedCount = Math.max(0, reviewedCount - 1);

				// Restore card to original position
				if (mode === 'scheduled' && undoResp.card) {
					// Re-insert the restored card at its original position
					const restoredCard: StudyCard = {
						...reviewToUndo.cardBefore,
						state: undoResp.card.state,
						due: undoResp.card.due,
						stability: undoResp.card.stability,
						difficulty: undoResp.card.difficulty,
						reps: undoResp.card.reps,
						lapses: undoResp.card.lapses
					};
					cardQueue = [
						...cardQueue.slice(0, reviewToUndo.insertPosition),
						restoredCard,
						...cardQueue.slice(reviewToUndo.insertPosition)
					];
				} else {
					// Practice mode: just re-insert the original card
					cardQueue = [
						...cardQueue.slice(0, reviewToUndo.insertPosition),
						reviewToUndo.cardBefore,
						...cardQueue.slice(reviewToUndo.insertPosition)
					];
				}

				// Move back to the restored card
				currentIndex = reviewToUndo.insertPosition;
				sessionComplete = false;
				isRevealed = false;
				reviewStartTime = Date.now();

				toast.success('Review undone');
			} else {
				const data = await res.json().catch(() => ({}));
				const message = data.message || 'Failed to undo review';
				toast.error(message);
			}
		} catch {
			toast.error('Failed to undo review');
		} finally {
			isUndoing = false;
		}
	}

	async function handleRate(rating: 1 | 2 | 3 | 4) {
		if (!currentCard || isSubmitting) return;

		const card = currentCard;
		const durationMs = Date.now() - reviewStartTime;
		isSubmitting = true;

		// Clear any existing undo state
		clearUndoState();

		// Capture state before advancing
		const insertPosition = currentIndex;

		// Fire-and-forget with retry
		const resultPromise = submitReview(card, rating, durationMs);

		// Toast on failure for non-scheduled (practice) mode where we don't .then() for re-queue
		if (mode !== 'scheduled') {
			resultPromise.then((result) => {
				if (!result) toast.error('Review failed to save. Please try again.');
			});
		}

		reviewedCount++;

		// Track if card was requeued for undo.
		// Note: For non-final cards, undo state is set before we know if requeue happened.
		// The .then() callback updates lastReview.wasRequeued after the fact. This has a
		// theoretical race if multiple reviews are submitted rapidly, but is acceptable
		// given the 8-second undo window and typical human review pace.
		let wasRequeued = false;

		// In scheduled mode, check if learning card should be re-queued
		if (mode === 'scheduled') {
			resultPromise.then((result) => {
				if (!result) {
					toast.error('Review failed to save. Please try again.');
					return;
				}
				const updatedDue = result.response.card.due;
				if (updatedDue) {
					const dueTime = new Date(updatedDue).getTime();
					const now = Date.now();
					const tenMinutes = 10 * 60 * 1000;
					if (dueTime - now < tenMinutes && dueTime > now) {
						// Re-insert card at end of queue with updated state
						const requeued: StudyCard = {
							...card,
							state: result.response.card.state,
							due: result.response.card.due,
							stability: result.response.card.stability,
							difficulty: result.response.card.difficulty,
							reps: result.response.card.reps,
							lapses: result.response.card.lapses
						};
						cardQueue = [...cardQueue, requeued];
						totalCards++;
						wasRequeued = true;

						// Update undo state with requeue info
						if (lastReview && lastReview.reviewId === result.reviewId) {
							lastReview = { ...lastReview, wasRequeued: true };
						}
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
			const result = await resultPromise;
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

			// Set up undo state after we know if requeue happened
			if (result) {
				lastReview = {
					reviewId: result.reviewId,
					cardId: card.id,
					cardBefore: card,
					insertPosition,
					wasRequeued
				};

				// Show undo toast
				undoTimeoutId = setTimeout(() => {
					lastReview = null;
					undoTimeoutId = null;
				}, UNDO_WINDOW_MS);

				toast('Review saved', {
					duration: UNDO_WINDOW_MS,
					action: {
						label: 'Undo',
						onClick: handleUndo
					}
				});
			}
			return;
		}

		// Set up undo state for non-final cards
		resultPromise.then((result) => {
			if (result) {
				lastReview = {
					reviewId: result.reviewId,
					cardId: card.id,
					cardBefore: card,
					insertPosition,
					wasRequeued: false
				};

				// Show undo toast
				undoTimeoutId = setTimeout(() => {
					lastReview = null;
					undoTimeoutId = null;
				}, UNDO_WINDOW_MS);

				toast('Review saved', {
					duration: UNDO_WINDOW_MS,
					action: {
						label: 'Undo',
						onClick: handleUndo
					}
				});
			}
		});
	}

	function handleClose() {
		clearUndoState();
		onClose?.();
	}

	let dialogEl = $state<HTMLDivElement>();

	onMount(() => {
		dialogEl?.focus();
		reviewStartTime = Date.now();

		function handleKeydown(e: KeyboardEvent) {
			// Z key for undo (case insensitive, no modifiers)
			if ((e.key === 'z' || e.key === 'Z') && !e.ctrlKey && !e.metaKey && !e.altKey) {
				if (lastReview && !isUndoing) {
					e.preventDefault();
					handleUndo();
					return;
				}
			}

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
		return () => {
			window.removeEventListener('keydown', handleKeydown);
			clearUndoState();
		};
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
			Space to flip &bull; 1-4 to rate &bull; Z to undo &bull; Esc to exit
		</div>
	{/if}
</div>
