<script lang="ts">
	import { cn } from '$lib/utils';

	type Tab = 'source' | 'cards' | 'summary' | 'chat';

	interface Props {
		activeTab?: Tab;
		onTabChange?: (tab: Tab) => void;
		class?: string;
	}

	let { activeTab = 'source', onTabChange, class: className }: Props = $props();

	const tabs: { id: Tab; label: string }[] = [
		{ id: 'source', label: 'Source' },
		{ id: 'cards', label: 'Cards' },
		{ id: 'summary', label: 'Summary' },
		{ id: 'chat', label: 'Chat' }
	];

	function selectTab(tab: Tab) {
		onTabChange?.(tab);
	}
</script>

<div class={cn('flex items-center gap-1 border-b border-slate-200 dark:border-slate-700', className)}>
	{#each tabs as tab (tab.id)}
		<button
			type="button"
			onclick={() => selectTab(tab.id)}
			class={cn(
				'px-4 py-2 text-sm font-medium transition-colors relative',
				activeTab === tab.id
					? 'text-sky-600 dark:text-sky-400'
					: 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'
			)}
		>
			{tab.label}
			{#if activeTab === tab.id}
				<div class="absolute bottom-0 left-0 right-0 h-0.5 bg-sky-500"></div>
			{/if}
		</button>
	{/each}
</div>
