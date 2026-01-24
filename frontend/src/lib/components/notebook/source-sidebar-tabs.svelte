<script lang="ts">
	import { cn } from '$lib/utils';

	export type SidebarTab = 'cards' | 'summary' | 'chat';

	interface Props {
		activeTab?: SidebarTab;
		onTabChange?: (tab: SidebarTab) => void;
		class?: string;
	}

	let { activeTab = 'cards', onTabChange, class: className }: Props = $props();

	const tabs: { id: SidebarTab; label: string }[] = [
		{ id: 'cards', label: 'Cards' },
		{ id: 'summary', label: 'Summary' },
		{ id: 'chat', label: 'Chat' }
	];

	function selectTab(tab: SidebarTab) {
		onTabChange?.(tab);
	}
</script>

<div
	class={cn(
		'flex items-center gap-1 border-b border-slate-200 px-2 py-1.5 dark:border-slate-700',
		className
	)}
>
	{#each tabs as tab (tab.id)}
		<button
			type="button"
			onclick={() => selectTab(tab.id)}
			class={cn(
				'rounded-md px-3 py-1.5 text-sm font-medium transition-colors',
				activeTab === tab.id
					? 'bg-sky-100 text-sky-900 dark:bg-sky-900/30 dark:text-sky-100'
					: 'text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white'
			)}
		>
			{tab.label}
		</button>
	{/each}
</div>
