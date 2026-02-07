# SvelteKit 2 + Svelte 5 URL, navigation, and history management (modern-only)

## 0) The "source of truth" rule
- **Durable + shareable state → URL** (path params + `?searchParams`). Survives SSR/reload and enables deep links.
- **Transient UI state → history state (`page.state`)** via `pushState/replaceState` (e.g., modals, drawers, selected row). Back/forward restores it.
- **Ephemeral "draft" DOM state you want to restore → snapshots** (notes editor drafts, chat composer, form inputs).
- **Session-persistent preferences → global store** (sidebar collapsed, theme). These persist across navigations but NOT tied to history (back/forward won't toggle them).

---

## 1) URL design for an "entities + panes" app
A clean baseline (example structure - adapt to your needs):
- `/app/notebooks`
- `/app/notebooks/[notebookId]` (overview)
- `/app/notebooks/[notebookId]/sources/[sourceId]?tab=transcript|summary`
- `/app/notebooks/[notebookId]/flashcards?due=today&sort=next`
- `/app/notebooks/[notebookId]/notes/[noteId]`
- `/app/notebooks/[notebookId]/chat/[threadId]`
- `/app/notebooks/[notebookId]/stats?range=30d`

Use **nested layouts** for your app shell (sidebar / notebook switcher) so the UI stays stable across navigations.

---

## 2) Reading the current URL (modern)
- In components, prefer `$app/state`:
  ```ts
  import { page, navigating } from '$app/state';
  // page.url, page.params, page.state, etc.
  ```

* **Important**: Changes to `page` require Svelte 5 runes. Legacy `$:` syntax will NOT react to changes:
  ```ts
  const id = $derived(page.params.id);     // ✅ Correct - updates reactively
  $: badId = page.params.id;               // ❌ Wrong - never updates after initial load
  ```

* In `load`, use the `url` argument (don't reach into component state).

---

## 3) Navigation: one API, predictable history

### Prefer links for normal navigation

Use `<a href>` for internal routes; customize behavior with `data-sveltekit-*` link options. ([Svelte][1])

### Use `goto()` for programmatic navigation

```ts
import { goto } from '$app/navigation';

await goto(`/app/notebooks/${id}`, {
  replaceState: false,
  keepFocus: false,
  noScroll: false
});
```

* Use `replaceState: true` for “filter/sort/tab” changes to avoid history spam.

---

## 4) URL search params for filters, tabs, sorting (workhorse pattern)

When controls change *view* state, put it in `searchParams` and typically **replace** history:

```ts
import { page } from '$app/state';
import { goto } from '$app/navigation';

function setSort(sort: string) {
  const url = new URL(page.url);
  url.searchParams.set('sort', sort);
  goto(url, { replaceState: true, keepFocus: true, noScroll: true });
}
```

---

## 5) Shallow routing + `page.state` for overlays (modals/drawers/details)

Use `pushState`/`replaceState` (not `history.pushState`) to make “Back closes modal” work cleanly. ([Svelte][2])

**Same-URL state change (modal open):**

```ts
import { pushState } from '$app/navigation';
import { page } from '$app/state';

function openQuickView(sourceId: string) {
  // first arg is the URL (relative). To stay on the current URL, use ''.
  pushState('', { quickView: { sourceId } });
}

// {#if page.state.quickView} ... {/if}
```

**Optional: change the URL shallowly (shareable, but still “overlay UX”):**

```ts
pushState(`/app/notebooks/${id}/sources/${sourceId}`, { quickView: { sourceId } });
```

(You can choose whether the URL changes; both are supported.) ([Svelte][2])

### Caveat: `page.state` is transient on full reload

There is a known issue where, after a full page refresh, `page.state` may be empty even though related data still exists inside `history.state`. Don’t rely on `page.state` as durable state; persist anything important in the URL (or reconstruct it on load if needed). ([GitHub][3])

---

## 6) Preserve drafts with snapshots (notes, chat composer, forms)

Snapshots must be exported from **`+page.svelte` or `+layout.svelte`** (not from an arbitrary child component). They’re stored as JSON (persisted in `sessionStorage`) and can be restored on **back/forward** and **page reload**, even when returning from another site. ([Svelte][4])

```ts
import type { Snapshot } from './$types';

let draft = $state('');

export const snapshot: Snapshot<string> = {
  capture: () => draft,
  restore: (v) => (draft = v)
};
```

---

## 7) Data refresh when the URL changes (and after mutations)

* Put URL-dependent fetching in `load` so navigation drives data.
* After create/edit/delete operations, use targeted invalidation (`invalidate(...)`) or `invalidateAll()` when needed.

---

## 8) Preloading: speed vs wasted work

`data-sveltekit-preload-data` can load data for navigations that never happen (“false positives”). For heavily real-time or expensive-load screens, consider:

* switching to `data-sveltekit-preload-data="tap"` (less eager), or
* using `data-sveltekit-preload-code` to preload *code* without prefetching *data*. ([Svelte][5])



[1]: https://svelte.dev/docs/kit/link-options?utm_source=chatgpt.com "Link options • SvelteKit Docs"
[2]: https://svelte.dev/docs/kit/shallow-routing "Shallow routing • SvelteKit Docs"
[3]: https://github.com/sveltejs/kit/issues/11956 "`$page.state` is lost after page refresh · Issue #11956 · sveltejs/kit · GitHub"
[4]: https://svelte.dev/docs/kit/snapshots "Snapshots • SvelteKit Docs"
[5]: https://svelte.dev/tutorial/kit/preload "Link options / Preloading • Svelte Tutorial"


