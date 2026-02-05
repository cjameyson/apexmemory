<script lang="ts">
	import type { FactType, FactDetail } from '$lib/types/fact';
	import type { JSONContent } from '@tiptap/core';
	import { untrack } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Layers, Braces, Image } from '@lucide/svelte';
	import FactTypeSelector from './fact-type-selector.svelte';
	import BasicFactEditor, { type BasicFactData } from './basic-fact-editor.svelte';

	const factTypeMeta: Record<FactType, { label: string; icon: typeof Layers }> = {
		basic: { label: 'Basic', icon: Layers },
		cloze: { label: 'Cloze', icon: Braces },
		image_occlusion: { label: 'Image Occlusion', icon: Image }
	};
	import ClozeFactEditor, { type ClozeFactData } from './cloze-fact-editor.svelte';
	import { ImageOcclusionEditor, type ImageOcclusionField } from '$lib/components/image-occlusion';
	import ConfirmDialog from '$lib/components/ui/confirm-dialog.svelte';

	export interface FactFormData {
		factType: FactType;
		content: {
			version: number;
			asset_ids?: string[];
			fields: { name: string; type: string; value: string | JSONContent; regions?: { id: string }[] }[];
		};
	}

	/** Walk TipTap JSON documents and collect all referenced asset IDs from image nodes. */
	function extractAssetIds(...docs: (JSONContent | null)[]): string[] {
		const ids: string[] = [];
		function walk(node: Record<string, unknown>) {
			if (node?.type === 'image') {
				const attrs = node.attrs as Record<string, unknown> | undefined;
				if (attrs?.asset_id && typeof attrs.asset_id === 'string') {
					ids.push(attrs.asset_id);
				}
			}
			if (Array.isArray(node?.content)) {
				(node.content as Record<string, unknown>[]).forEach(walk);
			}
		}
		for (const doc of docs) {
			if (doc) walk(doc as Record<string, unknown>);
		}
		return [...new Set(ids)];
	}

	/** Check if a TipTap JSON document is effectively empty (no text or image content). */
	function isRichTextEmpty(doc: JSONContent | null): boolean {
		if (!doc) return true;
		let hasContent = false;
		function walk(node: Record<string, unknown>) {
			if (hasContent) return;
			if (node.type === 'text' && typeof node.text === 'string' && node.text.trim()) {
				hasContent = true;
				return;
			}
			if (node.type === 'image') {
				hasContent = true;
				return;
			}
			if (Array.isArray(node.content)) {
				(node.content as Record<string, unknown>[]).forEach(walk);
			}
		}
		walk(doc as Record<string, unknown>);
		return !hasContent;
	}

	let {
		open = $bindable(false),
		notebookId,
		editFact,
		onclose,
		onsubmit
	}: {
		open: boolean;
		notebookId: string;
		editFact?: FactDetail | null;
		onclose: () => void;
		onsubmit: (data: FactFormData) => Promise<void>;
	} = $props();

	let isEditMode = $derived(!!editFact);

	let selectedType = $state<FactType>('basic');
	let submitting = $state(false);
	let submitError = $state<string | null>(null);

	// Independent state per type (preserved when switching)
	let basicData = $state<BasicFactData>({ front: null, back: null, backExtra: '', hint: '' });
	let clozeData = $state<ClozeFactData>({ text: '', backExtra: '' });

	// Validation errors
	let basicErrors = $state<Partial<Record<keyof BasicFactData, string>>>({});
	let clozeErrors = $state<Partial<Record<keyof ClozeFactData, string>>>({});
	let imageOcclusionErrors = $state<{ title?: boolean; regionLabels?: Set<string> }>({});

	// Card count from cloze editor
	let clozeCardCount = $state(0);
	// Image occlusion data
	let imageOcclusionData = $state<ImageOcclusionField | null>(null);
	// Incrementing formKey forces {#key} to destroy and recreate editor components,
	// clearing their internal state (e.g. input values) after "Save & New".
	let formKey = $state(0);

	let cardCount = $derived(
		selectedType === 'basic'
			? 1
			: selectedType === 'cloze'
				? clozeCardCount
				: selectedType === 'image_occlusion' && imageOcclusionData
					? imageOcclusionData.regions.length
					: 0
	);

	// Editor refs
	let basicEditor: BasicFactEditor | undefined = $state();
	let clozeEditor: ClozeFactEditor | undefined = $state();
	let imageOcclusionEditor: ImageOcclusionEditor | undefined = $state();

	// Discard confirmation for dirty image occlusion editor
	let discardConfirmOpen = $state(false);

	function validate(): boolean {
		if (selectedType === 'basic') {
			const errs: typeof basicErrors = {};
			if (isRichTextEmpty(basicData.front)) errs.front = 'Front is required';
			if (isRichTextEmpty(basicData.back)) errs.back = 'Back is required';
			basicErrors = errs;
			if (Object.keys(errs).length > 0) {
				basicEditor?.focus();
				return false;
			}
			return true;
		}
		if (selectedType === 'cloze') {
			const errs: typeof clozeErrors = {};
			if (!clozeData.text.trim()) {
				errs.text = 'Cloze text is required';
			} else if (!/\{\{c\d+::.+?\}\}/.test(clozeData.text)) {
				errs.text = 'Each cloze deletion must have content (e.g. {{c1::answer}})';
			}
			clozeErrors = errs;
			if (Object.keys(errs).length > 0) {
				clozeEditor?.focus();
				return false;
			}
			return true;
		}
		if (selectedType === 'image_occlusion') {
			const errs: typeof imageOcclusionErrors = {};
			const errors: string[] = [];

			if (!imageOcclusionData) {
				submitError = 'No image occlusion data';
				return false;
			}
			if (!imageOcclusionData.title.trim()) {
				errs.title = true;
				errors.push('Title is required');
			}
			if (imageOcclusionData.regions.length === 0) {
				errors.push('At least one region is required');
			}
			const missingLabels = imageOcclusionData.regions.filter((r) => !r.label.trim());
			if (missingLabels.length > 0) {
				errs.regionLabels = new Set(missingLabels.map((r) => r.id));
				errors.push(`${missingLabels.length} region(s) missing labels`);
			}
			imageOcclusionErrors = errs;
			if (errors.length > 0) {
				submitError = errors.join('. ');
				return false;
			}
			return true;
		}
		return false;
	}

	function buildContent(): FactFormData['content'] {
		if (selectedType === 'basic') {
			const assetIds = extractAssetIds(basicData.front, basicData.back);
			const fields: FactFormData['content']['fields'] = [
				{ name: 'front', type: 'rich_text', value: basicData.front! },
				{ name: 'back', type: 'rich_text', value: basicData.back! }
			];
			if (basicData.backExtra.trim()) {
				fields.push({ name: 'back_extra', type: 'plain_text', value: basicData.backExtra.trim() });
			}
			if (basicData.hint.trim()) {
				fields.push({ name: 'hint', type: 'plain_text', value: basicData.hint.trim() });
			}
			const content: FactFormData['content'] = { version: 1, fields };
			if (assetIds.length > 0) {
				content.asset_ids = assetIds;
			}
			return content;
		}
		if (selectedType === 'cloze') {
			const fields: FactFormData['content']['fields'] = [
				{ name: 'text', type: 'cloze_text', value: clozeData.text.trim() }
			];
			if (clozeData.backExtra.trim()) {
				fields.push({ name: 'back_extra', type: 'plain_text', value: clozeData.backExtra.trim() });
			}
			return { version: 1, fields };
		}
		// image_occlusion: title as separate plain_text field + occlusion data with regions array
		if (selectedType === 'image_occlusion' && imageOcclusionData) {
			const assetIds: string[] = [];
			if (imageOcclusionData.image.assetId) {
				assetIds.push(imageOcclusionData.image.assetId);
			}
			const content: FactFormData['content'] = {
				version: 1,
				fields: [
					{
						name: 'title',
						type: 'plain_text',
						value: imageOcclusionData.title.trim()
					},
					{
						name: 'image_occlusion',
						type: 'image_occlusion',
						value: imageOcclusionData,
						regions: imageOcclusionData.regions.map((r) => ({ id: r.id }))
					}
				]
			};
			if (assetIds.length > 0) {
				content.asset_ids = assetIds;
			}
			return content;
		}
		return { version: 1, fields: [] };
	}

	function resetActiveType() {
		if (selectedType === 'basic') {
			basicData = { front: null, back: null, backExtra: '', hint: '' };
			basicErrors = {};
		} else if (selectedType === 'cloze') {
			clozeData = { text: '', backExtra: '' };
			clozeErrors = {};
		} else if (selectedType === 'image_occlusion') {
			imageOcclusionData = null;
		}
		formKey++;
	}

	async function handleCreate(andNew = false) {
		if (!validate()) return;

		submitError = null;
		submitting = true;
		try {
			await onsubmit({ factType: selectedType, content: buildContent() });
			if (andNew) {
				resetActiveType();
				requestAnimationFrame(() => {
					if (selectedType === 'basic') basicEditor?.focus();
					else if (selectedType === 'cloze') clozeEditor?.focus();
				});
			} else {
				open = false;
			}
		} catch (err) {
			submitError =
				err instanceof Error
					? err.message
					: isEditMode
						? 'Failed to save fact'
						: 'Failed to create fact';
		} finally {
			submitting = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if ((e.metaKey || e.ctrlKey) && e.key === 's') {
			e.preventDefault();
			if (e.shiftKey && !isEditMode) {
				handleCreate(true);
			} else {
				handleCreate(false);
			}
		}
	}

	// Clear errors on data change
	function handleBasicChange(data: BasicFactData) {
		basicData = data;
		basicErrors = {};
	}

	function handleClozeChange(data: ClozeFactData) {
		clozeData = data;
		clozeErrors = {};
	}

	function handleImageOcclusionChange(data: ImageOcclusionField) {
		imageOcclusionData = data;
		imageOcclusionErrors = {};
		submitError = null;
	}

	function parseContentToEditorData(content: Record<string, unknown>, factType: FactType) {
		const fields = (content as { fields?: { name: string; type: string; value: string | JSONContent }[] }).fields ?? [];
		const fieldMap = new Map(fields.map((f) => [f.name, f]));

		if (factType === 'basic') {
			const frontField = fieldMap.get('front');
			const backField = fieldMap.get('back');
			const backExtraField = fieldMap.get('back_extra');
			const hintField = fieldMap.get('hint');

			// Support both rich_text (JSONContent) and legacy plain_text (string) values
			function toRichText(field: { type: string; value: string | JSONContent } | undefined): JSONContent | null {
				if (!field) return null;
				if (field.type === 'rich_text' && typeof field.value === 'object') {
					return field.value as JSONContent;
				}
				// Legacy plain_text: wrap in a minimal TipTap doc
				const text = typeof field.value === 'string' ? field.value : '';
				if (!text) return null;
				return {
					type: 'doc',
					content: [{ type: 'paragraph', content: [{ type: 'text', text }] }]
				};
			}

			basicData = {
				front: toRichText(frontField),
				back: toRichText(backField),
				backExtra: (backExtraField?.value as string) ?? '',
				hint: (hintField?.value as string) ?? ''
			};
			basicErrors = {};
		} else if (factType === 'cloze') {
			clozeData = {
				text: (fieldMap.get('text')?.value as string) ?? '',
				backExtra: (fieldMap.get('back_extra')?.value as string) ?? ''
			};
			clozeErrors = {};
		}
	}

	// Hydrate editor state when the modal opens.
	// Uses a prevOpen flag to detect the open transition, so we only hydrate once
	// per open (not while the modal stays open). This fixes the bug where
	// re-editing the same fact wouldn't re-populate fields.
	let prevOpen = false;
	$effect(() => {
		const isOpen = open;
		if (isOpen && !prevOpen) {
			// Modal just opened
			if (editFact) {
				const { factType, content } = editFact;
				untrack(() => {
					selectedType = factType;
					parseContentToEditorData(content, factType);
					formKey++;
				});
			} else {
				untrack(() => {
					resetActiveType();
				});
			}
		}
		prevOpen = isOpen;
	});

	// Auto-focus on type change
	$effect(() => {
		const type = selectedType;
		const isOpen = open;
		if (!isOpen) return;
		requestAnimationFrame(() => {
			if (type === 'basic') basicEditor?.focus();
			else if (type === 'cloze') clozeEditor?.focus();
		});
	});

	function handleClose() {
		if (selectedType === 'image_occlusion' && imageOcclusionEditor?.getIsDirty()) {
			// Re-open immediately since bind:open already set it to false
			open = true;
			discardConfirmOpen = true;
			return;
		}
		onclose();
	}

	function confirmDiscard() {
		discardConfirmOpen = false;
		onclose();
	}

	function cancelDiscard() {
		discardConfirmOpen = false;
	}
</script>

<svelte:window onkeydown={open ? handleKeydown : undefined} />

<Dialog.Root bind:open onOpenChange={(v) => !v && handleClose()}>
	<Dialog.Content
		class="top-[40px] flex translate-y-0 flex-col overflow-visible {selectedType ===
		'image_occlusion'
			? 'h-[calc(100vh-80px)] sm:max-w-[95vw]'
			: 'max-h-[calc(100vh-80px)] sm:max-w-[80vw]'}"
		showCloseButton={true}
	>
		<div class="shrink-0 space-y-4 pr-2">
			<Dialog.Header>
				{#if isEditMode}
					{@const meta = factTypeMeta[selectedType]}
					<Dialog.Title class="flex items-center gap-2">
						Edit Fact
						<span
							class="bg-muted text-muted-foreground inline-flex items-center gap-1.5 rounded-md px-2 py-0.5 text-xs font-medium"
						>
							<meta.icon class="h-3.5 w-3.5" />
							{meta.label}
						</span>
					</Dialog.Title>
				{:else}
					<Dialog.Title>
						<FactTypeSelector
							selected={selectedType}
							onchange={(t) => (selectedType = t)}
							disabled={submitting}
						/>
					</Dialog.Title>
				{/if}
				<Dialog.Description></Dialog.Description>
			</Dialog.Header>
		</div>

		{#if selectedType === 'image_occlusion'}
			<!-- Image occlusion editor fills the modal content area -->
			<div class="image-occlusion-editor min-h-0 flex-1 overflow-hidden">
				{#key formKey}
					<ImageOcclusionEditor bind:this={imageOcclusionEditor} onChange={handleImageOcclusionChange} errors={imageOcclusionErrors} />
				{/key}
			</div>
		{:else}
			<div class="flex-1 space-y-6 overflow-y-auto px-0.5 py-2">
				{#key formKey}
					{#if selectedType === 'basic'}
						<BasicFactEditor
							bind:this={basicEditor}
							initialData={basicData}
							onchange={handleBasicChange}
							errors={basicErrors}
						/>
					{:else if selectedType === 'cloze'}
						<ClozeFactEditor
							bind:this={clozeEditor}
							initialData={clozeData}
							onchange={handleClozeChange}
							errors={clozeErrors}
							oncardcount={(count) => (clozeCardCount = count)}
						/>
					{/if}
				{/key}
			</div>
		{/if}

		<Dialog.Footer class="w-full shrink-0 flex-col border-t pt-4">
			{#if submitError}
				<p class="text-destructive mb-2 text-sm">{submitError}</p>
			{/if}
			<div class="flex w-full items-center justify-between">
			<div class="text-muted-foreground flex items-center gap-3 text-xs">
				<span
					><kbd class="bg-muted rounded px-1.5 py-0.5 font-mono text-[10px]">&#8984;S</kbd> Save</span
				>
				{#if !isEditMode && selectedType !== 'image_occlusion'}
					<span
						><kbd class="bg-muted rounded px-1.5 py-0.5 font-mono text-[10px]">&#8984;&#8679;S</kbd> Save
						& New</span
					>
				{/if}
				<span
					><kbd class="bg-muted rounded px-1.5 py-0.5 font-mono text-[10px]">&#9099;</kbd> Cancel</span
				>
			</div>
			<div class="flex items-center gap-2">
				{#if cardCount > 0}
					<span class="bg-muted text-muted-foreground rounded-full px-2 py-0.5 text-xs font-medium">
						{cardCount}
						{cardCount === 1 ? 'card' : 'cards'}
					</span>
				{/if}
				{#if !isEditMode && selectedType !== 'image_occlusion'}
					<Button
						variant="secondary"
						size="sm"
						onclick={() => handleCreate(true)}
						disabled={submitting}
					>
						Save & New
					</Button>
				{/if}
				<Button size="sm" onclick={() => handleCreate(false)} disabled={submitting}>
					{isEditMode ? 'Save' : 'Create Fact'}
				</Button>
			</div>
			</div>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<ConfirmDialog
	bind:open={discardConfirmOpen}
	title="Discard changes?"
	description="You have unsaved changes to the image occlusion editor."
	confirmLabel="Discard"
	onconfirm={confirmDiscard}
	oncancel={cancelDiscard}
/>
