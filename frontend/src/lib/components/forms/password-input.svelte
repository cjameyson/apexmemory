<script lang="ts">
	import { cn } from '$lib/utils';
	import { EyeIcon, EyeOffIcon } from '@lucide/svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';

	interface Props extends Omit<HTMLInputAttributes, 'type'> {
		value?: string;
		error?: boolean;
	}

	let {
		value = $bindable(''),
		error = false,
		class: className,
		...restProps
	}: Props = $props();

	let showPassword = $state(false);
</script>

<div class="relative">
	<input
		type={showPassword ? 'text' : 'password'}
		bind:value
		aria-invalid={error}
		class={cn(
			'border-input bg-background selection:bg-primary dark:bg-input/30 selection:text-primary-foreground ring-offset-background placeholder:text-muted-foreground flex h-9 w-full min-w-0 rounded-md border px-3 py-1 pr-10 text-base shadow-xs transition-[color,box-shadow] outline-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
			'focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]',
			'aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive',
			className
		)}
		{...restProps}
	/>
	<button
		type="button"
		title={showPassword ? 'Hide password' : 'Show password'}
		class="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-muted-foreground hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 rounded-sm transition-colors"
		onclick={() => (showPassword = !showPassword)}
	>
		{#if showPassword}
			<EyeOffIcon class="size-4" />
		{:else}
			<EyeIcon class="size-4" />
		{/if}
	</button>
</div>
