<script lang="ts">
	import type { Fact, Card } from '$lib/types/fact';
	import type { ApiFactDetail } from '$lib/api/types';
	import { toCard } from '$lib/services/facts';
	import CardStateBadge from '$lib/components/ui/CardStateBadge.svelte';
	import { LoaderCircleIcon } from '@lucide/svelte';

	let {
		fact,
		notebookId
	}: {
		fact: Fact;
		notebookId: string;
	} = $props();

	let cards = $state<Card[] | null>(null);
	let loading = $state(false);
	let error = $state<string | null>(null);
	let fetched = $state(false);

	$effect(() => {
		if (fetched) return;
		fetched = true;

		loading = true;
		error = null;
		fetch(`/api/notebooks/${notebookId}/facts/${fact.id}`)
			.then(async (res) => {
				if (!res.ok) {
					error = 'Failed to load cards';
					return;
				}
				const data: ApiFactDetail = await res.json();
				cards = data.cards.map(toCard);
			})
			.catch(() => {
				error = 'Failed to load cards';
			})
			.finally(() => {
				loading = false;
			});
	});

	interface FieldDisplay {
		name: string;
		type: string;
		value: string;
	}

	const fields = $derived<FieldDisplay[]>(
		(fact.content as { fields?: FieldDisplay[] })?.fields ?? []
	);

	function formatFieldLabel(name: string): string {
		return name.charAt(0).toUpperCase() + name.slice(1);
	}

	function stripHtml(html: string | undefined): string {
		if (!html) return '';
		return html.replace(/<[^>]*>/g, '').trim();
	}

	function formatDue(due: string | null): string {
		if (!due) return '--';
		const d = new Date(due);
		const now = new Date();
		const diffMs = d.getTime() - now.getTime();
		const diffDays = Math.round(diffMs / (1000 * 60 * 60 * 24));
		if (diffDays < 0) return `${Math.abs(diffDays)}d ago`;
		if (diffDays === 0) return 'Today';
		return `${diffDays}d`;
	}

	/**
	 * Extract cloze answer text and hint for a given cloze ID (e.g. "c1").
	 * Pattern: {{c1::answer}} or {{c1::answer::hint}}
	 */
	function getClozeInfo(elementId: string): { text: string; hint: string } {
		const clozeField = fields.find((f) => f.type === 'cloze_text');
		if (!clozeField) return { text: '', hint: '' };

		const raw = stripHtml(clozeField.value);
		const num = elementId.replace('c', '');
		const re = new RegExp(`\\{\\{c${num}::([^}]*?)(?:::([^}]*))?\\}\\}`);
		const match = raw.match(re);
		if (!match) return { text: '', hint: '' };
		return { text: match[1] ?? '', hint: match[2] ?? '' };
	}
</script>

<div class="space-y-4">
	<!-- Fact Fields -->
	<div>
		<h4 class="mb-2 text-xs font-medium uppercase tracking-wider text-muted-foreground">Fields</h4>
		<div class="space-y-2">
			{#each fields as field}
				<div>
					<span class="text-xs font-medium text-muted-foreground">{formatFieldLabel(field.name)}</span>
					<p class="mt-0.5 text-sm text-foreground">{stripHtml(field.value) || '(empty)'}</p>
				</div>
			{/each}
		</div>
	</div>

	<!-- Cards -->
	<div>
		<h4 class="mb-2 text-xs font-medium uppercase tracking-wider text-muted-foreground">Cards</h4>
		{#if loading}
			<div class="flex items-center gap-2 py-2 text-sm text-muted-foreground">
				<LoaderCircleIcon class="h-4 w-4 animate-spin" />
				Loading cards...
			</div>
		{:else if error}
			<p class="text-sm text-destructive">{error}</p>
		{:else if cards && cards.length > 0}
			<table class="w-full text-sm">
				<thead>
					<tr class="border-b border-border text-xs text-muted-foreground">
						<th class="pb-1 text-left font-medium">State</th>
						{#if fact.factType === 'cloze'}
							<th class="pb-1 text-left font-medium">Cloze</th>
							<th class="pb-1 text-left font-medium">Text</th>
							<th class="pb-1 text-left font-medium">Hint</th>
						{:else if fact.factType === 'image_occlusion'}
							<th class="pb-1 text-left font-medium">Label</th>
							<th class="pb-1 text-left font-medium">Hint</th>
						{/if}
						<th class="pb-1 text-left font-medium">Due</th>
						<th class="pb-1 text-right font-medium">Reps</th>
						<th class="pb-1 text-right font-medium">Lapses</th>
					</tr>
				</thead>
				<tbody>
					{#each cards as card}
						{@const cloze = fact.factType === 'cloze' ? getClozeInfo(card.elementId) : null}
						<tr class="border-b border-border/50">
							<td class="py-1.5"><CardStateBadge state={card.state} /></td>
							{#if fact.factType === 'cloze'}
								<td class="py-1.5 text-muted-foreground">{card.elementId}</td>
								<td class="py-1.5">{cloze?.text || '--'}</td>
								<td class="py-1.5 text-muted-foreground">{cloze?.hint || '--'}</td>
							{:else if fact.factType === 'image_occlusion'}
								<td class="py-1.5">{card.elementId || '--'}</td>
								<td class="py-1.5 text-muted-foreground">--</td>
							{/if}
							<td class="py-1.5">{formatDue(card.due)}</td>
							<td class="py-1.5 text-right">{card.reps}</td>
							<td class="py-1.5 text-right">{card.lapses}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{:else}
			<p class="text-sm text-muted-foreground">No cards</p>
		{/if}
	</div>
</div>
