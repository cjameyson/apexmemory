<script lang="ts">
	import type { Region } from './types';
	import Input from '$lib/components/ui/input/input.svelte';
	import Textarea from '$lib/components/ui/textarea/textarea.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Trash2 } from '@lucide/svelte';

	interface Props {
		region: Region;
		index: number;
		isSelected: boolean;
		focusLabel?: boolean;
		labelError?: boolean;
		onSelect?: () => void;
		onLabelChange?: (value: string) => void;
		onHintChange?: (value: string) => void;
		onBackContentChange?: (value: string) => void;
		onDelete?: () => void;
	}

	let {
		region,
		index,
		isSelected,
		focusLabel = false,
		labelError = false,
		onSelect,
		onLabelChange,
		onHintChange,
		onBackContentChange,
		onDelete
	}: Props = $props();

	let labelInputRef: HTMLInputElement | null = $state(null);

	$effect(() => {
		if (focusLabel && isSelected && labelInputRef) {
			// Use tick to ensure DOM is updated after selected state renders inputs
			requestAnimationFrame(() => {
				labelInputRef?.focus();
				labelInputRef?.select();
			});
		}
	});
</script>

<div
	class="rounded-lg border p-2 transition-colors {isSelected
		? 'border-primary bg-primary/5'
		: labelError
			? 'border-destructive/50 bg-destructive/5'
			: 'border-transparent hover:bg-muted/50'}"
	onclick={onSelect}
	onkeydown={(e) => e.key === 'Enter' && onSelect?.()}
	role="button"
	tabindex="0"
>
	{#if isSelected}
		<!-- SELECTED: Full edit mode -->
		<div class="space-y-2">
			<!-- Top row: numbered badge + delete button -->
			<div class="flex items-center justify-between">
				<span
					class="flex h-5 w-5 shrink-0 items-center justify-center rounded bg-primary/10 text-xs font-medium text-primary"
				>
					{index}
				</span>
				<Button
					variant="ghost"
					size="icon-sm"
					class="h-6 w-6 text-muted-foreground hover:text-destructive"
					onclick={(e: MouseEvent) => {
						e.stopPropagation();
						onDelete?.();
					}}
				>
					<Trash2 class="h-3.5 w-3.5" />
				</Button>
			</div>

			<div>
				<span class="mb-1 block text-xs font-medium text-muted-foreground"
					>Label<span class="text-destructive"> *</span></span
				>
				<Input
					bind:ref={labelInputRef}
					value={region.label}
					placeholder="Label..."
					class="h-7 text-sm"
					aria-label="Region label"
					aria-invalid={labelError || undefined}
					onclick={(e: MouseEvent) => e.stopPropagation()}
					oninput={(e: Event) => onLabelChange?.((e.target as HTMLInputElement).value)}
				/>
			</div>

			<div>
				<span class="mb-1 block text-xs font-medium text-muted-foreground">Hint</span>
				<Input
					value={region.hint ?? ''}
					placeholder="Optional hint shown during review..."
					class="h-7 text-sm"
					aria-label="Region hint"
					onclick={(e: MouseEvent) => e.stopPropagation()}
					oninput={(e: Event) => onHintChange?.((e.target as HTMLInputElement).value)}
				/>
			</div>

			<div>
				<span class="mb-1 block text-xs font-medium text-muted-foreground">Back Content</span>
				<Textarea
					value={region.backContent ?? ''}
					placeholder="Additional content shown on answer..."
					class="min-h-[50px] resize-none bg-background text-sm"
					aria-label="Region back content"
					onclick={(e: MouseEvent) => e.stopPropagation()}
					oninput={(e: Event) =>
						onBackContentChange?.((e.target as HTMLTextAreaElement).value)}
				/>
			</div>
		</div>
	{:else}
		<!-- NOT SELECTED: Compact plain text display -->
		<div class="flex w-full items-center gap-2 text-left">
			<span
				class="flex h-5 w-5 shrink-0 items-center justify-center rounded bg-primary/10 text-xs font-medium text-primary"
			>
				{index}
			</span>
			<span class="flex-1 truncate text-sm text-foreground">
				{region.label || 'Untitled region'}
			</span>
		</div>
	{/if}
</div>
