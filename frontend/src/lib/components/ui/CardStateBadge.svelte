<script lang="ts">
	import type { CardState } from '$lib/types/fact';

	let { state }: { state: CardState } = $props();

	const config = $derived(
		({
			new: {
				label: 'New',
				classes: 'bg-info/15 text-info',
				tooltip: 'Card has never been reviewed. Will enter the learning phase on first review.'
			},
			learning: {
				label: 'Learning',
				classes: 'bg-warning/15 text-warning',
				tooltip: 'Card is being learned. Short intervals until you demonstrate retention.'
			},
			review: {
				label: 'Review',
				classes: 'bg-success/15 text-success',
				tooltip: 'Card graduated to review. Intervals grow with each successful review.'
			},
			relearning: {
				label: 'Relearning',
				classes: 'bg-destructive/15 text-destructive',
				tooltip: 'Card was forgotten and is being relearned. Will return to review once relearned.'
			}
		})[state]
	);
</script>

<span
	class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium cursor-help {config.classes}"
	title={config.tooltip}
>
	{config.label}
</span>
