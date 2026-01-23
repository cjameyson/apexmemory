<script lang="ts">
	import { cn } from '$lib/utils';
	import SourceIcon from '$lib/components/ui/source-icon.svelte';
	import type { Source } from '$lib/types';

	interface Props {
		source: Source;
		isSelected?: boolean;
		onclick?: () => void;
		class?: string;
	}

	let { source, isSelected = false, onclick, class: className }: Props = $props();
</script>

<button
	type="button"
	{onclick}
	class={cn(
		'w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors',
		isSelected
			? 'bg-sky-100 dark:bg-sky-900/30 text-sky-900 dark:text-sky-100'
			: 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300',
		className
	)}
>
	<SourceIcon type={source.type} />
	<div class="flex-1 min-w-0">
		<div class="text-sm font-medium truncate">{source.name}</div>
		{#if source.cards > 0}
			<div class="text-xs text-slate-500 dark:text-slate-400">
				{source.cards} cards
			</div>
		{/if}
	</div>
</button>
