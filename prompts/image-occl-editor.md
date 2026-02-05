# Image Occlusion Editor - AI Coding Agent Prompt

<context>
You are building a professional-grade image occlusion editor for Apex Memory, a spaced repetition learning application. This editor allows users to create flashcards by marking regions on images that will be hidden during review (similar to Anki's image occlusion feature, but modernized).

The editor is part of a larger SvelteKit application using:
- **Svelte 5** with runes (`$state`, `$derived`, `$effect`, `$props`)
- **TailwindCSS 4+**
- **shadcn-svelte** component library
- **TypeScript**

This is a **desktop-only** component. Do not implement mobile or tablet layouts.

Review the project files for additional context on conventions, schema patterns, and existing UI patterns.
</context>

<project_files>
Reference these files before implementation:
- backend related code and sql files for facts and cards
- `references/svelte-llms-small.md` - Quick Svelte 5 API reference
- `references/svelte-llms-full.md` - Full Svelte 5 API reference
</project_files>

<design_tokens>
Use the application's existing design tokens and color palette consistently.
- Primary action colors
- Selection/highlight states
- Background and border colors
- Text hierarchy colors
- Interactive state colors (hover, focus, active, disabled)

Do NOT hardcode color values. Use Tailwind classes that align with the existing app theme.
</design_tokens>

<data_schema>
The image occlusion fact type uses this content structure:

```typescript
interface ImageOcclusionField {
  name: "image";
  type: "image_occlusion";
  image: {
    asset_id: string;        // Backend asset reference, UUID
    url: string;             // CDN URL for display
    width: number;           // Original image width in pixels
    height: number;          // Original image height in pixels
    original_name: string;
    rotation?: 0 | 90 | 180 | 270;  // Applied rotation
  };
  settings: {
    mode: "hide_all_guess_one" | "hide_one_guess_one";
  };
  regions: Region[];
  annotations?: Annotation[];
}

interface Region {
  id: string;                // Format: m_<nanoid> (e.g., "m_k7Xp2mQ9Ab")
  shape: RectShape;          // Extensible for future shapes
  label: string;             // Required - the answer
  hint?: string;             // Optional hint shown during review
  back_content?: string;     // Optional extra content shown after answer
}

interface RectShape {
  type: "rect";
  x: number;      // Pixels from left edge of original image
  y: number;      // Pixels from top edge of original image
  width: number;  // Width in pixels (original image coordinates)
  height: number; // Height in pixels (original image coordinates)
}

// Future shape types (stub only, do not implement)
interface EllipseShape {
  type: "ellipse";
  cx: number; cy: number;  // Center in image pixels
  rx: number; ry: number;  // Radii in image pixels
}

interface PolygonShape {
  type: "polygon";
  points: Array<{ x: number; y: number }>;  // Each point in image pixels
}

interface Annotation {
  id: string;  // Format: a_<nanoid>
  type: 'arrow';  // Extensible for future: 'text', 'highlight', 'freehand'
  points: {
    start: { x: number; y: number };  // Image pixel coordinates
    end: { x: number; y: number };
  };
  style: {
    color: string;  // Hex color or design token reference
    thickness: 1 | 2 | 3;  // Line weight
  };
}
```

**Coordinate System**: All coordinates are stored in **original image pixel space**. When rendering at different display sizes, the editor must scale coordinates appropriately. This ensures:
- Coordinates remain meaningful regardless of display size
- Easier integration with backend OCR/AI that works in pixel coordinates
- No precision loss from percentage conversions
</data_schema>

<undo_redo_architecture>
Undo/Redo is a **core requirement**, not an afterthought. Design all state mutations with this in mind from Phase 1.

### Command Pattern Approach
```typescript
interface EditorCommand {
  type: string;
  execute(): void;
  undo(): void;
  // Optional: for command merging (e.g., continuous drag updates)
  mergeWith?(other: EditorCommand): EditorCommand | null;
}

interface EditorHistory {
  undoStack: EditorCommand[];
  redoStack: EditorCommand[];
  maxHistorySize: number;  // e.g., 50 commands
}

// Example commands:
// - CreateRegionCommand
// - DeleteRegionCommand
// - MoveRegionCommand
// - ResizeRegionCommand
// - UpdateRegionMetadataCommand
// - RotateImageCommand
// - CreateAnnotationCommand
// - DeleteAnnotationCommand
// - BatchCommand (for grouping multiple operations)
```

### State Management Implications
- All state mutations must go through command execution
- Direct state mutation is prohibited except within command.execute()
- Consider command merging for continuous operations (dragging should not create 100 undo steps)
- Undo/redo must restore full state including selection

### UI Requirements
- Toolbar buttons for Undo/Redo with disabled states
- Keyboard shortcuts: Ctrl+Z (undo), Ctrl+Y or Ctrl+Shift+Z (redo)
- Visual feedback when undo/redo stack is empty
</undo_redo_architecture>

<design_requirements>
**Visual Design Principles:**
- Clean, professional interface matching the existing app aesthetic
- Use established design tokens - do not introduce new colors
- Minimal visual noise - let the image be the focus
- Clear visual hierarchy between selected and unselected regions
- Smooth transitions and micro-interactions

**Region Styling:**
- Unselected regions: Semi-transparent fill, solid border (use app's secondary/muted colors)
- Selected region: Highlighted fill with marching ants animated border (use app's primary/accent colors)
- All regions must allow seeing the underlying image content

**Accessibility:**
- Full keyboard navigation support
- ARIA labels for all interactive elements
- Focus indicators that meet WCAG guidelines
</design_requirements>

---

## Phase 1: Design, Layout & State Architecture

<phase_1_objective>
Create the foundational layout, component structure, state architecture (including undo/redo infrastructure), and visual design system. Focus on structure, styling, and the command pattern foundation.
</phase_1_objective>

<phase_1_deliverables>
1. **Component file structure** in `$lib/components/image-occlusion/`
2. **Main editor layout** with two-panel design (image canvas + label sidebar)
3. **Toolbar component** with placeholder buttons (including undo/redo)
4. **State management architecture** with command pattern for undo/redo
5. **Coordinate transformation utilities** (image coords ↔ display coords)
6. **Static mockup** demonstrating the visual design with hardcoded sample data
7. **TypeScript interfaces** for all component props, state, and commands

(note: use an img element + svg overlay, not an actual <canvas> element)
</phase_1_deliverables>

<phase_1_specifications>
### File Structure
```
src/lib/components/image-occlusion/
├── ImageOcclusionEditor.svelte    # Main container component
├── EditorCanvas.svelte            # Left panel - image + regions + annotations
├── RegionOverlay.svelte           # Individual region rectangle
├── LabelPanel.svelte              # Right panel - region list + metadata
├── LabelPanelItem.svelte          # Single region's metadata row
├── EditorToolbar.svelte           # Top toolbar
├── ImageUploader.svelte           # Upload/paste interface (empty state)
├── types.ts                       # TypeScript interfaces
├── commands.ts                    # Command pattern implementations
├── history.svelte.ts              # Undo/redo state management
├── editor-state.svelte.ts         # Main editor state using runes
├── coordinates.ts                 # Coordinate transformation utilities
└── utils.ts                       # Helper functions (ID generation, etc.)
```

### Rendering Architecture

Use an `<img>` element for the base image with an SVG overlay for all interactive elements:

```svelte
<div class="editor-canvas relative overflow-hidden">
  <!-- Base image layer -->
  <img 
    src={image.url} 
    alt={image.original_name}
    style:transform="rotate({image.rotation}deg) scale({zoom})"
    class="select-none"
  />
  
  <!-- Interactive SVG overlay (regions, annotations, handles) -->
  <svg class="absolute inset-0 w-full h-full overflow-visible">
    <!-- Regions rendered as <rect> elements -->
    <!-- Annotations rendered as <line>, <path>, etc. -->
    <!-- Resize handles as small <rect> elements -->
  </svg>
</div>
```

**Rationale:**
- SVG provides native DOM events on each shape (no manual hit-testing)
- CSS animations work directly (marching ants via `stroke-dashoffset`)
- Resize handles are trivial SVG rects with drag handlers
- Zoom/rotation via CSS transforms with coordinate math in JS
- Better accessibility (semantic elements, screen reader support)
- If pixel access is needed later (AI/OCR preprocessing), render to an offscreen canvas separately

**Do NOT use a single `<canvas>` element** — the manual redraw loop, hit-testing, and animation management adds significant complexity without benefit for this use case.

### Layout Specifications (Desktop Only)
APPROXIMATE: ascii art is not perfect
```
┌─────────────────────────────────────────────────────────────────────────┐
│  EditorToolbar                                                          │
│  [Undo][Redo] | [Select][Draw Region][Annotate ▼] | [Rotate ↻][Zoom] | [Done] │
├─────────────────────────────────────────┬───────────────────────────────┤
│                                         │  LabelPanel                   │
│                                         │  ┌─────────────────────────┐  │
│                                         │  │ 🔍 Filter regions...    │  │
│                                         │  ├─────────────────────────┤  │
│  EditorCanvas (fills remaining)         │  │ Region 1 - Mitochondria │  │
│  ┌───────────────────────────────────┐  │  │ hint: Powerhouse...     │  │
│  │                                   │  │  │                    [🗑]  │  │
│  │      [Image with regions]         │  │  ├─────────────────────────┤  │
│  │                                   │  │  │ Region 2 - Nucleus      │  │
│  │   ┌──────┐                        │  │  │ hint: Control center    │  │
│  │   │region│  ┌────────┐            │  │  │                    [🗑]  │  │
│  │   └──────┘  │selected│ ← marching │  │  ├─────────────────────────┤  │
│  │             │ants    │   border   │  │  │ Region 3 - ...          │  │
│  │             └────────┘            │  │  └─────────────────────────┘  │
│  │                                   │  │                               │
│  └───────────────────────────────────┘  │  [+ Add Region Manually]      │
│                                         │                               │
├─────────────────────────────────────────┴───────────────────────────────┤
│  Status bar: "3 regions • Hide All, Guess One • 100%" | Undo: 5 Redo: 0 │
└─────────────────────────────────────────────────────────────────────────┘
```

### Coordinate Transformation System
```typescript
// coordinates.ts

interface DisplayContext {
  imageWidth: number;      // Original image dimensions
  imageHeight: number;
  displayWidth: number;    // Current rendered dimensions
  displayHeight: number;
  rotation: 0 | 90 | 180 | 270;
  zoom: number;
  panOffset: { x: number; y: number };
}

// Convert image coordinates to display coordinates for rendering
function imageToDisplay(
  point: { x: number; y: number },
  ctx: DisplayContext
): { x: number; y: number };

// Convert display coordinates (mouse events) to image coordinates for storage
function displayToImage(
  point: { x: number; y: number },
  ctx: DisplayContext
): { x: number; y: number };

// Convert a full region shape between coordinate systems
function regionToDisplay(region: RectShape, ctx: DisplayContext): RectShape;
function regionToImage(region: RectShape, ctx: DisplayContext): RectShape;
```

### Command Pattern Foundation
```typescript
// commands.ts

interface EditorCommand {
  readonly type: string;
  readonly description: string;  // For potential UI display
  execute(): void;
  undo(): void;
  // For merging continuous operations (e.g., dragging)
  canMergeWith?(other: EditorCommand): boolean;
  mergeWith?(other: EditorCommand): EditorCommand;
}

// Base class or factory for common command types
abstract class RegionCommand implements EditorCommand {
  constructor(
    protected editorState: EditorState,
    protected regionId: string
  ) {}
  abstract get type(): string;
  abstract get description(): string;
  abstract execute(): void;
  abstract undo(): void;
}

// history.svelte.ts
function createEditorHistory(maxSize = 50) {
  let undoStack = $state<EditorCommand[]>([]);
  let redoStack = $state<EditorCommand[]>([]);
  
  const canUndo = $derived(undoStack.length > 0);
  const canRedo = $derived(redoStack.length > 0);
  
  function execute(command: EditorCommand) {
    command.execute();
    // Handle merging for continuous operations
    const lastCommand = undoStack[undoStack.length - 1];
    if (lastCommand?.canMergeWith?.(command)) {
      undoStack[undoStack.length - 1] = lastCommand.mergeWith!(command);
    } else {
      undoStack = [...undoStack, command].slice(-maxSize);
    }
    redoStack = [];  // Clear redo stack on new action
  }
  
  function undo() { /* ... */ }
  function redo() { /* ... */ }
  function clear() { /* ... */ }
  
  return {
    get canUndo() { return canUndo; },
    get canRedo() { return canRedo; },
    get undoCount() { return undoStack.length; },
    get redoCount() { return redoStack.length; },
    execute,
    undo,
    redo,
    clear
  };
}
```

### Component Props API Design

```typescript
// types.ts

// Main editor props
interface ImageOcclusionEditorProps {
  initialValue?: ImageOcclusionField;  // For editing existing facts
  onSave: (value: ImageOcclusionField) => void;
  onCancel: () => void;
  uploadEndpoint: string;  // Backend endpoint for image upload
}

// Canvas props
interface EditorCanvasProps {
  image: ImageData | null;
  regions: Region[];
  annotations: Annotation[];
  selectedRegionId: string | null;
  activeTool: EditorTool;
  zoom: number;
  displayContext: DisplayContext;
  onRegionSelect: (id: string | null) => void;
  onRegionCreate: (region: Region) => void;
  onRegionMove: (id: string, newPosition: { x: number; y: number }) => void;
  onRegionResize: (id: string, newShape: RectShape) => void;
  onAnnotationCreate: (annotation: Annotation) => void;
}

// Label panel props
interface LabelPanelProps {
  regions: Region[];
  selectedRegionId: string | null;
  onSelect: (id: string) => void;
  onUpdate: (id: string, updates: Partial<Region>) => void;
  onDelete: (id: string) => void;
}

type EditorTool = 'select' | 'draw' | 'annotate-arrow';
```

### Styling Requirements

1. **Marching Ants Animation** - CSS animation for selected region borders:
```css
@keyframes marching-ants {
  0% { stroke-dashoffset: 0; }
  100% { stroke-dashoffset: 16; }
}

.marching-ants {
  stroke-dasharray: 8 8;
  animation: marching-ants 0.5s linear infinite;
}
```

2. **Design Token Usage** - Reference existing app patterns:
```svelte
<!-- Use app's established color classes, e.g.: -->
<rect class="fill-primary/40 stroke-primary" />           <!-- unselected -->
<rect class="fill-primary/60 stroke-primary-foreground" /> <!-- selected -->
```

3. **Panel Sizing** - Label panel: 320px fixed width, canvas fills remaining space
</phase_1_specifications>

<phase_1_implementation_notes>
- Use shadcn-svelte `Button`, `Input`, `ScrollArea`, and `Tooltip` components
- The canvas should use SVG overlays on top of the image for precise region rendering
- Implement image scaling that maintains aspect ratio and centers in available space
- Use Svelte 5 runes throughout - no legacy reactive syntax
- All state mutations must be designed to work with the command pattern
- Coordinate transforms are critical - test thoroughly with different image sizes
</phase_1_implementation_notes>

<phase_1_exit_criteria>
- [x] All component files created with proper TypeScript interfaces
- [x] Command pattern infrastructure in place (commands.ts, history.svelte.ts)
- [x] Coordinate transformation utilities implemented and tested
- [x] Static layout renders correctly with mock data
- [x] Marching ants animation working on "selected" region
- [x] Undo/Redo buttons visible in toolbar (disabled state)
- [x] Design tokens used consistently (no hardcoded colors)
- [x] Visual design approved by human reviewer
</phase_1_exit_criteria>

**PHASE 1 COMPLETE**

---

## Phase 2: Image Handling & Core State

<phase_2_objective>
Implement image upload, paste, display, rotation functionality, and wire up the core editor state with undo/redo support.
</phase_2_objective>

<phase_2_deliverables>
1. **Drag-and-drop upload** with visual feedback
2. **Clipboard paste** support for images
3. **Image preview** with proper scaling and centering
4. **Rotation controls** (90° increments) with undo support
5. **Backend upload integration** (return asset_id)
6. **Loading and error states**
7. **Core state management** wired to command pattern
</phase_2_deliverables>

<phase_2_specifications>
### ImageUploader Component (Empty State)
When no image is loaded, show an upload interface:
- Large drop zone with dashed border
- Icon + "Drop image here or click to upload"
- "Or paste from clipboard (Ctrl+V)"
- Accepted formats: PNG, JPG, WebP, GIF
- Max file size: 10MB
- Future placeholder: "Search web for images" (disabled button with "coming soon" indicator)

### Upload Flow
```
1. User drops/selects/pastes image
2. Show local preview immediately (optimistic UI)
3. Upload to backend endpoint: POST /api/assets/upload
   - Request: FormData with 'file' field
   - Response: { asset_id: uuid (string), url: string, width: number, height: number }
4. On success: Store asset_id, update image URL to CDN
5. On failure: Show error toast, allow retry
```

### Rotation Implementation
- Store rotation as metadata (0, 90, 180, 270 degrees)
- Apply rotation via CSS transform on display
- **Rotation is an undoable command**:
```typescript
class RotateImageCommand implements EditorCommand {
  type = 'rotate-image';
  description = 'Rotate image';
  
  constructor(
    private state: EditorState,
    private previousRotation: number,
    private newRotation: number
  ) {}
  
  execute() {
    this.state.image.rotation = this.newRotation;
  }
  
  undo() {
    this.state.image.rotation = this.previousRotation;
  }
}
```

### Coordinate Handling with Rotation
When rotation is applied, the display context changes but stored coordinates remain in original image space:
```typescript
// Rotation affects how we transform coordinates
// 90° CW: display (x, y) maps to image (y, height - x)
// 180°: display (x, y) maps to image (width - x, height - y)
// 270° CW: display (x, y) maps to image (width - y, x)

function getEffectiveImageDimensions(
  originalWidth: number,
  originalHeight: number,
  rotation: 0 | 90 | 180 | 270
): { width: number; height: number } {
  if (rotation === 90 || rotation === 270) {
    return { width: originalHeight, height: originalWidth };
  }
  return { width: originalWidth, height: originalHeight };
}
```

### Core State Structure
```typescript
// editor-state.svelte.ts
function createEditorState(initialValue?: ImageOcclusionField) {
  // Image state
  let image = $state<ImageData | null>(initialValue?.image ?? null);
  
  // Regions stored in image coordinates
  let regions = $state<Region[]>(initialValue?.regions ?? []);
  
  // Annotations stored in image coordinates
  let annotations = $state<Annotation[]>(initialValue?.annotations ?? []);
  
  // UI state (not persisted, not undoable)
  let selectedRegionId = $state<string | null>(null);
  let activeTool = $state<EditorTool>('select');
  let zoom = $state(1);
  let panOffset = $state({ x: 0, y: 0 });
  
  // Settings
  let mode = $state<'hide_all_guess_one' | 'hide_one_guess_one'>(
    initialValue?.settings?.mode ?? 'hide_all_guess_one'
  );
  
  // Computed display context
  const displayContext = $derived<DisplayContext>({
    imageWidth: image?.width ?? 0,
    imageHeight: image?.height ?? 0,
    displayWidth: /* computed from container */,
    displayHeight: /* computed from container */,
    rotation: image?.rotation ?? 0,
    zoom,
    panOffset
  });
  
  return {
    // State accessors
    get image() { return image; },
    get regions() { return regions; },
    get annotations() { return annotations; },
    get selectedRegionId() { return selectedRegionId; },
    // ... etc
    
    // Mutators (called by commands only)
    _setImage(img: ImageData) { image = img; },
    _setRegions(r: Region[]) { regions = r; },
    _addRegion(r: Region) { regions = [...regions, r]; },
    _updateRegion(id: string, updates: Partial<Region>) { /* ... */ },
    _removeRegion(id: string) { regions = regions.filter(r => r.id !== id); },
    // ... etc
    
    // UI state (directly mutable, not undoable)
    setSelectedRegionId(id: string | null) { selectedRegionId = id; },
    setActiveTool(tool: EditorTool) { activeTool = tool; },
    setZoom(z: number) { zoom = z; },
  };
}
```
</phase_2_specifications>

<phase_2_exit_criteria>
- [ ] Drag-and-drop upload working
- [ ] Paste from clipboard working
- [ ] Image displays correctly with proper aspect ratio
- [ ] Rotation controls work (all 4 orientations)
- [ ] Rotation is undoable (Ctrl+Z reverses it)
- [ ] Upload integrates with backend (can use mock endpoint for now)
- [ ] Error states handled gracefully
- [ ] Display context updates correctly on resize/zoom/rotation
</phase_2_exit_criteria>

**⏸️ STOP: Await human review before proceeding to Phase 3**

---

## Phase 3: Region Management

<phase_3_objective>
Implement region creation, selection, manipulation (drag/resize), and deletion - all with full undo/redo support.
</phase_3_objective>

<phase_3_deliverables>
1. **Region drawing tool** - Click and drag to create rectangles
2. **Selection** - Click to select, click away to deselect
3. **Drag to move** - Selected regions can be repositioned
4. **Resize handles** - 8-point resize (corners + edges)
5. **Delete** - Keyboard (Delete/Backspace) and UI button
6. **Full undo/redo** for all region operations
</phase_3_deliverables>

<phase_3_specifications>
### Drawing Tool Behavior
1. User clicks "Draw Region" button (or presses 'R' key)
2. Cursor changes to crosshair
3. Click and drag on canvas creates a rectangle
4. Minimum size enforced: 20x20 display pixels (convert to image coords)
5. On mouse up: 
   - Convert display coordinates to image coordinates
   - Execute CreateRegionCommand
   - Switch to select mode
   - Select new region
6. Auto-generate region ID: `m_${nanoid(10)}`

### Command Implementations
```typescript
class CreateRegionCommand implements EditorCommand {
  type = 'create-region';
  get description() { return `Create region "${this.region.label || 'untitled'}"`; }
  
  constructor(private state: EditorState, private region: Region) {}
  
  execute() {
    this.state._addRegion(this.region);
    this.state.setSelectedRegionId(this.region.id);
  }
  
  undo() {
    this.state._removeRegion(this.region.id);
    this.state.setSelectedRegionId(null);
  }
}

class MoveRegionCommand implements EditorCommand {
  type = 'move-region';
  get description() { return 'Move region'; }
  
  constructor(
    private state: EditorState,
    private regionId: string,
    private fromPosition: { x: number; y: number },
    private toPosition: { x: number; y: number }
  ) {}
  
  execute() {
    this.state._updateRegion(this.regionId, {
      shape: { 
        ...this.state.regions.find(r => r.id === this.regionId)!.shape,
        x: this.toPosition.x,
        y: this.toPosition.y
      }
    });
  }
  
  undo() {
    this.state._updateRegion(this.regionId, {
      shape: {
        ...this.state.regions.find(r => r.id === this.regionId)!.shape,
        x: this.fromPosition.x,
        y: this.fromPosition.y
      }
    });
  }
  
  // Merge consecutive moves into single undo step
  canMergeWith(other: EditorCommand): boolean {
    return other instanceof MoveRegionCommand && 
           other.regionId === this.regionId;
  }
  
  mergeWith(other: MoveRegionCommand): MoveRegionCommand {
    return new MoveRegionCommand(
      this.state,
      this.regionId,
      this.fromPosition,  // Keep original start position
      other.toPosition    // Use final end position
    );
  }
}

class ResizeRegionCommand implements EditorCommand { /* similar pattern */ }
class DeleteRegionCommand implements EditorCommand { /* similar pattern */ }
```

### Selection Behavior
- Click on region to select
- Click on canvas (not on region) to deselect
- Selected region shows resize handles and marching ants
- Only one region can be selected at a time
- Selection syncs between canvas and label panel
- Selection is NOT undoable (UI state only)

### Drag Behavior
- Only selected regions can be dragged
- Drag starts on mousedown inside region (not on resize handle)
- Constrain to image bounds (in image coordinate space)
- Visual feedback during drag (region follows cursor)
- Command executed on mouseup with start and end positions
- Continuous dragging merges into single undo step

### Resize Behavior
```
Resize handles layout:
    [nw]────[n]────[ne]
      │             │
    [w]           [e]
      │             │
    [sw]────[s]────[se]
```
- Handles are 8x8 display pixels, positioned at corners and edge midpoints
- Dragging a corner resizes freely
- Dragging an edge resizes in one dimension only
- Minimum size enforced during resize (in image coordinates)
- Command executed on mouseup

### Keyboard Shortcuts
| Key | Action |
|-----|--------|
| R | Activate draw tool |
| V / Escape | Activate select tool |
| Delete / Backspace | Delete selected region |
| Arrow keys | Nudge selected region by 1 image pixel |
| Shift + Arrow | Nudge by 10 image pixels |
| Ctrl+Z | Undo |
| Ctrl+Y / Ctrl+Shift+Z | Redo |

### Region Rendering (SVG)
All coordinates must be transformed from image space to display space:
```svelte
<svg class="absolute inset-0 pointer-events-none">
  {#each regions as region (region.id)}
    {@const displayShape = regionToDisplay(region.shape, displayContext)}
    <g 
      class="pointer-events-auto cursor-move"
      class:selected={region.id === selectedId}
    >
      <rect
        x={displayShape.x}
        y={displayShape.y}
        width={displayShape.width}
        height={displayShape.height}
        class="region-rect"
        class:region-selected={region.id === selectedId}
      />
      {#if region.id === selectedId}
        <!-- Resize handles at display coordinates -->
        {#each getResizeHandles(displayShape) as handle}
          <rect 
            class="resize-handle" 
            data-handle={handle.position}
            x={handle.x} y={handle.y}
            width={8} height={8}
          />
        {/each}
      {/if}
    </g>
  {/each}
</svg>
```
</phase_3_specifications>

<phase_3_exit_criteria>
- [ ] Can draw new rectangular regions
- [ ] Can select/deselect regions
- [ ] Can drag regions to new positions
- [ ] Can resize regions from any handle
- [ ] Regions constrained within image bounds
- [ ] Delete works via keyboard and UI
- [ ] Coordinates correctly stored in image space
- [ ] All operations undoable/redoable
- [ ] Continuous drag/resize merges into single undo step
- [ ] Performance acceptable with 20+ regions
</phase_3_exit_criteria>

**⏸️ STOP: Await human review before proceeding to Phase 4**

---

## Phase 4: Label Panel & Metadata Editing

<phase_4_objective>
Implement the label panel with region metadata editing, filtering, and deletion confirmation - with undo/redo for all edits.
</phase_4_objective>

<phase_4_deliverables>
1. **Region list** with all metadata fields displayed
2. **Inline editing** for label, hint, and back_content
3. **Filter/search** across all text fields
4. **Delete with confirmation** dialog
5. **Selection sync** between canvas and list
6. **Empty state** when no regions exist
7. **Undo/redo** for metadata changes
</phase_4_deliverables>

<phase_4_specifications>
### LabelPanelItem Layout
```
┌─────────────────────────────────────────┐
│ Region 1                        [🗑][⋮] │
│ ┌─────────────────────────────────────┐ │
│ │ Label*: [Mitochondria___________]   │ │
│ │ Hint:   [Powerhouse of the cell_]   │ │
│ │ Extra:  [Produces ATP through...]   │ │
│ └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

### Metadata Edit Commands
```typescript
class UpdateRegionMetadataCommand implements EditorCommand {
  type = 'update-metadata';
  get description() { return `Update ${this.field} for "${this.regionLabel}"`; }
  
  constructor(
    private state: EditorState,
    private regionId: string,
    private regionLabel: string,
    private field: 'label' | 'hint' | 'back_content',
    private previousValue: string,
    private newValue: string
  ) {}
  
  execute() {
    this.state._updateRegion(this.regionId, { [this.field]: this.newValue });
  }
  
  undo() {
    this.state._updateRegion(this.regionId, { [this.field]: this.previousValue });
  }
  
  // Merge rapid typing into single undo step
  canMergeWith(other: EditorCommand): boolean {
    return other instanceof UpdateRegionMetadataCommand &&
           other.regionId === this.regionId &&
           other.field === this.field;
  }
  
  mergeWith(other: UpdateRegionMetadataCommand): UpdateRegionMetadataCommand {
    return new UpdateRegionMetadataCommand(
      this.state,
      this.regionId,
      this.regionLabel,
      this.field,
      this.previousValue,  // Keep original value
      other.newValue       // Use final typed value
    );
  }
}
```

### Field Editing Behavior
- Click on field value to edit inline
- Use shadcn-svelte Input component
- Debounce changes (e.g., 500ms) before creating command
- On blur: finalize command if value changed
- On Enter: finalize and move to next field
- Label field is required - show validation error if empty on blur
- Hint and Extra are optional (can be empty)

### Filtering
- Single search input at top of panel
- Filters across label, hint, and back_content
- Case-insensitive partial match
- Show "N of M regions" count
- Clear button when filter active
- Filter is UI state only (not undoable)

### Delete Confirmation
Use shadcn-svelte AlertDialog:
```
┌────────────────────────────────────┐
│ Delete Region?                     │
│                                    │
│ This will permanently delete the   │
│ region "Mitochondria". This action │
│ can be undone with Ctrl+Z.         │
│                                    │
│         [Cancel]  [Delete]         │
└────────────────────────────────────┘
```

### Selection Sync
- Clicking a region in the canvas selects it in the panel (scroll into view)
- Clicking a region in the panel selects it on the canvas (visual highlight)
- Selected item in panel has different background color
</phase_4_specifications>

<phase_4_exit_criteria>
- [ ] All regions displayed in scrollable list
- [ ] Label, hint, extra fields editable inline
- [ ] Validation prevents empty labels
- [ ] Metadata changes are undoable
- [ ] Typing merges into single undo step
- [ ] Filter works across all text fields
- [ ] Delete shows confirmation dialog
- [ ] Selection syncs bidirectionally
- [ ] Keyboard navigation works in panel
</phase_4_exit_criteria>

**⏸️ STOP: Await human review before proceeding to Phase 5**

---

## Phase 5: Annotation Layer

<phase_5_objective>
Add an annotation layer for drawing arrows on the image, with full undo/redo support.
</phase_5_objective>

<phase_5_deliverables>
1. **Annotation tool mode** in toolbar
2. **Arrow drawing** - Click start point, drag to end point
3. **Arrow styling** - Color, thickness options
4. **Annotation list/management** - Delete annotations
5. **Undo/redo** for annotation operations
</phase_5_deliverables>

<phase_5_specifications>
### Arrow Tool Behavior
1. Select "Annotate → Arrow" from toolbar dropdown
2. Click on image to set start point (show small dot preview)
3. Drag to end point (show arrow preview with current style)
4. Release to create arrow (execute CreateAnnotationCommand)
5. Arrows are not selectable/movable (delete and redraw if needed)

### Annotation Commands
```typescript
class CreateAnnotationCommand implements EditorCommand {
  type = 'create-annotation';
  description = 'Add arrow annotation';
  
  constructor(private state: EditorState, private annotation: Annotation) {}
  
  execute() {
    this.state._addAnnotation(this.annotation);
  }
  
  undo() {
    this.state._removeAnnotation(this.annotation.id);
  }
}

class DeleteAnnotationCommand implements EditorCommand { /* similar */ }
class ClearAllAnnotationsCommand implements EditorCommand { /* batch delete */ }
```

### Arrow Rendering (SVG)
Coordinates transform from image space to display space:
```svelte
<defs>
  <marker id="arrowhead" markerWidth="10" markerHeight="7" 
          refX="9" refY="3.5" orient="auto">
    <polygon points="0 0, 10 3.5, 0 7" fill="currentColor" />
  </marker>
</defs>

{#each annotations as annotation (annotation.id)}
  {@const start = imageToDisplay(annotation.points.start, displayContext)}
  {@const end = imageToDisplay(annotation.points.end, displayContext)}
  <line 
    x1={start.x} y1={start.y}
    x2={end.x} y2={end.y}
    stroke={annotation.style.color}
    stroke-width={annotation.style.thickness * 2}
    marker-end="url(#arrowhead)"
    class="pointer-events-auto cursor-pointer"
    onclick={() => selectAnnotationForDeletion(annotation.id)}
  />
{/each}
```

### Toolbar Dropdown
```
[Annotate ▼]
├─ Arrow
├─ ──────── (separator)
├─ Clear All Annotations
```

### Style Options (Simple)
- Color picker: 3-4 preset colors (use design tokens)
- Thickness: 3 options (thin, medium, thick)
- Shown in toolbar when annotation tool active

### Future Stubs (Do Not Implement)
Add TypeScript interfaces for future annotation types:
```typescript
// types.ts - interfaces only, no implementation
interface TextAnnotation {
  id: string;
  type: 'text';
  position: { x: number; y: number };
  text: string;
  style: { color: string; fontSize: number };
}

interface HighlightAnnotation {
  id: string;
  type: 'highlight';
  shape: RectShape;
  style: { color: string; opacity: number };
}

interface FreehandAnnotation {
  id: string;
  type: 'freehand';
  points: Array<{ x: number; y: number }>;
  style: { color: string; thickness: number };
}
```
</phase_5_specifications>

<phase_5_exit_criteria>
- [ ] Arrow tool selectable from toolbar
- [ ] Can draw arrows on the image
- [ ] Arrows render with proper arrowhead
- [ ] Arrow creation is undoable
- [ ] Arrows can be deleted (undoable)
- [ ] Clear all annotations works (undoable)
- [ ] Annotations stored in image coordinates
- [ ] Annotations render correctly with image rotation/zoom
</phase_5_exit_criteria>

**⏸️ STOP: Await human review before proceeding to Phase 6**

---

## Phase 6: Integration, Polish & AI Stubs

<phase_6_objective>
Complete the editor with save/validation flow, zoom controls, AI processing stubs, and final polish.
</phase_6_objective>

<phase_6_deliverables>
1. **Save/Cancel flow** with proper data serialization and validation
2. **Zoom controls** for the canvas
3. **AI processing stubs** (hooks for future OCR/LLM integration)
4. **Final polish** - loading states, error handling, edge cases
5. **Documentation** - Component usage examples
</phase_6_deliverables>

<phase_6_specifications>
### Save Flow
```typescript
function handleSave() {
  // Validate: at least one region with a label
  if (state.regions.length === 0) {
    showError("Add at least one region");
    return;
  }
  
  const emptyLabels = state.regions.filter(r => !r.label.trim());
  if (emptyLabels.length > 0) {
    showError("All regions must have labels");
    highlightInvalidRegions(emptyLabels);
    return;
  }
  
  // Serialize to schema format (coordinates already in image space)
  const value: ImageOcclusionField = {
    name: "image",
    type: "image_occlusion",
    image: state.image,
    settings: { mode: state.mode },
    regions: state.regions,
    annotations: state.annotations.length > 0 ? state.annotations : undefined
  };
  
  onSave(value);
}
```

### Cancel Flow with Dirty Check
```typescript
function handleCancel() {
  if (history.undoCount > 0) {
    // Show confirmation dialog
    showConfirmation({
      title: "Discard changes?",
      message: "You have unsaved changes. Are you sure you want to cancel?",
      confirmLabel: "Discard",
      onConfirm: () => onCancel()
    });
  } else {
    onCancel();
  }
}
```

### Zoom Controls
- Zoom range: 50% - 300%
- Controls: +/- buttons in toolbar
- Scroll wheel with Ctrl/Cmd held
- "Fit to view" button (calculates optimal zoom)
- "100%" button (reset to actual size)
- Zoom level indicator in status bar

### Pan When Zoomed
- When zoomed > 100%, image may exceed container
- Space + drag to pan (standard image editor pattern)
- Pan limits constrained to keep image partially visible
- Pan offset is UI state (not undoable, not saved)

### AI Processing Stubs
Create placeholder infrastructure for future AI integration:

```typescript
// types.ts
interface AIAnalysisState {
  status: 'idle' | 'processing' | 'complete' | 'error';
  progress?: number;  // 0-100
  error?: string;
  results?: AIAnalysisResults;
}

interface AIAnalysisResults {
  imageDescription?: string;
  suggestedRegions?: Array<{
    shape: RectShape;          // In image coordinates
    suggestedLabel: string;
    confidence: number;        // 0-1
  }>;
  detectedText?: Array<{
    text: string;
    bounds: RectShape;         // In image coordinates
    confidence: number;
  }>;
}

// Stub function
async function requestAIAnalysis(assetId: string): Promise<void> {
  // TODO: Implement when backend endpoint ready
  // POST /api/assets/{assetId}/analyze
  // Returns analysis job ID, poll for results via:
  // GET /api/assets/{assetId}/analysis/{jobId}
  console.log('AI analysis not yet implemented for asset:', assetId);
}

// Hook for applying AI suggestions
function applyAISuggestedRegion(suggestion: AIAnalysisResults['suggestedRegions'][0]) {
  const region: Region = {
    id: `m_${nanoid(10)}`,
    shape: suggestion.shape,
    label: suggestion.suggestedLabel,
    hint: '',
    back_content: ''
  };
  history.execute(new CreateRegionCommand(state, region));
}
```

### AI UI Placeholder
Add disabled controls in toolbar:
```svelte
<DropdownMenu>
  <DropdownMenuTrigger asChild>
    <Button variant="outline" disabled>
      <Sparkles class="w-4 h-4 mr-2" />
      AI Assist
      <Badge variant="secondary" class="ml-2">Coming Soon</Badge>
    </Button>
  </DropdownMenuTrigger>
  <DropdownMenuContent>
    <DropdownMenuItem disabled>
      Auto-detect regions
    </DropdownMenuItem>
    <DropdownMenuItem disabled>
      Suggest labels from OCR
    </DropdownMenuItem>
    <DropdownMenuItem disabled>
      Describe image
    </DropdownMenuItem>
  </DropdownMenuContent>
</DropdownMenu>
```

### Status Bar
Bottom of editor shows contextual information:
- Region count: "5 regions"
- Mode indicator: "Hide All, Guess One"
- Zoom level: "100%"
- Undo/Redo state: "↶ 5 ↷ 2" or similar indicator
- Dirty indicator: "• Unsaved changes" when history has items

### Keyboard Shortcuts Summary
Display via help tooltip or modal:
| Shortcut | Action |
|----------|--------|
| R | Draw region tool |
| V | Select tool |
| A | Arrow annotation |
| Delete | Delete selected |
| Ctrl+Z | Undo |
| Ctrl+Y | Redo |
| Ctrl+S | Save |
| Escape | Cancel / Deselect / Exit tool |
| Ctrl++ | Zoom in |
| Ctrl+- | Zoom out |
| Ctrl+0 | Fit to view |
| Ctrl+1 | Actual size (100%) |
| Space+Drag | Pan (when zoomed) |
| Arrow keys | Nudge selected region |
</phase_6_specifications>

<phase_6_exit_criteria>
- [ ] Save produces valid schema JSON with image coordinates
- [ ] Validation prevents invalid saves
- [ ] Cancel shows confirmation if changes exist
- [ ] Zoom controls work smoothly
- [ ] Pan works when zoomed in
- [ ] AI stubs in place with disabled UI
- [ ] All keyboard shortcuts working
- [ ] Status bar shows accurate state
- [ ] No console errors in any flow
- [ ] Works in Chrome, Firefox, Safari
- [ ] Component documented with usage example
</phase_6_exit_criteria>

---

## Implementation Notes

<global_requirements>
### Code Quality
- Use TypeScript strict mode
- No `any` types - define proper interfaces
- Extract reusable logic into utility functions
- Keep components focused (< 200 lines preferred)
- All state mutations through command pattern

### Svelte 5 Patterns
```svelte
<script lang="ts">
  import type { ImageOcclusionEditorProps } from './types';
  
  // Props with $props rune
  let { initialValue, onSave, onCancel, uploadEndpoint }: ImageOcclusionEditorProps = $props();
  
  // Create state and history
  const state = createEditorState(initialValue);
  const history = createEditorHistory();
  
  // Derived values
  const canSave = $derived(
    state.image !== null && 
    state.regions.length > 0 &&
    state.regions.every(r => r.label.trim())
  );
  
  // Side effects
  $effect(() => {
    if (state.selectedRegionId) {
      scrollRegionIntoView(state.selectedRegionId);
    }
  });
  
  // Command-based mutations
  function handleRegionMove(id: string, newPosition: { x: number; y: number }) {
    const region = state.regions.find(r => r.id === id);
    if (!region) return;
    
    history.execute(new MoveRegionCommand(
      state,
      id,
      { x: region.shape.x, y: region.shape.y },
      newPosition
    ));
  }
</script>
```

### Performance Considerations
- Use SVG for region rendering (vector precision at any zoom)
- Debounce continuous operations, execute command on completion
- Use requestAnimationFrame for drag/resize visual updates
- Virtualize region list if > 50 items (unlikely but good practice)

### Testing Considerations
- Components should accept mock data via props
- Commands are pure and testable in isolation
- Coordinate transforms should have unit tests
- Mock history for testing component behavior
</global_requirements>

<shadcn_components_to_use>
- `Button` - All toolbar and action buttons
- `Input` - Label/hint/extra fields, filter input
- `ScrollArea` - Label panel scrolling
- `Tooltip` - Toolbar button hints, keyboard shortcut hints
- `AlertDialog` - Delete confirmation, unsaved changes warning
- `DropdownMenu` - Annotate tool menu, mode selector
- `Separator` - Visual dividers
- `Badge` - "Coming soon" indicators
- `Popover` - Color/style pickers
</shadcn_components_to_use>

<do_not_implement>
- Image occlusion review/study renderer (separate component)
- Actual AI/OCR processing (stub infrastructure only)
- Web image search
- Non-rectangle shapes (TypeScript interfaces only)
- Multi-user collaboration
- Version history beyond undo/redo
- Mobile/tablet layouts
</do_not_implement>