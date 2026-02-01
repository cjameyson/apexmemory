<script lang="ts">
	import type { CardState } from '$lib/types/fact';

	let { state }: { state: CardState } = $props();

	const config = $derived(
		({
			new: { label: 'New', classes: 'bg-info/15 text-info' },
			learning: { label: 'Learning', classes: 'bg-warning/15 text-warning' },
			review: { label: 'Review', classes: 'bg-success/15 text-success' },
			relearning: { label: 'Relearning', classes: 'bg-destructive/15 text-destructive' }
		})[state]
	);
</script>

<span class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium {config.classes}">
	{config.label}
</span>
