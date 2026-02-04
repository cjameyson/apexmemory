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
	import ConfirmDialog from '$lib/components/ui/confirm-dialog.svelte';
	import { studyCounts } from '$lib/stores/study-counts.svelte';

	let { stats, notebookId, notebook }: { stats: FactStats; notebookId: string; notebook: Notebook } = $props();

	// Use store for live due count (stays in sync after reviews)
	let dueCount = $derived(studyCounts.getDueCount(notebookId));

	let modalOpen = $state(false);
	let editingFact = $state<FactDetail | null>(null);
	let fetchError = $state<string | null>(null);
	let editRequestId = 0;
	let isLoadingReview = $state(false);
	let isLoadingPractice = $state(false);

	// Delete modal state
	let deleteModalOpen = $state(false);
	let deletingFactId = $state<string | null>(null);
	let deletingFactDisplay = $state<string>('');
	let isDeleting = $state(false);

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

	export function openDelete(factId: string, displayText: string) {
		deletingFactId = factId;
		deletingFactDisplay = displayText;
		deleteModalOpen = true;
	}

	async function handleDeleteConfirm() {
		if (!deletingFactId) return;
		isDeleting = true;
		try {
			const res = await fetch(`/api/notebooks/${notebookId}/facts/${deletingFactId}`, {
				method: 'DELETE'
			});
			if (!res.ok) {
				const err = await res.json().catch(() => ({ message: 'Failed to delete fact' }));
				throw new Error(err.message ?? 'Failed to delete fact');
			}
			deleteModalOpen = false;
			deletingFactId = null;
			await studyCounts.refresh();
			await invalidateAll();
		} catch (err) {
			fetchError = err instanceof Error ? err.message : 'Failed to delete fact';
		} finally {
			isDeleting = false;
		}
	}

	function handleDeleteCancel() {
		deleteModalOpen = false;
		deletingFactId = null;
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

		// Refresh study counts (new cards may be due)
		await studyCounts.refresh();
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
				{stats.totalFacts} facts &middot; {stats.totalCards} cards &middot; {dueCount} due for review
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
				disabled={isLoadingReview || dueCount === 0}
				class="bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50 inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors"
			>
				{#if isLoadingReview}
					<LoaderCircleIcon class="h-4 w-4 animate-spin" />
				{:else}
					<ZapIcon class="h-4 w-4" />
				{/if}
				Review ({dueCount})
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

<ConfirmDialog
	bind:open={deleteModalOpen}
	title="Delete Fact"
	description="Are you sure you want to delete this fact? All associated cards will be permanently deleted. Review history will be preserved but unlinked."
	loading={isDeleting}
	onconfirm={handleDeleteConfirm}
	oncancel={handleDeleteCancel}
/>
