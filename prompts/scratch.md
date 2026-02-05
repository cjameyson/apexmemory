
## Prompt: Image Occlusion Editor for Flashcard App

Create a professional level svelte component of an image occlusion editor for a flashcard application. This editor allows users to draw rectangular regions over an image that will be hidden during review.

Out of scope:
- A separate render component will be created to display the occluded image during review.

### Layout
- **Left panel (70%)**: Canvas area displaying the image with draggable/resizable occlusion rectangles
- **Right panel (30%)**: Sidebar with a filterable list of all occlusion regions

### Image
Use a placeholder image of a biological cell diagram (you can use a public domain cell diagram URL or a placeholder service). The image should fill the canvas area while maintaining aspect ratio.

### Canvas Behavior
- Each rectangle should have a semi-transparent fill with a colored border
- Show resize handles on corners when a region appears "selected"
- Include a toolbar above the canvas with: "Add Region" button, "Clear All" button, zoom controls
### Sidebar: Region List
- **Search/filter input** at the top to filter regions by label
- **Scrollable list** of regions each containing:
  - Editable label field (inline text input)
  - Editable hint field (smaller, muted text input below label)
  - Trash icon button to delete
- Clicking a region card should visually indicate it's "selected" (highlight the card and corresponding rectangle)
### Sample Data
Pre-populate with regions like:
1. Label: "Nucleus", Hint: "Control center of the cell"
2. Label: "Mitochondria", Hint: "Powerhouse of the cell"
3. Label: "Cell Membrane", Hint: none
4. Label: "Ribosome", Hint: "Protein synthesis"
5. Label: "Endoplasmic Reticulum", Hint: "Transport network"

All regions should use the same color, but the selected region could have a marching ants border and slightly altered shade to indicate selection.

### Styling
- Clean, modern UI with tailwindcss 4+
- Use a neutral background (light gray) for the editor area
- Mobile-friendly isn't required—this is a desktop-focused editor
- Edward Tufte level of thought and design
### Interactions to Demonstrate
- One region should appear in a "selected" state
- The filter input can have placeholder text "Filter regions..."

