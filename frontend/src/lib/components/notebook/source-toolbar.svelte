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
		'flex items-center justify-between gap-2 border-b border-border bg-muted px-4 py-2',
		className
	)}
>
	<!-- Type-specific controls -->
	<div class="flex items-center gap-1">
		{#if type === 'pdf'}
			<Button variant="ghost" size="icon" class="size-8" title="Zoom out">
				<ZoomOutIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" title="Zoom in">
				<ZoomInIcon class="size-4" />
			</Button>
			<div class="mx-1 h-4 w-px bg-border"></div>
			<Button variant="ghost" size="icon" class="size-8" title="Previous page">
				<ChevronLeftIcon class="size-4" />
			</Button>
			<span class="px-2 text-sm text-muted-foreground">1 / 87</span>
			<Button variant="ghost" size="icon" class="size-8" title="Next page">
				<ChevronRightIcon class="size-4" />
			</Button>
			<div class="mx-1 h-4 w-px bg-border"></div>
			<Button variant="ghost" size="icon" class="size-8" title="Search">
				<SearchIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" title="Table of contents">
				<ListIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" title="Download">
				<DownloadIcon class="size-4" />
			</Button>
		{:else if type === 'youtube'}
			<Button variant="ghost" size="icon" class="size-8" title="Play">
				<PlayIcon class="size-4" />
			</Button>
			<span class="px-2 text-sm text-muted-foreground">0:00 / 14:32</span>
			<Button variant="ghost" size="icon" class="size-8" title="Volume">
				<Volume2Icon class="size-4" />
			</Button>
			<div class="mx-1 h-4 w-px bg-border"></div>
			<Button variant="ghost" size="sm" class="h-7 text-xs">1x</Button>
			<Button variant="ghost" size="icon" class="size-8" title="Bookmark">
				<BookmarkIcon class="size-4" />
			</Button>
		{:else if type === 'audio'}
			<Button variant="ghost" size="icon" class="size-8" title="Play">
				<PlayIcon class="size-4" />
			</Button>
			<span class="px-2 text-sm text-muted-foreground">0:00 / 28:15</span>
			<Button variant="ghost" size="icon" class="size-8" title="Volume">
				<Volume2Icon class="size-4" />
			</Button>
			<div class="mx-1 h-4 w-px bg-border"></div>
			<Button variant="ghost" size="sm" class="h-7 text-xs">1x</Button>
		{:else if type === 'url'}
			<Button variant="ghost" size="icon" class="size-8" title="Open original">
				<ExternalLinkIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" title="Refresh">
				<RefreshCwIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" title="Reader mode">
				<BookOpenIcon class="size-4" />
			</Button>
		{:else if type === 'notes'}
			<Button variant="ghost" size="icon" class="size-8" title="Edit">
				<EditIcon class="size-4" />
			</Button>
			<Button variant="ghost" size="icon" class="size-8" title="Download">
				<DownloadIcon class="size-4" />
			</Button>
		{/if}
	</div>

	<!-- Generate cards button -->
	<Button
		variant="default"
		size="sm"
		class="hidden gap-1.5"
		onclick={onGenerateCards}
	>
		<SparklesIcon class="size-4" />
		<span class="sm:inline">Generate Cards</span>
	</Button>
</div>
