<script lang="ts">
	import type { ImageData, Region, DisplayContext, EditorTool, Point, RectShape, ResizeHandlePosition } from './types';
	import { calculateScale, calculateCenteredOffset, regionToDisplay, displayToImage, constrainShapeToImageBounds, getResizeCursor } from './coordinates';
	import RegionOverlay from './RegionOverlay.svelte';

	type InteractionMode = 'idle' | 'drawing' | 'moving' | 'resizing';

	interface Props {
		image: ImageData | null;
		regions: Region[];
		selectedRegionId: string | null;
		displayContext: DisplayContext;
		activeTool: EditorTool;
		onSelectRegion?: (id: string | null) => void;
		onDblClickRegion?: (id: string) => void;
		onPanChange?: (offset: Point) => void;
		onZoomChange?: (zoom: number) => void;
		onRegionCreate?: (displayShape: RectShape) => void;
		onRegionMove?: (id: string, newShape: RectShape) => void;
		onRegionMoveEnd?: (id: string, originalShape: RectShape, newPosition: Point) => void;
		onRegionResize?: (id: string, newShape: RectShape) => void;
		onRegionResizeEnd?: (id: string, originalShape: RectShape, newShape: RectShape) => void;
		onContainerResize?: (width: number, height: number) => void;
	}

	let {
		image,
		regions,
		selectedRegionId,
		displayContext,
		activeTool,
		onSelectRegion,
		onDblClickRegion,
		onPanChange,
		onZoomChange,
		onRegionCreate,
		onRegionMove,
		onRegionMoveEnd,
		onRegionResize,
		onRegionResizeEnd,
		onContainerResize
	}: Props = $props();

	let containerRef: HTMLDivElement | undefined = $state();

	// Interaction state
	let interactionMode = $state<InteractionMode>('idle');

	// Panning state
	let isPanning = $state(false);
	let panStartMouse = $state<Point>({ x: 0, y: 0 });
	let panStartOffset = $state<Point>({ x: 0, y: 0 });
	let isSpaceHeld = $state(false);

	// Drawing state
	let drawStart = $state<Point | null>(null);
	let drawCurrent = $state<Point | null>(null);
	let justFinishedDrawing = $state(false);

	// Move state
	let moveRegionId = $state<string | null>(null);
	let moveOriginalShape = $state<RectShape | null>(null);
	let moveMouseOffset = $state<Point>({ x: 0, y: 0 });
	let lastMovePosition = $state<Point | null>(null);

	// Resize state
	let resizeRegionId = $state<string | null>(null);
	let resizeOriginalShape = $state<RectShape | null>(null);
	let resizeHandle = $state<ResizeHandlePosition | null>(null);
	let resizeStartMouse = $state<Point>({ x: 0, y: 0 });
	let lastResizeShape = $state<RectShape | null>(null);

	// Draw preview in display coordinates
	let drawPreview = $derived<RectShape | null>(
		drawStart && drawCurrent
			? {
					x: Math.min(drawStart.x, drawCurrent.x),
					y: Math.min(drawStart.y, drawCurrent.y),
					width: Math.abs(drawCurrent.x - drawStart.x),
					height: Math.abs(drawCurrent.y - drawStart.y)
				}
			: null
	);

	// Set up ResizeObserver
	$effect(() => {
		if (!containerRef) return;

		const resizeObserver = new ResizeObserver((entries) => {
			for (const entry of entries) {
				const { width, height } = entry.contentRect;
				onContainerResize?.(width, height);
			}
		});

		resizeObserver.observe(containerRef);

		return () => resizeObserver.disconnect();
	});

	// Space key tracking for pan mode
	$effect(() => {
		function onKeyDown(e: KeyboardEvent) {
			if (e.key === ' ' && !e.repeat) {
				const target = e.target as HTMLElement;
				if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA') return;
				e.preventDefault();
				isSpaceHeld = true;
			}
		}
		function onKeyUp(e: KeyboardEvent) {
			if (e.key === ' ') {
				isSpaceHeld = false;
			}
		}
		window.addEventListener('keydown', onKeyDown);
		window.addEventListener('keyup', onKeyUp);
		return () => {
			window.removeEventListener('keydown', onKeyDown);
			window.removeEventListener('keyup', onKeyUp);
		};
	});

	// Window-level mouse handlers during drawing
	$effect(() => {
		if (interactionMode !== 'drawing') return;

		function onMouseMove(e: MouseEvent) {
			drawCurrent = getContainerPoint(e);
		}

		function onMouseUp() {
			const preview = drawPreview;
			// Reset state
			interactionMode = 'idle';
			drawStart = null;
			drawCurrent = null;
			// Prevent the subsequent click event from deselecting
			justFinishedDrawing = true;

			if (!preview) return;
			// Minimum 20x20 display pixels
			if (preview.width < 20 || preview.height < 20) return;

			onRegionCreate?.(preview);
		}

		window.addEventListener('mousemove', onMouseMove);
		window.addEventListener('mouseup', onMouseUp);
		return () => {
			window.removeEventListener('mousemove', onMouseMove);
			window.removeEventListener('mouseup', onMouseUp);
		};
	});

	// Window-level mouse handlers while panning
	$effect(() => {
		if (!isPanning) return;

		function onMouseMove(e: MouseEvent) {
			const dx = e.clientX - panStartMouse.x;
			const dy = e.clientY - panStartMouse.y;
			onPanChange?.({ x: panStartOffset.x + dx, y: panStartOffset.y + dy });
		}
		function onMouseUp() {
			isPanning = false;
		}

		window.addEventListener('mousemove', onMouseMove);
		window.addEventListener('mouseup', onMouseUp);
		return () => {
			window.removeEventListener('mousemove', onMouseMove);
			window.removeEventListener('mouseup', onMouseUp);
		};
	});

	// Window-level mouse handlers during move
	$effect(() => {
		if (interactionMode !== 'moving' || !moveRegionId || !moveOriginalShape) return;

		const regionId = moveRegionId;
		const originalShape = moveOriginalShape;
		const offset = moveMouseOffset;

		function onMouseMove(e: MouseEvent) {
			const mousePos = getContainerPoint(e);
			const newDisplayX = mousePos.x - offset.x;
			const newDisplayY = mousePos.y - offset.y;

			// Convert to image coords
			const newImagePos = displayToImage({ x: newDisplayX, y: newDisplayY }, displayContext);

			// Clamp to image bounds (accounting for region size)
			const clampedX = Math.max(0, Math.min(displayContext.imageWidth - originalShape.width, newImagePos.x));
			const clampedY = Math.max(0, Math.min(displayContext.imageHeight - originalShape.height, newImagePos.y));

			const pos = { x: clampedX, y: clampedY };
			lastMovePosition = pos;

			// Live visual update
			onRegionMove?.(regionId, {
				...originalShape,
				x: clampedX,
				y: clampedY
			});
		}

		function onMouseUp() {
			const finalPos = lastMovePosition;
			if (finalPos && (finalPos.x !== originalShape.x || finalPos.y !== originalShape.y)) {
				onRegionMoveEnd?.(regionId, originalShape, finalPos);
			}
			interactionMode = 'idle';
			moveRegionId = null;
			moveOriginalShape = null;
			lastMovePosition = null;
		}

		function onKeyDown(e: KeyboardEvent) {
			if (e.key === 'Escape') {
				e.preventDefault();
				e.stopPropagation();
				// Cancel: restore original position
				onRegionMove?.(regionId, originalShape);
				interactionMode = 'idle';
				moveRegionId = null;
				moveOriginalShape = null;
				lastMovePosition = null;
			}
		}

		window.addEventListener('mousemove', onMouseMove);
		window.addEventListener('mouseup', onMouseUp);
		window.addEventListener('keydown', onKeyDown);
		return () => {
			window.removeEventListener('mousemove', onMouseMove);
			window.removeEventListener('mouseup', onMouseUp);
			window.removeEventListener('keydown', onKeyDown);
		};
	});

	// Window-level mouse handlers during resize
	$effect(() => {
		if (interactionMode !== 'resizing' || !resizeRegionId || !resizeOriginalShape || !resizeHandle) return;

		const regionId = resizeRegionId;
		const originalShape = resizeOriginalShape;
		const handle = resizeHandle;
		const startMouse = resizeStartMouse;

		function onMouseMove(e: MouseEvent) {
			const currentMouse = getContainerPoint(e);
			const newShape = computeResizedShape(originalShape, handle, startMouse, currentMouse, displayContext);
			lastResizeShape = newShape;
			onRegionResize?.(regionId, newShape);
		}

		function onMouseUp() {
			const finalShape = lastResizeShape;
			if (finalShape && originalShape) {
				const changed =
					finalShape.x !== originalShape.x ||
					finalShape.y !== originalShape.y ||
					finalShape.width !== originalShape.width ||
					finalShape.height !== originalShape.height;
				if (changed) {
					onRegionResizeEnd?.(regionId, originalShape, finalShape);
				}
			}
			interactionMode = 'idle';
			resizeRegionId = null;
			resizeOriginalShape = null;
			resizeHandle = null;
			lastResizeShape = null;
		}

		function onKeyDown(e: KeyboardEvent) {
			if (e.key === 'Escape') {
				e.preventDefault();
				e.stopPropagation();
				onRegionResize?.(regionId, originalShape);
				interactionMode = 'idle';
				resizeRegionId = null;
				resizeOriginalShape = null;
				resizeHandle = null;
				lastResizeShape = null;
			}
		}

		window.addEventListener('mousemove', onMouseMove);
		window.addEventListener('mouseup', onMouseUp);
		window.addEventListener('keydown', onKeyDown);
		return () => {
			window.removeEventListener('mousemove', onMouseMove);
			window.removeEventListener('mouseup', onMouseUp);
			window.removeEventListener('keydown', onKeyDown);
		};
	});

	function handleMoveStart(regionId: string, e: MouseEvent) {
		if (interactionMode !== 'idle') return;
		const region = regions.find(r => r.id === regionId);
		if (!region) return;

		interactionMode = 'moving';
		moveRegionId = regionId;
		moveOriginalShape = { ...region.shape };

		// Offset between mouse and region origin (in display coords)
		const displayShape = regionToDisplay(region.shape, displayContext);
		const mousePos = getContainerPoint(e);
		moveMouseOffset = {
			x: mousePos.x - displayShape.x,
			y: mousePos.y - displayShape.y
		};
	}

	function handleResizeStart(regionId: string, e: MouseEvent, handle: ResizeHandlePosition) {
		if (interactionMode !== 'idle') return;
		const region = regions.find(r => r.id === regionId);
		if (!region) return;

		e.preventDefault();
		interactionMode = 'resizing';
		resizeRegionId = regionId;
		resizeOriginalShape = { ...region.shape };
		resizeHandle = handle;
		resizeStartMouse = getContainerPoint(e);
	}

	function computeResizedShape(
		original: RectShape,
		handle: ResizeHandlePosition,
		startMouse: Point,
		currentMouse: Point,
		ctx: DisplayContext
	): RectShape {
		const scale = calculateScale(ctx) * ctx.zoom;
		// Convert mouse delta from display to image pixels
		const dx = (currentMouse.x - startMouse.x) / scale;
		const dy = (currentMouse.y - startMouse.y) / scale;

		let { x, y, width, height } = original;

		// Adjust edges based on handle position
		switch (handle) {
			case 'nw': x += dx; y += dy; width -= dx; height -= dy; break;
			case 'n':  y += dy; height -= dy; break;
			case 'ne': y += dy; width += dx; height -= dy; break;
			case 'w':  x += dx; width -= dx; break;
			case 'e':  width += dx; break;
			case 'sw': x += dx; width -= dx; height += dy; break;
			case 's':  height += dy; break;
			case 'se': width += dx; height += dy; break;
		}

		// Handle flipping (drag past opposite edge)
		if (width < 0) { x += width; width = -width; }
		if (height < 0) { y += height; height = -height; }

		// Enforce minimum size (20 image pixels)
		const MIN_SIZE = 20;
		if (width < MIN_SIZE) width = MIN_SIZE;
		if (height < MIN_SIZE) height = MIN_SIZE;

		// Constrain to image bounds
		return constrainShapeToImageBounds({ x, y, width, height }, ctx.imageWidth, ctx.imageHeight);
	}

	// Calculate image transform with transform-origin: 0 0
	let imageTransform = $derived(
		(() => {
			if (!image || displayContext.containerWidth === 0) return '';

			const scale = calculateScale(displayContext) * displayContext.zoom;
			const offset = calculateCenteredOffset(displayContext);

			const transforms: string[] = [];

			// Position image at centered offset
			transforms.push(`translate(${offset.x}px, ${offset.y}px)`);

			// Scale from top-left (transform-origin: 0 0)
			transforms.push(`scale(${scale})`);

			// Rotation around image center (in unscaled image coords, after scale is applied)
			if (displayContext.rotation !== 0) {
				const cx = image.width / 2;
				const cy = image.height / 2;
				transforms.push(`translate(${cx}px, ${cy}px)`);
				transforms.push(`rotate(${displayContext.rotation}deg)`);
				transforms.push(`translate(${-cx}px, ${-cy}px)`);
			}

			return transforms.join(' ');
		})()
	);

	// Image dimensions for SVG viewBox
	let svgViewBox = $derived(
		!displayContext.containerWidth || !displayContext.containerHeight
			? '0 0 100 100'
			: `0 0 ${displayContext.containerWidth} ${displayContext.containerHeight}`
	);

	function handleWheel(e: WheelEvent) {
		if (!e.ctrlKey && !e.metaKey) return;
		e.preventDefault();
		const zoomFactor = 1 - e.deltaY * 0.003;
		const newZoom = Math.max(0.1, Math.min(5, displayContext.zoom * zoomFactor));
		onZoomChange?.(newZoom);
	}

	function getContainerPoint(e: MouseEvent): Point {
		if (!containerRef) return { x: 0, y: 0 };
		const rect = containerRef.getBoundingClientRect();
		return { x: e.clientX - rect.left, y: e.clientY - rect.top };
	}

	function handleCanvasMouseDown(e: MouseEvent) {
		// Start drawing when draw tool is active
		if (e.button === 0 && activeTool === 'draw_region' && interactionMode === 'idle' && !isSpaceHeld) {
			e.preventDefault();
			interactionMode = 'drawing';
			const point = getContainerPoint(e);
			drawStart = point;
			drawCurrent = point;
			return;
		}

		// Middle mouse button or space + left click = pan
		if (e.button === 1 || (e.button === 0 && isSpaceHeld)) {
			e.preventDefault();
			isPanning = true;
			panStartMouse = { x: e.clientX, y: e.clientY };
			panStartOffset = { ...displayContext.panOffset };
			return;
		}
	}

	function handleCanvasClick(e: MouseEvent) {
		// Don't deselect if we were panning
		if (isSpaceHeld) return;

		// Don't deselect after finishing a draw (click fires after mouseup)
		if (justFinishedDrawing) {
			justFinishedDrawing = false;
			return;
		}

		// Only deselect if clicking directly on the canvas (not on a region)
		if (e.target === e.currentTarget || (e.target as Element).tagName === 'svg') {
			onSelectRegion?.(null);
		}
	}

	function handleRegionClick(regionId: string) {
		if (isSpaceHeld) return;
		onSelectRegion?.(regionId);
	}

	// Cursor based on active tool and pan state
	let canvasCursor = $derived(
		isPanning ? 'grabbing'
		: isSpaceHeld ? 'grab'
		: interactionMode === 'resizing' && resizeHandle ? getResizeCursor(resizeHandle)
		: interactionMode === 'drawing' ? 'crosshair'
		: interactionMode === 'moving' ? 'move'
		: activeTool === 'draw_region' ? 'crosshair'
		: activeTool === 'draw_arrow' ? 'crosshair'
		: 'default'
	);
</script>

<div
	bind:this={containerRef}
	class="relative h-full w-full overflow-hidden bg-muted/30"
	style="cursor: {canvasCursor}"
	onmousedown={handleCanvasMouseDown}
	onclick={handleCanvasClick}
	onwheel={handleWheel}
	onkeydown={(e) => e.key === 'Escape' && interactionMode === 'idle' && onSelectRegion?.(null)}
	tabindex="0"
	role="application"
	aria-label="Image occlusion editor canvas"
>
	{#if image}
		<!-- Image layer with explicit dimensions and transform-origin: 0 0 -->
		<img
			src={image.url}
			alt="Occlusion target"
			class="pointer-events-none absolute left-0 top-0 max-w-none"
			style="width: {image.width}px; height: {image.height}px; transform: {imageTransform}; transform-origin: 0 0;"
			draggable="false"
		/>

		<!-- SVG overlay for regions -->
		<svg
			class="absolute inset-0 h-full w-full"
			viewBox={svgViewBox}
			preserveAspectRatio="none"
		>
			{#each regions as region (region.id)}
				<RegionOverlay
					{region}
					isSelected={region.id === selectedRegionId}
					{displayContext}
					onClick={() => handleRegionClick(region.id)}
					onDblClick={() => onDblClickRegion?.(region.id)}
					onMoveStart={(e) => handleMoveStart(region.id, e)}
					onResizeStart={(e, pos) => handleResizeStart(region.id, e, pos)}
				/>
			{/each}
			{#if drawPreview}
				<rect
					x={drawPreview.x}
					y={drawPreview.y}
					width={drawPreview.width}
					height={drawPreview.height}
					fill="oklch(from var(--primary) l c h / 0.2)"
					stroke="var(--primary)"
					stroke-width="2"
					stroke-dasharray="6 3"
					pointer-events="none"
				/>
			{/if}
		</svg>
	{:else}
		<!-- Empty state placeholder -->
		<div class="flex h-full w-full items-center justify-center">
			<div class="text-center text-muted-foreground">
				<p class="text-sm">No image loaded</p>
			</div>
		</div>
	{/if}
</div>
