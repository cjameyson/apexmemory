<script lang="ts">
	import type { FactStats, FactDetail } from '$lib/types/fact';
	import type { ApiFactDetail } from '$lib/api/types';
	import type { FactFormData } from '$lib/components/facts/create-fact-modal.svelte';
	import { toFactDetail } from '$lib/services/facts';
	import { PlusIcon } from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import QuickStats from './QuickStats.svelte';
	import CreateFactModal from '$lib/components/facts/create-fact-modal.svelte';

	let { stats, notebookId }: { stats: FactStats; notebookId: string } = $props();

	let modalOpen = $state(false);
	let editingFact = $state<FactDetail | null>(null);
	let fetchError = $state<string | null>(null);
	let editRequestId = 0;

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

	function openCreate() {
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
			<a
				href="/notebooks/{notebookId}/review"
				class="bg-primary text-primary-foreground hover:bg-primary/90 inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium"
			>
				Review ({stats.totalDue})
			</a>
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
