# Occlusion Mask Effects — Design

**Date:** 2026-02-05
**Goal:** Experiment with visual effects on image occlusion masks during review, find the best combination, then bake the winner into the component and remove the debug tooling.

## Approach

Build a dev-only debug panel in focus mode that exposes independent toggles for mask color, animations, and transitions. All effects are CSS-driven — no JS animation logic. The panel is tree-shaken from production builds via `import.meta.env.DEV`.

## Debug State Store

**File:** `src/lib/stores/debug-occlusion.ts`

```typescript
type MaskColor = 'blue' | 'amber' | 'violet' | 'rose';
type PulseMode = 'off' | 'subtle' | 'pronounced';
type MarchingAnts = 'off' | 'slow' | 'medium';
type EntranceMode = 'none' | 'fade-in' | 'scale-up';
type RevealMode = 'none' | 'fade-cross' | 'dissolve' | 'slide-away';

interface OcclusionDebugState {
  enabled: boolean;
  maskColor: MaskColor;
  showQuestionMark: boolean;
  pulse: PulseMode;
  marchingAnts: MarchingAnts;
  entrance: EntranceMode;
  reveal: RevealMode;
}
```

Defaults match current production behavior: `blue`, question mark on, everything else off. When `enabled` is `false`, the card component ignores all debug values — zero regression risk.

## Toggles

| Control | Options | Default |
|---------|---------|---------|
| Mask color | `blue`, `amber`, `violet`, `rose` | `blue` |
| Question mark | on / off | on |
| Pulse | `off`, `subtle`, `pronounced` | `off` |
| Marching ants | `off`, `slow`, `medium` | `off` |
| Card entrance | `none`, `fade-in`, `scale-up` | `none` |
| Reveal transition | `none`, `fade-cross`, `dissolve`, `slide-away` | `none` |

All toggles are independent — any combination is valid.

## Effect Specifications

### Color Presets

Map to Tailwind color families for both the mask background and border:

| Preset | Background | Border |
|--------|-----------|--------|
| blue | `bg-sky-500` | `border-sky-300` |
| amber | `bg-amber-500` | `border-amber-300` |
| violet | `bg-violet-500` | `border-violet-300` |
| rose | `bg-rose-500` | `border-rose-300` |

Non-target masks remain `bg-slate-600 / border-slate-500` regardless of color preset.

### Pulse

CSS `@keyframes` on the target mask, pre-reveal only.

- **subtle:** opacity 0.85 to 1.0 over 2s ease-in-out, infinite.
- **pronounced:** opacity 0.7 to 1.0 + scale(1.0 to 1.02) over 1.5s ease-in-out, infinite.

### Marching Ants

Pure CSS using four `background-image` linear gradients (one per edge) with animated `background-position`. Applied as a pseudo-element or additional layer on the target mask, pre-reveal only.

- **slow:** 2s linear infinite cycle.
- **medium:** 0.8s linear infinite cycle.

Dashes are white/semi-transparent, 2px thick, ~8px dash + 8px gap.

### Card Entrance

Svelte `in:` transition on the region overlay wrapper when the card mounts.

- **fade-in:** opacity 0 to 1 over 400ms ease-out.
- **scale-up:** scale(0.9) + opacity(0) to scale(1) + opacity(1) over 300ms ease-out.

### Reveal Transition

Requires restructuring the region overlay from conditional `{#if}/{:else if}` blocks to a **single persistent div** that changes classes based on state. This lets CSS transitions interpolate between masked and revealed states.

- **fade-cross:** 400ms ease-out opacity crossfade. Mask color fades out while green highlight fades in.
- **dissolve:** 500ms. Mask background-color transitions to transparent, revealing the image beneath. Green border fades in after a 200ms delay.
- **slide-away:** 300ms. Mask translateY(-100%) with opacity fade-out. Green highlight fades in immediately.

## Structural Change: Persistent Region Div

Current (`{#if}` conditional — DOM nodes destroyed/created):
```svelte
{#if shouldShowMask(isTarget)}
  <div class="mask-styles">...</div>
{:else if isTarget && isRevealed}
  <div class="reveal-styles">...</div>
{/if}
```

New (single div, class-driven transitions):
```svelte
<div class="w-full h-full rounded-sm transition-all duration-300
  {getMaskClasses(isTarget, isRevealed, debugState)}">
  {#if !isRevealed && isTarget}
    <span class="?-indicator">?</span>
  {/if}
</div>
```

A helper function `getMaskClasses()` returns the appropriate Tailwind classes based on the current state and debug settings. The div is always present, so `transition-all` works as expected.

## Debug Panel

**File:** `src/lib/components/debug/occlusion-debug-panel.svelte`

- Rendered inside `focus-mode.svelte`, only when `import.meta.env.DEV` is true and the current card is `image_occlusion` type.
- Fixed position, top-right corner of focus mode overlay, `z-50`.
- Collapsed by default — small gear icon. Toggle with `D` key.
- ~240px wide when expanded, `bg-black/80 backdrop-blur-sm`, rounded.
- Each control: label on left, cycle-button on right (click to advance).
- Color options show a small colored dot preview.
- "Reset defaults" button snaps all values back to production defaults.
- `D` key is ignored during post-reveal rating state to avoid conflicts with future shortcuts.

## Files to Create/Modify

| File | Action |
|------|--------|
| `src/lib/stores/debug-occlusion.ts` | Create — debug state store |
| `src/lib/components/debug/occlusion-debug-panel.svelte` | Create — floating panel UI |
| `src/lib/components/cards/image-occlusion-card.svelte` | Modify — persistent div structure, read debug state, apply effect classes |
| `src/lib/components/overlays/focus-mode.svelte` | Modify — render debug panel, add `D` keydown handler |
| `src/app.css` | Modify — add `@keyframes` for pulse, marching ants, entrance/reveal animations |

## After Experimentation

Once a winning combination is chosen:
1. Bake the chosen effect values directly into the card component as defaults.
2. Delete `debug-occlusion.ts` and `occlusion-debug-panel.svelte`.
3. Remove debug panel rendering from `focus-mode.svelte`.
4. Remove the `D` keydown handler.
5. Clean up any unused `@keyframes` from `app.css`.
