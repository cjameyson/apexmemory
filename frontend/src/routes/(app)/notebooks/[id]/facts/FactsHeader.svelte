<script lang="ts">
	import type { FactStats } from '$lib/types/fact';
	import type { FactFormData } from '$lib/components/facts/create-fact-modal.svelte';
	import { PlusIcon } from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import QuickStats from './QuickStats.svelte';
	import CreateFactModal from '$lib/components/facts/create-fact-modal.svelte';

	let { stats, notebookId }: { stats: FactStats; notebookId: string } = $props();

	let createModalOpen = $state(false);

	async function handleSubmit(data: FactFormData) {
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

		await invalidateAll();
	}
</script>

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
				onclick={() => (createModalOpen = true)}
			>
				<PlusIcon class="h-4 w-4" />
				Create Fact
			</button>
		</div>
	</div>
	<QuickStats {stats} />
</div>

<CreateFactModal
	bind:open={createModalOpen}
	{notebookId}
	onclose={() => (createModalOpen = false)}
	onsubmit={handleSubmit}
/>
