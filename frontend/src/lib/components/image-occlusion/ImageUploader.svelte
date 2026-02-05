<script lang="ts">
	import { Upload } from '@lucide/svelte';
	import { uploadAsset, assetUrl } from '$lib/api/client';

	interface Props {
		onImageLoad?: (url: string, width: number, height: number, assetId: string) => void;
	}

	let { onImageLoad }: Props = $props();

	let isDragging = $state(false);
	let uploading = $state(false);
	let errorMsg = $state<string | null>(null);

	async function handleFile(file: File) {
		if (!file.type.startsWith('image/')) {
			errorMsg = 'Please select an image file (JPEG, PNG, WebP, or GIF)';
			return;
		}

		errorMsg = null;
		uploading = true;

		try {
			const asset = await uploadAsset(file);
			const url = assetUrl(asset.id);
			const width = asset.metadata?.width ?? 0;
			const height = asset.metadata?.height ?? 0;
			onImageLoad?.(url, width, height, asset.id);
		} catch (err) {
			errorMsg = err instanceof Error ? err.message : 'Upload failed';
		} finally {
			uploading = false;
		}
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
		const file = e.dataTransfer?.files[0];
		if (file) handleFile(file);
	}

	function handleFileInput(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		if (file) handleFile(file);
		input.value = '';
	}

	let fileInput: HTMLInputElement;
</script>

<div
	class="flex h-full w-full flex-col items-center justify-center p-8 {isDragging ? 'bg-primary/5' : ''}"
	ondragover={(e) => {
		e.preventDefault();
		isDragging = true;
	}}
	ondragleave={() => (isDragging = false)}
	ondrop={handleDrop}
	role="button"
	tabindex="0"
	onclick={() => fileInput?.click()}
	onkeydown={(e) => {
		if (e.key === 'Enter' || e.key === ' ') fileInput?.click();
	}}
>
	<input
		bind:this={fileInput}
		type="file"
		accept="image/jpeg,image/png,image/webp,image/gif"
		class="hidden"
		onchange={handleFileInput}
	/>

	<div
		class="flex max-w-md flex-col items-center rounded-xl border-2 border-dashed border-border p-12 text-center"
	>
		{#if uploading}
			<div class="text-muted-foreground">Uploading...</div>
		{:else}
			<Upload class="mb-4 h-12 w-12 text-muted-foreground" />
			<p class="mb-2 font-medium text-foreground">Drop an image here or click to browse</p>
			<p class="text-sm text-muted-foreground">JPEG, PNG, WebP, or GIF up to 10MB</p>
		{/if}

		{#if errorMsg}
			<p class="mt-4 text-sm text-destructive">{errorMsg}</p>
		{/if}
	</div>
</div>
