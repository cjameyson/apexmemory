<script lang="ts">
	import { cn } from '$lib/utils';

	type Tab = 'source' | 'cards' | 'summary' | 'chat';

	interface Props {
		activeTab?: Tab;
		onTabChange?: (tab: Tab) => void;
		class?: string;
	}

	let { activeTab = 'source', onTabChange, class: className }: Props = $props();

	const tabs: { id: Tab; label: string; icon: string }[] = [
		{ id: 'source', label: 'Source', icon: 'file' },
		{ id: 'cards', label: 'Cards', icon: 'cards' },
		{ id: 'summary', label: 'Summary', icon: 'summary' },
		{ id: 'chat', label: 'Chat', icon: 'chat' }
	];

	function selectTab(tab: Tab) {
		onTabChange?.(tab);
	}
</script>

<div class={cn('flex items-center gap-1 rounded-lg bg-white p-1 dark:bg-slate-800', className)}>
	{#each tabs as tab (tab.id)}
		<button
			type="button"
			onclick={() => selectTab(tab.id)}
			class={cn(
				'rounded-md px-3 py-1.5 text-sm font-medium transition-colors',
				activeTab === tab.id
					? 'bg-slate-100 text-slate-900 shadow-none dark:bg-slate-700 dark:text-white'
					: 'text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white'
			)}
		>
			{tab.label}
		</button>
	{/each}
</div>
