<script lang="ts">
	import type { FactType } from '$lib/types/fact';
	import { Layers, TextCursorInput, Image, Braces, Layers3 } from '@lucide/svelte';

	let {
		selected,
		onchange,
		disabled = false
	}: {
		selected: FactType;
		onchange: (type: FactType) => void;
		disabled?: boolean;
	} = $props();

	const types: { value: FactType; label: string; description: string; icon: typeof Layers }[] = [
		{ value: 'basic', label: 'Basic', description: 'Front & back', icon: Layers },
		{ value: 'cloze', label: 'Cloze', description: 'Fill in the blank', icon: Braces },
		{
			value: 'image_occlusion',
			label: 'Image Occlusion',
			description: 'Hide image regions',
			icon: Image
		}
	];

	function handleKeydown(e: KeyboardEvent) {
		if (disabled) return;
		const currentIndex = types.findIndex((t) => t.value === selected);
		let nextIndex = currentIndex;

		if (e.key === 'ArrowRight' || e.key === 'ArrowDown') {
			e.preventDefault();
			nextIndex = (currentIndex + 1) % types.length;
		} else if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') {
			e.preventDefault();
			nextIndex = (currentIndex - 1 + types.length) % types.length;
		}

		if (nextIndex !== currentIndex) {
			onchange(types[nextIndex].value);
		}
	}
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div role="radiogroup" aria-label="Fact type" class="flex gap-3" onkeydown={handleKeydown}>
	{#each types as type (type.value)}
		{@const isSelected = selected === type.value}
		<button
			role="radio"
			aria-checked={isSelected}
			{disabled}
			tabindex={isSelected ? 0 : -1}
			class="flex flex-1 items-center gap-2.5 rounded-lg border px-3 py-2 text-left transition-all
				{isSelected
				? 'ring-primary bg-primary/5 border-primary/50 ring-2'
				: 'border-border hover:border-primary/50'}
				{disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'}"
			onclick={() => !disabled && onchange(type.value)}
		>
			<type.icon class="text-muted-foreground h-4 w-4 shrink-0" />
			<div class="min-w-0">
				<div class="text-sm leading-tight font-medium">{type.label}</div>
				<div class="text-muted-foreground text-xs leading-tight">{type.description}</div>
			</div>
		</button>
	{/each}
</div>
