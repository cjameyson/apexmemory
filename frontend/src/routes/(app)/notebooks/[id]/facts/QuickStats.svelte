<script lang="ts">
	import type { FactStats } from '$lib/types/fact';
	import { LayersIcon, TypeIcon, ImageIcon, ClockIcon } from '@lucide/svelte';

	let { stats }: { stats: FactStats } = $props();
</script>

<div class="flex hidden flex-wrap gap-6 px-6 py-3">
	<div class="flex items-center gap-3">
		<div class="bg-muted flex h-8 w-8 items-center justify-center rounded-md">
			<LayersIcon class="text-muted-foreground h-4 w-4" />
		</div>
		<div>
			<p class="text-muted-foreground text-xs">Basic</p>
			<p class="text-sm font-semibold">{stats.byType.basic}</p>
		</div>
	</div>

	<div class="flex items-center gap-3">
		<div class="bg-cloze/15 flex h-8 w-8 items-center justify-center rounded-md">
			<TypeIcon class="text-cloze h-4 w-4" />
		</div>
		<div>
			<p class="text-muted-foreground text-xs">Cloze</p>
			<p class="text-sm font-semibold">{stats.byType.cloze}</p>
		</div>
	</div>

	<div class="flex items-center gap-3">
		<div class="bg-warning/15 flex h-8 w-8 items-center justify-center rounded-md">
			<ImageIcon class="text-warning h-4 w-4" />
		</div>
		<div>
			<p class="text-muted-foreground text-xs">Image</p>
			<p class="text-sm font-semibold">{stats.byType.imageOcclusion}</p>
		</div>
	</div>

	<div class="flex items-center gap-3">
		<div class="bg-info/15 flex h-8 w-8 items-center justify-center rounded-md">
			<ClockIcon class="text-info h-4 w-4" />
		</div>
		<div>
			<p class="text-muted-foreground text-xs">Due Today</p>
			<p class="text-info text-sm font-semibold">{stats.totalDue}</p>
		</div>
	</div>
</div>
