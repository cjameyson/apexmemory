<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import { ChevronRightIcon } from '@lucide/svelte';
	import type { Component, Snippet } from 'svelte';

	interface Props {
		title: string;
		icon: Component;
		count?: number;
		isOpen?: boolean;
		actions?: Snippet;
		children: Snippet;
		class?: string;
	}

	let {
		title,
		icon: Icon,
		count,
		isOpen = $bindable(true),
		actions,
		children,
		class: className
	}: Props = $props();

	function toggle() {
		isOpen = !isOpen;
	}
</script>

<div class={cn('', className)}>
	<!-- Section header -->
	<div class="flex items-center gap-1 px-2 py-1.5">
		<button
			type="button"
			onclick={toggle}
			class="flex items-center gap-2 flex-1 min-w-0 text-left hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md px-2 py-1 transition-colors"
		>
			<ChevronRightIcon
				class={cn(
					'size-4 text-slate-400 transition-transform duration-200',
					isOpen && 'rotate-90'
				)}
			/>
			<Icon class="size-4 text-slate-500 dark:text-slate-400" />
			<span class="text-sm font-medium text-slate-700 dark:text-slate-300 truncate">
				{title}
			</span>
			{#if count !== undefined}
				<span class="text-xs text-slate-400 dark:text-slate-500 ml-auto">
					{count}
				</span>
			{/if}
		</button>

		{#if actions}
			<div class="flex items-center">
				{@render actions()}
			</div>
		{/if}
	</div>

	<!-- Section content -->
	{#if isOpen}
		<div class="pl-4">
			{@render children()}
		</div>
	{/if}
</div>
