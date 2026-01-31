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
		'bg-card rounded-2xl p-5 border border-border',
		className
	)}
>
	<div class="flex items-start gap-3">
		{#if Icon}
			<div class="w-10 h-10 bg-muted rounded-xl flex items-center justify-center shrink-0">
				<Icon class="size-5 text-muted-foreground" />
			</div>
		{/if}
		<div class="min-w-0">
			<div class="text-sm text-muted-foreground mb-1">{label}</div>
			<div class="text-2xl font-bold text-foreground">{value}</div>
		</div>
	</div>
</div>
