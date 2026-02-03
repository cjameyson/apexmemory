<script lang="ts">
	import { page } from '$app/state';
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import NotebooksDropdown from '$lib/components/navigation/notebooks-dropdown.svelte';
	import SearchTrigger from '$lib/components/navigation/search-trigger.svelte';
	import ReviewLauncher from '$lib/components/navigation/review-launcher.svelte';
	import { HomeIcon, BrainIcon, SettingsIcon, LogOutIcon } from '@lucide/svelte';
	import type { Notebook, ReviewScope, StudyCard } from '$lib/types';
	import type { User } from '$lib/api/types';

	interface Props {
		user: User;
		notebooks: Notebook[];
		currentNotebook?: Notebook;
		onStartFocusMode?: (scope: ReviewScope, cards: StudyCard[]) => void;
		onOpenSearch?: () => void;
		onCreateNotebook?: () => void;
	}

	let {
		user,
		notebooks,
		currentNotebook,
		onStartFocusMode,
		onOpenSearch,
		onCreateNotebook
	}: Props = $props();

	let isHome = $derived(page.url.pathname === '/home');
	let isInNotebook = $derived(page.url.pathname.startsWith('/notebooks/'));

	function handleOpenSearch() {
		onOpenSearch?.();
	}

	function handleStartReview(scope: ReviewScope, cards: StudyCard[]) {
		onStartFocusMode?.(scope, cards);
	}

	function getInitials(u: User): string {
		if (u.username) {
			return u.username.slice(0, 2).toUpperCase();
		}
		return u.email.slice(0, 2).toUpperCase();
	}

	async function handleLogout() {
		await fetch('/api/auth/logout', { method: 'POST' });
		window.location.href = '/login';
	}
</script>

<header
	class="shrink-0 border-b border-border bg-card px-2 py-2"
>
	<div class="flex items-center justify-between">
		<!-- Left group -->
		<div class="flex items-center gap-1">
			<!-- Logo -->
			<a
				href="/home"
				class="flex items-center gap-2 rounded-lg px-2 py-1.5 text-foreground transition-colors hover:bg-accent"
			>
				<BrainIcon class="text-primary size-6" />
			</a>

			<!-- Notebooks dropdown -->
			<NotebooksDropdown {notebooks} current={currentNotebook} {isInNotebook} {onCreateNotebook} />
		</div>

		<!-- Right group -->
		<div class="flex items-center gap-2">
			<!-- Search trigger -->
			<SearchTrigger onclick={handleOpenSearch} />

			<!-- Review launcher -->
			<ReviewLauncher {notebooks} {currentNotebook} onStartReview={handleStartReview} />

			<!-- User avatar dropdown -->
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<button
							{...props}
							class="bg-primary hover:bg-primary/80 flex size-8 items-center justify-center rounded-full text-sm font-medium text-white transition-colors"
						>
							{getInitials(user)}
						</button>
					{/snippet}
				</DropdownMenu.Trigger>

				<DropdownMenu.Content align="end" class="w-48">
					<DropdownMenu.Label class="font-normal">
						<div class="flex flex-col space-y-1">
							<p class="text-sm font-medium">{user.username || 'User'}</p>
							<p class="truncate text-xs text-muted-foreground">{user.email}</p>
						</div>
					</DropdownMenu.Label>
					<DropdownMenu.Separator />
					<DropdownMenu.Item class="cursor-pointer gap-2">
						<SettingsIcon class="size-4" />
						Settings
					</DropdownMenu.Item>
					<DropdownMenu.Separator />
					<DropdownMenu.Item
						class="cursor-pointer gap-2 text-destructive"
						onclick={handleLogout}
					>
						<LogOutIcon class="size-4" />
						Logout
					</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
	</div>
</header>
