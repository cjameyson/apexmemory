<script lang="ts">
	import { cn } from '$lib/utils';
	import { Skeleton } from '$lib/components/ui/skeleton';

	interface Props {
		variant?: 'card' | 'list' | 'text' | 'avatar';
		count?: number;
		class?: string;
	}

	let { variant = 'text', count = 1, class: className }: Props = $props();
</script>

{#if variant === 'card'}
	<div class={cn('space-y-4', className)}>
		{#each Array(count) as _, i (i)}
			<div class="rounded-lg border border-border p-4 space-y-3">
				<Skeleton class="h-4 w-3/4" />
				<Skeleton class="h-4 w-1/2" />
				<Skeleton class="h-20 w-full" />
			</div>
		{/each}
	</div>
{:else if variant === 'list'}
	<div class={cn('space-y-3', className)}>
		{#each Array(count) as _, i (i)}
			<div class="flex items-center gap-3">
				<Skeleton class="h-10 w-10 rounded-full" />
				<div class="space-y-2 flex-1">
					<Skeleton class="h-4 w-1/3" />
					<Skeleton class="h-3 w-2/3" />
				</div>
			</div>
		{/each}
	</div>
{:else if variant === 'avatar'}
	<div class={cn('flex gap-2', className)}>
		{#each Array(count) as _, i (i)}
			<Skeleton class="h-10 w-10 rounded-full" />
		{/each}
	</div>
{:else}
	<div class={cn('space-y-2', className)}>
		{#each Array(count) as _, i (i)}
			<Skeleton class="h-4 w-full" />
		{/each}
	</div>
{/if}
