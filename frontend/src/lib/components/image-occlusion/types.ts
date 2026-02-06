/**
 * Image Occlusion Editor Type Definitions
 *
 * These types define the data structures for the image occlusion editor,
 * supporting rectangular mask regions over images for spaced repetition cards.
 */

// ============================================================================
// Core Data Types
// ============================================================================

/** 2D point in image or display coordinates */
export interface Point {
	x: number;
	y: number;
}

/** Rectangular shape definition */
export interface RectShape {
	x: number;
	y: number;
	width: number;
	height: number;
}

/** A mask region that occludes part of the image */
export interface Region {
	id: string;
	shape: RectShape;
	label: string;
	hint?: string;
	backExtra?: string;
}

/** Arrow annotation for pointing at regions */
export interface ArrowAnnotation {
	type: 'arrow';
	id: string;
	start: Point;
	end: Point;
	color?: string;
}

/** Union of all annotation types */
export type Annotation = ArrowAnnotation;

/** Image data including URL and dimensions */
export interface ImageData {
	url: string;
	assetId?: string;
	width: number;
	height: number;
	rotation: 0 | 90 | 180 | 270;
}

/** How regions are revealed during review */
export type OcclusionMode = 'hide_all_guess_one' | 'hide_one_guess_one';

/** How the target region is displayed after reveal */
export type RevealStyle = 'show_label' | 'image_only';

/** The complete field structure stored in fact content */
export interface ImageOcclusionField {
	version: 1;
	title: string;
	image: ImageData;
	regions: Region[];
	annotations: Annotation[];
	mode: OcclusionMode;
	revealStyle?: RevealStyle;
}

// ============================================================================
// Editor State Types
// ============================================================================

/** Available editor tools */
export type EditorTool = 'select' | 'draw_region' | 'draw_arrow';

/** Resize handle positions on a selected region */
export type ResizeHandlePosition =
	| 'nw' | 'n' | 'ne'
	| 'w'  |      'e'
	| 'sw' | 's' | 'se';

/** A resize handle with position and bounds */
export interface ResizeHandle {
	position: ResizeHandlePosition;
	x: number;
	y: number;
	size: number;
}

/** Context for coordinate transformations */
export interface DisplayContext {
	/** Natural image width */
	imageWidth: number;
	/** Natural image height */
	imageHeight: number;
	/** Current rotation */
	rotation: 0 | 90 | 180 | 270;
	/** Current zoom level (1.0 = 100%) */
	zoom: number;
	/** Pan offset in display coordinates */
	panOffset: Point;
	/** Container dimensions */
	containerWidth: number;
	containerHeight: number;
}

// ============================================================================
// Command Pattern Types
// ============================================================================

/** Unique identifier for command types (for merging) */
export type CommandType =
	| 'create_region'
	| 'delete_region'
	| 'move_region'
	| 'resize_region'
	| 'update_region_metadata'
	| 'rotate_image';

/** Interface for undoable/redoable commands */
export interface EditorCommand {
	/** Command type identifier */
	readonly type: CommandType;
	/** Execute the command */
	execute(): void;
	/** Undo the command */
	undo(): void;
	/** Optional: merge with another command of same type (for continuous drag) */
	merge?(other: EditorCommand): boolean;
	/** Unique ID for tracking (e.g., region ID for moves) */
	readonly targetId?: string;
}

// ============================================================================
// Component Props Types
// ============================================================================

/** Props for ImageOcclusionEditor (main container) */
export interface ImageOcclusionEditorProps {
	/** Initial field value (for editing existing facts) */
	initialValue?: ImageOcclusionField;
	/** Called when user clicks Done with valid data */
	onSave?: (field: ImageOcclusionField) => void;
	/** Called when user cancels */
	onCancel?: () => void;
}

/** Props for EditorCanvas */
export interface EditorCanvasProps {
	/** Image to display */
	image: ImageData | null;
	/** Regions to render */
	regions: Region[];
	/** Currently selected region ID */
	selectedRegionId: string | null;
	/** Display context for coordinate transforms */
	displayContext: DisplayContext;
	/** Active tool */
	activeTool: EditorTool;
	/** Callback when region is selected */
	onSelectRegion?: (id: string | null) => void;
	/** Callback when container size changes */
	onContainerResize?: (width: number, height: number) => void;
}

/** Props for RegionOverlay */
export interface RegionOverlayProps {
	/** The region to render */
	region: Region;
	/** Whether this region is selected */
	isSelected: boolean;
	/** Display context for coordinate transforms */
	displayContext: DisplayContext;
	/** Callback when clicked */
	onClick?: () => void;
}

/** Props for LabelPanel */
export interface LabelPanelProps {
	/** All regions */
	regions: Region[];
	/** Currently selected region ID */
	selectedRegionId: string | null;
	/** Callback when a region is selected */
	onSelectRegion?: (id: string) => void;
	/** Callback when region metadata is updated */
	onUpdateRegion?: (id: string, updates: Partial<Pick<Region, 'label' | 'hint' | 'backExtra'>>) => void;
	/** Callback when region is deleted */
	onDeleteRegion?: (id: string) => void;
}

/** Props for LabelPanelItem */
export interface LabelPanelItemProps {
	/** The region to display */
	region: Region;
	/** Display index (1-based) */
	index: number;
	/** Whether this region is selected */
	isSelected: boolean;
	/** Callback when clicked */
	onSelect?: () => void;
	/** Callback when label changes */
	onLabelChange?: (value: string) => void;
	/** Callback when hint changes */
	onHintChange?: (value: string) => void;
	/** Callback when back content changes */
	onBackExtraChange?: (value: string) => void;
	/** Callback when delete is clicked */
	onDelete?: () => void;
}

/** Props for EditorToolbar */
export interface EditorToolbarProps {
	/** Current active tool */
	activeTool: EditorTool;
	/** Whether undo is available */
	canUndo: boolean;
	/** Whether redo is available */
	canRedo: boolean;
	/** Current zoom level */
	zoom: number;
	/** Callback when tool changes */
	onToolChange?: (tool: EditorTool) => void;
	/** Callback for undo */
	onUndo?: () => void;
	/** Callback for redo */
	onRedo?: () => void;
	/** Callback for rotate */
	onRotate?: () => void;
	/** Callback for zoom change */
	onZoomChange?: (zoom: number) => void;
	/** Callback for zoom to fit */
	onZoomFit?: () => void;
	/** Callback for done */
	onDone?: () => void;
}

/** Props for StatusBar */
export interface StatusBarProps {
	/** Number of regions */
	regionCount: number;
	/** Current occlusion mode */
	mode: OcclusionMode;
	/** Current zoom percentage */
	zoom: number;
}
