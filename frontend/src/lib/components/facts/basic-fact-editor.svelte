<script lang="ts">
	import { Info, Lightbulb } from '@lucide/svelte';
	import { cn } from '$lib/utils';
	import MiniToolbar from './mini-toolbar.svelte';
	import IconInput from './icon-input.svelte';

	export interface BasicFactData {
		front: string;
		back: string;
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
	let front = $state(initialData?.front ?? '');
	// svelte-ignore state_referenced_locally
	let back = $state(initialData?.back ?? '');
	// svelte-ignore state_referenced_locally
	let backExtra = $state(initialData?.backExtra ?? '');
	// svelte-ignore state_referenced_locally
	let hint = $state(initialData?.hint ?? '');

	const textareaBase = 'border-input bg-background text-foreground placeholder:text-muted-foreground focus:ring-ring w-full resize-y rounded-md border p-3 text-sm focus:ring-2 focus:outline-none';
	const errorClasses = 'border-destructive focus:ring-destructive';

	let frontRef: HTMLTextAreaElement | undefined = $state();

	function notify() {
		onchange({ front, back, backExtra, hint });
	}

	export function focus() {
		frontRef?.focus();
	}
</script>

<div class="space-y-4">
	<div class="space-y-1">
		<label class="text-sm font-medium" for="basic-front">Front</label>
		<MiniToolbar />
		<textarea
			id="basic-front"
			bind:this={frontRef}
			rows={2}
			class={cn(textareaBase, errors?.front && errorClasses)}
			placeholder="Question or prompt..."
			bind:value={front}
			oninput={notify}
		></textarea>
		{#if errors?.front}
			<p class="text-destructive text-xs">{errors.front}</p>
		{/if}
	</div>

	<div class="space-y-1">
		<label class="text-sm font-medium" for="basic-back">Back</label>
		<MiniToolbar />
		<textarea
			id="basic-back"
			rows={2}
			class={cn(textareaBase, errors?.back && errorClasses)}
			placeholder="Answer..."
			bind:value={back}
			oninput={notify}
		></textarea>
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
