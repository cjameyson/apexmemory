<script lang="ts">
	import type { Region, DisplayContext, RectShape, ResizeHandlePosition } from './types';
	import { regionToDisplay, getResizeHandles } from './coordinates';

	interface Props {
		region: Region;
		isSelected: boolean;
		displayContext: DisplayContext;
		onClick?: () => void;
		onDblClick?: () => void;
		onMoveStart?: (e: MouseEvent) => void;
		onResizeStart?: (e: MouseEvent, position: ResizeHandlePosition) => void;
	}

	let { region, isSelected, displayContext, onClick, onDblClick, onMoveStart, onResizeStart }: Props = $props();

	// Transform region to display coordinates
	let displayShape = $derived<RectShape>(regionToDisplay(region.shape, displayContext));

	// Resize handles (only shown when selected)
	let handles = $derived(isSelected ? getResizeHandles(displayShape) : []);

	// Region index for display
	let regionIndex = $derived(
		(() => {
			// This would ideally come from parent, but we show label for now
			return region.label || '?';
		})()
	);
</script>

<g
	class="region-overlay"
	role="button"
	tabindex="0"
	onclick={onClick}
	ondblclick={onDblClick}
	onkeydown={(e) => e.key === 'Enter' && onClick?.()}
>
	<!-- Region fill (semi-transparent primary color) -->
	<rect
		x={displayShape.x}
		y={displayShape.y}
		width={displayShape.width}
		height={displayShape.height}
		class="region-fill"
		class:selected={isSelected}
		onmousedown={(e) => {
			if (isSelected && e.button === 0) {
				e.stopPropagation();
				onMoveStart?.(e);
			}
		}}
	/>

	<!-- Region border -->
	{#if isSelected}
		<!-- Marching ants animation for selected region -->
		<rect
			x={displayShape.x}
			y={displayShape.y}
			width={displayShape.width}
			height={displayShape.height}
			class="region-border-ants"
		/>
	{:else}
		<!-- Solid border for unselected -->
		<rect
			x={displayShape.x}
			y={displayShape.y}
			width={displayShape.width}
			height={displayShape.height}
			class="region-border"
		/>
	{/if}

	<!-- Region label in center -->
	<text
		x={displayShape.x + displayShape.width / 2}
		y={displayShape.y + displayShape.height / 2}
		class="region-label"
		text-anchor="middle"
		dominant-baseline="central"
	>
		{regionIndex}
	</text>

	<!-- Resize handles (when selected) -->
	{#if isSelected}
		{#each handles as handle}
			<rect
				x={handle.x}
				y={handle.y}
				width={handle.size}
				height={handle.size}
				class="resize-handle"
				data-position={handle.position}
				onmousedown={(e) => {
					e.stopPropagation();
					onResizeStart?.(e, handle.position);
				}}
			/>
		{/each}
	{/if}
</g>

<style>
	.region-fill {
		fill: oklch(from var(--primary) l c h / 0.3);
		cursor: pointer;
		transition: fill 0.15s ease;
	}

	.region-fill:hover {
		fill: oklch(from var(--primary) l c h / 0.4);
	}

	.region-fill.selected {
		fill: oklch(from var(--primary) l c h / 0.35);
	}

	.region-border {
		fill: none;
		stroke: var(--primary);
		stroke-width: 2;
		pointer-events: none;
	}

	.region-border-ants {
		fill: none;
		stroke: var(--primary);
		stroke-width: 2;
		stroke-dasharray: 8 4;
		pointer-events: none;
		animation: marching-ants 0.5s linear infinite;
	}

	@keyframes marching-ants {
		0% {
			stroke-dashoffset: 0;
		}
		100% {
			stroke-dashoffset: 12;
		}
	}

	.region-label {
		fill: var(--primary-foreground);
		font-size: 14px;
		font-weight: 600;
		pointer-events: none;
		user-select: none;
		text-shadow:
			0 1px 2px oklch(from var(--primary) calc(l - 0.3) c h / 0.8),
			0 0 8px oklch(from var(--primary) l c h / 0.5);
	}

	.resize-handle {
		fill: var(--primary);
		stroke: var(--primary-foreground);
		stroke-width: 1;
		cursor: pointer;
	}

	/* Cursor styles for resize handles */
	.resize-handle[data-position='nw'],
	.resize-handle[data-position='se'] {
		cursor: nwse-resize;
	}

	.resize-handle[data-position='ne'],
	.resize-handle[data-position='sw'] {
		cursor: nesw-resize;
	}

	.resize-handle[data-position='n'],
	.resize-handle[data-position='s'] {
		cursor: ns-resize;
	}

	.resize-handle[data-position='e'],
	.resize-handle[data-position='w'] {
		cursor: ew-resize;
	}
</style>
