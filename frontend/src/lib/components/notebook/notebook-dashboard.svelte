<script lang="ts">
	import type { Notebook, Source, Card } from '$lib/types';
	import { ConstructionIcon } from '@lucide/svelte';

	interface Props {
		notebook: Notebook;
		sources: Source[];
		cards: Card[];
	}

	let { notebook, sources, cards }: Props = $props();
</script>

<div class="flex flex-1 flex-col items-center justify-center p-8">
	<div class="text-center max-w-md">
		<div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-muted">
			<ConstructionIcon class="h-8 w-8 text-muted-foreground" />
		</div>
		<h1 class="text-2xl font-bold text-foreground mb-2">Dashboard</h1>
		<p class="text-muted-foreground">
			Coming soon. Use the sidebar to manage your facts and sources.
		</p>
	</div>
</div>
