# Image Occlusion Editor - Phase 1 Implementation Plan

## Overview

Phase 1 creates the foundational layout, component structure, state architecture (including undo/redo infrastructure), and visual design system for the image occlusion editor.

## Deviations from Prompt

| Prompt Assumption | Actual Codebase | Resolution |
|-------------------|-----------------|------------|
| `asset_id` + CDN URLs | No assets table exists; images stored as external URLs | Use `url` field only, stub `asset_id` as optional |
| Region `label`, `hint`, `back_content` validated | Backend only extracts region IDs | Store in JSONB content (passes through unvalidated) |
| `mode` setting validated | Not validated by backend | Store in content, handle in frontend |

## File Structure

```
frontend/src/lib/components/image-occlusion/
├── index.ts                      # Barrel exports
├── types.ts                      # TypeScript interfaces
├── ImageOcclusionEditor.svelte   # Main container
├── EditorCanvas.svelte           # Left panel - image + SVG overlay
├── RegionOverlay.svelte          # Individual region SVG group
├── LabelPanel.svelte             # Right panel - region list
├── LabelPanelItem.svelte         # Single region metadata row
├── EditorToolbar.svelte          # Top toolbar
├── ImageUploader.svelte          # Empty state (placeholder for Phase 2)
├── StatusBar.svelte              # Bottom status bar
├── commands.ts                   # Command pattern implementations
├── history.svelte.ts             # Undo/redo state management
├── editor-state.svelte.ts        # Main editor state
├── coordinates.ts                # Coordinate transformation utilities
└── utils.ts                      # Helper functions (ID generation, debounce)
```

## Implementation Tasks

### Task 1: TypeScript Interfaces (`types.ts`)

Define all interfaces:
- `RectShape`, `Region`, `ArrowAnnotation`, `Annotation`
- `ImageData`, `ImageOcclusionField`, `OcclusionMode`
- `EditorTool`, `DisplayContext`, `Point`
- Component props: `ImageOcclusionEditorProps`, `EditorCanvasProps`, `RegionOverlayProps`, `LabelPanelProps`, `LabelPanelItemProps`, `EditorToolbarProps`, `StatusBarProps`
- `EditorCommand` interface, `ResizeHandle`, `ResizeHandlePosition`

### Task 2: Coordinate Utilities (`coordinates.ts`)

Implement:
- `imageToDisplay(point, ctx)` - Convert image coords to display coords
- `displayToImage(point, ctx)` - Convert display coords to image coords
- `regionToDisplay(shape, ctx)` - Convert full region shape
- `calculateScale(ctx)` - Compute scale factor for fit
- `getEffectiveImageDimensions(w, h, rotation)` - Handle rotation swap
- `calculateCenteredOffset(ctx)` - Initial centering
- `getResizeHandles(shape, handleSize)` - 8 resize handle positions
- `constrainToImageBounds(point, w, h)` - Clamp to image bounds

### Task 3: State Management (`editor-state.svelte.ts`)

Factory function pattern using Svelte 5 runes:
- Image state: `image`, `regions`, `annotations`, `mode`
- UI state: `selectedRegionId`, `activeTool`, `zoom`, `panOffset`, `containerSize`
- Derived: `displayContext`, `selectedRegion`
- Internal mutators (prefixed `_`): `_setImage`, `_addRegion`, `_updateRegion`, `_removeRegion`, etc.
- Public setters for UI state: `setSelectedRegionId`, `setActiveTool`, `setZoom`, etc.

### Task 4: History Management (`history.svelte.ts`)

Factory function with:
- `undoStack`, `redoStack` as `$state`
- `canUndo`, `canRedo` as `$derived`
- `execute(command)` - Run and push to undo stack, handle merging
- `undo()`, `redo()`, `clear()`
- `undoCount`, `redoCount` getters

### Task 5: Command Pattern (`commands.ts`)

Implement commands (Phase 1 stubs, full logic in Phase 3):
- `CreateRegionCommand` - Add region, select it
- `DeleteRegionCommand` - Remove region, handle selection
- `MoveRegionCommand` - Update position, support merging
- `ResizeRegionCommand` - Update shape, support merging
- `UpdateRegionMetadataCommand` - Update label/hint/back_content, support merging
- `RotateImageCommand` - Change rotation

### Task 6: Utilities (`utils.ts`)

- `generateRegionId()` - Returns `m_<nanoid(10)>`
- `generateAnnotationId()` - Returns `a_<nanoid(10)>`
- `debounce(fn, delay)` - For continuous operations

### Task 7: Main Container (`ImageOcclusionEditor.svelte`)

Layout:
```
┌─────────────────────────────────────────────────────────────────┐
│  EditorToolbar                                                   │
├───────────────────────────────────────────┬─────────────────────┤
│  EditorCanvas (flex-1)                    │  LabelPanel (320px)  │
├───────────────────────────────────────────┴─────────────────────┤
│  StatusBar                                                       │
└─────────────────────────────────────────────────────────────────┘
```

CSS Grid: `grid-template-rows: auto 1fr auto`, `grid-template-columns: 1fr 320px`

Props: `initialValue`, `onSave`, `onCancel`

Creates `editorState` and `history`, passes to children.

### Task 8: Editor Canvas (`EditorCanvas.svelte`)

Structure:
- Container `<div>` with ResizeObserver
- `<img>` with CSS transform for rotation/zoom
- SVG overlay (absolute positioned)
- Renders regions via `RegionOverlay` components

Phase 1: Static display only (interactions in Phase 3)

### Task 9: Region Overlay (`RegionOverlay.svelte`)

SVG `<g>` containing:
- `<rect>` for region (fill + stroke using design tokens)
- Marching ants border when selected (CSS animation)
- Resize handles when selected (8 small rects)

CSS keyframes for marching ants:
```css
@keyframes marching-ants {
  0% { stroke-dashoffset: 0; }
  100% { stroke-dashoffset: 16; }
}
```

### Task 10: Label Panel (`LabelPanel.svelte`)

- Filter input (shadcn `Input`)
- Scrollable list (shadcn `ScrollArea`)
- Empty state when no regions
- Fixed 320px width

### Task 11: Label Panel Item (`LabelPanelItem.svelte`)

- Region header with number
- Label input (required indicator)
- Hint input (optional)
- Back content textarea (optional)
- Delete button (icon)
- Selected state highlight

### Task 12: Editor Toolbar (`EditorToolbar.svelte`)

Button groups:
1. Undo/Redo (disabled when empty)
2. Tools: Select, Draw Region (Annotate dropdown placeholder)
3. Rotate
4. Zoom: -/+/Fit
5. Done (primary button)

Uses shadcn `Button`, `Separator`, `Tooltip`

### Task 13: Status Bar (`StatusBar.svelte`)

Displays:
- Region count: "N regions"
- Mode: "Hide All, Guess One"
- Zoom: "100%"
- Undo/Redo: "↶ N ↷ M"
- Dirty indicator when history has items

### Task 14: Image Uploader Placeholder (`ImageUploader.svelte`)

Empty state component (full implementation in Phase 2):
- Drop zone visual
- "Drop image or click to upload"
- Paste hint

### Task 15: Barrel Exports (`index.ts`)

Export all public components and types.

### Task 16: Mock Data & Static Demo

Add mock data for testing:
- Sample image (external URL)
- 3 sample regions with labels/hints
- One region pre-selected to show marching ants

## Styling Approach

- Use semantic design tokens: `bg-card`, `text-foreground`, `border-border`, `bg-primary`, etc.
- Dark mode via existing CSS variables (no new tokens)
- Region colors: `oklch(from var(--primary) l c h / 0.3)` for fill
- Marching ants: `var(--primary-foreground)` stroke with dash animation
- Resize handles: `bg-primary` small squares

## Dependencies

- `nanoid` - Already in project (check), or add if missing
- shadcn components: `Button`, `Input`, `Textarea`, `ScrollArea`, `Separator`, `Tooltip`

## Integration Point

**Modal approach** - Replace `ImageOcclusionPlaceholder` in `create-fact-modal.svelte`:
- File: `frontend/src/lib/components/facts/create-fact-modal.svelte`
- The existing placeholder shows "Coming Soon" message
- Replace with `ImageOcclusionEditor` when fact type is `image_occlusion`
- Editor should fill the modal content area

## Demo Mode

**Mock data only** for Phase 1:
- Hardcoded sample image URL (Wikipedia commons image)
- 3 pre-defined regions with sample labels
- One region pre-selected to demonstrate marching ants
- No image upload functionality yet (Phase 2)

## Exit Criteria

- [ ] All component files created with proper TypeScript interfaces
- [ ] Command pattern infrastructure in place
- [ ] Coordinate transformation utilities implemented
- [ ] Static layout renders correctly with mock data
- [ ] Marching ants animation working on selected region
- [ ] Undo/Redo buttons visible in toolbar (disabled state)
- [ ] Design tokens used consistently (no hardcoded colors)
- [ ] Visual design ready for review

## Verification

1. Start dev server: `make dev.up`
2. Navigate to notebook facts page
3. Click "Create Fact" and select "Image Occlusion"
4. Verify editor renders with mock data (sample image + 3 regions)
5. Verify selected region shows marching ants animation
6. Verify toolbar buttons show correct disabled states
7. Verify label panel shows all regions with metadata fields
8. Verify dark mode works correctly
