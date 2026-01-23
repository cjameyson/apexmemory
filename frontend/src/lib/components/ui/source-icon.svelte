<script lang="ts">
	import { cn } from '$lib/utils';
	import { FileTextIcon, YoutubeIcon, LinkIcon, MicIcon, StickyNoteIcon } from '@lucide/svelte';
	import type { SourceType } from '$lib/types';

	interface Props {
		type: SourceType;
		class?: string;
	}

	let { type, class: className }: Props = $props();

	const iconMap = {
		pdf: FileTextIcon,
		youtube: YoutubeIcon,
		url: LinkIcon,
		audio: MicIcon,
		notes: StickyNoteIcon
	};

	const colorMap = {
		pdf: 'text-red-500',
		youtube: 'text-red-600',
		url: 'text-blue-500',
		audio: 'text-purple-500',
		notes: 'text-amber-500'
	};

	let Icon = $derived(iconMap[type] || FileTextIcon);
	let colorClass = $derived(colorMap[type] || 'text-slate-500');
</script>

<Icon class={cn('size-4', colorClass, className)} />
