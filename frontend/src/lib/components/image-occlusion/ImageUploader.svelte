<script lang="ts">
	/**
	 * Image Uploader - Placeholder for Phase 2
	 *
	 * This component will handle:
	 * - Drag and drop file upload
	 * - Click to browse
	 * - Paste from clipboard
	 * - URL input
	 */

	import { Upload, Image, Link } from '@lucide/svelte';

	interface Props {
		onImageLoad?: (url: string, width: number, height: number) => void;
	}

	let { onImageLoad }: Props = $props();

	// Placeholder - full implementation in Phase 2
	let isDragging = $state(false);
</script>

<div
	class="flex h-full w-full flex-col items-center justify-center p-8 {isDragging ? 'bg-primary/5' : ''}"
	ondragover={(e) => {
		e.preventDefault();
		isDragging = true;
	}}
	ondragleave={() => (isDragging = false)}
	ondrop={(e) => {
		e.preventDefault();
		isDragging = false;
		// Phase 2: Handle file drop
	}}
	role="button"
	tabindex="0"
>
	<div
		class="flex max-w-md flex-col items-center rounded-xl border-2 border-dashed p-12 text-center transition-colors {isDragging ? 'border-primary' : 'border-border'}"
	>
		<div class="mb-4 rounded-full bg-muted p-4">
			<Upload class="h-8 w-8 text-muted-foreground" />
		</div>

		<h3 class="mb-2 text-lg font-semibold text-foreground">Upload an Image</h3>

		<p class="mb-4 text-sm text-muted-foreground">
			Drag and drop an image here, or click to browse
		</p>

		<div class="mb-4 flex items-center gap-2 text-xs text-muted-foreground">
			<Image class="h-4 w-4" />
			<span>PNG, JPG, GIF, WebP</span>
		</div>

		<div class="flex items-center gap-2 text-xs text-muted-foreground">
			<Link class="h-4 w-4" />
			<span>Or paste an image URL</span>
		</div>

		<p class="mt-6 rounded bg-muted px-3 py-1.5 text-xs text-muted-foreground">
			<kbd class="font-mono">Cmd+V</kbd> to paste from clipboard
		</p>
	</div>
</div>
