<script lang="ts">
	import { cn } from '$lib/utils';
	import * as Popover from '$lib/components/ui/popover';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import { SearchIcon, XIcon, ChevronDownIcon } from '@lucide/svelte';
	import {
		notebookEmojis,
		emojiCategories,
		searchEmojis,
		type NotebookEmoji
	} from '$lib/data/notebook-emojis';

	interface Props {
		value?: string;
		defaultValue?: string;
		placeholder?: string;
		class?: string;
	}

	let { value = $bindable(''), defaultValue = '📓', placeholder = 'Select', class: className }: Props = $props();

	// Initialize value if empty (one-time initialization on mount is intentional)
	// svelte-ignore state_referenced_locally
	if (!value && defaultValue) {
		value = defaultValue;
	}

	let open = $state(false);
	let searchQuery = $state('');
	let customInput = $state('');

	// Filter emojis based on search
	let filteredEmojis = $derived(searchEmojis(searchQuery));

	// Group filtered emojis by category
	let groupedEmojis = $derived.by(() => {
		const groups = new Map<string, NotebookEmoji[]>();
		for (const emoji of filteredEmojis) {
			const existing = groups.get(emoji.category) ?? [];
			existing.push(emoji);
			groups.set(emoji.category, existing);
		}
		return groups;
	});

	// Order categories as defined
	let orderedCategories = $derived(
		emojiCategories.filter((cat) => groupedEmojis.has(cat))
	);

	function selectEmoji(emoji: string) {
		value = emoji;
		open = false;
		searchQuery = '';
		customInput = '';
	}

	function clearEmoji() {
		value = '';
		open = false;
	}

	function handleCustomInput(e: Event) {
		const input = (e.target as HTMLInputElement).value;
		customInput = input;
		// Check if it's a valid emoji (rough check - at least one emoji character)
		const emojiRegex = /\p{Emoji}/u;
		if (emojiRegex.test(input)) {
			// Extract just the first emoji
			const match = input.match(/\p{Emoji_Presentation}|\p{Emoji}\uFE0F?/u);
			if (match) {
				value = match[0];
			}
		}
	}

	function handleOpenChange(isOpen: boolean) {
		open = isOpen;
		if (!isOpen) {
			searchQuery = '';
			customInput = '';
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			open = false;
		}
	}
</script>

<Popover.Root bind:open onOpenChange={handleOpenChange}>
	<Popover.Trigger>
		{#snippet child({ props })}
			<button
				{...props}
				type="button"
				class={cn(
					'inline-flex h-9 items-center justify-center gap-1 rounded-md border border-input bg-transparent px-2 text-xl shadow-xs transition-colors hover:bg-accent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
					className
				)}
			>
				{#if value}
					<span>{value}</span>
				{:else}
					<span class="text-sm text-muted-foreground">{placeholder}</span>
				{/if}
				<ChevronDownIcon class="size-3 text-muted-foreground" />
			</button>
		{/snippet}
	</Popover.Trigger>

	<Popover.Content class="w-80 p-0" align="start" sideOffset={4}>
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<!-- svelte-ignore a11y_interactive_supports_focus -->
		<div class="flex flex-col" role="dialog" aria-label="Emoji picker" onkeydown={handleKeydown}>
			<!-- Search input -->
			<div class="border-b p-2">
				<div class="relative">
					<SearchIcon class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
					<Input
						type="text"
						placeholder="Search emojis..."
						class="h-8 pl-8 text-sm"
						bind:value={searchQuery}
						onclick={(e) => e.stopPropagation()}
						onkeydown={(e) => e.stopPropagation()}
					/>
				</div>
			</div>

			<!-- Custom emoji input -->
			<div class="flex items-center gap-2 border-b px-3 py-2">
				<span class="text-xs text-muted-foreground">Custom:</span>
				<Input
					type="text"
					placeholder="Paste emoji"
					class="h-7 w-20 px-2 text-center text-base"
					value={customInput}
					oninput={handleCustomInput}
					onclick={(e) => e.stopPropagation()}
					onkeydown={(e) => e.stopPropagation()}
				/>
				{#if value}
					<Button
						variant="ghost"
						size="sm"
						class="ml-auto h-7 gap-1 px-2 text-xs text-muted-foreground"
						onclick={clearEmoji}
					>
						<XIcon class="size-3" />
						Clear
					</Button>
				{/if}
			</div>

			<!-- Emoji grid -->
			<div class="max-h-64 overflow-y-auto p-2">
				{#if orderedCategories.length === 0}
					<p class="py-4 text-center text-sm text-muted-foreground">No emojis match "{searchQuery}"</p>
				{:else}
					{#each orderedCategories as category}
						{@const emojis = groupedEmojis.get(category) ?? []}
						<div class="mb-3 last:mb-0">
							<h4 class="mb-1.5 px-1 text-xs font-medium text-muted-foreground">{category}</h4>
							<div class="grid grid-cols-8 gap-0.5">
								{#each emojis as emoji}
									<button
										type="button"
										class={cn(
											'flex size-8 items-center justify-center rounded text-lg transition-colors hover:bg-accent',
											value === emoji.emoji && 'bg-primary/10 ring-1 ring-primary'
										)}
										onclick={() => selectEmoji(emoji.emoji)}
										title={emoji.keywords.slice(0, 3).join(', ')}
									>
										{emoji.emoji}
									</button>
								{/each}
							</div>
						</div>
					{/each}
				{/if}
			</div>
		</div>
	</Popover.Content>
</Popover.Root>
