<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { BookOpenIcon, ClockIcon, CheckCircleIcon } from '@lucide/svelte';
</script>

<svelte:head>
	<title>Study - Apex Memory</title>
</svelte:head>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold">Study</h1>
		<p class="text-muted-foreground">Review your flashcards using spaced repetition.</p>
	</div>

	<!-- Study Queue Summary -->
	<div class="grid gap-4 sm:grid-cols-3">
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">New</Card.Title>
				<BookOpenIcon class="size-4 text-blue-500" />
			</Card.Header>
			<Card.Content>
				<p class="text-2xl font-bold">0</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Learning</Card.Title>
				<ClockIcon class="size-4 text-orange-500" />
			</Card.Header>
			<Card.Content>
				<p class="text-2xl font-bold">0</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Review</Card.Title>
				<CheckCircleIcon class="size-4 text-green-500" />
			</Card.Header>
			<Card.Content>
				<p class="text-2xl font-bold">0</p>
			</Card.Content>
		</Card.Root>
	</div>

	<!-- Study Session Card -->
	<Card.Root class="max-w-2xl">
		<Card.Header>
			<Card.Title>Ready to Study?</Card.Title>
			<Card.Description>
				You have no cards due for review right now. Create some decks and add cards to get started.
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="flex flex-col items-center justify-center py-12 text-center">
				<BookOpenIcon class="size-12 text-muted-foreground mb-4" />
				<p class="text-muted-foreground mb-4">No cards to study</p>
				<Button href="/decks">Create a Deck</Button>
			</div>
		</Card.Content>
	</Card.Root>
</div>
