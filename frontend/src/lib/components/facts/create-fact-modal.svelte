<script lang="ts">
	import type { FactType } from '$lib/types/fact';
	import * as Dialog from '$lib/components/ui/dialog';
	import Button from '$lib/components/ui/button/button.svelte';
	import FactTypeSelector from './fact-type-selector.svelte';
	import BasicFactEditor, { type BasicFactData } from './basic-fact-editor.svelte';
	import ClozeFactEditor, { type ClozeFactData } from './cloze-fact-editor.svelte';
	import ImageOcclusionPlaceholder from './image-occlusion-placeholder.svelte';

	export interface FactFormData {
		factType: FactType;
		content: { version: number; fields: { name: string; type: string; value: string }[] };
	}

	let {
		open = $bindable(false),
		notebookId,
		onclose,
		onsubmit
	}: {
		open: boolean;
		notebookId: string;
		onclose: () => void;
		onsubmit: (data: FactFormData) => Promise<void>;
	} = $props();

	let selectedType = $state<FactType>('basic');
	let submitting = $state(false);
	let submitError = $state<string | null>(null);

	// Independent state per type (preserved when switching)
	let basicData = $state<BasicFactData>({ front: '', back: '', backExtra: '', hint: '' });
	let clozeData = $state<ClozeFactData>({ text: '', backExtra: '' });

	// Validation errors
	let basicErrors = $state<Partial<Record<keyof BasicFactData, string>>>({});
	let clozeErrors = $state<Partial<Record<keyof ClozeFactData, string>>>({});

	// Card count from cloze editor
	let clozeCardCount = $state(0);
	// Incrementing formKey forces {#key} to destroy and recreate editor components,
	// clearing their internal state (e.g. input values) after "Save & New".
	let formKey = $state(0);

	let cardCount = $derived(
		selectedType === 'basic' ? 1 : selectedType === 'cloze' ? clozeCardCount : 0
	);

	// Editor refs
	let basicEditor: BasicFactEditor | undefined = $state();
	let clozeEditor: ClozeFactEditor | undefined = $state();

	function validate(): boolean {
		if (selectedType === 'basic') {
			const errs: typeof basicErrors = {};
			if (!basicData.front.trim()) errs.front = 'Front is required';
			if (!basicData.back.trim()) errs.back = 'Back is required';
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
		return false;
	}

	function buildContent(): FactFormData['content'] {
		if (selectedType === 'basic') {
			const fields: { name: string; type: string; value: string }[] = [
				{ name: 'front', type: 'plain_text', value: basicData.front.trim() },
				{ name: 'back', type: 'plain_text', value: basicData.back.trim() }
			];
			if (basicData.backExtra.trim()) {
				fields.push({ name: 'back_extra', type: 'plain_text', value: basicData.backExtra.trim() });
			}
			if (basicData.hint.trim()) {
				fields.push({ name: 'hint', type: 'plain_text', value: basicData.hint.trim() });
			}
			return { version: 1, fields };
		}
		// cloze
		const fields: { name: string; type: string; value: string }[] = [
			{ name: 'text', type: 'cloze_text', value: clozeData.text.trim() }
		];
		if (clozeData.backExtra.trim()) {
			fields.push({ name: 'back_extra', type: 'plain_text', value: clozeData.backExtra.trim() });
		}
		return { version: 1, fields };
	}

	function resetActiveType() {
		if (selectedType === 'basic') {
			basicData = { front: '', back: '', backExtra: '', hint: '' };
			basicErrors = {};
		} else if (selectedType === 'cloze') {
			clozeData = { text: '', backExtra: '' };
			clozeErrors = {};
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
			submitError = err instanceof Error ? err.message : 'Failed to create fact';
		} finally {
			submitting = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if ((e.metaKey || e.ctrlKey) && e.key === 's') {
			e.preventDefault();
			if (e.shiftKey) {
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
</script>

<svelte:window onkeydown={open ? handleKeydown : undefined} />

<Dialog.Root bind:open onOpenChange={(v) => !v && onclose()}>
	<Dialog.Content
		class="top-[40px] flex max-h-[calc(100vh-80px)] translate-y-0 flex-col overflow-visible sm:max-w-[80vw]"
		showCloseButton={true}
	>
		<div class="space-y-4">
			<Dialog.Header>
				<Dialog.Title>Create Fact</Dialog.Title>
				<Dialog.Description></Dialog.Description>
			</Dialog.Header>
			<FactTypeSelector
				selected={selectedType}
				onchange={(t) => (selectedType = t)}
				disabled={submitting}
			/>
		</div>

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
				{:else}
					<ImageOcclusionPlaceholder />
				{/if}
			{/key}
		</div>

		{#if submitError}
			<p class="text-destructive px-0.5 text-sm">{submitError}</p>
		{/if}

		<Dialog.Footer class="w-full flex-row items-center !justify-between border-t pt-4">
			<div class="text-muted-foreground flex items-center gap-3 text-xs">
				<span
					><kbd class="bg-muted rounded px-1.5 py-0.5 font-mono text-[10px]">&#8984;S</kbd> Save</span
				>
				<span
					><kbd class="bg-muted rounded px-1.5 py-0.5 font-mono text-[10px]">&#8984;&#8679;S</kbd> Save
					& New</span
				>
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
				<Button
					variant="secondary"
					size="sm"
					onclick={() => handleCreate(true)}
					disabled={submitting || selectedType === 'image_occlusion'}
				>
					Save & New
				</Button>
				<Button
					size="sm"
					onclick={() => handleCreate(false)}
					disabled={submitting || selectedType === 'image_occlusion'}
				>
					Create Fact
				</Button>
			</div>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
