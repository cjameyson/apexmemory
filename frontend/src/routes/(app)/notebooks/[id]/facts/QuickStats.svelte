<script lang="ts">
	import type { FactStats } from '$lib/types/fact';
	import { LayersIcon, TypeIcon, ImageIcon, ClockIcon } from '@lucide/svelte';

	let { stats }: { stats: FactStats } = $props();
</script>

<div class="flex flex-wrap gap-6 px-6 py-3">
	<div class="flex items-center gap-3">
		<div class="flex h-8 w-8 items-center justify-center rounded-md bg-muted">
			<LayersIcon class="h-4 w-4 text-muted-foreground" />
		</div>
		<div>
			<p class="text-xs text-muted-foreground">Basic</p>
			<p class="text-sm font-semibold">{stats.byType.basic}</p>
		</div>
	</div>

	<div class="flex items-center gap-3">
		<div class="flex h-8 w-8 items-center justify-center rounded-md bg-cloze/15">
			<TypeIcon class="h-4 w-4 text-cloze" />
		</div>
		<div>
			<p class="text-xs text-muted-foreground">Cloze</p>
			<p class="text-sm font-semibold">{stats.byType.cloze}</p>
		</div>
	</div>

	<div class="flex items-center gap-3">
		<div class="flex h-8 w-8 items-center justify-center rounded-md bg-warning/15">
			<ImageIcon class="h-4 w-4 text-warning" />
		</div>
		<div>
			<p class="text-xs text-muted-foreground">Image</p>
			<p class="text-sm font-semibold">{stats.byType.imageOcclusion}</p>
		</div>
	</div>

	<div class="flex items-center gap-3">
		<div class="flex h-8 w-8 items-center justify-center rounded-md bg-info/15">
			<ClockIcon class="h-4 w-4 text-info" />
		</div>
		<div>
			<p class="text-xs text-muted-foreground">Due Today</p>
			<p class="text-sm font-semibold text-info">{stats.totalDue}</p>
		</div>
	</div>
</div>
