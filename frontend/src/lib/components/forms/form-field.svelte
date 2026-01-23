<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	interface Props {
		label: string;
		name: string;
		error?: string;
		description?: string;
		required?: boolean;
		class?: string;
		children: Snippet;
	}

	let {
		label,
		name,
		error,
		description,
		required = false,
		class: className,
		children,
	}: Props = $props();
</script>

<div class={cn('space-y-2', className)}>
	<label for={name} class="text-sm font-medium leading-none">
		{label}
		{#if required}
			<span class="text-destructive">*</span>
		{/if}
	</label>

	{@render children()}

	{#if description}
		<p class="text-sm text-muted-foreground">{description}</p>
	{/if}
	{#if error}
		<p class="text-sm text-destructive" role="alert" aria-live="polite">{error}</p>
	{/if}
</div>
