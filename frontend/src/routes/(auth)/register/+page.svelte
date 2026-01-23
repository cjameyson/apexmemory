<script lang="ts">
	import { enhance } from '$app/forms';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Card from '$lib/components/ui/card';
	import { FormField, PasswordInput } from '$lib/components/forms';
	import { createFormEnhance } from '$lib/forms/enhance';
	import { toast } from 'svelte-sonner';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();
	let loading = $state(false);

	const formEnhance = createFormEnhance((v) => (loading = v));

	$effect(() => {
		if (form?.error && !form.fieldErrors) {
			toast.error(form.error);
		}
	});
</script>

<svelte:head>
	<title>Create Account - Apex Memory</title>
</svelte:head>

<Card.Root>
	<Card.Header class="space-y-1">
		<Card.Title class="text-2xl">Create an account</Card.Title>
		<Card.Description>Enter your details to get started</Card.Description>
	</Card.Header>

	<Card.Content>
		<form method="POST" use:enhance={formEnhance} class="space-y-4">
			<FormField
				label="Email"
				name="email"
				error={form?.fieldErrors?.email}
				required
			>
				<Input
					type="email"
					name="email"
					id="email"
					placeholder="you@example.com"
					autocomplete="email"
					value={form?.values?.email ?? ''}
					aria-invalid={!!form?.fieldErrors?.email}
					disabled={loading}
					required
				/>
			</FormField>

			<FormField
				label="Username"
				name="username"
				error={form?.fieldErrors?.username}
				description="3-30 characters, letters, numbers, underscores, and hyphens"
				required
			>
				<Input
					type="text"
					name="username"
					id="username"
					placeholder="johndoe"
					autocomplete="username"
					value={form?.values?.username ?? ''}
					aria-invalid={!!form?.fieldErrors?.username}
					disabled={loading}
					required
				/>
			</FormField>

			<FormField
				label="Password"
				name="password"
				error={form?.fieldErrors?.password}
				description="At least 8 characters"
				required
			>
				<PasswordInput
					name="password"
					id="password"
					autocomplete="new-password"
					error={!!form?.fieldErrors?.password}
					disabled={loading}
					required
				/>
			</FormField>

			<FormField
				label="Confirm Password"
				name="confirmPassword"
				error={form?.fieldErrors?.confirmPassword}
				required
			>
				<PasswordInput
					name="confirmPassword"
					id="confirmPassword"
					autocomplete="new-password"
					error={!!form?.fieldErrors?.confirmPassword}
					disabled={loading}
					required
				/>
			</FormField>

			<Button type="submit" class="w-full" disabled={loading}>
				{#if loading}
					Creating account...
				{:else}
					Create account
				{/if}
			</Button>
		</form>
	</Card.Content>

	<Card.Footer class="flex justify-center">
		<p class="text-sm text-muted-foreground">
			Already have an account?
			<a href="/login" class="text-primary hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 rounded-sm">Sign in</a>
		</p>
	</Card.Footer>
</Card.Root>
