<script lang="ts">
	import { TagIcon, PauseIcon, Trash2Icon, XIcon } from '@lucide/svelte';

	let {
		selectedCount,
		onClear
	}: {
		selectedCount: number;
		onClear: () => void;
	} = $props();
</script>

<div class="flex items-center justify-between border-b border-primary/10 bg-primary/5 px-4 py-2">
	<div class="flex items-center gap-2 text-sm">
		<span class="font-medium text-foreground">{selectedCount} selected</span>
		<button
			onclick={onClear}
			class="inline-flex items-center gap-1 rounded px-2 py-0.5 text-xs text-muted-foreground hover:bg-accent hover:text-foreground"
		>
			<XIcon class="h-3 w-3" />
			Clear
		</button>
	</div>
	<div class="flex items-center gap-1">
		<button class="inline-flex items-center gap-1.5 rounded px-3 py-1.5 text-sm text-muted-foreground hover:bg-accent hover:text-foreground">
			<TagIcon class="h-4 w-4" />
			Add Tags
		</button>
		<button class="inline-flex items-center gap-1.5 rounded px-3 py-1.5 text-sm text-muted-foreground hover:bg-accent hover:text-foreground">
			<PauseIcon class="h-4 w-4" />
			Suspend
		</button>
		<button class="inline-flex items-center gap-1.5 rounded px-3 py-1.5 text-sm text-destructive hover:bg-destructive/10">
			<Trash2Icon class="h-4 w-4" />
			Delete
		</button>
	</div>
</div>
