<script lang="ts">
	import { untrack } from 'svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Button } from '$lib/components/ui/button';
	import EmojiSelector from '$lib/components/ui/emoji-selector.svelte';
	import type { Notebook } from '$lib/types';

	interface Props {
		open: boolean;
		onOpenChange: (open: boolean) => void;
		onSuccess?: (notebook: Notebook) => void;
	}

	let { open = $bindable(false), onOpenChange, onSuccess }: Props = $props();

	// Form state
	let emoji = $state('📓');
	let name = $state('');
	let description = $state('');

	// UI state
	let submitting = $state(false);
	let errors = $state<{ name?: string; description?: string }>({});

	// Track open state to reset form
	let prevOpen = false;
	$effect(() => {
		const isOpen = open;
		if (isOpen && !prevOpen) {
			// Modal just opened - reset form
			untrack(() => {
				emoji = '📓';
				name = '';
				description = '';
				errors = {};
			});
		}
		prevOpen = isOpen;
	});

	function validate(): boolean {
		const newErrors: typeof errors = {};
		const trimmedName = name.trim();

		if (!trimmedName) {
			newErrors.name = 'Name is required';
		} else if (trimmedName.length > 255) {
			newErrors.name = 'Name must be 255 characters or less';
		}

		if (description.length > 10000) {
			newErrors.description = 'Description must be 10,000 characters or less';
		}

		errors = newErrors;
		return Object.keys(newErrors).length === 0;
	}

	async function handleSubmit() {
		if (!validate()) return;

		submitting = true;
		errors = {};

		try {
			const response = await fetch('/api/notebooks', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					name: name.trim(),
					description: description.trim() || undefined,
					emoji: emoji || undefined
				})
			});

			if (!response.ok) {
				const data = await response.json();
				if (data.fieldErrors) {
					errors = data.fieldErrors;
				} else {
					errors = { name: data.message || 'Failed to create notebook' };
				}
				return;
			}

			const notebook = await response.json();
			onSuccess?.(notebook);
			open = false;
		} catch (err) {
			errors = { name: 'Network error. Please try again.' };
		} finally {
			submitting = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') {
			e.preventDefault();
			handleSubmit();
		}
	}

	function handleOpenChange(isOpen: boolean) {
		if (!isOpen) {
			onOpenChange(false);
		}
	}

	// Clear field errors when user types
	function handleNameInput() {
		if (errors.name) {
			errors = { ...errors, name: undefined };
		}
	}

	function handleDescriptionInput() {
		if (errors.description) {
			errors = { ...errors, description: undefined };
		}
	}
</script>

<svelte:window onkeydown={open ? handleKeydown : undefined} />

<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content class="sm:max-w-md" showCloseButton={true}>
		<Dialog.Header>
			<Dialog.Title>Create Notebook</Dialog.Title>
			<Dialog.Description class="sr-only">Create a new notebook to organize your flashcards</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
			<!-- Emoji + Name row -->
			<div class="space-y-2">
				<label for="notebook-name" class="text-sm font-medium">
					Name <span class="text-destructive">*</span>
				</label>
				<div class="flex gap-2">
					<EmojiSelector bind:value={emoji} defaultValue="📓" />
					<Input
						id="notebook-name"
						type="text"
						placeholder="Chemistry 101"
						class="flex-1"
						bind:value={name}
						oninput={handleNameInput}
						aria-invalid={!!errors.name}
						disabled={submitting}
					/>
				</div>
				{#if errors.name}
					<p class="text-sm text-destructive">{errors.name}</p>
				{/if}
			</div>

			<!-- Description -->
			<div class="space-y-2">
				<label for="notebook-description" class="text-sm font-medium">Description</label>
				<Textarea
					id="notebook-description"
					placeholder="What will you study? Add context to help AI generate better flashcards..."
					class="min-h-24 resize-none"
					bind:value={description}
					oninput={handleDescriptionInput}
					aria-invalid={!!errors.description}
					disabled={submitting}
				/>
				{#if errors.description}
					<p class="text-sm text-destructive">{errors.description}</p>
				{/if}
			</div>
		</form>

		<Dialog.Footer class="flex-row gap-2 sm:justify-end">
			<Button variant="outline" onclick={() => (open = false)} disabled={submitting}>
				Cancel
			</Button>
			<Button onclick={handleSubmit} disabled={submitting}>
				{#if submitting}
					Creating...
				{:else}
					Create
				{/if}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
