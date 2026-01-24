<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import {
		ZoomInIcon,
		ZoomOutIcon,
		ChevronLeftIcon,
		ChevronRightIcon,
		SearchIcon,
		ListIcon,
		DownloadIcon,
		PlayIcon,
		PauseIcon,
		Volume2Icon,
		BookmarkIcon,
		ExternalLinkIcon,
		RefreshCwIcon,
		BookOpenIcon,
		EditIcon,
		SparklesIcon
	} from '@lucide/svelte';
	import type { SourceType } from '$lib/types';

	interface Props {
		type: SourceType;
		onGenerateCards?: () => void;
		class?: string;
	}

	let { type, onGenerateCards, class: className }: Props = $props();
</script>

<div
	class={cn(
		'flex items-center justify-between gap-2 border-b border-slate-200 bg-slate-50 px-4 py-2 dark:border-slate-700 dark:bg-slate-800',
		className
	)}
>
	<!-- Type-specific controls -->
	<div class="flex items-center gap-1">
		{#if type === 'pdf'}
			<Button variant="ghost" size="icon" class="size-8" aria-label="Zoom out">
				<ZoomOutIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Zoom in">
				<ZoomInIcon class="size-4" />
			</Button>
			<div class="mx-1 h-4 w-px bg-slate-300 dark:bg-slate-600"></div>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Previous page">
				<ChevronLeftIcon class="size-4" />
			</Button>
			<span class="px-2 text-sm text-slate-500 dark:text-slate-400">1 / 87</span>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Next page">
				<ChevronRightIcon class="size-4" />
			</Button>
			<div class="mx-1 h-4 w-px bg-slate-300 dark:bg-slate-600"></div>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Search">
				<SearchIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Table of contents">
				<ListIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Download">
				<DownloadIcon class="size-4" />
			</Button>
		{:else if type === 'youtube'}
			<Button variant="ghost" size="icon" class="size-8" aria-label="Play">
				<PlayIcon class="size-4" />
			</Button>
			<span class="px-2 text-sm text-slate-500 dark:text-slate-400">0:00 / 14:32</span>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Volume">
				<Volume2Icon class="size-4" />
			</Button>
			<div class="mx-1 h-4 w-px bg-slate-300 dark:bg-slate-600"></div>
			<Button variant="ghost" size="sm" class="h-7 text-xs">1x</Button>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Bookmark">
				<BookmarkIcon class="size-4" />
			</Button>
		{:else if type === 'audio'}
			<Button variant="ghost" size="icon" class="size-8" aria-label="Play">
				<PlayIcon class="size-4" />
			</Button>
			<span class="px-2 text-sm text-slate-500 dark:text-slate-400">0:00 / 28:15</span>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Volume">
				<Volume2Icon class="size-4" />
			</Button>
			<div class="mx-1 h-4 w-px bg-slate-300 dark:bg-slate-600"></div>
			<Button variant="ghost" size="sm" class="h-7 text-xs">1x</Button>
		{:else if type === 'url'}
			<Button variant="ghost" size="icon" class="size-8" aria-label="Open original">
				<ExternalLinkIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Refresh">
				<RefreshCwIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Reader mode">
				<BookOpenIcon class="size-4" />
			</Button>
		{:else if type === 'notes'}
			<Button variant="ghost" size="icon" class="size-8" aria-label="Edit">
				<EditIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" aria-label="Download">
				<DownloadIcon class="size-4" />
			</Button>
		{/if}
	</div>

	<!-- Generate cards button -->
	<Button
		variant="default"
		size="sm"
		class="hidden gap-1.5 bg-sky-500 hover:bg-sky-600"
		onclick={onGenerateCards}
	>
		<SparklesIcon class="size-4" />
		<span class="sm:inline">Generate Cards</span>
	</Button>
</div>
