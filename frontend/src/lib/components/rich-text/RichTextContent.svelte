<script lang="ts">
	import { generateHTML } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Underline from '@tiptap/extension-underline';
	import Image from '@tiptap/extension-image';
	import { assetUrl } from '$lib/api/client';
	import type { JSONContent } from '@tiptap/core';

	interface Props {
		content: JSONContent;
		class?: string;
	}

	let { content, class: className = '' }: Props = $props();

	/**
	 * Custom Image extension that resolves asset_id to proxy URLs during HTML generation.
	 * If asset_id is present, it overrides src with the resolved URL.
	 * This ensures images always point to the correct proxy route regardless of
	 * what src was stored in the JSON content.
	 */
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
		},
		renderHTML({ HTMLAttributes }) {
			const attrs = { ...HTMLAttributes };
			if (attrs['data-asset-id']) {
				attrs.src = assetUrl(attrs['data-asset-id'] as string);
			}
			return ['img', attrs];
		}
	});

	const extensions = [StarterKit, Underline, AssetImage];

	let html = $derived(generateHTML(content, extensions));
</script>

<div class="rich-text-content prose prose-sm dark:prose-invert max-w-none {className}">
	{@html html}
</div>

<style>
	/* Constrain images in rendered content */
	.rich-text-content :global(img) {
		max-width: 100%;
		height: auto;
		border-radius: var(--radius-md);
	}
</style>
