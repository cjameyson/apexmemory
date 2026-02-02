<script lang="ts">
	import { Bold, Italic, Underline, Code, Link, List, ListOrdered } from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';

	let { class: className = '' }: { class?: string } = $props();

	const groups = [
		[
			{ icon: Bold, label: 'Bold' },
			{ icon: Italic, label: 'Italic' },
			{ icon: Underline, label: 'Underline' }
		],
		[
			{ icon: Code, label: 'Code' },
			{ icon: Link, label: 'Link' }
		],
		[
			{ icon: List, label: 'Bullet list' },
			{ icon: ListOrdered, label: 'Numbered list' }
		]
	];
</script>

<div class="flex items-center gap-0.5 {className}">
	{#each groups as group, gi}
		{#if gi > 0}
			<div class="bg-border mx-1 h-4 w-px"></div>
		{/if}
		{#each group as item}
			<!-- TODO: Wire to TipTap commands -->
			<Button
				variant="ghost"
				size="icon"
				class="h-7 w-7"
				title={item.label}
				disabled
			>
				<item.icon class="h-3.5 w-3.5" />
			</Button>
		{/each}
	{/each}
</div>
