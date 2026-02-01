<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { ChevronsLeftIcon, ChevronLeftIcon, ChevronRightIcon, ChevronsRightIcon } from '@lucide/svelte';

	interface Props {
		page: number;
		totalPages: number;
		total: number;
		pageSize: number;
	}

	let { page: currentPage, totalPages, total, pageSize }: Props = $props();

	let start = $derived(Math.min((currentPage - 1) * pageSize + 1, total));
	let end = $derived(Math.min(currentPage * pageSize, total));
	let isFirst = $derived(currentPage <= 1);
	let isLast = $derived(currentPage >= totalPages);

	let visiblePages = $derived.by(() => {
		const pages: number[] = [];
		let startPage = Math.max(1, currentPage - 2);
		let endPage = Math.min(totalPages, startPage + 4);
		startPage = Math.max(1, endPage - 4);
		for (let i = startPage; i <= endPage; i++) {
			pages.push(i);
		}
		return pages;
	});

	function navigate(p: number) {
		const url = new URL($page.url);
		if (p <= 1) {
			url.searchParams.delete('page');
		} else {
			url.searchParams.set('page', String(p));
		}
		goto(url.toString(), { replaceState: true });
	}
</script>

{#if totalPages > 1}
	<div class="flex items-center justify-between border-t border-border px-6 py-3">
		<span class="text-sm text-muted-foreground">
			Showing {start} to {end} of {total} facts
		</span>

		<div class="flex items-center gap-1">
			<button
				onclick={() => navigate(1)}
				disabled={isFirst}
				class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground hover:bg-muted disabled:pointer-events-none disabled:opacity-50"
				aria-label="First page"
			>
				<ChevronsLeftIcon class="h-4 w-4" />
			</button>
			<button
				onclick={() => navigate(currentPage - 1)}
				disabled={isFirst}
				class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground hover:bg-muted disabled:pointer-events-none disabled:opacity-50"
				aria-label="Previous page"
			>
				<ChevronLeftIcon class="h-4 w-4" />
			</button>

			{#each visiblePages as p}
				<button
					onclick={() => navigate(p)}
					class="flex h-8 w-8 items-center justify-center rounded-md text-sm {p === currentPage
						? 'bg-primary text-primary-foreground'
						: 'text-muted-foreground hover:bg-muted'}"
				>
					{p}
				</button>
			{/each}

			<button
				onclick={() => navigate(currentPage + 1)}
				disabled={isLast}
				class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground hover:bg-muted disabled:pointer-events-none disabled:opacity-50"
				aria-label="Next page"
			>
				<ChevronRightIcon class="h-4 w-4" />
			</button>
			<button
				onclick={() => navigate(totalPages)}
				disabled={isLast}
				class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground hover:bg-muted disabled:pointer-events-none disabled:opacity-50"
				aria-label="Last page"
			>
				<ChevronsRightIcon class="h-4 w-4" />
			</button>
		</div>
	</div>
{/if}
