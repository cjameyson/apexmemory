<script lang="ts">
	import type { OcclusionMode } from './types';
	import { formatZoom, formatOcclusionMode } from './utils';
	import { Undo2, Redo2, Circle } from '@lucide/svelte';

	interface Props {
		regionCount: number;
		mode: OcclusionMode;
		zoom: number;
		undoCount: number;
		redoCount: number;
		isDirty: boolean;
	}

	let { regionCount, mode, zoom, undoCount, redoCount, isDirty }: Props = $props();
</script>

<div
	class="flex items-center justify-between border-t border-border bg-muted/50 px-3 py-1 text-xs text-muted-foreground"
>
	<!-- Left section: Region count and mode -->
	<div class="flex items-center gap-4">
		<span>
			{regionCount}
			{regionCount === 1 ? 'region' : 'regions'}
		</span>
		<span class="text-border">|</span>
		<span>{formatOcclusionMode(mode)}</span>
	</div>

	<!-- Right section: Zoom and history info -->
	<div class="flex items-center gap-4">
		<!-- Dirty indicator -->
		{#if isDirty}
			<span class="flex items-center gap-1 text-amber-500" title="Unsaved changes">
				<Circle class="h-2 w-2 fill-current" />
				Modified
			</span>
		{/if}

		<!-- Undo/Redo counts -->
		<span class="flex items-center gap-2" title="Undo/Redo history">
			<span class="flex items-center gap-0.5">
				<Undo2 class="h-3 w-3" />
				{undoCount}
			</span>
			<span class="flex items-center gap-0.5">
				<Redo2 class="h-3 w-3" />
				{redoCount}
			</span>
		</span>

		<span class="text-border">|</span>

		<!-- Zoom level -->
		<span>{formatZoom(zoom)}</span>
	</div>
</div>
