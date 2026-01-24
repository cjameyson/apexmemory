<script lang="ts">
	import { page } from '$app/stores';
	import { cn } from '$lib/utils';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import NotebooksDropdown from '$lib/components/navigation/notebooks-dropdown.svelte';
	import SearchTrigger from '$lib/components/navigation/search-trigger.svelte';
	import ReviewLauncher from '$lib/components/navigation/review-launcher.svelte';
	import { HomeIcon, BrainIcon, SettingsIcon, LogOutIcon } from '@lucide/svelte';
	import type { Notebook, ReviewScope } from '$lib/types';
	import type { User } from '$lib/api/types';
	import { appState } from '$lib/stores/app.svelte';

	interface Props {
		user: User;
		notebooks: Notebook[];
		currentNotebook?: Notebook;
	}

	let { user, notebooks, currentNotebook }: Props = $props();

	let isHome = $derived($page.url.pathname === '/home');
	let isInNotebook = $derived($page.url.pathname.startsWith('/notebooks/'));

	function handleOpenSearch() {
		appState.openCommandPalette();
	}

	function handleStartReview(scope: ReviewScope) {
		appState.startFocusMode(scope);
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
	class="shrink-0 border-b border-slate-200 bg-white px-2 py-2 dark:border-slate-800 dark:bg-slate-900"
>
	<div class="flex items-center justify-between">
		<!-- Left group -->
		<div class="flex items-center gap-1">
			<!-- Logo -->
			<a
				href="/home"
				class="flex items-center gap-2 rounded-lg px-2 py-1.5 text-slate-900 transition-colors hover:bg-slate-100 dark:text-white dark:hover:bg-slate-800"
			>
				<BrainIcon class="size-6 text-sky-500" />
			</a>

			<!-- Notebooks dropdown -->
			<NotebooksDropdown {notebooks} current={currentNotebook} {isInNotebook} />
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
							class="flex size-8 items-center justify-center rounded-full bg-sky-500 text-sm font-medium text-white transition-colors hover:bg-sky-600"
						>
							{getInitials(user)}
						</button>
					{/snippet}
				</DropdownMenu.Trigger>

				<DropdownMenu.Content align="end" class="w-48">
					<DropdownMenu.Label class="font-normal">
						<div class="flex flex-col space-y-1">
							<p class="text-sm font-medium">{user.username || 'User'}</p>
							<p class="truncate text-xs text-slate-500">{user.email}</p>
						</div>
					</DropdownMenu.Label>
					<DropdownMenu.Separator />
					<DropdownMenu.Item class="cursor-pointer gap-2">
						<SettingsIcon class="size-4" />
						Settings
					</DropdownMenu.Item>
					<DropdownMenu.Separator />
					<DropdownMenu.Item
						class="cursor-pointer gap-2 text-red-600 dark:text-red-400"
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
