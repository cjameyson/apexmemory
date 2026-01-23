<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	interface Props extends HTMLAttributes<HTMLDivElement> {
		children: Snippet;
		gap?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl';
		align?: 'start' | 'center' | 'end' | 'stretch';
		class?: string;
	}

	let { children, gap = 'md', align = 'stretch', class: className, ...rest }: Props = $props();

	const gapClasses = {
		none: 'gap-0',
		xs: 'gap-1',
		sm: 'gap-2',
		md: 'gap-4',
		lg: 'gap-6',
		xl: 'gap-8'
	};

	const alignClasses = {
		start: 'items-start',
		center: 'items-center',
		end: 'items-end',
		stretch: 'items-stretch'
	};
</script>

<div class={cn('flex flex-col', gapClasses[gap], alignClasses[align], className)} {...rest}>
	{@render children()}
</div>
