<script lang="ts">
	import type { JSONContent } from '@tiptap/core';
	import { Info, Lightbulb } from '@lucide/svelte';
	import { RichTextEditor } from '$lib/components/rich-text';
	import IconInput from './icon-input.svelte';

	export interface BasicFactData {
		front: JSONContent | null;
		back: JSONContent | null;
		backExtra: string;
		hint: string;
	}

	let {
		initialData,
		onchange,
		errors
	}: {
		initialData?: Partial<BasicFactData>;
		onchange: (data: BasicFactData) => void;
		errors?: Partial<Record<keyof BasicFactData, string>>;
	} = $props();

	// initialData is captured at mount time only.
	// Parent must destroy/recreate this component (via {#if} or {#key}) to reset.
	// svelte-ignore state_referenced_locally
	let front: JSONContent | null = $state(initialData?.front ?? null);
	// svelte-ignore state_referenced_locally
	let back: JSONContent | null = $state(initialData?.back ?? null);
	// svelte-ignore state_referenced_locally
	let backExtra = $state(initialData?.backExtra ?? '');
	// svelte-ignore state_referenced_locally
	let hint = $state(initialData?.hint ?? '');

	let frontEditor: RichTextEditor | undefined = $state();

	function notify() {
		onchange({ front, back, backExtra, hint });
	}

	export function focus() {
		frontEditor?.focus();
	}
</script>

<div class="space-y-4">
	<div class="space-y-1">
		<!-- svelte-ignore a11y_label_has_associated_control -->
		<label class="text-sm font-medium">Front</label>
		<RichTextEditor
			bind:this={frontEditor}
			content={front}
			onchange={(c) => { front = c; notify(); }}
			placeholder="Question or prompt..."
			class={errors?.front ? 'border-destructive focus-within:ring-destructive' : ''}
		/>
		{#if errors?.front}
			<p class="text-destructive text-xs">{errors.front}</p>
		{/if}
	</div>

	<div class="space-y-1">
		<!-- svelte-ignore a11y_label_has_associated_control -->
		<label class="text-sm font-medium">Back</label>
		<RichTextEditor
			content={back}
			onchange={(c) => { back = c; notify(); }}
			placeholder="Answer..."
			class={errors?.back ? 'border-destructive focus-within:ring-destructive' : ''}
		/>
		{#if errors?.back}
			<p class="text-destructive text-xs">{errors.back}</p>
		{/if}
	</div>

	<IconInput
		icon={Info}
		placeholder="Additional info shown on the back..."
		bind:value={backExtra}
		resizable
		oninput={notify}
		error={errors?.backExtra}
	/>

	<IconInput
		icon={Lightbulb}
		placeholder="Hint (shown on request)..."
		bind:value={hint}
		oninput={notify}
		error={errors?.hint}
	/>
</div>
