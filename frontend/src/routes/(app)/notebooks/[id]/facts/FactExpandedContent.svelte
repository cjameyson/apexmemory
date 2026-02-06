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

	$effect(() => {
		const factId = fact.id;
		const currentNotebookId = notebookId;
		const controller = new AbortController();

		loading = true;
		error = null;
		cards = null;

		const loadCards = async () => {
			try {
				const res = await fetch(`/api/notebooks/${currentNotebookId}/facts/${factId}`, {
					signal: controller.signal
				});
				if (!res.ok) {
					error = 'Failed to load cards';
					return;
				}

				const data: ApiFactDetail = await res.json();
				cards = data.cards.map(toCard);
			} catch {
				if (controller.signal.aborted) return;
				error = 'Failed to load cards';
			} finally {
				if (!controller.signal.aborted) {
					loading = false;
				}
			}
		};

		void loadCards();

		return () => {
			controller.abort();
		};
	});

	interface FieldDisplay {
		name: string;
		type: string;
		value: string | Record<string, unknown>;
	}

	const allFields = $derived<FieldDisplay[]>(
		(fact.content as { fields?: FieldDisplay[] })?.fields ?? []
	);

	// Hide image_occlusion JSON blob from field display
	const fields = $derived(
		allFields.filter((f) => f.type !== 'image_occlusion')
	);

	// Parse image occlusion regions for label/hint lookup
	const imageOcclusionRegions = $derived.by(() => {
		const field = allFields.find((f) => f.type === 'image_occlusion');
		if (!field) return new Map<string, { label: string; hint: string }>();
		try {
			const parsed = (typeof field.value === 'string' ? JSON.parse(field.value) : field.value) as {
				regions?: { id: string; label: string; hint?: string }[];
			};
			return new Map(
				(parsed.regions ?? []).map((r) => [r.id, { label: r.label, hint: r.hint ?? '' }])
			);
		} catch {
			return new Map<string, { label: string; hint: string }>();
		}
	});

	function formatFieldLabel(name: string): string {
		return name.charAt(0).toUpperCase() + name.slice(1);
	}

	function extractText(value: string | Record<string, unknown> | undefined): string {
		if (!value) return '';
		if (typeof value === 'string') {
			return value.replace(/<[^>]*>/g, '').trim();
		}
		// TipTap JSON: walk text nodes
		const parts: string[] = [];
		function walk(node: Record<string, unknown>) {
			if (node.type === 'text' && typeof node.text === 'string') {
				parts.push(node.text);
			}
			if (Array.isArray(node.content)) {
				for (const child of node.content) {
					walk(child as Record<string, unknown>);
				}
			}
		}
		walk(value);
		return parts.join(' ').replace(/\s+/g, ' ').trim();
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

		const raw = extractText(clozeField.value);
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
				{#each fields as field (`${field.name}-${field.type}`)}
					<div>
						<span class="text-xs font-medium text-muted-foreground">{formatFieldLabel(field.name)}</span>
						<p class="mt-0.5 text-sm text-foreground">{extractText(field.value) || '(empty)'}</p>
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
						{#each cards as card (card.id)}
							{@const cloze = fact.factType === 'cloze' ? getClozeInfo(card.elementId) : null}
							<tr class="border-b border-border/50">
							<td class="py-1.5"><CardStateBadge state={card.state} /></td>
							{#if fact.factType === 'cloze'}
								<td class="py-1.5 text-muted-foreground">{card.elementId}</td>
								<td class="py-1.5">{cloze?.text || '--'}</td>
								<td class="py-1.5 text-muted-foreground">{cloze?.hint || '--'}</td>
							{:else if fact.factType === 'image_occlusion'}
								{@const region = imageOcclusionRegions.get(card.elementId)}
								<td class="py-1.5">{region?.label || card.elementId || '--'}</td>
								<td class="py-1.5 text-muted-foreground">{region?.hint || '--'}</td>
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
