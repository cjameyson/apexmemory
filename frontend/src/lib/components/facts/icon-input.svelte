<script lang="ts">
	import type { Component } from 'svelte';

	let {
		icon,
		placeholder = '',
		value = $bindable(''),
		multiline = false,
		resizable = false,
		oninput,
		error
	}: {
		icon: Component<{ class?: string }>;
		placeholder?: string;
		value?: string;
		multiline?: boolean;
		resizable?: boolean;
		oninput?: (value: string) => void;
		error?: string;
	} = $props();

	const inputBase = 'border-input bg-background text-foreground placeholder:text-muted-foreground focus:ring-ring w-full rounded-md border text-sm focus:ring-2 focus:outline-none';
	const textareaClasses = `${inputBase} resize-y py-2 pr-3 pl-9`;
	const inputClasses = `${inputBase} h-9 py-2 pr-3 pl-9`;

	function handleInput(e: Event) {
		const target = e.target as HTMLInputElement | HTMLTextAreaElement;
		value = target.value;
		oninput?.(value);
	}
</script>

<div class="space-y-1">
	<div class="relative">
		<div class="text-muted-foreground pointer-events-none absolute top-2.5 left-2.5">
			{@render iconSlot()}
		</div>
		{#if multiline}
			<textarea
				rows={2}
				class="{textareaClasses} {error ? 'border-destructive' : ''}"
				{placeholder}
				{value}
				oninput={handleInput}
			></textarea>
		{:else if resizable}
			<textarea
				rows={1}
				class="{textareaClasses} {error ? 'border-destructive' : ''}"
				{placeholder}
				{value}
				oninput={handleInput}
			></textarea>
		{:else}
			<input
				type="text"
				class="{inputClasses} {error ? 'border-destructive' : ''}"
				{placeholder}
				{value}
				oninput={handleInput}
			/>
		{/if}
	</div>
	{#if error}
		<p class="text-destructive text-xs">{error}</p>
	{/if}
</div>

{#snippet iconSlot()}
	{@const Icon = icon}
	<Icon class="h-4 w-4" />
{/snippet}
