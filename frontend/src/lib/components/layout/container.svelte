<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	interface Props extends HTMLAttributes<HTMLDivElement> {
		children: Snippet;
		size?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
		class?: string;
	}

	let { children, size = 'lg', class: className, ...rest }: Props = $props();

	const sizeClasses = {
		sm: 'max-w-xl',
		md: 'max-w-3xl',
		lg: 'max-w-5xl',
		xl: 'max-w-7xl',
		full: 'max-w-full'
	};
</script>

<div class={cn('w-full mx-auto px-4', sizeClasses[size], className)} {...rest}>
	{@render children()}
</div>
