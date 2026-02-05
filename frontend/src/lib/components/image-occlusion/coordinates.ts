/**
 * Coordinate transformation utilities for the image occlusion editor.
 *
 * The editor works with two coordinate systems:
 * 1. Image coordinates: Natural pixel coordinates of the source image
 * 2. Display coordinates: Screen pixel coordinates after zoom/pan
 *
 * Regions are stored in image coordinates and transformed to display coordinates for rendering.
 * Note: Rotation only affects the image display, not the annotation coordinates.
 */

import type { Point, RectShape, DisplayContext, ResizeHandle, ResizeHandlePosition } from './types';

/**
 * Calculate the scale factor to fit the image within the container.
 * Uses the original image dimensions (rotation doesn't affect annotation scale).
 */
export function calculateScale(ctx: DisplayContext): number {
	const scaleX = ctx.containerWidth / ctx.imageWidth;
	const scaleY = ctx.containerHeight / ctx.imageHeight;
	return Math.min(scaleX, scaleY, 1); // Never scale up
}

/**
 * Calculate the offset to center the image within the container.
 */
export function calculateCenteredOffset(ctx: DisplayContext): Point {
	const scale = calculateScale(ctx) * ctx.zoom;
	const scaledWidth = ctx.imageWidth * scale;
	const scaledHeight = ctx.imageHeight * scale;

	return {
		x: (ctx.containerWidth - scaledWidth) / 2 + ctx.panOffset.x,
		y: (ctx.containerHeight - scaledHeight) / 2 + ctx.panOffset.y
	};
}

/**
 * Transform a point from image coordinates to display coordinates.
 * Does NOT apply rotation - rotation is a visual-only transform on the image element.
 */
export function imageToDisplay(point: Point, ctx: DisplayContext): Point {
	const scale = calculateScale(ctx) * ctx.zoom;
	const offset = calculateCenteredOffset(ctx);

	return {
		x: point.x * scale + offset.x,
		y: point.y * scale + offset.y
	};
}

/**
 * Transform a point from display coordinates to image coordinates.
 */
export function displayToImage(point: Point, ctx: DisplayContext): Point {
	const scale = calculateScale(ctx) * ctx.zoom;
	const offset = calculateCenteredOffset(ctx);

	return {
		x: (point.x - offset.x) / scale,
		y: (point.y - offset.y) / scale
	};
}

/**
 * Transform a region shape from image coordinates to display coordinates.
 */
export function regionToDisplay(shape: RectShape, ctx: DisplayContext): RectShape {
	const scale = calculateScale(ctx) * ctx.zoom;
	const offset = calculateCenteredOffset(ctx);

	return {
		x: shape.x * scale + offset.x,
		y: shape.y * scale + offset.y,
		width: shape.width * scale,
		height: shape.height * scale
	};
}

/**
 * Transform a region shape from display coordinates to image coordinates.
 */
export function displayToRegion(shape: RectShape, ctx: DisplayContext): RectShape {
	const scale = calculateScale(ctx) * ctx.zoom;
	const offset = calculateCenteredOffset(ctx);

	return {
		x: (shape.x - offset.x) / scale,
		y: (shape.y - offset.y) / scale,
		width: shape.width / scale,
		height: shape.height / scale
	};
}

/**
 * Get the 8 resize handle positions for a shape in display coordinates.
 */
export function getResizeHandles(shape: RectShape, handleSize: number = 8): ResizeHandle[] {
	const halfSize = handleSize / 2;
	const positions: ResizeHandlePosition[] = ['nw', 'n', 'ne', 'w', 'e', 'sw', 's', 'se'];

	const getPosition = (pos: ResizeHandlePosition): Point => {
		const centerX = shape.x + shape.width / 2;
		const centerY = shape.y + shape.height / 2;

		switch (pos) {
			case 'nw':
				return { x: shape.x, y: shape.y };
			case 'n':
				return { x: centerX, y: shape.y };
			case 'ne':
				return { x: shape.x + shape.width, y: shape.y };
			case 'w':
				return { x: shape.x, y: centerY };
			case 'e':
				return { x: shape.x + shape.width, y: centerY };
			case 'sw':
				return { x: shape.x, y: shape.y + shape.height };
			case 's':
				return { x: centerX, y: shape.y + shape.height };
			case 'se':
				return { x: shape.x + shape.width, y: shape.y + shape.height };
		}
	};

	return positions.map((pos) => {
		const point = getPosition(pos);
		return {
			position: pos,
			x: point.x - halfSize,
			y: point.y - halfSize,
			size: handleSize
		};
	});
}

/**
 * Constrain a point to be within image bounds.
 */
export function constrainToImageBounds(point: Point, imageWidth: number, imageHeight: number): Point {
	return {
		x: Math.max(0, Math.min(imageWidth, point.x)),
		y: Math.max(0, Math.min(imageHeight, point.y))
	};
}

/**
 * Constrain a shape to be within image bounds, adjusting position and size.
 */
export function constrainShapeToImageBounds(
	shape: RectShape,
	imageWidth: number,
	imageHeight: number
): RectShape {
	let { x, y, width, height } = shape;

	// Constrain position
	x = Math.max(0, Math.min(imageWidth - width, x));
	y = Math.max(0, Math.min(imageHeight - height, y));

	// If shape is larger than image, constrain size
	if (width > imageWidth) {
		width = imageWidth;
		x = 0;
	}
	if (height > imageHeight) {
		height = imageHeight;
		y = 0;
	}

	return { x, y, width, height };
}

/**
 * Check if a point (in display coordinates) is within the image bounds.
 */
export function isPointInImage(point: Point, ctx: DisplayContext): boolean {
	const imagePoint = displayToImage(point, ctx);
	return (
		imagePoint.x >= 0 &&
		imagePoint.x <= ctx.imageWidth &&
		imagePoint.y >= 0 &&
		imagePoint.y <= ctx.imageHeight
	);
}

/**
 * Check if a point (in display coordinates) is within a shape (also in display coordinates).
 */
export function isPointInShape(point: Point, shape: RectShape): boolean {
	return (
		point.x >= shape.x &&
		point.x <= shape.x + shape.width &&
		point.y >= shape.y &&
		point.y <= shape.y + shape.height
	);
}

/**
 * Get the cursor style for a resize handle position.
 */
export function getResizeCursor(position: ResizeHandlePosition): string {
	switch (position) {
		case 'nw':
		case 'se':
			return 'nwse-resize';
		case 'ne':
		case 'sw':
			return 'nesw-resize';
		case 'n':
		case 's':
			return 'ns-resize';
		case 'e':
		case 'w':
			return 'ew-resize';
	}
}
