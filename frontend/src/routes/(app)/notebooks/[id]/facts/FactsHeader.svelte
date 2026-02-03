<script lang="ts">
	import type { FactStats, FactDetail } from '$lib/types/fact';
	import type { ApiFactDetail } from '$lib/api/types';
	import type { FactFormData } from '$lib/components/facts/create-fact-modal.svelte';
	import type { Notebook } from '$lib/types';
	import { toFactDetail } from '$lib/services/facts';
	import { fetchStudyCards, fetchPracticeCards } from '$lib/services/reviews';
	import { PlusIcon, ZapIcon, RepeatIcon, LoaderCircleIcon } from '@lucide/svelte';
	import { invalidateAll, pushState } from '$app/navigation';
	import QuickStats from './QuickStats.svelte';
	import CreateFactModal from '$lib/components/facts/create-fact-modal.svelte';

	let { stats, notebookId, notebook }: { stats: FactStats; notebookId: string; notebook: Notebook } = $props();

	let modalOpen = $state(false);
	let editingFact = $state<FactDetail | null>(null);
	let fetchError = $state<string | null>(null);
	let editRequestId = 0;
	let isLoadingReview = $state(false);
	let isLoadingPractice = $state(false);

	async function startReview() {
		isLoadingReview = true;
		try {
			const cards = await fetchStudyCards(notebookId);
			pushState('', {
				focusMode: {
					type: 'notebook',
					mode: 'scheduled',
					notebookId,
					notebookName: notebook.name,
					notebookEmoji: notebook.emoji,
					cards
				}
			});
		} finally {
			isLoadingReview = false;
		}
	}

	async function startPractice() {
		isLoadingPractice = true;
		try {
			const cards = await fetchPracticeCards(notebookId);
			pushState('', {
				focusMode: {
					type: 'notebook',
					mode: 'practice',
					notebookId,
					notebookName: notebook.name,
					notebookEmoji: notebook.emoji,
					cards
				}
			});
		} finally {
			isLoadingPractice = false;
		}
	}

	export async function openEdit(factId: string) {
		fetchError = null;
		const reqId = ++editRequestId;
		let res: Response;
		try {
			res = await fetch(`/api/notebooks/${notebookId}/facts/${factId}`);
		} catch {
			fetchError = 'Network error — could not load fact.';
			return;
		}
		if (reqId !== editRequestId) return; // stale
		if (!res.ok) {
			fetchError = 'Failed to load fact for editing.';
			return;
		}
		const raw: ApiFactDetail = await res.json();
		editingFact = toFactDetail(raw);
		modalOpen = true;
	}

	$effect(() => {
		if (!fetchError) return;
		const timer = setTimeout(() => (fetchError = null), 4000);
		return () => clearTimeout(timer);
	});

	export function openCreate() {
		editingFact = null;
		modalOpen = true;
	}

	async function handleSubmit(data: FactFormData) {
		if (editingFact) {
			const res = await fetch(`/api/notebooks/${notebookId}/facts/${editingFact.id}`, {
				method: 'PATCH',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ content: data.content })
			});
			if (!res.ok) {
				const err = await res.json().catch(() => ({ message: 'Failed to save fact' }));
				throw new Error(err.message ?? 'Failed to save fact');
			}
		} else {
			const res = await fetch(`/api/notebooks/${notebookId}/facts`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					fact_type: data.factType,
					content: data.content
				})
			});
			if (!res.ok) {
				const err = await res.json().catch(() => ({ message: 'Failed to create fact' }));
				throw new Error(err.message ?? 'Failed to create fact');
			}
		}

		await invalidateAll();
	}
</script>

{#if fetchError}
	<div class="bg-destructive/10 text-destructive border-destructive/20 border-b px-6 py-2 text-sm">
		{fetchError}
	</div>
{/if}
<div class="border-border border-b">
	<div class="flex items-center justify-between px-6 py-4">
		<div>
			<h2 class="text-xl font-bold">Facts & Cards</h2>
			<p class="text-muted-foreground text-sm">
				{stats.totalFacts} facts &middot; {stats.totalCards} cards &middot; {stats.totalDue} due for review
			</p>
		</div>
		<div class="flex items-center gap-2">
			<button
				onclick={startPractice}
				disabled={isLoadingPractice || stats.totalCards === 0}
				class="border-border text-foreground hover:bg-accent disabled:opacity-50 inline-flex items-center gap-1.5 rounded-md border px-3 py-1.5 text-sm font-medium transition-colors"
			>
				{#if isLoadingPractice}
					<LoaderCircleIcon class="h-4 w-4 animate-spin" />
				{:else}
					<RepeatIcon class="h-4 w-4" />
				{/if}
				Practice
			</button>
			<button
				onclick={startReview}
				disabled={isLoadingReview || stats.totalDue === 0}
				class="bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50 inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors"
			>
				{#if isLoadingReview}
					<LoaderCircleIcon class="h-4 w-4 animate-spin" />
				{:else}
					<ZapIcon class="h-4 w-4" />
				{/if}
				Review ({stats.totalDue})
			</button>
			<button
				class="bg-foreground text-background hover:bg-foreground/90 inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium"
				onclick={openCreate}
			>
				<PlusIcon class="h-4 w-4" />
				Create Fact
			</button>
		</div>
	</div>
	<QuickStats {stats} />
</div>

<CreateFactModal
	bind:open={modalOpen}
	{notebookId}
	editFact={editingFact}
	onclose={() => { modalOpen = false; }}
	onsubmit={handleSubmit}
/>
