<script lang="ts">
	import type { ImageData, Region, DisplayContext, EditorTool, Point, RectShape, ResizeHandlePosition } from './types';
	import { calculateScale, calculateCenteredOffset, regionToDisplay, displayToImage } from './coordinates';
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

	// Move state
	let moveRegionId = $state<string | null>(null);
	let moveOriginalShape = $state<RectShape | null>(null);
	let moveMouseOffset = $state<Point>({ x: 0, y: 0 });

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

			// Live visual update
			onRegionMove?.(regionId, {
				...originalShape,
				x: clampedX,
				y: clampedY
			});
		}

		function onMouseUp() {
			// Find current region state to get final position
			const region = regions.find(r => r.id === regionId);
			if (region && (region.shape.x !== originalShape.x || region.shape.y !== originalShape.y)) {
				onRegionMoveEnd?.(regionId, originalShape, { x: region.shape.x, y: region.shape.y });
			}
			interactionMode = 'idle';
			moveRegionId = null;
			moveOriginalShape = null;
		}

		function onKeyDown(e: KeyboardEvent) {
			if (e.key === 'Escape') {
				// Cancel: restore original position
				onRegionMove?.(regionId, originalShape);
				interactionMode = 'idle';
				moveRegionId = null;
				moveOriginalShape = null;
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
		// Will be implemented in Task 5
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
	onkeydown={(e) => e.key === 'Escape' && onSelectRegion?.(null)}
	tabindex="0"
	role="application"
	aria-label="Image occlusion editor canvas"
>
	{#if image}
		<!-- Image layer with explicit dimensions and transform-origin: 0 0 -->
		<img
			src={image.url}
			alt="Occlusion target"
			class="pointer-events-none absolute left-0 top-0"
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
