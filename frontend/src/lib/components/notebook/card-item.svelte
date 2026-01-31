<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Card, Source } from '$lib/types';

	interface Props {
		card: Card;
		source?: Source;
		onclick?: () => void;
		class?: string;
	}

	let { card, source, onclick, class: className }: Props = $props();
</script>

<button
	type="button"
	{onclick}
	class={cn(
		'w-full text-left bg-muted hover:bg-accent rounded-xl p-4 border border-border transition-colors',
		className
	)}
>
	<!-- Card front (question) -->
	<div class="text-sm font-medium text-foreground mb-2 line-clamp-2">
		{card.front}
	</div>

	<!-- Card back (answer) preview -->
	<div class="text-xs text-muted-foreground line-clamp-2 mb-3">
		{card.back}
	</div>

	<!-- Card metadata -->
	<div class="flex items-center justify-between text-xs">
		{#if card.due}
			<span class="text-primary font-medium">Due</span>
		{:else}
			<span class="text-muted-foreground">Next: {card.interval}</span>
		{/if}

		<div class="flex items-center gap-2">
			{#if card.tags.length > 0}
				<span class="text-muted-foreground truncate max-w-24">
					{card.tags.slice(0, 2).join(', ')}
				</span>
			{/if}
		</div>
	</div>
</button>
