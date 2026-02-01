<script lang="ts">
	import { cn } from '$lib/utils';
	import SourceIcon from '$lib/components/ui/source-icon.svelte';
	import type { Source } from '$lib/types';

	interface Props {
		source: Source;
		isSelected?: boolean;
		href?: string;
		onclick?: () => void;
		class?: string;
	}

	let { source, isSelected = false, href, onclick, class: className }: Props = $props();

	const classes = $derived(
		cn(
			'w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors',
			isSelected ? 'bg-primary/10 text-primary' : 'hover:bg-accent text-foreground',
			className
		)
	);
</script>

{#snippet content()}
	<SourceIcon type={source.type} />
	<div class="flex-1 min-w-0">
		<div class="text-sm font-medium truncate">{source.name}</div>
		{#if source.cards > 0}
			<div class="text-xs text-muted-foreground">
				{source.cards} cards
			</div>
		{/if}
	</div>
{/snippet}

{#if href}
	<a {href} class={classes} data-sveltekit-noscroll>
		{@render content()}
	</a>
{:else}
	<button type="button" {onclick} class={classes}>
		{@render content()}
	</button>
{/if}
