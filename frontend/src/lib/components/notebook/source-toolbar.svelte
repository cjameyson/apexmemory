<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
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

<Tooltip.Provider delayDuration={300}>
	<div
		class={cn(
			'flex items-center justify-between gap-2 border-b border-slate-200 bg-slate-50 px-4 py-2 dark:border-slate-700 dark:bg-slate-800',
			className
		)}
	>
		<!-- Type-specific controls -->
		<div class="flex items-center gap-1">
			{#if type === 'pdf'}
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<ZoomOutIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Zoom out</Tooltip.Content>
				</Tooltip.Root>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<ZoomInIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Zoom in</Tooltip.Content>
				</Tooltip.Root>
				<div class="mx-1 h-4 w-px bg-slate-300 dark:bg-slate-600"></div>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<ChevronLeftIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Previous page</Tooltip.Content>
				</Tooltip.Root>
				<span class="px-2 text-sm text-slate-500 dark:text-slate-400">1 / 87</span>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<ChevronRightIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Next page</Tooltip.Content>
				</Tooltip.Root>
				<div class="mx-1 h-4 w-px bg-slate-300 dark:bg-slate-600"></div>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<SearchIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Search</Tooltip.Content>
				</Tooltip.Root>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<ListIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Table of contents</Tooltip.Content>
				</Tooltip.Root>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<DownloadIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Download</Tooltip.Content>
				</Tooltip.Root>
			{:else if type === 'youtube'}
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<PlayIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Play</Tooltip.Content>
				</Tooltip.Root>
				<span class="px-2 text-sm text-slate-500 dark:text-slate-400">0:00 / 14:32</span>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<Volume2Icon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Volume</Tooltip.Content>
				</Tooltip.Root>
				<div class="mx-1 h-4 w-px bg-slate-300 dark:bg-slate-600"></div>
				<Button variant="ghost" size="sm" class="h-7 text-xs">1x</Button>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<BookmarkIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Bookmark</Tooltip.Content>
				</Tooltip.Root>
			{:else if type === 'audio'}
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<PlayIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Play</Tooltip.Content>
				</Tooltip.Root>
				<span class="px-2 text-sm text-slate-500 dark:text-slate-400">0:00 / 28:15</span>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<Volume2Icon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Volume</Tooltip.Content>
				</Tooltip.Root>
				<div class="mx-1 h-4 w-px bg-slate-300 dark:bg-slate-600"></div>
				<Button variant="ghost" size="sm" class="h-7 text-xs">1x</Button>
			{:else if type === 'url'}
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<ExternalLinkIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Open original</Tooltip.Content>
				</Tooltip.Root>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<RefreshCwIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Refresh</Tooltip.Content>
				</Tooltip.Root>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<BookOpenIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Reader mode</Tooltip.Content>
				</Tooltip.Root>
			{:else if type === 'notes'}
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<EditIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Edit</Tooltip.Content>
				</Tooltip.Root>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon" class="size-8" {...props}>
								<DownloadIcon class="size-4" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Download</Tooltip.Content>
				</Tooltip.Root>
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
</Tooltip.Provider>
