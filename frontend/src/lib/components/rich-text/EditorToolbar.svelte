<script lang="ts">
	import type { Editor } from '@tiptap/core';
	import {
		Bold,
		Italic,
		Underline,
		Heading2,
		List,
		ListOrdered,
		ImagePlus,
		Undo2,
		Redo2
	} from '@lucide/svelte';

	interface Props {
		editor: Editor;
		/** Incremented on each editor transaction to force toolbar reactivity. */
		transactionKey: number;
		onImageInsert?: () => void;
	}

	let { editor, transactionKey, onImageInsert }: Props = $props();

	// Consume transactionKey so Svelte treats it as a dependency.
	// This forces $derived re-evaluation when the editor state changes.
	// svelte-ignore state_referenced_locally
	let _tick = $derived(transactionKey);

	function btn(active: boolean): string {
		return active
			? 'rounded p-1.5 bg-muted text-foreground'
			: 'rounded p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground';
	}
</script>

<div class="flex flex-wrap items-center gap-0.5 border-b border-border bg-muted/40 px-2 py-1">
	<!-- Text formatting -->
	<button
		type="button"
		class={btn(editor.isActive('bold'))}
		onclick={() => editor.chain().focus().toggleBold().run()}
		disabled={!editor.can().chain().focus().toggleBold().run()}
		aria-label="Bold"
		aria-pressed={editor.isActive('bold')}
		title="Bold (Ctrl+B)"
	>
		<Bold class="h-4 w-4" />
	</button>

	<button
		type="button"
		class={btn(editor.isActive('italic'))}
		onclick={() => editor.chain().focus().toggleItalic().run()}
		disabled={!editor.can().chain().focus().toggleItalic().run()}
		aria-label="Italic"
		aria-pressed={editor.isActive('italic')}
		title="Italic (Ctrl+I)"
	>
		<Italic class="h-4 w-4" />
	</button>

	<button
		type="button"
		class={btn(editor.isActive('underline'))}
		onclick={() => editor.chain().focus().toggleUnderline().run()}
		disabled={!editor.can().chain().focus().toggleUnderline().run()}
		aria-label="Underline"
		aria-pressed={editor.isActive('underline')}
		title="Underline (Ctrl+U)"
	>
		<Underline class="h-4 w-4" />
	</button>

	<span class="mx-1 h-5 w-px bg-border"></span>

	<!-- Heading -->
	<button
		type="button"
		class={btn(editor.isActive('heading', { level: 2 }))}
		onclick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
		aria-label="Heading 2"
		aria-pressed={editor.isActive('heading', { level: 2 })}
		title="Heading 2"
	>
		<Heading2 class="h-4 w-4" />
	</button>

	<span class="mx-1 h-5 w-px bg-border"></span>

	<!-- Lists -->
	<button
		type="button"
		class={btn(editor.isActive('bulletList'))}
		onclick={() => editor.chain().focus().toggleBulletList().run()}
		aria-label="Bullet list"
		aria-pressed={editor.isActive('bulletList')}
		title="Bullet List"
	>
		<List class="h-4 w-4" />
	</button>

	<button
		type="button"
		class={btn(editor.isActive('orderedList'))}
		onclick={() => editor.chain().focus().toggleOrderedList().run()}
		aria-label="Ordered list"
		aria-pressed={editor.isActive('orderedList')}
		title="Ordered List"
	>
		<ListOrdered class="h-4 w-4" />
	</button>

	<span class="mx-1 h-5 w-px bg-border"></span>

	<!-- Image -->
	<button
		type="button"
		class={btn(false)}
		onclick={() => onImageInsert?.()}
		aria-label="Insert image"
		title="Insert Image"
	>
		<ImagePlus class="h-4 w-4" />
	</button>

	<span class="mx-1 h-5 w-px bg-border"></span>

	<!-- Undo / Redo -->
	<button
		type="button"
		class={btn(false)}
		onclick={() => editor.chain().focus().undo().run()}
		disabled={!editor.can().chain().focus().undo().run()}
		aria-label="Undo"
		title="Undo (Ctrl+Z)"
	>
		<Undo2 class="h-4 w-4" />
	</button>

	<button
		type="button"
		class={btn(false)}
		onclick={() => editor.chain().focus().redo().run()}
		disabled={!editor.can().chain().focus().redo().run()}
		aria-label="Redo"
		title="Redo (Ctrl+Shift+Z)"
	>
		<Redo2 class="h-4 w-4" />
	</button>
</div>
