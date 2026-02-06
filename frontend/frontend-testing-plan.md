# Frontend Testing Plan

## Testing Layers for SvelteKit

There are three distinct layers. Use the right tool at each.

### 1. Pure Logic — Vitest, Node environment

Unit tests for service functions, transforms, stores, utilities — anything that doesn't touch the DOM. This is what `reviews.test.ts`, `hooks.server.test.ts`, and `theme.test.ts` do today.

- Runs fast (sub-second)
- No browser needed
- `jsdom` is fine here for store tests that touch `document`/`window` lightly

### 2. Component Testing — Vitest Browser Mode

Svelte 5 runes and `mount()` don't work correctly in jsdom — you get "mount(...) is not available on the server" errors, and simulated DOM behavior diverges from real browsers in subtle ways (focus, accessibility, event propagation).

The modern approach uses **Vitest Browser Mode** with Playwright as the provider and `vitest-browser-svelte` for rendering.

#### Setup

```bash
npm install -D @vitest/browser vitest-browser-svelte playwright
npx playwright install chromium
```

#### Multi-project vitest config

```ts
// vitest.config.ts
export default defineConfig({
  test: {
    projects: [
      {
        // Component tests — run in real Chromium
        test: {
          name: 'client',
          environment: 'browser',
          browser: {
            enabled: true,
            provider: 'playwright',
            instances: [{ browser: 'chromium' }],
          },
          include: ['src/**/*.svelte.{test,spec}.ts'],
        },
      },
      {
        // Server/logic tests — run in Node
        test: {
          name: 'server',
          environment: 'node',
          include: ['src/**/*.{test,spec}.ts'],
        },
      },
    ],
  },
});
```

#### Component test example

```ts
// MyComponent.svelte.test.ts  (note: .svelte.test.ts enables runes)
import { render } from 'vitest-browser-svelte';
import { page } from '@vitest/browser/context';
import MyComponent from './MyComponent.svelte';

it('renders and handles clicks', async () => {
  render(MyComponent, { props: { label: 'Click me' } });

  const button = page.getByRole('button');
  await expect.element(button).toHaveTextContent('Click me');
  await button.click();
  // assert outcomes...
});
```

#### Key differences from `@testing-library/svelte` + jsdom

- **Real browser APIs** — focus, events, accessibility all work correctly
- **`page.getBy*()` locators** — auto-retry, no flaky `querySelector` calls
- **`flushSync()`** needed for external `$state` updates to reflect in DOM
- **`.svelte.test.ts` extension** — enables runes in test files

### 3. End-to-End — Playwright (full app)

For testing complete user flows (login -> create notebook -> add facts -> review), Playwright runs against the actual dev server.

```bash
npx playwright test
```

SvelteKit scaffolds this with `npx sv create`. Tests hit real pages, real API proxy routes, real backend. Target flows:

- Auth flow (login, session persistence, logout)
- Navigation and layout behavior
- Review session from start to finish
- Mobile viewport behavior

## Priority Components for Browser-Based Tests

| Component | Why |
|---|---|
| `RichTextEditor` | TipTap integration, paste/drop handlers, external content sync — jsdom can't meaningfully test this |
| `FactExpandedContent` | AbortController cleanup, loading states, DOM updates on prop change |
| `command-palette` | Keyboard interactions, scope switching |
| Focus mode / review session | Multi-step user flow with state transitions |

## Strategy

Keep `jsdom` for server/logic tests (hooks, services). Add Vitest Browser Mode for component tests where real DOM behavior matters. Add Playwright E2E for critical user flows.

## References

- [Testing - Svelte Docs](https://svelte.dev/docs/svelte/testing)
- [Vitest Component Testing Guide](https://vitest.dev/guide/browser/component-testing)
- [From JSDOM to Real Browsers: Testing Svelte with Vitest Browser Mode](https://scottspence.com/posts/testing-with-vitest-browser-svelte-guide)
- [Migrating from @testing-library/svelte to vitest-browser-svelte](https://scottspence.com/posts/migrating-from-testing-library-svelte-to-vitest-browser-svelte)
