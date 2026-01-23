<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { navigating, page } from '$app/stores';
	import { theme } from '$lib/stores/theme.svelte';
	import ThemeToggle from '$lib/components/ui/theme-toggle.svelte';
	import { Toaster } from '$lib/components/ui/sonner';
	import { Button } from '$lib/components/ui/button';

	let { children } = $props();

	onMount(() => {
		theme.init();
	});

	let loggingOut = $state(false);

	async function handleLogout() {
		loggingOut = true;
		const form = document.createElement('form');
		form.method = 'POST';
		form.action = '/api/auth/logout';
		document.body.appendChild(form);
		form.submit();
	}
</script>

<svelte:head>
	<title>Apex Memory</title>
	<meta name="description" content="Simple and effective spaced repetition learning" />
</svelte:head>

<!-- Skip to main content link for accessibility -->
<a href="#main-content" class="skip-link">Skip to main content</a>

<!-- Navigation loading indicator -->
{#if $navigating}
	<div class="fixed top-0 left-0 right-0 h-[3px] z-[9999] bg-muted overflow-hidden" role="progressbar" aria-label="Loading page">
		<div class="h-full bg-primary animate-progress"></div>
	</div>
{/if}

<div class="flex flex-col min-h-dvh safe-x">
	<header class="sticky top-0 z-[150] bg-background border-b border-border safe-top">
		<div class="flex items-center justify-between max-w-5xl mx-auto px-4 py-3">
			<a href="/" class="flex items-center gap-2 text-foreground hover:text-primary no-underline transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 rounded-sm">
				<span class="text-lg font-semibold">Apex Memory</span>
			</a>
			<nav class="flex items-center gap-2">
				<ThemeToggle />
				{#if $page.data.user}
					<Button variant="ghost" size="sm" href="/dashboard">Dashboard</Button>
					<Button
						variant="outline"
						size="sm"
						onclick={handleLogout}
						disabled={loggingOut}
					>
						{loggingOut ? 'Signing out...' : 'Sign out'}
					</Button>
				{:else}
					<Button variant="ghost" size="sm" href="/login">Sign in</Button>
					<Button size="sm" href="/register">Sign up</Button>
				{/if}
			</nav>
		</div>
	</header>

	<main id="main-content" class="flex-1 w-full max-w-5xl mx-auto px-4 py-6" tabindex="-1">
		{@render children()}
	</main>

	<footer class="border-t border-border px-4 py-4 safe-bottom text-center">
		<p class="text-sm text-muted-foreground">Apex Memory</p>
	</footer>
</div>

<!-- Toast notifications -->
<Toaster richColors position="top-center" />

<style>
	@keyframes progress {
		0% {
			width: 0%;
			margin-left: 0%;
		}
		50% {
			width: 50%;
			margin-left: 25%;
		}
		100% {
			width: 0%;
			margin-left: 100%;
		}
	}

	.animate-progress {
		animation: progress 1s ease-in-out infinite;
	}

	@media (prefers-reduced-motion: reduce) {
		.animate-progress {
			animation: none;
			width: 100%;
		}
	}

	main:focus {
		outline: none;
	}
</style>
