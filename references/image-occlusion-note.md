# Image Occlusion Field Schema

This document defines the JSON schema for `image_occlusion` type notes.

## Overview

Image occlusion notes allow users to create flashcards from a single image by defining mask regions. During review, all masks are shown and the user must identify the content under the target mask.

**MVP Scope:**
- Single image per note
- Rectangle shapes only
- No region grouping
- One review mode: hide-all-reveal-one

## Note Content Structure

```json
{
  "version": 1,
  "fields": [
    {
      "name": "image",
      "type": "image_occlusion",
      "image": {
        "url": "https://storage.example.com/user_abc/img_xyz.png",
        "width": 1200,
        "height": 800,
        "original_name": "cell_diagram.png"
      },
      "settings": {
        "mode": "hide_all_reveal_one"
      },
      "regions": [
        {
          "id": "m_k7Xp2mQ9",
          "shape": {
            "type": "rect",
            "x": 150,
            "y": 200,
            "width": 80,
            "height": 40
          },
          "label": "Mitochondria",
          "hint": "Powerhouse of the cell",
          "back_content": "Produces ATP through cellular respiration"
        },
        {
          "id": "m_nR4wL8vK",
          "shape": {
            "type": "rect",
            "x": 330,
            "y": 150,
            "width": 60,
            "height": 90
          },
          "label": "Nucleus",
          "hint": null,
          "back_content": "Contains genetic material (DNA)"
        }
      ]
    }
  ]
}
```

## Schema Reference

### Field Object

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | string | ✓ | Field identifier, typically `"image"` |
| `type` | string | ✓ | Must be `"image_occlusion"` |
| `image` | object | ✓ | Image metadata |
| `settings` | object | | Review settings (defaults applied if omitted) |
| `regions` | array | ✓ | Mask region definitions (1-128 regions) |

### Image Object

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `url` | string | ✓ | Stored image URL |
| `width` | integer | ✓ | Original image width in pixels |
| `height` | integer | ✓ | Original image height in pixels |
| `original_name` | string | | Original filename for display |

The `width` and `height` define the coordinate space for all regions. These are captured when the image is uploaded and must not change.

### Settings Object

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| `mode` | string | `"hide_all_reveal_one"` | Review mode |

**Mode Options:**

| Mode | Description | MVP |
|------|-------------|-----|
| `hide_all_reveal_one` | All regions masked, user identifies target | ✓ |
| `hide_one_reveal_all` | Only target masked, context visible | Future |

### Region Object

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `id` | string | ✓ | Stable identifier: `m_` + nanoid (8-16 chars) |
| `shape` | object | ✓ | Shape definition |
| `label` | string | | Text displayed on the mask (plain text) |
| `hint` | string | | Hint shown during review (plain text) |
| `back_content` | string | | Answer text shown after reveal (plain text) |

**Text field guidance:**
- `label`: Short identifier shown on the mask itself (e.g., "?", "A", "1", or a term)
- `hint`: Clue to help recall, shown near the image before reveal (e.g., "Powerhouse of the cell")
- `back_content`: Full answer or explanation shown after reveal

### Shape Object — Rectangle

```json
{
  "type": "rect",
  "x": 150,
  "y": 200,
  "width": 80,
  "height": 40
}
```

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `type` | string | ✓ | Must be `"rect"` for MVP |
| `x` | integer | ✓ | Left edge in pixels |
| `y` | integer | ✓ | Top edge in pixels |
| `width` | integer | ✓ | Width in pixels |
| `height` | integer | ✓ | Height in pixels |

## Coordinate System

**Origin:** Top-left corner of the image is (0, 0).

**Units:** Absolute pixels relative to the original image dimensions (`image.width` × `image.height`).

**Scaling:** The client scales coordinates proportionally when rendering at different sizes.

```javascript
// Rendering at any display size
function renderRegion(region, image, displayWidth) {
  const scale = displayWidth / image.width;
  
  return {
    x: region.shape.x * scale,
    y: region.shape.y * scale,
    width: region.shape.width * scale,
    height: region.shape.height * scale,
  };
}
```

### Design Note: Why Absolute Pixels?

We considered normalized coordinates (0-1 range) but chose absolute pixels for several reasons:

| Consideration | Absolute Pixels | Normalized 0-1 |
|---------------|-----------------|----------------|
| Editor UX | What user sees = what's stored | Requires conversion |
| Debugging | Coordinates match image viewers | Abstract values |
| Precision | Exact integers | Potential float rounding |
| Industry precedent | Anki, Figma, Photoshop, SVG | Some maps/games |

The render-time math is equivalent either way — we're choosing where the division happens. Storing pixels keeps the editor simple and matches user expectations when they inspect the data or debug issues.

**Anki precedent:** Anki's Image Occlusion stores SVG with pixel coordinates. The `viewBox` attribute handles scaling:

```svg
<svg viewBox="0 0 1200 800">
  <rect x="150" y="200" width="80" height="40" />
</svg>
```

## Card Generation

Each region generates one card with `element_id` set to the region's `id`.

```
Input regions:
  [{ id: "m_k7Xp2mQ9", ... }, { id: "m_nR4wL8vK", ... }]

Generated cards:
  [{ element_id: "m_k7Xp2mQ9" }, { element_id: "m_nR4wL8vK" }]
```

The `element_id` is stable — editing the note to add/remove other regions does not affect existing cards or their review history.

## Review Behavior

### Question Side (hide_all_reveal_one)

1. Display image scaled to fit viewport
2. Render all region masks as colored rectangles
3. Highlight target mask (e.g., pulsing border, different color)
4. Display target's `label` centered on the mask (if set)
5. Display target's `hint` below the image (if set)

### Answer Side

1. Remove target mask (reveal that region)
2. Keep other masks in place
3. Display target's `back_content` below the image (if set)

### Design Note: Why "Hide All, Reveal One"?

This mode is the most effective for learning identification tasks:

- **Active recall:** User must retrieve the answer from memory, not just recognize it
- **Contextual learning:** Surrounding masked regions provide structural context
- **Prevents shortcuts:** User can't identify by elimination if all regions look the same

The alternative mode ("hide one, reveal all") is useful for different scenarios but is deferred to post-MVP.

## Validation Rules

| Rule | Constraint |
|------|------------|
| Regions count | 1 ≤ count ≤ 128 |
| Region ID format | `^m_[a-zA-Z0-9_-]{6,24}$` |
| Region ID uniqueness | All IDs must be unique within note |
| Shape type | Must be `"rect"` (MVP) |
| Shape coordinates | Non-negative integers |
| Shape dimensions | Positive integers (width > 0, height > 0) |
| Image URL | Non-empty string |
| Image dimensions | Positive integers |

### Design Note: 128 Region Limit

The 128 limit balances flexibility with sanity:

- Most use cases need 2-20 regions
- Complex diagrams (anatomy, circuits) rarely exceed 50
- 128 accommodates edge cases without enabling abuse
- Consistent with our max-cards-per-note constraint
- Well under database and JSON size concerns

## Examples

### Anatomy Diagram — Multiple Organs

```json
{
  "version": 1,
  "fields": [
    {
      "name": "image",
      "type": "image_occlusion",
      "image": {
        "url": "https://cdn.example.com/uploads/heart_diagram.png",
        "width": 800,
        "height": 600,
        "original_name": "heart_anatomy.png"
      },
      "regions": [
        {
          "id": "m_aLv8kQ2x",
          "shape": { "type": "rect", "x": 240, "y": 120, "width": 120, "height": 60 },
          "label": "?",
          "back_content": "Right Atrium — receives deoxygenated blood from the body"
        },
        {
          "id": "m_bNw9rT3y",
          "shape": { "type": "rect", "x": 440, "y": 120, "width": 120, "height": 60 },
          "label": "?",
          "back_content": "Left Atrium — receives oxygenated blood from the lungs"
        },
        {
          "id": "m_cPx0sU4z",
          "shape": { "type": "rect", "x": 240, "y": 240, "width": 144, "height": 120 },
          "label": "?",
          "back_content": "Right Ventricle — pumps blood to the lungs"
        },
        {
          "id": "m_dQy1tV5a",
          "shape": { "type": "rect", "x": 416, "y": 240, "width": 144, "height": 120 },
          "label": "?",
          "back_content": "Left Ventricle — pumps blood to the body (thickest wall)"
        }
      ]
    }
  ]
}
```

### Code Screenshot — Identify Syntax

Useful for learning programming languages, keyboard shortcuts, or UI elements.

```json
{
  "version": 1,
  "fields": [
    {
      "name": "image",
      "type": "image_occlusion",
      "image": {
        "url": "https://cdn.example.com/uploads/python_code.png",
        "width": 600,
        "height": 400
      },
      "regions": [
        {
          "id": "m_eRz2uW6b",
          "shape": { "type": "rect", "x": 30, "y": 40, "width": 72, "height": 32 },
          "label": "keyword",
          "hint": "Used to define a function",
          "back_content": "def"
        },
        {
          "id": "m_fSa3vX7c",
          "shape": { "type": "rect", "x": 420, "y": 120, "width": 150, "height": 32 },
          "label": "expression",
          "hint": "Pythonic way to build a list",
          "back_content": "List comprehension: [x**2 for x in range(10)]"
        }
      ]
    }
  ]
}
```

### Geography — Country Identification

```json
{
  "version": 1,
  "fields": [
    {
      "name": "image",
      "type": "image_occlusion",
      "image": {
        "url": "https://cdn.example.com/uploads/europe_map.png",
        "width": 1000,
        "height": 800
      },
      "regions": [
        {
          "id": "m_gTb4wY8d",
          "shape": { "type": "rect", "x": 420, "y": 280, "width": 80, "height": 100 },
          "hint": "Known for pasta and Rome",
          "back_content": "Italy"
        },
        {
          "id": "m_hUc5xZ9e",
          "shape": { "type": "rect", "x": 300, "y": 200, "width": 120, "height": 80 },
          "hint": "Known for wine and Paris",
          "back_content": "France"
        }
      ]
    }
  ]
}
```

### Minimal Example — Single Region

The simplest valid image occlusion note.

```json
{
  "version": 1,
  "fields": [
    {
      "name": "image",
      "type": "image_occlusion",
      "image": {
        "url": "https://cdn.example.com/uploads/flag.png",
        "width": 400,
        "height": 267
      },
      "regions": [
        {
          "id": "m_jWd6yA0f",
          "shape": { "type": "rect", "x": 40, "y": 27, "width": 320, "height": 213 },
          "back_content": "Brazil"
        }
      ]
    }
  ]
}
```

## Client Implementation Notes

### Region ID Generation

Use a nanoid library to generate IDs when creating new regions:

```javascript
import { nanoid } from 'nanoid';

function createRegion(shape) {
  return {
    id: `m_${nanoid(10)}`,  // e.g., "m_V1StGXR8_Z"
    shape,
    label: null,
    hint: null,
    back_content: null,
  };
}
```

### Image Upload Flow

1. User selects image file
2. Client reads dimensions from the file
3. Client uploads to storage, receives URL
4. Client initializes field with `image: { url, width, height, original_name }`
5. User draws regions using editor
6. On save, note content is validated and sent to API

### Rendering Considerations

- **Aspect ratio:** Always maintain original aspect ratio when scaling
- **Minimum size:** Ensure masks remain tappable on mobile (suggest ≥44px rendered)
- **High DPI:** Use `image.width`/`image.height` for coordinate math, not rendered size
- **Accessibility:** Consider high-contrast mask colors and clear visual focus indicators

## Future Considerations

### Additional Shapes (Post-MVP)

```typescript
type Shape =
  | { type: "rect"; x: number; y: number; width: number; height: number }
  | { type: "ellipse"; cx: number; cy: number; rx: number; ry: number }
  | { type: "polygon"; points: [number, number][] }
```

Ellipse is useful for circular structures (cells, planets). Polygon enables precise outlines for irregular shapes (countries, organs).

### Region Grouping (Post-MVP)

Groups allow multiple regions to act as a single card:
**POST MVP**: Defer implementation

```json
{
  "groups": [
    {
      "id": "g_arm_bones",
      "name": "Arm Bones",
      "region_ids": ["m_abc", "m_def", "m_ghi"],
      "back_content": "Humerus, Radius, and Ulna"
    }
  ]
}
```

Use cases: related anatomy (all bones in a limb), grouped concepts, sequential reveals.

### Custom Mask Styling (Post-MVP)

```json
{
  "settings": {
    "default_mask_style": {
      "fill": "#3b82f6",
      "fill_opacity": 0.85,
      "stroke": "#1e40af",
      "stroke_width": 2
    }
  },
  "regions": [
    {
      "id": "m_special",
      "shape": { "..." : "..." },
      "style": {
        "fill": "#ef4444"
      }
    }
  ]
}
```

Useful for color-coding regions by category or difficulty.

### Alternative Review Modes (Post-MVP)

| Mode | Behavior | Use Case |
|------|----------|----------|
| `hide_all_reveal_one` | All masked, reveal target | Identification (MVP) |
| `hide_one_reveal_all` | Only target masked | Context-heavy learning (krebs cycle)|