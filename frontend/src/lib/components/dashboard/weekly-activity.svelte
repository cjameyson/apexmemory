<script lang="ts">
	import { cn } from '$lib/utils';

	interface Props {
		data: number[]; // 7 days
		labels?: string[];
		todayIndex?: number; // Which day is today (0-6)
		class?: string;
	}

	let {
		data,
		labels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
		todayIndex = new Date().getDay() === 0 ? 6 : new Date().getDay() - 1, // Convert Sunday=0 to index 6
		class: className
	}: Props = $props();

	let maxValue = $derived(Math.max(...data, 1));
</script>

<div
	class={cn(
		'bg-white dark:bg-slate-900 rounded-2xl p-5 border border-slate-200 dark:border-slate-800',
		className
	)}
>
	<h2 class="text-lg font-semibold text-slate-900 dark:text-white mb-4">Weekly Activity</h2>
	<div class="flex items-end gap-2 h-24">
		{#each data as count, i}
			<div class="flex-1 flex flex-col items-center gap-1">
				<div class="w-full flex-1 flex flex-col justify-end">
					<div
						class={cn(
							'w-full rounded-t-md transition-all duration-300',
							i === todayIndex
								? 'bg-sky-500'
								: count > 0
									? 'bg-sky-300 dark:bg-sky-700'
									: 'bg-slate-200 dark:bg-slate-700'
						)}
						style="height: {Math.max((count / maxValue) * 100, 4)}%"
					></div>
				</div>
				<div
					class={cn(
						'text-xs',
						i === todayIndex
							? 'text-sky-600 dark:text-sky-400 font-medium'
							: 'text-slate-500 dark:text-slate-400'
					)}
				>
					{labels[i]}
				</div>
			</div>
		{/each}
	</div>
</div>
