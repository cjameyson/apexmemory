<script lang="ts">
	import type { OcclusionMode } from './types';
	import { formatZoom, formatOcclusionMode } from './utils';

	interface Props {
		regionCount: number;
		mode: OcclusionMode;
		zoom: number;
	}

	let { regionCount, mode, zoom }: Props = $props();
</script>

<div class="flex items-center gap-3 px-3 py-1 text-xs text-muted-foreground">
	<span>
		{regionCount}
		{regionCount === 1 ? 'region' : 'regions'}
	</span>
	<span class="text-border">|</span>
	<span>{formatOcclusionMode(mode)}</span>
	<span class="text-border">|</span>
	<span>{formatZoom(zoom)}</span>
</div>
