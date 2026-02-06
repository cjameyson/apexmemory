<script lang="ts">
	import type { Region, DisplayContext, RectShape, ResizeHandlePosition } from './types';
	import { regionToDisplay, getResizeHandles } from './coordinates';

	interface Props {
		region: Region;
		index: number;
		isSelected: boolean;
		showLabels?: boolean;
		displayContext: DisplayContext;
		onClick?: () => void;
		onDblClick?: () => void;
		onMoveStart?: (e: MouseEvent) => void;
		onResizeStart?: (e: MouseEvent, position: ResizeHandlePosition) => void;
	}

	let { region, index, isSelected, showLabels = true, displayContext, onClick, onDblClick, onMoveStart, onResizeStart }: Props = $props();

	// Transform region to display coordinates
	let displayShape = $derived<RectShape>(regionToDisplay(region.shape, displayContext));

	// Resize handles (only shown when selected)
	let handles = $derived(isSelected ? getResizeHandles(displayShape) : []);

	let isUnlabeled = $derived(!region.label.trim());
	let regionLabel = $derived(region.label || '?');

	// Dynamic color: selected = success (green), unlabeled = warning (amber), normal = primary (blue)
	let regionColor = $derived(
		isSelected ? 'var(--success)' : isUnlabeled ? 'var(--warning)' : 'var(--primary)'
	);
	let regionColorFg = $derived(
		isSelected
			? 'var(--success-foreground)'
			: isUnlabeled
				? 'var(--warning-foreground)'
				: 'var(--primary-foreground)'
	);

	// Index badge centered on top-left corner of region rect
	const BADGE_R = 9;
	let badgeCx = $derived(displayShape.x);
	let badgeCy = $derived(displayShape.y);
</script>

<g
	class="region-overlay"
	style="--rc: {regionColor}; --rc-fg: {regionColorFg}"
	role="button"
	tabindex="0"
	onclick={onClick}
	ondblclick={onDblClick}
	onkeydown={(e) => e.key === 'Enter' && onClick?.()}
>
	<!-- Region fill (semi-transparent) -->
	<rect
		x={displayShape.x}
		y={displayShape.y}
		width={displayShape.width}
		height={displayShape.height}
		class="region-fill"
		class:selected={isSelected}
		onmousedown={(e) => {
			if (e.button === 0) {
				// Always stop propagation to prevent canvas draw handler
				e.stopPropagation();
				if (isSelected) {
					onMoveStart?.(e);
				}
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

	{#if showLabels}
		<!-- Region label in center -->
		<text
			x={displayShape.x + displayShape.width / 2}
			y={displayShape.y + displayShape.height / 2}
			class="region-label"
			text-anchor="middle"
			dominant-baseline="central"
		>
			{regionLabel}
		</text>
	{/if}

	<!-- Index badge (centered on top-left corner, always visible) -->
	<circle cx={badgeCx} cy={badgeCy} r={BADGE_R} class="badge-circle" />
	<text
		x={badgeCx}
		y={badgeCy}
		class="badge-text"
		text-anchor="middle"
		dominant-baseline="central"
	>
		{index}
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
		fill: oklch(from var(--rc) l c h / 0.3);
		cursor: pointer;
		transition: fill 0.15s ease;
	}

	.region-fill:hover {
		fill: oklch(from var(--rc) l c h / 0.4);
	}

	.region-fill.selected {
		fill: oklch(from var(--rc) l c h / 0.35);
	}

	.region-border {
		fill: none;
		stroke: var(--rc);
		stroke-width: 2;
		pointer-events: none;
	}

	.region-border-ants {
		fill: none;
		stroke: var(--rc);
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
		fill: white;
		font-size: 14px;
		font-weight: 600;
		pointer-events: none;
		user-select: none;
		text-shadow:
			0 1px 2px rgba(0, 0, 0, 0.6),
			0 0 3px rgba(0, 0, 0, 0.3);
	}

	.badge-circle {
		fill: var(--rc);
		opacity: 0.9;
	}

	.badge-text {
		fill: var(--rc-fg);
		font-size: 11px;
		font-weight: 700;
		pointer-events: none;
		user-select: none;
	}

	.resize-handle {
		fill: var(--rc);
		stroke: var(--rc-fg);
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
