/**
 * Editor state management using Svelte 5 runes.
 *
 * Uses a factory function pattern to create reactive state with proper encapsulation.
 * Internal mutators (prefixed with _) are used by commands for undo/redo.
 * Public setters are used for direct UI state changes.
 */

import type {
	ImageData,
	Region,
	Annotation,
	OcclusionMode,
	EditorTool,
	DisplayContext,
	Point,
	RectShape
} from './types';
import type { EditorStateMutators } from './commands';

/** Default display context when no image is loaded */
const DEFAULT_DISPLAY_CONTEXT: DisplayContext = {
	imageWidth: 0,
	imageHeight: 0,
	rotation: 0,
	zoom: 1,
	panOffset: { x: 0, y: 0 },
	containerWidth: 0,
	containerHeight: 0
};

export interface EditorState extends EditorStateMutators {
	// ============================================================================
	// Image Data (readonly - modified through commands)
	// ============================================================================
	/** The loaded image data */
	readonly image: ImageData | null;
	/** Mask regions */
	readonly regions: Region[];
	/** Annotations (arrows, etc) */
	readonly annotations: Annotation[];
	/** Occlusion mode */
	readonly mode: OcclusionMode;

	// ============================================================================
	// UI State (read/write)
	// ============================================================================
	/** Currently selected region ID */
	readonly selectedRegionId: string | null;
	/** Active editor tool */
	readonly activeTool: EditorTool;
	/** Current zoom level (1.0 = fit to container) */
	readonly zoom: number;
	/** Pan offset from center */
	readonly panOffset: Point;
	/** Container dimensions */
	readonly containerWidth: number;
	readonly containerHeight: number;

	// ============================================================================
	// Derived State
	// ============================================================================
	/** Complete display context for coordinate transforms */
	readonly displayContext: DisplayContext;
	/** The currently selected region, or null */
	readonly selectedRegion: Region | null;
	/** Whether there's an image loaded */
	readonly hasImage: boolean;

	// ============================================================================
	// Public Setters (UI state)
	// ============================================================================
	setSelectedRegionId(id: string | null): void;
	setActiveTool(tool: EditorTool): void;
	setZoom(zoom: number): void;
	setPanOffset(offset: Point): void;
	setContainerSize(width: number, height: number): void;
	setMode(mode: OcclusionMode): void;

	// ============================================================================
	// Initialization
	// ============================================================================
	/** Initialize with an image and optional existing data */
	initialize(
		image: ImageData,
		regions?: Region[],
		annotations?: Annotation[],
		mode?: OcclusionMode
	): void;
	/** Reset to empty state */
	reset(): void;
}

/**
 * Create the editor state manager.
 */
export function createEditorState(): EditorState {
	// ============================================================================
	// Reactive State
	// ============================================================================

	// svelte-ignore state_referenced_locally
	let image = $state<ImageData | null>(null);
	// svelte-ignore state_referenced_locally
	let regions = $state<Region[]>([]);
	// svelte-ignore state_referenced_locally
	let annotations = $state<Annotation[]>([]);
	// svelte-ignore state_referenced_locally
	let mode = $state<OcclusionMode>('hide_all_guess_one');

	// UI State
	// svelte-ignore state_referenced_locally
	let selectedRegionId = $state<string | null>(null);
	// svelte-ignore state_referenced_locally
	let activeTool = $state<EditorTool>('select');
	// svelte-ignore state_referenced_locally
	let zoom = $state(1);
	// svelte-ignore state_referenced_locally
	let panOffset = $state<Point>({ x: 0, y: 0 });
	// svelte-ignore state_referenced_locally
	let containerWidth = $state(0);
	// svelte-ignore state_referenced_locally
	let containerHeight = $state(0);

	// ============================================================================
	// Derived State
	// ============================================================================

	const displayContext = $derived<DisplayContext>(
		image
			? {
					imageWidth: image.width,
					imageHeight: image.height,
					rotation: image.rotation,
					zoom,
					panOffset,
					containerWidth,
					containerHeight
				}
			: DEFAULT_DISPLAY_CONTEXT
	);

	const selectedRegion = $derived<Region | null>(
		selectedRegionId ? regions.find((r) => r.id === selectedRegionId) ?? null : null
	);

	const hasImage = $derived(image !== null);

	// ============================================================================
	// Internal Mutators (for commands)
	// ============================================================================

	function _addRegion(region: Region): void {
		regions = [...regions, region];
	}

	function _updateRegion(id: string, updates: Partial<Region>): void {
		regions = regions.map((r) => {
			if (r.id !== id) return r;

			// Special handling for shape updates - preserve unspecified dimensions
			if (updates.shape) {
				const existingShape = r.shape;
				const newShape: RectShape = {
					x: updates.shape.x ?? existingShape.x,
					y: updates.shape.y ?? existingShape.y,
					width: updates.shape.width ?? existingShape.width,
					height: updates.shape.height ?? existingShape.height
				};
				return { ...r, ...updates, shape: newShape };
			}

			return { ...r, ...updates };
		});
	}

	function _removeRegion(id: string): void {
		regions = regions.filter((r) => r.id !== id);
	}

	function _setRotation(rotation: 0 | 90 | 180 | 270): void {
		if (image) {
			image = { ...image, rotation };
		}
	}

	function _setSelectedRegionId(id: string | null): void {
		selectedRegionId = id;
	}

	// ============================================================================
	// Public Setters
	// ============================================================================

	function setSelectedRegionId(id: string | null): void {
		selectedRegionId = id;
	}

	function setActiveTool(tool: EditorTool): void {
		activeTool = tool;
	}

	function setZoom(newZoom: number): void {
		// Clamp zoom between 0.1 (10%) and 5 (500%)
		zoom = Math.max(0.1, Math.min(5, newZoom));
	}

	function setPanOffset(offset: Point): void {
		panOffset = offset;
	}

	function setContainerSize(width: number, height: number): void {
		containerWidth = width;
		containerHeight = height;
	}

	function setMode(newMode: OcclusionMode): void {
		mode = newMode;
	}

	// ============================================================================
	// Initialization
	// ============================================================================

	function initialize(
		newImage: ImageData,
		newRegions: Region[] = [],
		newAnnotations: Annotation[] = [],
		newMode: OcclusionMode = 'hide_all_guess_one'
	): void {
		image = newImage;
		regions = newRegions;
		annotations = newAnnotations;
		mode = newMode;
		selectedRegionId = null;
		activeTool = 'select';
		zoom = 1;
		panOffset = { x: 0, y: 0 };
	}

	function reset(): void {
		image = null;
		regions = [];
		annotations = [];
		mode = 'hide_all_guess_one';
		selectedRegionId = null;
		activeTool = 'select';
		zoom = 1;
		panOffset = { x: 0, y: 0 };
	}

	// ============================================================================
	// Return Public Interface
	// ============================================================================

	return {
		// Image data (getters)
		get image() {
			return image;
		},
		get regions() {
			return regions;
		},
		get annotations() {
			return annotations;
		},
		get mode() {
			return mode;
		},

		// UI state (getters)
		get selectedRegionId() {
			return selectedRegionId;
		},
		get activeTool() {
			return activeTool;
		},
		get zoom() {
			return zoom;
		},
		get panOffset() {
			return panOffset;
		},
		get containerWidth() {
			return containerWidth;
		},
		get containerHeight() {
			return containerHeight;
		},

		// Derived state
		get displayContext() {
			return displayContext;
		},
		get selectedRegion() {
			return selectedRegion;
		},
		get hasImage() {
			return hasImage;
		},

		// Internal mutators (for commands)
		_addRegion,
		_updateRegion,
		_removeRegion,
		_setRotation,
		_setSelectedRegionId,

		// Public setters
		setSelectedRegionId,
		setActiveTool,
		setZoom,
		setPanOffset,
		setContainerSize,
		setMode,

		// Initialization
		initialize,
		reset
	};
}
