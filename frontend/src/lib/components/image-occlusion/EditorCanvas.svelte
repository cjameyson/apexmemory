<script lang="ts">
	import type { ImageData, Region, DisplayContext, EditorTool, Point } from './types';
	import { calculateScale, calculateCenteredOffset } from './coordinates';
	import RegionOverlay from './RegionOverlay.svelte';

	interface Props {
		image: ImageData | null;
		regions: Region[];
		selectedRegionId: string | null;
		displayContext: DisplayContext;
		activeTool: EditorTool;
		onSelectRegion?: (id: string | null) => void;
		onDblClickRegion?: (id: string) => void;
		onPanChange?: (offset: Point) => void;
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
		onContainerResize
	}: Props = $props();

	let containerRef: HTMLDivElement | undefined = $state();

	// Panning state
	let isPanning = $state(false);
	let panStartMouse = $state<Point>({ x: 0, y: 0 });
	let panStartOffset = $state<Point>({ x: 0, y: 0 });
	let isSpaceHeld = $state(false);

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

	// Calculate image transform with transform-origin: 0 0
	let imageTransform = $derived(() => {
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
	});

	// Image dimensions for SVG viewBox
	let svgViewBox = $derived(() => {
		if (!displayContext.containerWidth || !displayContext.containerHeight) {
			return '0 0 100 100';
		}
		return `0 0 ${displayContext.containerWidth} ${displayContext.containerHeight}`;
	});

	function handleCanvasMouseDown(e: MouseEvent) {
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
	let canvasCursor = $derived(() => {
		if (isPanning) return 'grabbing';
		if (isSpaceHeld) return 'grab';
		switch (activeTool) {
			case 'draw_region':
				return 'crosshair';
			case 'draw_arrow':
				return 'crosshair';
			default:
				return 'default';
		}
	});
</script>

<div
	bind:this={containerRef}
	class="relative h-full w-full overflow-hidden bg-muted/30"
	style="cursor: {canvasCursor()}"
	onmousedown={handleCanvasMouseDown}
	onclick={handleCanvasClick}
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
			style="width: {image.width}px; height: {image.height}px; transform: {imageTransform()}; transform-origin: 0 0;"
			draggable="false"
		/>

		<!-- SVG overlay for regions -->
		<svg
			class="absolute inset-0 h-full w-full"
			viewBox={svgViewBox()}
			preserveAspectRatio="none"
		>
			{#each regions as region (region.id)}
				<RegionOverlay
					{region}
					isSelected={region.id === selectedRegionId}
					{displayContext}
					onClick={() => handleRegionClick(region.id)}
					onDblClick={() => onDblClickRegion?.(region.id)}
				/>
			{/each}
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
