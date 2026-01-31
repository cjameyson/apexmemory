<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { BrainIcon, ZapIcon, BookOpenIcon, SparklesIcon } from '@lucide/svelte';
</script>

<!-- Landing page for unauthenticated users -->
<div class="min-h-screen bg-gradient-to-br from-accent to-primary/5 dark:from-background dark:to-background">
	<!-- Header -->
	<header class="px-6 py-4">
		<div class="max-w-6xl mx-auto flex items-center justify-between">
			<div class="flex items-center gap-2">
				<BrainIcon class="size-8 text-primary" />
				<span class="text-xl font-bold text-foreground">Apex Memory</span>
			</div>
			<div class="flex items-center gap-2">
				<Button variant="ghost" href="/login">Sign in</Button>
				<Button href="/register">Get Started</Button>
			</div>
		</div>
	</header>

	<!-- Hero -->
	<main class="px-6 py-20">
		<div class="max-w-4xl mx-auto text-center">
			<h1 class="text-5xl font-bold text-foreground mb-6">
				Remember everything you learn
			</h1>
			<p class="text-xl text-muted-foreground mb-10 max-w-2xl mx-auto">
				Apex Memory uses AI-enhanced flashcards and proven spaced repetition algorithms
				to help you retain knowledge for life.
			</p>
			<div class="flex items-center justify-center gap-4">
				<Button size="lg" href="/register" class="gap-2">
					<ZapIcon class="size-5" />
					Start learning free
				</Button>
				<Button size="lg" variant="outline" href="/demo">
					View demo
				</Button>
			</div>
		</div>

		<!-- Features -->
		<div class="max-w-5xl mx-auto mt-24 grid md:grid-cols-3 gap-8">
			<div class="bg-card rounded-2xl p-6 shadow-sm border border-border">
				<div class="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center mb-4">
					<SparklesIcon class="size-6 text-primary" />
				</div>
				<h3 class="text-lg font-semibold text-foreground mb-2">
					AI Card Generation
				</h3>
				<p class="text-muted-foreground">
					Upload PDFs, videos, or notes. Our AI creates high-quality flashcards with citations.
				</p>
			</div>

			<div class="bg-card rounded-2xl p-6 shadow-sm border border-border">
				<div class="w-12 h-12 bg-success/10 rounded-xl flex items-center justify-center mb-4">
					<BookOpenIcon class="size-6 text-success" />
				</div>
				<h3 class="text-lg font-semibold text-foreground mb-2">
					FSRS Algorithm
				</h3>
				<p class="text-muted-foreground">
					The latest evidence-based scheduling algorithm optimizes your review intervals.
				</p>
			</div>

			<div class="bg-card rounded-2xl p-6 shadow-sm border border-border">
				<div class="w-12 h-12 bg-warning/10 rounded-xl flex items-center justify-center mb-4">
					<ZapIcon class="size-6 text-warning" />
				</div>
				<h3 class="text-lg font-semibold text-foreground mb-2">
					Focus Mode
				</h3>
				<p class="text-muted-foreground">
					Distraction-free review sessions with keyboard shortcuts for maximum efficiency.
				</p>
			</div>
		</div>
	</main>
</div>
