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

	let Icon = $derived(iconMap[type] || FileTextIcon);
</script>

<Icon class={cn('size-4 text-slate-500 dark:text-slate-400', className)} />
