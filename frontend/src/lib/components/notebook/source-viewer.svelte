<script lang="ts">
	import { cn } from '$lib/utils';
	import SourceToolbar from './source-toolbar.svelte';
	import type { Source } from '$lib/types';

	interface Props {
		source: Source;
		onSectionChange?: (section: string) => void;
		onGenerateCards?: () => void;
		class?: string;
	}

	let { source, onSectionChange, onGenerateCards, class: className }: Props = $props();
</script>

<div class={cn('flex flex-1 flex-col overflow-hidden', className)}>
	<!-- Toolbar -->
	<SourceToolbar type={source.type} {onGenerateCards} />

	<!-- Source content -->
	<div class="flex-1 overflow-auto p-6">
		{#if source.type === 'pdf'}
			<div class="rounded-lg bg-muted p-8 text-center">
				<div class="mb-2 text-muted-foreground">PDF Preview</div>
				<div class="text-sm text-muted-foreground">
					{source.pages} pages
				</div>
			</div>
		{:else if source.type === 'youtube'}
			<div class="flex aspect-video items-center justify-center rounded-lg bg-foreground">
				<div class="text-background/50">YouTube Player Placeholder</div>
			</div>
		{:else if source.type === 'audio'}
			<div class="rounded-lg bg-muted p-8">
				<div class="mb-4 text-center text-muted-foreground">Audio Player</div>
				<div class="h-2 rounded-full bg-border">
					<div class="h-full w-0 rounded-full bg-primary"></div>
				</div>
				<div class="mt-2 flex justify-between text-xs text-muted-foreground">
					<span>0:00</span>
					<span>{source.duration}</span>
				</div>
			</div>
		{:else if source.type === 'url'}
			<div class="prose max-w-none dark:prose-invert">
				<p>{source.excerpt}</p>
				<p class="text-sm text-muted-foreground">
					[Web content preview would appear here]
				</p>
			</div>
		{:else if source.type === 'notes'}
			<div class="prose max-w-none dark:prose-invert">
				<p>{source.excerpt}</p>
			</div>
		{/if}
	</div>
</div>
