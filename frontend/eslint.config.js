import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import svelte from 'eslint-plugin-svelte';
import globals from 'globals';

export default tseslint.config(
	eslint.configs.recommended,
	...tseslint.configs.recommended,
	...svelte.configs['flat/recommended'],
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node
			}
		}
	},
	{
		files: ['**/*.svelte'],
		languageOptions: {
			parserOptions: {
				parser: tseslint.parser
			}
		}
	},
	{
		rules: {
			// TypeScript
			'@typescript-eslint/no-unused-vars': [
				'warn',
				{
					argsIgnorePattern: '^_',
					varsIgnorePattern: '^_'
				}
			],
			'@typescript-eslint/no-explicit-any': 'warn',

			// Svelte - disable strict navigation rule for internal links
			'svelte/no-at-html-tags': 'warn',
			'svelte/no-navigation-without-resolve': 'off',

			// General
			'no-console': ['warn', { allow: ['warn', 'error'] }]
		}
	},
	{
		// Disable parsing for .svelte.ts files (Svelte 5 runes)
		files: ['**/*.svelte.ts'],
		rules: {
			// These files use Svelte 5 runes which aren't standard TS
		}
	},
	{
		ignores: [
			'.svelte-kit/**',
			'build/**',
			'dist/**',
			'node_modules/**',
			'*.config.js',
			'*.config.ts',
			'**/*.svelte.ts' // Ignore runes files for now
		]
	}
);
