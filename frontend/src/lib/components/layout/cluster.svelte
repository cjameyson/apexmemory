<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	interface Props extends HTMLAttributes<HTMLDivElement> {
		children: Snippet;
		gap?: 'none' | 'xs' | 'sm' | 'md' | 'lg' | 'xl';
		justify?: 'start' | 'center' | 'end' | 'between' | 'around';
		align?: 'start' | 'center' | 'end' | 'stretch' | 'baseline';
		wrap?: boolean;
		class?: string;
	}

	let {
		children,
		gap = 'md',
		justify = 'start',
		align = 'center',
		wrap = true,
		class: className,
		...rest
	}: Props = $props();

	const gapClasses = {
		none: 'gap-0',
		xs: 'gap-1',
		sm: 'gap-2',
		md: 'gap-4',
		lg: 'gap-6',
		xl: 'gap-8'
	};

	const justifyClasses = {
		start: 'justify-start',
		center: 'justify-center',
		end: 'justify-end',
		between: 'justify-between',
		around: 'justify-around'
	};

	const alignClasses = {
		start: 'items-start',
		center: 'items-center',
		end: 'items-end',
		stretch: 'items-stretch',
		baseline: 'items-baseline'
	};
</script>

<div
	class={cn('flex', gapClasses[gap], justifyClasses[justify], alignClasses[align], wrap && 'flex-wrap', className)}
	{...rest}
>
	{@render children()}
</div>
