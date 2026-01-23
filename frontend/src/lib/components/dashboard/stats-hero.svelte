<script lang="ts">
	import ProgressRing from '$lib/components/ui/progress-ring.svelte';

	interface Props {
		totalDue: number;
		reviewedToday: number;
		currentStreak: number;
		averageRetention: number;
	}

	let { totalDue, reviewedToday, currentStreak, averageRetention }: Props = $props();
</script>

<div class="bg-gradient-to-br from-sky-500 to-cyan-600 rounded-2xl p-6 text-white">
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 items-center">
		<!-- Due count -->
		<div>
			<div class="text-sm font-medium text-white/80 mb-1">Cards due today</div>
			<div class="text-4xl font-bold">{totalDue}</div>
			<div class="text-sm text-white/70 mt-1">
				{reviewedToday} reviewed so far
			</div>
		</div>

		<!-- Streak -->
		<div class="text-center">
			<div class="text-4xl mb-1">🔥</div>
			<div class="text-2xl font-bold">{currentStreak}</div>
			<div class="text-sm text-white/70">day streak</div>
		</div>

		<!-- Retention ring -->
		<div class="flex flex-col items-center">
			<div class="relative">
				<ProgressRing
					progress={averageRetention}
					size={80}
					stroke={4}
					trackColor="white"
					progressColor="white"
				/>
				<div class="absolute inset-0 flex items-center justify-center">
					<span class="text-lg font-bold">{averageRetention}%</span>
				</div>
			</div>
			<div class="text-sm text-white/70 mt-2">retention</div>
		</div>
	</div>
</div>
