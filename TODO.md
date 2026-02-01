# TODO

## Frontend - Facts Table

- [ ] **image_occlusion region counting may be incorrect** - `fact-display.ts` line 61 filters fields by `f.type === 'image_occlusion'` to count regions. Verify against actual image_occlusion content shape from backend (the regions are nested inside the field, not separate fields). The current logic likely returns 1 (the image field itself) rather than the actual region count.
- [ ] **`stripHtml` is naive** - `fact-display.ts` uses `/<[^>]*>/g` which breaks on edge cases like `<img alt="a>b">` or malformed HTML. Acceptable for display-only truncation but consider a proper sanitizer if user-authored HTML is allowed in fields.
- [ ] **Expanded row `<tr>` has no keyed association** - In `FactTableRow.svelte`, the `{#if expanded}` block renders a second `<tr>` outside the `{#each}` key scope. If Svelte reorders rows during keyed updates, the expanded row could theoretically detach. Not a practical issue with current structure but worth noting.
