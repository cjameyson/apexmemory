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
		'flex items-center gap-1 border-none border-border px-2 py-1.5',
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
					? 'bg-primary/10 text-primary'
					: 'text-muted-foreground hover:text-foreground'
			)}
		>
			{tab.label}
		</button>
	{/each}
</div>
