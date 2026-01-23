<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Notebook } from '$lib/types';

	interface Props {
		notebook: Notebook;
		class?: string;
	}

	let { notebook, class: className }: Props = $props();
</script>

<a
	href="/notebooks/{notebook.id}"
	class={cn(
		'block bg-white dark:bg-slate-900 hover:bg-slate-50 dark:hover:bg-slate-800 rounded-2xl p-5 border border-slate-200 dark:border-slate-800 transition-colors',
		className
	)}
>
	<div class="flex items-start gap-4">
		<span class="text-3xl">{notebook.emoji}</span>
		<div class="flex-1 min-w-0">
			<div class="font-semibold text-slate-900 dark:text-white mb-1">
				{notebook.name}
			</div>
			<div class="text-sm text-slate-500 dark:text-slate-400 mb-2">
				{notebook.totalCards} cards
			</div>
			<div class="flex items-center gap-3 text-xs flex-wrap">
				{#if notebook.dueCount > 0}
					<span class="text-sky-600 dark:text-sky-400 font-medium">
						{notebook.dueCount} due
					</span>
				{/if}
				<span class="text-slate-400 dark:text-slate-500">
					{notebook.retention}% retention
				</span>
				{#if notebook.streak > 0}
					<span class="text-slate-400 dark:text-slate-500">
						🔥 {notebook.streak}
					</span>
				{/if}
			</div>
		</div>
	</div>
</a>
