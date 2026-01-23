<script lang="ts">
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button';
	import { Container, Stack } from '$lib/components/layout';
</script>

<Container size="sm" class="py-16">
	<Stack gap="lg" align="center" class="text-center">
		<div class="space-y-2">
			<p class="text-6xl font-bold text-primary">{$page.status}</p>
			<h1 class="text-2xl font-semibold text-foreground">
				{#if $page.status === 404}
					Page not found
				{:else if $page.status >= 500}
					Something went wrong
				{:else}
					An error occurred
				{/if}
			</h1>
		</div>

		<p class="text-muted-foreground max-w-md">
			{#if $page.status === 404}
				The page you're looking for doesn't exist or has been moved.
			{:else}
				{$page.error?.message || 'An unexpected error occurred. Please try again later.'}
			{/if}
		</p>

		<div class="flex gap-3">
			<Button href="/">Go home</Button>
			<Button variant="outline" onclick={() => history.back()}>Go back</Button>
		</div>
	</Stack>
</Container>
