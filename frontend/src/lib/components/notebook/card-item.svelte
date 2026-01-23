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
		'w-full text-left bg-slate-50 dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-xl p-4 border border-slate-200 dark:border-slate-700 transition-colors',
		className
	)}
>
	<!-- Card front (question) -->
	<div class="text-sm font-medium text-slate-900 dark:text-white mb-2 line-clamp-2">
		{card.front}
	</div>

	<!-- Card back (answer) preview -->
	<div class="text-xs text-slate-500 dark:text-slate-400 line-clamp-2 mb-3">
		{card.back}
	</div>

	<!-- Card metadata -->
	<div class="flex items-center justify-between text-xs">
		{#if card.due}
			<span class="text-sky-600 dark:text-sky-400 font-medium">Due</span>
		{:else}
			<span class="text-slate-400 dark:text-slate-500">Next: {card.interval}</span>
		{/if}

		<div class="flex items-center gap-2">
			{#if card.tags.length > 0}
				<span class="text-slate-400 dark:text-slate-500 truncate max-w-24">
					{card.tags.slice(0, 2).join(', ')}
				</span>
			{/if}
		</div>
	</div>
</button>
