<script lang="ts">
	import type { EditorTool } from './types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Separator } from '$lib/components/ui/separator';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import {
		Undo2,
		Redo2,
		MousePointer2,
		Square,
		RotateCw,
		ZoomIn,
		ZoomOut,
		Maximize,
		Sparkles,
		Type
	} from '@lucide/svelte';

	export type ToolbarPosition = 'left' | 'center' | 'right';

	interface Props {
		activeTool: EditorTool;
		canUndo: boolean;
		canRedo: boolean;
		zoom: number;
		showLabels?: boolean;
		position?: ToolbarPosition;
		onToolChange?: (tool: EditorTool) => void;
		onToggleLabels?: () => void;
		onUndo?: () => void;
		onRedo?: () => void;
		onRotate?: () => void;
		onZoomChange?: (zoom: number) => void;
		onZoomFit?: () => void;
		onPositionChange?: (position: ToolbarPosition) => void;
	}

	let {
		activeTool,
		canUndo,
		canRedo,
		zoom,
		showLabels = true,
		position = 'center',
		onToolChange,
		onToggleLabels,
		onUndo,
		onRedo,
		onRotate,
		onZoomChange,
		onZoomFit,
		onPositionChange
	}: Props = $props();

	function handleZoomIn() {
		onZoomChange?.(Math.min(5, zoom + 0.25));
	}

	function handleZoomOut() {
		onZoomChange?.(Math.max(0.1, zoom - 0.25));
	}

	let zoomPercentage = $derived(Math.round(zoom * 100));
</script>

<Tooltip.Provider ignoreNonKeyboardFocus={true}>
	<div
		role="toolbar"
		data-toolbar-name="Image Occlusion Editor Toolbar"
		class="border-border bg-card/95 inline-flex items-center gap-1 rounded-md border px-2 py-1 shadow-sm backdrop-blur-sm"
	>
		<!-- Undo/Redo group -->
		<div class="flex items-center">
			<Tooltip.Root>
				<Tooltip.Trigger>
					<Button
						variant="ghost"
						size="icon-sm"
						disabled={!canUndo}
						onclick={onUndo}
						aria-label="Undo"
						tabindex={-1}
					>
						<Undo2 class="h-4 w-4" />
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content>
					<p>Undo (Cmd+Z)</p>
				</Tooltip.Content>
			</Tooltip.Root>

			<Tooltip.Root>
				<Tooltip.Trigger>
					<Button
						variant="ghost"
						size="icon-sm"
						disabled={!canRedo}
						onclick={onRedo}
						aria-label="Redo"
					>
						<Redo2 class="h-4 w-4" />
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content>
					<p>Redo (Cmd+Shift+Z)</p>
				</Tooltip.Content>
			</Tooltip.Root>
		</div>

		<Separator orientation="vertical" class="mx-1 h-6" />

		<!-- Tools group -->
		<div class="flex items-center">
			<Tooltip.Root>
				<Tooltip.Trigger>
					<Button
						variant={activeTool === 'select' ? 'default' : 'ghost'}
						size="icon-sm"
						onclick={() => onToolChange?.('select')}
						aria-label="Select tool"
						aria-pressed={activeTool === 'select'}
					>
						<MousePointer2 class="h-4 w-4" />
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content>
					<p>Select (V)</p>
				</Tooltip.Content>
			</Tooltip.Root>

			<Tooltip.Root>
				<Tooltip.Trigger>
					<Button
						variant={activeTool === 'draw_region' ? 'default' : 'ghost'}
						size="icon-sm"
						onclick={() => onToolChange?.('draw_region')}
						aria-label="Draw region tool"
						aria-pressed={activeTool === 'draw_region'}
					>
						<Square class="h-4 w-4" />
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content>
					<p>Draw Region (R)</p>
				</Tooltip.Content>
			</Tooltip.Root>
		</div>

		<Separator orientation="vertical" class="mx-1 h-6" />

		<!-- Rotate -->
		<Tooltip.Root>
			<Tooltip.Trigger>
				<Button variant="ghost" size="icon-sm" onclick={onRotate} aria-label="Rotate 90 degrees">
					<RotateCw class="h-4 w-4" />
				</Button>
			</Tooltip.Trigger>
			<Tooltip.Content>
				<p>Rotate 90°</p>
			</Tooltip.Content>
		</Tooltip.Root>

		<!-- Toggle labels -->
		<Tooltip.Root>
			<Tooltip.Trigger>
				<Button
					variant="ghost"
					size="icon-sm"
					class={showLabels ? 'bg-primary/15 text-primary hover:bg-primary/25' : ''}
					onclick={onToggleLabels}
					aria-label="Toggle region labels"
					aria-pressed={showLabels}
				>
					<Type class="h-4 w-4" />
				</Button>
			</Tooltip.Trigger>
			<Tooltip.Content>
				<p>{showLabels ? 'Hide' : 'Show'} Labels (L)</p>
			</Tooltip.Content>
		</Tooltip.Root>

		<Separator orientation="vertical" class="mx-1 h-6" />

		<!-- Zoom group -->
		<div class="flex items-center gap-0.5">
			<Tooltip.Root>
				<Tooltip.Trigger>
					<Button
						variant="ghost"
						size="icon-sm"
						onclick={handleZoomOut}
						disabled={zoom <= 0.1}
						aria-label="Zoom out"
					>
						<ZoomOut class="h-4 w-4" />
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content>
					<p>Zoom Out</p>
				</Tooltip.Content>
			</Tooltip.Root>

			<span class="text-muted-foreground min-w-[48px] text-center text-xs font-medium">
				{zoomPercentage}%
			</span>

			<Tooltip.Root>
				<Tooltip.Trigger>
					<Button
						variant="ghost"
						size="icon-sm"
						onclick={handleZoomIn}
						disabled={zoom >= 5}
						aria-label="Zoom in"
					>
						<ZoomIn class="h-4 w-4" />
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content>
					<p>Zoom In</p>
				</Tooltip.Content>
			</Tooltip.Root>

			<Tooltip.Root>
				<Tooltip.Trigger>
					<Button variant="ghost" size="icon-sm" onclick={onZoomFit} aria-label="Fit to screen">
						<Maximize class="h-4 w-4" />
					</Button>
				</Tooltip.Trigger>
				<Tooltip.Content>
					<p>Fit to Screen</p>
				</Tooltip.Content>
			</Tooltip.Root>
		</div>

		<!-- AI Assist (hidden for now)
		<Separator orientation="vertical" class="mx-1 h-6" />
		<Tooltip.Root>
			<Tooltip.Trigger>
				<Button
					variant="ghost"
					size="sm"
					disabled={true}
					aria-label="AI Assist (coming soon)"
					class="relative gap-1.5 opacity-60"
				>
					<Sparkles class="h-4 w-4" />
					<span class="text-xs">AI</span>
					<span
						class="bg-muted text-muted-foreground absolute -right-1 -top-1 rounded px-1 text-[9px] font-medium leading-tight"
					>
						Soon
					</span>
				</Button>
			</Tooltip.Trigger>
			<Tooltip.Content>
				<p>AI-assisted region detection (coming soon)</p>
			</Tooltip.Content>
		</Tooltip.Root>
		-->
	</div>
</Tooltip.Provider>
