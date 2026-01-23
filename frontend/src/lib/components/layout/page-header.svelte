<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	interface Props extends HTMLAttributes<HTMLElement> {
		title: string;
		description?: string;
		actions?: Snippet;
		class?: string;
	}

	let { title, description, actions, class: className, ...rest }: Props = $props();
</script>

<header class={cn('flex flex-col gap-1 mb-6 sm:flex-row sm:items-center sm:justify-between sm:gap-4', className)} {...rest}>
	<div class="space-y-1">
		<h1 class="text-2xl font-semibold tracking-tight text-foreground sm:text-3xl">{title}</h1>
		{#if description}
			<p class="text-muted-foreground">{description}</p>
		{/if}
	</div>
	{#if actions}
		<div class="flex items-center gap-2 mt-4 sm:mt-0">
			{@render actions()}
		</div>
	{/if}
</header>
