<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { PlusIcon, FolderIcon, SearchIcon } from '@lucide/svelte';

	let searchQuery = $state('');
</script>

<svelte:head>
	<title>Decks - Apex Memory</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
		<div>
			<h1 class="text-2xl font-bold">Decks</h1>
			<p class="text-muted-foreground">Organize your flashcards into decks.</p>
		</div>
		<Button disabled>
			<PlusIcon class="size-4 mr-2" />
			New Deck
		</Button>
	</div>

	<!-- Search -->
	<div class="relative max-w-sm">
		<SearchIcon class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
		<Input
			type="search"
			placeholder="Search decks..."
			class="pl-9"
			bind:value={searchQuery}
			disabled
		/>
	</div>

	<!-- Empty State -->
	<Card.Root>
		<Card.Content class="py-12">
			<div class="flex flex-col items-center justify-center text-center">
				<FolderIcon class="size-12 text-muted-foreground mb-4" />
				<h3 class="text-lg font-semibold mb-2">No decks yet</h3>
				<p class="text-muted-foreground mb-4 max-w-sm">
					Create your first deck to start organizing your flashcards. Decks can be nested to create a hierarchy.
				</p>
				<Button disabled>
					<PlusIcon class="size-4 mr-2" />
					Create Your First Deck
				</Button>
			</div>
		</Card.Content>
	</Card.Root>
</div>
