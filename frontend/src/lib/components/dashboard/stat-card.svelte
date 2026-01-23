<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Component } from 'svelte';

	interface Props {
		icon?: Component;
		label: string;
		value: string | number;
		class?: string;
	}

	let { icon: Icon, label, value, class: className }: Props = $props();
</script>

<div
	class={cn(
		'bg-white dark:bg-slate-900 rounded-2xl p-5 border border-slate-200 dark:border-slate-800',
		className
	)}
>
	<div class="flex items-start gap-3">
		{#if Icon}
			<div class="w-10 h-10 bg-slate-100 dark:bg-slate-800 rounded-xl flex items-center justify-center shrink-0">
				<Icon class="size-5 text-slate-500 dark:text-slate-400" />
			</div>
		{/if}
		<div class="min-w-0">
			<div class="text-sm text-slate-500 dark:text-slate-400 mb-1">{label}</div>
			<div class="text-2xl font-bold text-slate-900 dark:text-white">{value}</div>
		</div>
	</div>
</div>
