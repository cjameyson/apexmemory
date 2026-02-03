<script lang="ts">
	import { cn } from '$lib/utils';

	type Rating = 1 | 2 | 3 | 4;

	interface Props {
		onRate?: (rating: Rating) => void;
		intervals?: {
			again: string;
			hard: string;
			good: string;
			easy: string;
		};
		disabled?: boolean;
		class?: string;
	}

	let {
		onRate,
		intervals = { again: '<1m', hard: '<10m', good: '1d', easy: '4d' },
		disabled = false,
		class: className
	}: Props = $props();

	const buttons: { rating: Rating; label: string; key: string; color: string }[] = [
		{ rating: 1, label: 'Again', key: '1', color: 'bg-again hover:bg-again/90 focus:ring-again/50' },
		{ rating: 2, label: 'Hard', key: '2', color: 'bg-hard hover:bg-hard/90 focus:ring-hard/50' },
		{ rating: 3, label: 'Good', key: '3', color: 'bg-good hover:bg-good/90 focus:ring-good/50' },
		{ rating: 4, label: 'Easy', key: '4', color: 'bg-easy hover:bg-easy/90 focus:ring-easy/50' }
	];

	function getInterval(rating: Rating): string {
		switch (rating) {
			case 1: return intervals.again;
			case 2: return intervals.hard;
			case 3: return intervals.good;
			case 4: return intervals.easy;
		}
	}

	function handleRate(rating: Rating) {
		if (!disabled) {
			onRate?.(rating);
		}
	}
</script>

<div class={cn('flex items-center gap-3', className)}>
	{#each buttons as btn (btn.rating)}
		<div class="flex-1 flex flex-col items-center gap-2">
			<button
				type="button"
				onclick={() => handleRate(btn.rating)}
				{disabled}
				class={cn(
					'w-full flex flex-col items-center gap-1 px-4 py-3 rounded-xl text-white font-medium transition-all focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-900',
					btn.color,
					disabled && 'opacity-50 cursor-not-allowed'
				)}
			>
				<span class="text-sm">{btn.label}</span>
				<span class="text-xs opacity-80">{getInterval(btn.rating)}</span>
			</button>
			<kbd class="px-1.5 py-0.5 bg-white/10 rounded text-xs text-white/50">{btn.key}</kbd>
		</div>
	{/each}
</div>
