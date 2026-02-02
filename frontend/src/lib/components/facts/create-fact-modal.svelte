<script lang="ts">
	import type { FactType } from '$lib/types/fact';
	import * as Dialog from '$lib/components/ui/dialog';
	import Button from '$lib/components/ui/button/button.svelte';
	import FactTypeSelector from './fact-type-selector.svelte';

	export interface FactFormData {
		factType: FactType;
		content: Record<string, unknown>;
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

	let cardCount = $derived(selectedType === 'basic' ? 1 : selectedType === 'cloze' ? 0 : 0);

	async function handleCreate(andNew = false) {
		submitting = true;
		try {
			await onsubmit({ factType: selectedType, content: {} });
			if (andNew) {
				selectedType = 'basic';
			} else {
				open = false;
			}
		} finally {
			submitting = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.metaKey && e.key === 's') {
			e.preventDefault();
			if (e.shiftKey) {
				handleCreate(true);
			} else {
				handleCreate(false);
			}
		}
	}
</script>

<svelte:window onkeydown={open ? handleKeydown : undefined} />

<Dialog.Root bind:open onOpenChange={(v) => !v && onclose()}>
	<Dialog.Content
		class="top-[40px] flex max-h-[calc(100vh-80px)] translate-y-0 flex-col overflow-visible sm:max-w-[80vw]"
		showCloseButton={true}
	>
		<div class="">
			<Dialog.Header>
				<Dialog.Title>Create Fact</Dialog.Title>
				<Dialog.Description></Dialog.Description>
			</Dialog.Header>
		</div>

		<div class="px-0.5">
			<FactTypeSelector
				selected={selectedType}
				onchange={(t) => (selectedType = t)}
				disabled={submitting}
			/>
		</div>

		<div class="flex-1 space-y-6 overflow-y-auto py-2">
			<!-- Editor content area (phase 2+) -->
			<div class="border-border rounded-lg border border-dashed p-8 text-center">
				<p class="text-muted-foreground text-sm">
					Editor for <span class="font-medium">{selectedType}</span> facts coming in phase 2.
				</p>
			</div>
		</div>

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
				<Button
					variant="secondary"
					size="sm"
					onclick={() => handleCreate(true)}
					disabled={submitting}
				>
					Save & New
				</Button>
				<Button size="sm" onclick={() => handleCreate(false)} disabled={submitting}>
					Create Fact
				</Button>
			</div>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
