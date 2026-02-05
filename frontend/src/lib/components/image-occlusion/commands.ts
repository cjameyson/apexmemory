/**
 * Command pattern implementations for undoable editor operations.
 *
 * Phase 1: Stub implementations that prepare the infrastructure.
 * Full interaction logic will be added in Phase 3.
 */

import type { EditorCommand, CommandType, Region, RectShape, Point } from './types';

/** State mutators interface that commands use to modify editor state */
export interface EditorStateMutators {
	_addRegion(region: Region): void;
	_updateRegion(id: string, updates: Partial<Region>): void;
	_removeRegion(id: string): void;
	_setRotation(rotation: 0 | 90 | 180 | 270): void;
	_setSelectedRegionId(id: string | null): void;
}

/**
 * Command to create a new region.
 */
export class CreateRegionCommand implements EditorCommand {
	readonly type: CommandType = 'create_region';
	readonly targetId: string;

	constructor(
		private readonly mutators: EditorStateMutators,
		private readonly region: Region
	) {
		this.targetId = region.id;
	}

	execute(): void {
		this.mutators._addRegion(this.region);
		this.mutators._setSelectedRegionId(this.region.id);
	}

	undo(): void {
		this.mutators._removeRegion(this.region.id);
		this.mutators._setSelectedRegionId(null);
	}
}

/**
 * Command to delete a region.
 */
export class DeleteRegionCommand implements EditorCommand {
	readonly type: CommandType = 'delete_region';
	readonly targetId: string;
	private previousSelectedId: string | null = null;

	constructor(
		private readonly mutators: EditorStateMutators,
		private readonly region: Region,
		currentSelectedId: string | null
	) {
		this.targetId = region.id;
		this.previousSelectedId = currentSelectedId;
	}

	execute(): void {
		this.mutators._removeRegion(this.region.id);
		if (this.previousSelectedId === this.region.id) {
			this.mutators._setSelectedRegionId(null);
		}
	}

	undo(): void {
		this.mutators._addRegion(this.region);
		this.mutators._setSelectedRegionId(this.previousSelectedId);
	}
}

/**
 * Command to move a region to a new position.
 * Supports merging for continuous drag operations.
 */
export class MoveRegionCommand implements EditorCommand {
	readonly type: CommandType = 'move_region';
	readonly targetId: string;
	private finalPosition: Point;

	constructor(
		private readonly mutators: EditorStateMutators,
		private readonly regionId: string,
		private readonly originalShape: RectShape,
		newPosition: Point
	) {
		this.targetId = regionId;
		this.finalPosition = newPosition;
	}

	execute(): void {
		this.mutators._updateRegion(this.regionId, {
			shape: {
				...this.originalShape,
				x: this.finalPosition.x,
				y: this.finalPosition.y
			}
		});
	}

	undo(): void {
		this.mutators._updateRegion(this.regionId, {
			shape: this.originalShape
		});
	}

	merge(other: EditorCommand): boolean {
		if (
			other instanceof MoveRegionCommand &&
			other.regionId === this.regionId
		) {
			// Update final position to the newer command's position
			this.finalPosition = other.finalPosition;
			return true;
		}
		return false;
	}
}

/**
 * Command to resize a region.
 * Supports merging for continuous drag operations.
 */
export class ResizeRegionCommand implements EditorCommand {
	readonly type: CommandType = 'resize_region';
	readonly targetId: string;
	private finalShape: RectShape;

	constructor(
		private readonly mutators: EditorStateMutators,
		private readonly regionId: string,
		private readonly originalShape: RectShape,
		newShape: RectShape
	) {
		this.targetId = regionId;
		this.finalShape = newShape;
	}

	execute(): void {
		this.mutators._updateRegion(this.regionId, {
			shape: this.finalShape
		});
	}

	undo(): void {
		this.mutators._updateRegion(this.regionId, {
			shape: this.originalShape
		});
	}

	merge(other: EditorCommand): boolean {
		if (
			other instanceof ResizeRegionCommand &&
			other.regionId === this.regionId
		) {
			// Update final shape to the newer command's shape
			this.finalShape = other.finalShape;
			return true;
		}
		return false;
	}
}

/**
 * Command to update region metadata (label, hint, backContent).
 * Supports merging for continuous typing.
 */
export class UpdateRegionMetadataCommand implements EditorCommand {
	readonly type: CommandType = 'update_region_metadata';
	readonly targetId: string;
	private finalUpdates: Partial<Pick<Region, 'label' | 'hint' | 'backContent'>>;

	constructor(
		private readonly mutators: EditorStateMutators,
		private readonly regionId: string,
		private readonly originalValues: Partial<Pick<Region, 'label' | 'hint' | 'backContent'>>,
		updates: Partial<Pick<Region, 'label' | 'hint' | 'backContent'>>
	) {
		this.targetId = regionId;
		this.finalUpdates = updates;
	}

	execute(): void {
		this.mutators._updateRegion(this.regionId, this.finalUpdates);
	}

	undo(): void {
		this.mutators._updateRegion(this.regionId, this.originalValues);
	}

	merge(other: EditorCommand): boolean {
		if (
			other instanceof UpdateRegionMetadataCommand &&
			other.regionId === this.regionId
		) {
			// Merge the updates
			this.finalUpdates = { ...this.finalUpdates, ...other.finalUpdates };
			return true;
		}
		return false;
	}
}

/**
 * Command to rotate the image.
 */
export class RotateImageCommand implements EditorCommand {
	readonly type: CommandType = 'rotate_image';

	constructor(
		private readonly mutators: EditorStateMutators,
		private readonly originalRotation: 0 | 90 | 180 | 270,
		private readonly newRotation: 0 | 90 | 180 | 270
	) {}

	execute(): void {
		this.mutators._setRotation(this.newRotation);
	}

	undo(): void {
		this.mutators._setRotation(this.originalRotation);
	}
}
