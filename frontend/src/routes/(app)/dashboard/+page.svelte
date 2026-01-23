<script lang="ts">
	import { page } from '$app/stores';
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { BookOpenIcon, FolderIcon, TrendingUpIcon } from '@lucide/svelte';
</script>

<svelte:head>
	<title>Dashboard - Apex Memory</title>
</svelte:head>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold">Welcome back, {$page.data.user?.username}</h1>
		<p class="text-muted-foreground">Here's an overview of your learning progress.</p>
	</div>

	<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Due Today</Card.Title>
				<BookOpenIcon class="size-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<p class="text-3xl font-bold">0</p>
				<p class="text-xs text-muted-foreground">cards ready for review</p>
			</Card.Content>
			<Card.Footer>
				<Button href="/study" class="w-full">Start Studying</Button>
			</Card.Footer>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">Total Decks</Card.Title>
				<FolderIcon class="size-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<p class="text-3xl font-bold">0</p>
				<p class="text-xs text-muted-foreground">decks in your library</p>
			</Card.Content>
			<Card.Footer>
				<Button href="/decks" variant="outline" class="w-full">Manage Decks</Button>
			</Card.Footer>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between pb-2">
				<Card.Title class="text-sm font-medium">This Week</Card.Title>
				<TrendingUpIcon class="size-4 text-muted-foreground" />
			</Card.Header>
			<Card.Content>
				<p class="text-3xl font-bold">0</p>
				<p class="text-xs text-muted-foreground">cards reviewed</p>
			</Card.Content>
			<Card.Footer>
				<Button href="/settings" variant="outline" class="w-full">View Stats</Button>
			</Card.Footer>
		</Card.Root>
	</div>

	<!-- Recent Activity -->
	<Card.Root>
		<Card.Header>
			<Card.Title>Recent Activity</Card.Title>
			<Card.Description>Your latest study sessions</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="flex items-center justify-center py-8 text-muted-foreground">
				<p>No recent activity. Start studying to see your progress here.</p>
			</div>
		</Card.Content>
	</Card.Root>
</div>
