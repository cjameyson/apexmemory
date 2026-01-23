<script lang="ts">
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button';
	import {
		LayoutDashboardIcon,
		BookOpenIcon,
		FolderIcon,
		SettingsIcon,
	} from '@lucide/svelte';

	let { children } = $props();

	const navItems = [
		{ href: '/dashboard', label: 'Dashboard', icon: LayoutDashboardIcon },
		{ href: '/study', label: 'Study', icon: BookOpenIcon },
		{ href: '/decks', label: 'Decks', icon: FolderIcon },
		{ href: '/settings', label: 'Settings', icon: SettingsIcon },
	];

	function isActive(href: string): boolean {
		return $page.url.pathname === href || $page.url.pathname.startsWith(href + '/');
	}
</script>

<div class="flex flex-col md:flex-row gap-6">
	<!-- Sidebar Navigation -->
	<aside class="md:w-56 shrink-0">
		<nav class="flex md:flex-col gap-1">
			{#each navItems as item}
				<Button
					href={item.href}
					variant={isActive(item.href) ? 'secondary' : 'ghost'}
					class="justify-start gap-2 w-full"
					size="sm"
				>
					<item.icon class="size-4" />
					<span class="hidden md:inline">{item.label}</span>
				</Button>
			{/each}
		</nav>
	</aside>

	<!-- Main Content -->
	<div class="flex-1 min-w-0">
		{@render children()}
	</div>
</div>
