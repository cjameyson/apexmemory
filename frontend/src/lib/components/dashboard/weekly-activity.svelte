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
		'bg-card rounded-2xl p-5 border border-border',
		className
	)}
>
	<h2 class="text-lg font-semibold text-foreground mb-4">Weekly Activity</h2>
	<div class="flex items-end gap-2 h-24">
		{#each data as count, i}
			<div class="flex-1 flex flex-col items-center gap-1">
				<div class="w-full flex-1 flex flex-col justify-end">
					<div
						class={cn(
							'w-full rounded-t-md transition-all duration-300',
							i === todayIndex
								? 'bg-primary'
								: count > 0
									? 'bg-primary/50'
									: 'bg-border'
						)}
						style="height: {Math.max((count / maxValue) * 100, 4)}%"
					></div>
				</div>
				<div
					class={cn(
						'text-xs',
						i === todayIndex
							? 'text-primary font-medium'
							: 'text-muted-foreground'
					)}
				>
					{labels[i]}
				</div>
			</div>
		{/each}
	</div>
</div>
