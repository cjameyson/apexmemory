<script lang="ts">
	import { cn } from '$lib/utils';
	import CardItem from './card-item.svelte';
	import type { Source, Card } from '$lib/types';

	type Tab = 'source' | 'cards' | 'summary' | 'chat';

	interface Props {
		source: Source;
		cards: Card[];
		activeTab?: Tab;
		class?: string;
	}

	let { source, cards, activeTab = 'source', class: className }: Props = $props();

	// Filter cards for this source
	let sourceCards = $derived(cards.filter((c) => c.sourceId === source.id));
</script>

<div class={cn('flex-1 overflow-auto', className)}>
	{#if activeTab === 'source'}
		<!-- Source preview -->
		<div class="p-6">
			{#if source.type === 'pdf'}
				<div class="bg-slate-100 dark:bg-slate-800 rounded-lg p-8 text-center">
					<div class="text-slate-500 dark:text-slate-400 mb-2">PDF Preview</div>
					<div class="text-sm text-slate-400 dark:text-slate-500">
						{source.pages} pages
					</div>
				</div>
			{:else if source.type === 'youtube'}
				<div class="aspect-video bg-slate-900 rounded-lg flex items-center justify-center">
					<div class="text-white/50">YouTube Player Placeholder</div>
				</div>
			{:else if source.type === 'audio'}
				<div class="bg-slate-100 dark:bg-slate-800 rounded-lg p-8">
					<div class="text-center text-slate-500 dark:text-slate-400 mb-4">Audio Player</div>
					<div class="h-2 bg-slate-200 dark:bg-slate-700 rounded-full">
						<div class="h-full w-0 bg-sky-500 rounded-full"></div>
					</div>
					<div class="flex justify-between mt-2 text-xs text-slate-400">
						<span>0:00</span>
						<span>{source.duration}</span>
					</div>
				</div>
			{:else if source.type === 'url'}
				<div class="prose dark:prose-invert max-w-none">
					<p>{source.excerpt}</p>
					<p class="text-slate-400 dark:text-slate-500 text-sm">
						[Web content preview would appear here]
					</p>
				</div>
			{:else if source.type === 'notes'}
				<div class="prose dark:prose-invert max-w-none">
					<p>{source.excerpt}</p>
				</div>
			{/if}
		</div>

	{:else if activeTab === 'cards'}
		<!-- Cards from this source -->
		<div class="p-6">
			{#if sourceCards.length > 0}
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
					{#each sourceCards as card (card.id)}
						<CardItem {card} />
					{/each}
				</div>
			{:else}
				<div class="text-center py-12 text-slate-500 dark:text-slate-400">
					No cards generated from this source yet.
				</div>
			{/if}
		</div>

	{:else if activeTab === 'summary'}
		<!-- AI Summary -->
		<div class="p-6">
			<div class="bg-slate-50 dark:bg-slate-800 rounded-lg p-6">
				<h3 class="font-semibold text-slate-900 dark:text-white mb-3">AI Summary</h3>
				<p class="text-slate-600 dark:text-slate-400">
					{source.excerpt}
				</p>
				<p class="text-slate-400 dark:text-slate-500 text-sm mt-4">
					[AI-generated summary would appear here]
				</p>
			</div>
		</div>

	{:else if activeTab === 'chat'}
		<!-- Chat with source -->
		<div class="flex flex-col h-full">
			<div class="flex-1 p-6">
				<div class="text-center py-12 text-slate-500 dark:text-slate-400">
					Ask questions about this source...
				</div>
			</div>
			<div class="p-4 border-t border-slate-200 dark:border-slate-700">
				<input
					type="text"
					placeholder="Ask a question about this source..."
					class="w-full px-4 py-2 bg-slate-100 dark:bg-slate-800 rounded-lg border-0 focus:ring-2 focus:ring-sky-500 text-slate-900 dark:text-white placeholder:text-slate-400"
				/>
			</div>
		</div>
	{/if}
</div>
