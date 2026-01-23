<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { BrainIcon, ZapIcon, BookOpenIcon, SparklesIcon } from '@lucide/svelte';

	// If user is logged in, this page won't be shown (see +page.server.ts)
</script>

<!-- Landing page for unauthenticated users -->
<div class="min-h-screen bg-gradient-to-br from-slate-50 to-sky-50 dark:from-slate-950 dark:to-slate-900">
	<!-- Header -->
	<header class="px-6 py-4">
		<div class="max-w-6xl mx-auto flex items-center justify-between">
			<div class="flex items-center gap-2">
				<BrainIcon class="size-8 text-sky-500" />
				<span class="text-xl font-bold text-slate-900 dark:text-white">Apex Memory</span>
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
			<h1 class="text-5xl font-bold text-slate-900 dark:text-white mb-6">
				Remember everything you learn
			</h1>
			<p class="text-xl text-slate-600 dark:text-slate-400 mb-10 max-w-2xl mx-auto">
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
			<div class="bg-white dark:bg-slate-900 rounded-2xl p-6 shadow-sm border border-slate-200 dark:border-slate-800">
				<div class="w-12 h-12 bg-sky-100 dark:bg-sky-900/30 rounded-xl flex items-center justify-center mb-4">
					<SparklesIcon class="size-6 text-sky-600 dark:text-sky-400" />
				</div>
				<h3 class="text-lg font-semibold text-slate-900 dark:text-white mb-2">
					AI Card Generation
				</h3>
				<p class="text-slate-600 dark:text-slate-400">
					Upload PDFs, videos, or notes. Our AI creates high-quality flashcards with citations.
				</p>
			</div>

			<div class="bg-white dark:bg-slate-900 rounded-2xl p-6 shadow-sm border border-slate-200 dark:border-slate-800">
				<div class="w-12 h-12 bg-emerald-100 dark:bg-emerald-900/30 rounded-xl flex items-center justify-center mb-4">
					<BookOpenIcon class="size-6 text-emerald-600 dark:text-emerald-400" />
				</div>
				<h3 class="text-lg font-semibold text-slate-900 dark:text-white mb-2">
					FSRS Algorithm
				</h3>
				<p class="text-slate-600 dark:text-slate-400">
					The latest evidence-based scheduling algorithm optimizes your review intervals.
				</p>
			</div>

			<div class="bg-white dark:bg-slate-900 rounded-2xl p-6 shadow-sm border border-slate-200 dark:border-slate-800">
				<div class="w-12 h-12 bg-amber-100 dark:bg-amber-900/30 rounded-xl flex items-center justify-center mb-4">
					<ZapIcon class="size-6 text-amber-600 dark:text-amber-400" />
				</div>
				<h3 class="text-lg font-semibold text-slate-900 dark:text-white mb-2">
					Focus Mode
				</h3>
				<p class="text-slate-600 dark:text-slate-400">
					Distraction-free review sessions with keyboard shortcuts for maximum efficiency.
				</p>
			</div>
		</div>
	</main>
</div>
