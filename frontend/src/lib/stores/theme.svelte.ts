/**
 * Theme store for managing light/dark theme switching
 * Uses Svelte 5 runes for reactivity
 * Compatible with shadcn-svelte's .dark class approach
 */

export type Theme = 'light' | 'dark';

const STORAGE_KEY = 'theme';

function createThemeStore() {
	let current = $state<Theme>('light');

	// Initialize from localStorage or system preference
	function init() {
		if (typeof window === 'undefined') return;

		const stored = localStorage.getItem(STORAGE_KEY) as Theme | null;
		if (stored === 'light' || stored === 'dark') {
			current = stored;
		} else {
			const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
			current = prefersDark ? 'dark' : 'light';
		}

		applyTheme(current);

		// Listen for system preference changes when no stored preference
		window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
			if (!localStorage.getItem(STORAGE_KEY)) {
				const newTheme = e.matches ? 'dark' : 'light';
				current = newTheme;
				applyTheme(newTheme);
			}
		});
	}

	function applyTheme(theme: Theme) {
		if (typeof document === 'undefined') return;

		if (theme === 'dark') {
			document.documentElement.classList.add('dark');
		} else {
			document.documentElement.classList.remove('dark');
		}
	}

	function set(theme: Theme) {
		if (theme !== 'light' && theme !== 'dark') return;
		current = theme;
		localStorage.setItem(STORAGE_KEY, theme);
		applyTheme(theme);
	}

	function toggle() {
		set(current === 'dark' ? 'light' : 'dark');
	}

	function clear() {
		localStorage.removeItem(STORAGE_KEY);
		const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		const newTheme = prefersDark ? 'dark' : 'light';
		current = newTheme;
		applyTheme(newTheme);
	}

	return {
		get current() {
			return current;
		},
		get isDark() {
			return current === 'dark';
		},
		get isLight() {
			return current === 'light';
		},
		init,
		set,
		toggle,
		clear
	};
}

export const theme = createThemeStore();
