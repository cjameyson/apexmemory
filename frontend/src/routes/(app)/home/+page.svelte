<script lang="ts">
	import { getAllNotebooks, getGlobalStats } from '$lib/mocks';
	import { StatsHero, StatCard, WeeklyActivity, NotebookCard } from '$lib/components/dashboard';
	import { LayersIcon, CalendarIcon, TrendingUpIcon, FlameIcon } from '@lucide/svelte';
	import type { Notebook } from '$lib/types';

	// Get mock data
	let notebooks: Notebook[] = $state(getAllNotebooks());
	let stats = $state(getGlobalStats());

	// Calculate totals
	let totalDue = $derived(notebooks.reduce((sum, nb) => sum + nb.dueCount, 0));
</script>

<div class="flex-1 overflow-auto">
	<div class="max-w-6xl mx-auto p-6 lg:p-8">
		<!-- Welcome header -->
		<div class="mb-8">
			<h1 class="text-3xl font-bold text-foreground mb-2">
				Welcome back
			</h1>
			<p class="text-muted-foreground">
				You have {totalDue} cards due for review today.
			</p>
		</div>

		<!-- Stats Hero -->
		<div class="mb-6">
			<StatsHero
				totalDue={stats.upcomingDue.today}
				reviewedToday={stats.cardsReviewedToday}
				currentStreak={stats.currentStreak}
				averageRetention={stats.averageRetention}
			/>
		</div>

		<!-- Stats cards grid -->
		<div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
			<StatCard
				icon={LayersIcon}
				label="Total Cards"
				value={stats.totalCards}
			/>
			<StatCard
				icon={CalendarIcon}
				label="Tomorrow"
				value={stats.upcomingDue.tomorrow}
			/>
			<StatCard
				icon={TrendingUpIcon}
				label="This Week"
				value={stats.upcomingDue.thisWeek}
			/>
			<StatCard
				icon={FlameIcon}
				label="Best Streak"
				value="{stats.longestStreak} days"
			/>
		</div>

		<!-- Weekly activity chart -->
		<div class="mb-6">
			<WeeklyActivity
				data={stats.reviewsThisWeek}
				todayIndex={5}
			/>
		</div>

		<!-- Notebooks grid -->
		<div>
			<h2 class="text-lg font-semibold text-foreground mb-4">Your Notebooks</h2>
			<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
				{#each notebooks as notebook (notebook.id)}
					<NotebookCard {notebook} />
				{/each}
			</div>
		</div>
	</div>
</div>
