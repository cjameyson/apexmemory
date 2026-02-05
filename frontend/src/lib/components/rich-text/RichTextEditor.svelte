<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor, type JSONContent } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Underline from '@tiptap/extension-underline';
	import Image from '@tiptap/extension-image';
	import { uploadAsset, assetUrl } from '$lib/api/client';
	import EditorToolbar from './EditorToolbar.svelte';

	interface Props {
		content?: JSONContent | null;
		onchange?: (json: JSONContent) => void;
		placeholder?: string;
		class?: string;
	}

	let { content = null, onchange, placeholder = '', class: className = '' }: Props = $props();

	let editorElement: HTMLDivElement | undefined = $state();
	let editor: Editor | undefined = $state();
	// svelte-ignore state_referenced_locally
	let transactionKey: number = $state(0);
	let uploading: boolean = $state(false);

	/** Custom Image node extended with asset_id attribute. */
	const AssetImage = Image.extend({
		addAttributes() {
			return {
				...this.parent?.(),
				asset_id: {
					default: null,
					parseHTML: (element: HTMLElement) => element.getAttribute('data-asset-id'),
					renderHTML: (attributes: Record<string, unknown>) => {
						if (!attributes.asset_id) return {};
						return { 'data-asset-id': attributes.asset_id };
					}
				}
			};
		}
	});

	onMount(() => {
		if (!editorElement) return;

		editor = new Editor({
			extensions: [
				StarterKit,
				Underline,
				AssetImage.configure({ allowBase64: false })
			],
			content: content ?? undefined,
			onUpdate: ({ editor: ed }) => {
				onchange?.(ed.getJSON());
			},
			onTransaction: () => {
				transactionKey++;
			},
			editorProps: {
				attributes: {
					class:
						'prose prose-sm max-w-none focus:outline-none min-h-[80px] px-3 py-2',
					'data-placeholder': placeholder
				},
				handleDrop: (_view, event, _slice, moved) => {
					if (moved || !event.dataTransfer?.files.length) return false;
					const images = Array.from(event.dataTransfer.files).filter((f) =>
						f.type.startsWith('image/')
					);
					if (!images.length) return false;
					event.preventDefault();
					for (const file of images) {
						insertImage(file);
					}
					return true;
				},
				handlePaste: (_view, event) => {
					const items = event.clipboardData?.items;
					if (!items) return false;
					const imageItems = Array.from(items).filter((item) =>
						item.type.startsWith('image/')
					);
					if (!imageItems.length) return false;
					event.preventDefault();
					for (const item of imageItems) {
						const file = item.getAsFile();
						if (file) insertImage(file);
					}
					return true;
				}
			}
		});

		editor.mount(editorElement);
	});

	onDestroy(() => {
		editor?.destroy();
	});

	/** Upload an image file and insert it into the editor. */
	async function insertImage(file: File) {
		if (!editor) return;
		uploading = true;
		try {
			const asset = await uploadAsset(file);
			const src = assetUrl(asset.id);
			editor
				.chain()
				.focus()
				.insertContent({
					type: 'image',
					attrs: {
						src,
						alt: asset.filename,
						asset_id: asset.id
					}
				})
				.run();
		} catch (err) {
			console.error('Failed to upload image:', err);
		} finally {
			uploading = false;
		}
	}

	/** Focus the editor. */
	export function focus() {
		editor?.commands.focus();
	}

	/** Open file picker and insert selected images. */
	function handleImageInsert() {
		const input = document.createElement('input');
		input.type = 'file';
		input.accept = 'image/*';
		input.multiple = true;
		input.onchange = () => {
			if (!input.files) return;
			for (const file of Array.from(input.files)) {
				insertImage(file);
			}
		};
		input.click();
	}
</script>

<div class="rounded-md border border-input bg-card text-foreground {className}">
	{#if editor}
		<EditorToolbar {editor} {transactionKey} onImageInsert={handleImageInsert} />
	{/if}

	<div bind:this={editorElement}></div>

	{#if uploading}
		<div class="border-t border-border px-3 py-1 text-xs text-muted-foreground">
			Uploading image...
		</div>
	{/if}
</div>

<style>
	/* Placeholder styling for empty editor */
	:global(.ProseMirror p.is-editor-empty:first-child::before) {
		content: attr(data-placeholder);
		float: left;
		color: var(--muted-foreground);
		pointer-events: none;
		height: 0;
	}

	/* Placeholder via editor-level data attribute */
	:global(.ProseMirror.is-empty::before) {
		content: attr(data-placeholder);
		float: left;
		color: var(--muted-foreground);
		pointer-events: none;
		height: 0;
	}

	/* Ensure images inside the editor are constrained */
	:global(.ProseMirror img) {
		max-width: 100%;
		height: auto;
		border-radius: var(--radius-md);
	}

	/* Selected image highlight */
	:global(.ProseMirror img.ProseMirror-selectednode) {
		outline: 2px solid var(--ring);
		outline-offset: 2px;
	}
</style>
