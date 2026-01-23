import { test, expect } from '@playwright/test';

test.describe('Home Page', () => {
	test('should load the home page', async ({ page }) => {
		await page.goto('/');
		await expect(page).toHaveTitle(/Apex Memory/);
	});

	test('should toggle theme', async ({ page }) => {
		await page.goto('/');

		// Find theme toggle button
		const themeToggle = page.getByRole('button', { name: /switch to/i });
		await expect(themeToggle).toBeVisible();

		// Check initial state (light mode)
		const html = page.locator('html');
		await expect(html).not.toHaveClass(/dark/);

		// Click to switch to dark mode
		await themeToggle.click();
		await expect(html).toHaveClass(/dark/);

		// Click to switch back to light mode
		await themeToggle.click();
		await expect(html).not.toHaveClass(/dark/);
	});

	test('should have skip link for accessibility', async ({ page }) => {
		await page.goto('/');

		const skipLink = page.getByRole('link', { name: 'Skip to main content' });
		await expect(skipLink).toBeAttached();

		// Skip link should be visible when focused
		await skipLink.focus();
		await expect(skipLink).toBeVisible();
	});
});
