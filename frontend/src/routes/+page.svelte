<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { PageHeader, Stack, Cluster, Container } from '$lib/components/layout';
	import Loading from '$lib/components/ui/loading.svelte';
	import { toast } from 'svelte-sonner';

	let dialogOpen = $state(false);
	let showLoading = $state(false);

	function handleShowToast() {
		toast.success('Toast notification working!', {
			description: 'This is a success message.'
		});
	}

	function toggleLoading() {
		showLoading = !showLoading;
	}
</script>

<Container>
	<Stack gap="xl">
		<PageHeader
			title="Apex Memory"
			description="Design System Foundation"
		>
			{#snippet actions()}
				<Button variant="outline" onclick={handleShowToast}>Show Toast</Button>
			{/snippet}
		</PageHeader>

		<!-- FSRS Rating Colors -->
		<Card.Root>
			<Card.Header>
				<Card.Title>FSRS Rating Colors</Card.Title>
				<Card.Description>Spaced repetition answer rating colors</Card.Description>
			</Card.Header>
			<Card.Content>
				<Cluster gap="sm">
					<div class="flex items-center gap-2">
						<div class="w-8 h-8 rounded-md bg-again"></div>
						<span class="text-sm">Again</span>
					</div>
					<div class="flex items-center gap-2">
						<div class="w-8 h-8 rounded-md bg-hard"></div>
						<span class="text-sm">Hard</span>
					</div>
					<div class="flex items-center gap-2">
						<div class="w-8 h-8 rounded-md bg-good"></div>
						<span class="text-sm">Good</span>
					</div>
					<div class="flex items-center gap-2">
						<div class="w-8 h-8 rounded-md bg-easy"></div>
						<span class="text-sm">Easy</span>
					</div>
				</Cluster>
			</Card.Content>
		</Card.Root>

		<!-- Button Variants -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Buttons</Card.Title>
				<Card.Description>Button component variants</Card.Description>
			</Card.Header>
			<Card.Content>
				<Stack gap="md">
					<Cluster gap="sm">
						<Button>Primary</Button>
						<Button variant="secondary">Secondary</Button>
						<Button variant="outline">Outline</Button>
						<Button variant="ghost">Ghost</Button>
						<Button variant="destructive">Destructive</Button>
						<Button variant="link">Link</Button>
					</Cluster>
					<Cluster gap="sm">
						<Button size="sm">Small</Button>
						<Button size="default">Default</Button>
						<Button size="lg">Large</Button>
						<Button size="icon">
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M5 12h14"/></svg>
						</Button>
					</Cluster>
				</Stack>
			</Card.Content>
		</Card.Root>

		<!-- Input -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Input</Card.Title>
				<Card.Description>Text input component</Card.Description>
			</Card.Header>
			<Card.Content>
				<Stack gap="sm" class="max-w-sm">
					<Input placeholder="Default input" />
					<Input type="email" placeholder="Email input" />
					<Input disabled placeholder="Disabled input" />
				</Stack>
			</Card.Content>
		</Card.Root>

		<!-- Dialog -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Dialog</Card.Title>
				<Card.Description>Modal dialog component</Card.Description>
			</Card.Header>
			<Card.Content>
				<Dialog.Root bind:open={dialogOpen}>
					<Dialog.Trigger>
						{#snippet child({ props })}
							<Button {...props}>Open Dialog</Button>
						{/snippet}
					</Dialog.Trigger>
					<Dialog.Content>
						<Dialog.Header>
							<Dialog.Title>Dialog Title</Dialog.Title>
							<Dialog.Description>
								This is a dialog description. Dialogs are used for important information or actions.
							</Dialog.Description>
						</Dialog.Header>
						<Dialog.Footer>
							<Button variant="outline" onclick={() => (dialogOpen = false)}>Cancel</Button>
							<Button onclick={() => (dialogOpen = false)}>Confirm</Button>
						</Dialog.Footer>
					</Dialog.Content>
				</Dialog.Root>
			</Card.Content>
		</Card.Root>

		<!-- Dropdown Menu -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Dropdown Menu</Card.Title>
				<Card.Description>Dropdown menu component</Card.Description>
			</Card.Header>
			<Card.Content>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Button variant="outline" {...props}>Open Menu</Button>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content>
						<DropdownMenu.Item>Profile</DropdownMenu.Item>
						<DropdownMenu.Item>Settings</DropdownMenu.Item>
						<DropdownMenu.Separator />
						<DropdownMenu.Item>Log out</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Card.Content>
		</Card.Root>

		<!-- Loading States -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Loading States</Card.Title>
				<Card.Description>Skeleton loading components</Card.Description>
			</Card.Header>
			<Card.Content>
				<Stack gap="lg">
					<Button variant="outline" onclick={toggleLoading}>
						{showLoading ? 'Hide' : 'Show'} Loading States
					</Button>

					{#if showLoading}
						<div class="grid gap-6 md:grid-cols-2">
							<div>
								<p class="text-sm text-muted-foreground mb-2">Card Loading</p>
								<Loading variant="card" count={2} />
							</div>
							<div>
								<p class="text-sm text-muted-foreground mb-2">List Loading</p>
								<Loading variant="list" count={3} />
							</div>
						</div>
					{/if}
				</Stack>
			</Card.Content>
		</Card.Root>

		<!-- Typography -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Typography</Card.Title>
				<Card.Description>Fluid typography scale</Card.Description>
			</Card.Header>
			<Card.Content>
				<Stack gap="sm">
					<p class="text-5xl font-bold">Heading 5XL</p>
					<p class="text-4xl font-bold">Heading 4XL</p>
					<p class="text-3xl font-semibold">Heading 3XL</p>
					<p class="text-2xl font-semibold">Heading 2XL</p>
					<p class="text-xl font-medium">Heading XL</p>
					<p class="text-lg">Text Large</p>
					<p class="text-base">Text Base</p>
					<p class="text-sm text-muted-foreground">Text Small (muted)</p>
					<p class="text-xs text-muted-foreground">Text Extra Small (muted)</p>
				</Stack>
			</Card.Content>
		</Card.Root>

		<!-- Prose -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Prose</Card.Title>
				<Card.Description>Markdown-style content rendering</Card.Description>
			</Card.Header>
			<Card.Content>
				<div class="prose">
					<h2>Sample Content</h2>
					<p>This is a paragraph demonstrating the prose styling for rendered markdown content in flashcards.</p>
					<ul>
						<li>First item</li>
						<li>Second item</li>
						<li>Third item</li>
					</ul>
					<blockquote>This is a blockquote for important information.</blockquote>
					<p>Inline <code>code</code> looks like this.</p>
				</div>
			</Card.Content>
		</Card.Root>
	</Stack>
</Container>
