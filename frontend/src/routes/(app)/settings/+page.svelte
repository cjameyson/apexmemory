<script lang="ts">
	import { page } from '$app/stores';
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { FormField } from '$lib/components/forms';
	import { UserIcon, BellIcon, PaletteIcon, DatabaseIcon } from '@lucide/svelte';
</script>

<svelte:head>
	<title>Settings - Apex Memory</title>
</svelte:head>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold">Settings</h1>
		<p class="text-muted-foreground">Manage your account and preferences.</p>
	</div>

	<!-- Profile Section -->
	<Card.Root>
		<Card.Header>
			<div class="flex items-center gap-2">
				<UserIcon class="size-5" />
				<Card.Title>Profile</Card.Title>
			</div>
			<Card.Description>Your account information</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<FormField label="Email" name="email">
				<Input
					type="email"
					name="email"
					value={$page.data.user?.email ?? ''}
					disabled
				/>
			</FormField>
			<FormField label="Username" name="username">
				<Input
					type="text"
					name="username"
					value={$page.data.user?.username ?? ''}
					disabled
				/>
			</FormField>
		</Card.Content>
		<Card.Footer>
			<Button variant="outline" disabled>Update Profile</Button>
		</Card.Footer>
	</Card.Root>

	<!-- Study Preferences -->
	<Card.Root>
		<Card.Header>
			<div class="flex items-center gap-2">
				<DatabaseIcon class="size-5" />
				<Card.Title>Study Preferences</Card.Title>
			</div>
			<Card.Description>Customize your learning experience</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<FormField label="New Cards per Day" name="newCardsPerDay" description="Maximum new cards to introduce daily">
				<Input
					type="number"
					name="newCardsPerDay"
					value="20"
					disabled
				/>
			</FormField>
			<FormField label="Maximum Reviews per Day" name="maxReviewsPerDay" description="Maximum reviews per day (0 = unlimited)">
				<Input
					type="number"
					name="maxReviewsPerDay"
					value="200"
					disabled
				/>
			</FormField>
		</Card.Content>
		<Card.Footer>
			<Button variant="outline" disabled>Save Preferences</Button>
		</Card.Footer>
	</Card.Root>

	<!-- Notifications -->
	<Card.Root>
		<Card.Header>
			<div class="flex items-center gap-2">
				<BellIcon class="size-5" />
				<Card.Title>Notifications</Card.Title>
			</div>
			<Card.Description>Configure how you receive reminders</Card.Description>
		</Card.Header>
		<Card.Content>
			<p class="text-sm text-muted-foreground">Notification settings coming soon.</p>
		</Card.Content>
	</Card.Root>

	<!-- Appearance -->
	<Card.Root>
		<Card.Header>
			<div class="flex items-center gap-2">
				<PaletteIcon class="size-5" />
				<Card.Title>Appearance</Card.Title>
			</div>
			<Card.Description>Customize the look and feel</Card.Description>
		</Card.Header>
		<Card.Content>
			<p class="text-sm text-muted-foreground">
				Use the theme toggle in the header to switch between light and dark modes.
			</p>
		</Card.Content>
	</Card.Root>

	<!-- Danger Zone -->
	<Card.Root class="border-destructive/50">
		<Card.Header>
			<Card.Title class="text-destructive">Danger Zone</Card.Title>
			<Card.Description>Irreversible actions</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="flex items-center justify-between">
				<div>
					<p class="font-medium">Delete Account</p>
					<p class="text-sm text-muted-foreground">Permanently delete your account and all data</p>
				</div>
				<Button variant="destructive" disabled>Delete Account</Button>
			</div>
		</Card.Content>
	</Card.Root>
</div>
