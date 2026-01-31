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
		'block bg-card hover:bg-accent rounded-2xl p-5 border border-border transition-colors',
		className
	)}
>
	<div class="flex items-start gap-4">
		<span class="text-3xl">{notebook.emoji}</span>
		<div class="flex-1 min-w-0">
			<div class="font-semibold text-foreground mb-1">
				{notebook.name}
			</div>
			<div class="text-sm text-muted-foreground mb-2">
				{notebook.totalCards} cards
			</div>
			<div class="flex items-center gap-3 text-xs flex-wrap">
				{#if notebook.dueCount > 0}
					<span class="text-primary font-medium">
						{notebook.dueCount} due
					</span>
				{/if}
				<span class="text-muted-foreground">
					{notebook.retention}% retention
				</span>
				{#if notebook.streak > 0}
					<span class="text-muted-foreground">
						🔥 {notebook.streak}
					</span>
				{/if}
			</div>
		</div>
	</div>
</a>
