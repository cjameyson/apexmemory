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
			<div class="rounded-lg bg-slate-100 p-8 text-center dark:bg-slate-800">
				<div class="mb-2 text-slate-500 dark:text-slate-400">PDF Preview</div>
				<div class="text-sm text-slate-400 dark:text-slate-500">
					{source.pages} pages
				</div>
			</div>
		{:else if source.type === 'youtube'}
			<div class="flex aspect-video items-center justify-center rounded-lg bg-slate-900">
				<div class="text-white/50">YouTube Player Placeholder</div>
			</div>
		{:else if source.type === 'audio'}
			<div class="rounded-lg bg-slate-100 p-8 dark:bg-slate-800">
				<div class="mb-4 text-center text-slate-500 dark:text-slate-400">Audio Player</div>
				<div class="h-2 rounded-full bg-slate-200 dark:bg-slate-700">
					<div class="h-full w-0 rounded-full bg-sky-500"></div>
				</div>
				<div class="mt-2 flex justify-between text-xs text-slate-400">
					<span>0:00</span>
					<span>{source.duration}</span>
				</div>
			</div>
		{:else if source.type === 'url'}
			<div class="prose max-w-none dark:prose-invert">
				<p>{source.excerpt}</p>
				<p class="text-sm text-slate-400 dark:text-slate-500">
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
