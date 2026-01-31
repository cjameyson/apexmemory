<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { navigating } from '$app/stores';
	import { theme } from '$lib/stores/theme.svelte';
	import { Toaster } from '$lib/components/ui/sonner';

	let { children } = $props();

	onMount(() => {
		theme.init();
	});
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

{@render children()}

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
</style>
