import { describe, it, expect, beforeEach, vi } from 'vitest';
import { theme } from '$lib/stores/theme.svelte';

describe('Theme Store', () => {
	beforeEach(() => {
		// Reset DOM
		document.documentElement.classList.remove('dark');
		// Reset localStorage mock
		vi.mocked(localStorage.getItem).mockReturnValue(null);
		vi.mocked(localStorage.setItem).mockClear();
	});

	it('should default to light theme when no stored preference', () => {
		vi.mocked(localStorage.getItem).mockReturnValue(null);
		vi.mocked(window.matchMedia).mockImplementation(() => ({
			matches: false,
			media: '',
			onchange: null,
			addListener: vi.fn(),
			removeListener: vi.fn(),
			addEventListener: vi.fn(),
			removeEventListener: vi.fn(),
			dispatchEvent: vi.fn()
		}));

		theme.init();
		expect(theme.current).toBe('light');
		expect(theme.isDark).toBe(false);
		expect(theme.isLight).toBe(true);
	});

	it('should respect stored dark theme preference', () => {
		vi.mocked(localStorage.getItem).mockReturnValue('dark');

		theme.init();
		expect(theme.current).toBe('dark');
		expect(theme.isDark).toBe(true);
		expect(document.documentElement.classList.contains('dark')).toBe(true);
	});

	it('should toggle between light and dark', () => {
		vi.mocked(localStorage.getItem).mockReturnValue(null);
		theme.init();

		theme.toggle();
		expect(theme.current).toBe('dark');
		expect(localStorage.setItem).toHaveBeenCalledWith('theme', 'dark');

		theme.toggle();
		expect(theme.current).toBe('light');
		expect(localStorage.setItem).toHaveBeenCalledWith('theme', 'light');
	});

	it('should set specific theme', () => {
		theme.set('dark');
		expect(theme.current).toBe('dark');
		expect(localStorage.setItem).toHaveBeenCalledWith('theme', 'dark');
		expect(document.documentElement.classList.contains('dark')).toBe(true);

		theme.set('light');
		expect(theme.current).toBe('light');
		expect(document.documentElement.classList.contains('dark')).toBe(false);
	});

	it('should replace existing matchMedia listener on repeated init', () => {
		const addEventListener = vi.fn();
		const removeEventListener = vi.fn();

		vi.mocked(window.matchMedia).mockImplementation(() => ({
			matches: false,
			media: '(prefers-color-scheme: dark)',
			onchange: null,
			addListener: vi.fn(),
			removeListener: vi.fn(),
			addEventListener,
			removeEventListener,
			dispatchEvent: vi.fn()
		}));

		theme.init();
		theme.init();

		expect(addEventListener).toHaveBeenCalledTimes(2);
		expect(removeEventListener).toHaveBeenCalledTimes(1);
	});
});
