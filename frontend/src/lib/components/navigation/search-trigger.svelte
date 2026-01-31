<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import { SearchIcon } from '@lucide/svelte';

	interface Props {
		onclick?: () => void;
		class?: string;
	}

	let { onclick, class: className }: Props = $props();
</script>

<Button
	variant="outline"
	size="sm"
	class={cn('gap-2 text-muted-foreground', className)}
	{onclick}
>
	<SearchIcon class="size-4" />
	<span class="hidden sm:inline">Search</span>
	<kbd class="hidden md:inline-flex items-center gap-0.5 text-xs bg-muted px-1.5 py-0.5 rounded">
		<span class="text-[10px]">⌘</span>K
	</kbd>
</Button>
