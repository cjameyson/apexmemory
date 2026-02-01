<script lang="ts">
	import { LayersIcon, TypeIcon, ImageIcon } from '@lucide/svelte';
	import type { FactType } from '$lib/types/fact';

	let { factType }: { factType: FactType } = $props();

	const config = $derived(
		factType === 'cloze'
			? { icon: TypeIcon, label: 'Cloze', classes: 'bg-cloze/15 text-cloze' }
			: factType === 'image_occlusion'
				? { icon: ImageIcon, label: 'Image', classes: 'bg-warning/15 text-warning' }
				: { icon: LayersIcon, label: 'Basic', classes: 'bg-muted text-muted-foreground' }
	);
</script>

<span class="inline-flex items-center gap-1 rounded px-2 py-0.5 text-xs font-medium {config.classes}">
	<config.icon class="h-3 w-3" />
	{config.label}
</span>
