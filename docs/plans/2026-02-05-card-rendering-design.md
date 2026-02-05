# Card Rendering Design: Rich Text & Image Occlusion in Focus Mode

**Date:** 2026-02-05
**Scope:** Study/review focus mode rendering only (no facts list or detail views)

## Overview

The study focus mode currently renders basic and cloze cards as plain text. This design adds:
1. Rich text rendering for basic card fields (TipTap JSON to HTML)
2. Image occlusion card rendering with region masking and reveal

## Architecture

### Data Flow

```
StudyCard.factContent (raw JSONB)
  -> extractCardDisplay() parses by factType
  -> CardDisplay (discriminated union)
  -> focus-mode.svelte routes to correct renderer
  -> ImageOcclusionCard.svelte or RichTextContent.svelte
```

### File Manifest

| Action | File | Purpose |
|--------|------|---------|
| Create | `src/lib/components/cards/image-occlusion-card.svelte` | Study-mode image occlusion renderer |
| Create | `src/lib/components/cards/rich-text-content.svelte` | TipTap JSON to rendered HTML |
| Modify | `src/lib/types/review.ts` | Discriminated union CardDisplay type |
| Modify | `src/lib/services/reviews.ts` | extractCardDisplay() for all 3 types |
| Modify | `src/lib/components/overlays/focus-mode.svelte` | Route to correct renderer |
| Modify | `src/lib/components/image-occlusion/types.ts` | Add RevealStyle type |

No backend changes required.

## CardDisplay Discriminated Union

Replace the current flat `CardDisplay` type with a discriminated union:

```typescript
interface BasicCardDisplay {
  type: 'basic';
  front: TipTapDocument | string;
  back: TipTapDocument | string;
}

interface ClozeCardDisplay {
  type: 'cloze';
  front: string;
  clozeAnswer: string;
}

interface ImageOcclusionCardDisplay {
  type: 'image_occlusion';
  title: string;
  image: ImageData;
  regions: Region[];
  targetRegionId: string;
  mode: OcclusionMode;
  revealStyle: RevealStyle;
}

type CardDisplay = BasicCardDisplay | ClozeCardDisplay | ImageOcclusionCardDisplay;
```

### extractCardDisplay() Changes

- **Basic:** Passes TipTap JSON documents through as-is (no longer strips to plain text)
- **Cloze:** Same logic, adds `type: 'cloze'` discriminant
- **Image occlusion:** Parses the `image_occlusion` field, maps `card.elementId` to `targetRegionId`, resolves image URL via `assetUrl()`, defaults `revealStyle` to `'show_label'`

## ImageOcclusionCard.svelte

### Props
- `display: ImageOcclusionCardDisplay`
- `isRevealed: boolean`
- `onReveal: () => void`

### Layout
- Image rendered at natural aspect ratio, constrained by max-width/max-height of available card area, centered
- Regions rendered as absolutely-positioned divs using percentage-based coordinates (region x/y/width/height relative to image dimensions)
- Image rotation via CSS `transform: rotate()`
- No zoom/pan in study mode

### Occlusion Modes

**hide_all_guess_one (before reveal):**
- All regions show as opaque colored rectangles
- Target region has distinct treatment (different color or pulsing border) to indicate "guess this one"

**hide_one_guess_one (before reveal):**
- Only the target region is masked
- All other regions are invisible (no overlay)

### Reveal Behavior

**show_label mode (default):**
- Target region mask removed or made semi-transparent
- Label text rendered centered inside region bounds
- White text on semi-transparent dark background for readability

**image_only mode:**
- Target region mask simply disappears, revealing underlying image

**Both modes:**
- In hide_all_guess_one, non-target regions stay masked after reveal
- Only the target region changes

### Hints
- If target region has a `hint` field, a small "Show hint" button appears below the image pre-reveal
- Tapping shows hint text inline
- Non-intrusive, on-request only

### Reveal Button
- "Tap or press space to reveal" positioned below the image inside the component
- Calls parent's `reveal()` function

## RichTextContent.svelte

### Approach
- Uses TipTap's `generateHTML(doc, extensions)` -- pure function, no editor instance
- Same extension set as the editor (StarterKit, Image, etc.)
- Renders via `{@html}` with scoped typography styles

### Props
- `doc: TipTapDocument | string` -- TipTap JSON or plain string fallback

### Styling
- Scoped CSS for headings, paragraphs, bold, italic, lists, images, code blocks
- Inherits focus-mode white-on-dark color scheme
- Images get `max-width: 100%` and rounded corners

### Fallback
- Plain string values rendered directly as text
- Malformed JSON falls back to `extractTextFromDoc()` plain text extraction

## RevealStyle Type Addition

```typescript
export type RevealStyle = 'show_label' | 'image_only';

export interface ImageOcclusionField {
  // ... existing fields
  revealStyle?: RevealStyle;  // optional, defaults to 'show_label'
}
```

Optional field -- no migration needed, existing facts default to `show_label`. Editor toggle to set this value is a separate future task.

## Focus Mode Integration

Template changes to `focus-mode.svelte`:
- Switch from `if (display.isCloze)` boolean to `switch (display.type)` discriminant
- New branch for `image_occlusion` renders `<ImageOcclusionCard>`
- Basic branch updated to use `<RichTextContent>` instead of plain text interpolation
- Rating buttons, undo, session stats, keyboard shortcuts -- all unchanged

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Labeled vs unlabeled diagrams | Per-fact `revealStyle` toggle, default `show_label` | Covers both use cases, zero config for default |
| Image layout | Full card area, no side panel | Immersive, matches focus mode aesthetic |
| Hints | On-request only (button) | Non-intrusive, preserves difficulty |
| Zoom/pan in study | Not included | Fit-to-area sufficient, revisit for mobile later |
| Region rendering | CSS absolute positioning (%) | Responsive, simpler than canvas/SVG |
