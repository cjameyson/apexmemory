/**
 * Image Occlusion Editor - Barrel Exports
 */

// Main component
export { default as ImageOcclusionEditor } from './ImageOcclusionEditor.svelte';

// Sub-components (for advanced customization)
export { default as EditorCanvas } from './EditorCanvas.svelte';
export { default as EditorToolbar } from './EditorToolbar.svelte';
export { default as LabelPanel } from './LabelPanel.svelte';
export { default as LabelPanelItem } from './LabelPanelItem.svelte';
export { default as RegionOverlay } from './RegionOverlay.svelte';
export { default as StatusBar } from './StatusBar.svelte';
export { default as ImageUploader } from './ImageUploader.svelte';

// State management
export { createEditorState, type EditorState } from './editor-state.svelte';
export { createHistoryManager, type HistoryManager } from './history.svelte';

// Commands
export {
	CreateRegionCommand,
	DeleteRegionCommand,
	MoveRegionCommand,
	ResizeRegionCommand,
	UpdateRegionMetadataCommand,
	RotateImageCommand,
	type EditorStateMutators
} from './commands';

// Utilities
export {
	generateRegionId,
	generateAnnotationId,
	debounce,
	throttle,
	clamp,
	formatZoom,
	formatOcclusionMode
} from './utils';

// Coordinate utilities
export {
	imageToDisplay,
	displayToImage,
	regionToDisplay,
	displayToRegion,
	calculateScale,
	calculateCenteredOffset,
	getResizeHandles,
	constrainToImageBounds,
	constrainShapeToImageBounds,
	isPointInImage,
	isPointInShape,
	getResizeCursor
} from './coordinates';

// Types
export type {
	Point,
	RectShape,
	Region,
	ArrowAnnotation,
	Annotation,
	ImageData,
	ImageOcclusionField,
	OcclusionMode,
	EditorTool,
	DisplayContext,
	ResizeHandle,
	ResizeHandlePosition,
	EditorCommand,
	CommandType,
	ImageOcclusionEditorProps,
	EditorCanvasProps,
	RegionOverlayProps,
	LabelPanelProps,
	LabelPanelItemProps,
	EditorToolbarProps,
	StatusBarProps
} from './types';
