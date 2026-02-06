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
		onBackExtraChange?: (value: string) => void;
		onDelete?: () => void;
		onAdvanceNext?: () => void;
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
		onBackExtraChange,
		onDelete,
		onAdvanceNext
	}: Props = $props();

	let labelInputRef: HTMLInputElement | null = $state(null);

	let isUnlabeled = $derived(!region.label.trim());

	// Color classes based on state
	let badgeClasses = $derived(
		isSelected
			? 'bg-success text-success-foreground'
			: isUnlabeled
				? 'bg-warning text-warning-foreground'
				: 'bg-primary text-primary-foreground'
	);

	let wrapperClasses = $derived(
		isSelected
			? 'border-success bg-success/5'
			: labelError
				? 'border-destructive/50 bg-destructive/5'
				: isUnlabeled
					? 'border-warning/50 bg-warning/5'
					: 'border-transparent hover:bg-muted/50'
	);

	$effect(() => {
		if (focusLabel && isSelected && labelInputRef) {
			requestAnimationFrame(() => {
				labelInputRef?.focus();
				labelInputRef?.select();
			});
		}
	});

	function handleLabelKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && e.shiftKey) {
			e.preventDefault();
			e.stopPropagation();
			onAdvanceNext?.();
		}
	}
</script>

<div
	class="rounded-lg border p-2 transition-colors {wrapperClasses}"
	onclick={onSelect}
	onkeydown={(e) => e.key === 'Enter' && onSelect?.()}
	role="button"
	tabindex="0"
>
	{#if isSelected}
		<!-- SELECTED: Full edit mode -->
		<div class="space-y-1.5">
			<!-- Label field with badge -->
			<div class="flex items-center gap-2">
				<span
					class="flex h-5 w-5 shrink-0 items-center justify-center rounded text-[10px] font-bold {badgeClasses}"
				>
					{index}
				</span>
				<span class="w-12 shrink-0 text-right text-xs font-medium text-muted-foreground"
					>Label<span class="text-destructive"> *</span></span
				>
				<Input
					bind:ref={labelInputRef}
					value={region.label}
					placeholder="Label..."
					class="h-7 flex-1 text-sm"
					aria-label="Region label"
					aria-invalid={labelError || undefined}
					onclick={(e: MouseEvent) => e.stopPropagation()}
					oninput={(e: Event) => onLabelChange?.((e.target as HTMLInputElement).value)}
					onkeydown={handleLabelKeydown}
				/>
			</div>

			<!-- Hint field (indented to align with Label input) -->
			<div class="flex items-center gap-2 pl-7">
				<span class="w-12 shrink-0 text-right text-xs font-medium text-muted-foreground"
					>Hint</span
				>
				<Input
					value={region.hint ?? ''}
					placeholder="Hint for review..."
					class="h-7 flex-1 text-sm"
					aria-label="Region hint"
					onclick={(e: MouseEvent) => e.stopPropagation()}
					oninput={(e: Event) => onHintChange?.((e.target as HTMLInputElement).value)}
				/>
			</div>

			<!-- Back Extra field (indented to align) -->
			<div class="flex items-start gap-2 pl-7">
				<span class="w-12 shrink-0 pt-1.5 text-right text-xs font-medium text-muted-foreground"
					>Extra</span
				>
				<Textarea
					value={region.backExtra ?? ''}
					placeholder="Additional answer content..."
					class="min-h-[40px] flex-1 resize-none bg-background text-sm"
					aria-label="Region back extra"
					onclick={(e: MouseEvent) => e.stopPropagation()}
					oninput={(e: Event) =>
						onBackExtraChange?.((e.target as HTMLTextAreaElement).value)}
				/>
			</div>

			<!-- Delete button (bottom-right) -->
			<div class="flex justify-end pl-7">
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
		</div>
	{:else}
		<!-- NOT SELECTED: Compact inline display -->
		<div class="flex w-full items-center gap-2 text-left">
			<span
				class="flex h-5 w-5 shrink-0 items-center justify-center rounded text-[10px] font-bold {badgeClasses}"
			>
				{index}
			</span>
			<span class="flex-1 truncate text-sm {isUnlabeled ? 'italic text-warning' : 'text-foreground'}">
				{region.label || 'Untitled region'}
			</span>
		</div>
	{/if}
</div>
