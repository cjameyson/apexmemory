<script lang="ts">
	import { Info } from '@lucide/svelte';
	import { cn } from '$lib/utils';
	import MiniToolbar from './mini-toolbar.svelte';
	import IconInput from './icon-input.svelte';
	import Button from '$lib/components/ui/button/button.svelte';

	export interface ClozeFactData {
		text: string;
		backExtra: string;
	}

	let {
		initialData,
		onchange,
		oncardcount,
		errors
	}: {
		initialData?: Partial<ClozeFactData>;
		onchange: (data: ClozeFactData) => void;
		oncardcount?: (count: number) => void;
		errors?: Partial<Record<keyof ClozeFactData, string>>;
	} = $props();

	// initialData is captured at mount time only.
	// Parent must destroy/recreate this component (via {#if} or {#key}) to reset.
	let text = $state(initialData?.text ?? '');
	let backExtra = $state(initialData?.backExtra ?? '');

	const textareaBase = 'border-input bg-background text-foreground placeholder:text-muted-foreground focus:ring-ring w-full resize-y rounded-md border p-3 text-sm focus:ring-2 focus:outline-none';

	let textareaRef: HTMLTextAreaElement | undefined = $state();

	// Parse cloze numbers from text
	let usedNumbers = $derived.by(() => {
		const nums = new Set<number>();
		const regex = /\{\{c(\d+)::/g;
		let match;
		while ((match = regex.exec(text)) !== null) {
			nums.add(parseInt(match[1], 10));
		}
		return nums;
	});

	$effect(() => {
		oncardcount?.(usedNumbers.size);
	});

	// Buttons c1-c3 plus dynamic extras
	let clozeButtons = $derived.by(() => {
		const maxUsed = usedNumbers.size > 0 ? Math.max(...usedNumbers) : 0;
		const count = Math.max(3, maxUsed);
		return Array.from({ length: count }, (_, i) => i + 1);
	});

	function notify() {
		onchange({ text, backExtra });
	}

	function insertCloze(num: number) {
		const el = textareaRef;
		if (!el) return;

		el.focus();
		const start = el.selectionStart;
		const end = el.selectionEnd;
		const selected = text.substring(start, end);

		if (selected) {
			const replacement = `{{c${num}::${selected}}}`;
			text = text.substring(0, start) + replacement + text.substring(end);
			// Position cursor after the replacement
			const newPos = start + replacement.length;
			requestAnimationFrame(() => {
				el.setSelectionRange(newPos, newPos);
			});
		} else {
			const insertion = `{{c${num}::}}`;
			text = text.substring(0, start) + insertion + text.substring(end);
			// Position cursor between :: and }}
			const cursorPos = start + insertion.length - 2;
			requestAnimationFrame(() => {
				el.setSelectionRange(cursorPos, cursorPos);
			});
		}
		notify();
	}

	function insertNextCloze() {
		let next = 1;
		while (usedNumbers.has(next)) next++;
		insertCloze(next);
	}

	function handleKeydown(e: KeyboardEvent) {
		if ((e.metaKey || e.ctrlKey) && e.key >= '1' && e.key <= '9') {
			e.preventDefault();
			insertCloze(parseInt(e.key, 10));
		}
	}

	export function focus() {
		textareaRef?.focus();
	}
</script>

<div class="space-y-4">
	<div class="space-y-1">
		<label class="text-sm font-medium" for="cloze-text">Cloze Text</label>
		<div class="flex flex-wrap items-center gap-1">
			<MiniToolbar />
			<div class="bg-border mx-1 h-4 w-px"></div>
			{#each clozeButtons as num}
				{@const isUsed = usedNumbers.has(num)}
				<Button
					variant={isUsed ? 'secondary' : 'outline'}
					size="sm"
					class="h-7 px-2 font-mono text-xs {isUsed ? 'text-muted-foreground' : 'text-primary'}"
					onclick={() => insertCloze(num)}
				>
					c{num}
				</Button>
			{/each}
			<Button
				variant="outline"
				size="sm"
				class="h-7 px-2 text-xs"
				onclick={insertNextCloze}
				title="Insert next cloze"
			>
				+
			</Button>
		</div>
		<textarea
			id="cloze-text"
			bind:this={textareaRef}
			rows={2}
			class={cn(textareaBase, 'font-mono', errors?.text && 'border-destructive focus:ring-destructive')}
			placeholder={'The {{c1::answer}} is hidden...'}
			bind:value={text}
			oninput={notify}
			onkeydown={handleKeydown}
		></textarea>
		{#if errors?.text}
			<p class="text-destructive text-xs">{errors.text}</p>
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
</div>
