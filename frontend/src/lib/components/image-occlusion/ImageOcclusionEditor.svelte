<script lang="ts">
	import type { ImageOcclusionField, Region, RectShape, Point } from './types';
	import { createEditorState } from './editor-state.svelte';
	import { createHistoryManager } from './history.svelte';
	import {
		CreateRegionCommand,
		DeleteRegionCommand,
		MoveRegionCommand,
		ResizeRegionCommand,
		UpdateRegionMetadataCommand,
		RotateImageCommand
	} from './commands';
	import { displayToRegion, constrainShapeToImageBounds } from './coordinates';
	import EditorCanvas from './EditorCanvas.svelte';
	import EditorToolbar, { type ToolbarPosition } from './EditorToolbar.svelte';
	import LabelPanel from './LabelPanel.svelte';
	import StatusBar from './StatusBar.svelte';
	import ImageUploader from './ImageUploader.svelte';
	import ConfirmDialog from '$lib/components/ui/confirm-dialog.svelte';
	import { generateRegionId } from './utils';
	import { uploadAsset, assetUrl } from '$lib/api/client';

	interface Props {
		initialValue?: ImageOcclusionField;
		onChange?: (field: ImageOcclusionField) => void;
		errors?: { title?: boolean; regionLabels?: Set<string> };
	}

	let { initialValue, onChange, errors }: Props = $props();

	// Create editor state and history managers
	const editor = createEditorState();
	const history = createHistoryManager();

	// Panel width state
	let panelWidth = $state(380);
	let isDragging = $state(false);
	let containerRef: HTMLDivElement | undefined = $state();

	// Focus label signal: set to a region id to trigger label input focus+select
	let focusLabelRegionId = $state<string | null>(null);

	// Delete confirmation state
	let deleteConfirmOpen = $state(false);
	let deletingRegionId = $state<string | null>(null);
	let deletingRegionLabel = $derived(
		deletingRegionId ? (editor.regions.find((r) => r.id === deletingRegionId)?.label ?? '') : ''
	);

	// Dirty state (exposed for parent to use in close confirmation)
	export function getIsDirty() {
		return history.undoCount > 0;
	}

	// Card title
	let cardTitle = $state('');

	// Toolbar position
	let toolbarPosition = $state<ToolbarPosition>('left');

	// Filter: set of visible region IDs (null = show all)
	let visibleRegionIds = $state<Set<string> | null>(null);
	let canvasRegions = $derived(
		visibleRegionIds ? editor.regions.filter((r) => visibleRegionIds!.has(r.id)) : editor.regions
	);
	const toolbarJustify: Record<ToolbarPosition, string> = {
		left: 'justify-start',
		center: 'justify-center',
		right: 'justify-end'
	};

	// Initialize from provided initial value
	$effect(() => {
		if (initialValue) {
			editor.initialize(
				initialValue.image,
				initialValue.regions,
				initialValue.annotations,
				initialValue.mode
			);
			cardTitle = initialValue.title ?? '';
		}
	});

	// Notify parent of changes
	function notifyChange() {
		if (!editor.image) return;

		const field: ImageOcclusionField = {
			version: 1,
			title: cardTitle,
			image: editor.image,
			regions: editor.regions,
			annotations: editor.annotations,
			mode: editor.mode
		};

		onChange?.(field);
	}

	// Call notifyChange whenever regions or title change
	$effect(() => {
		// Track regions and title for changes
		const _ = editor.regions;
		const _t = cardTitle;
		notifyChange();
	});

	// ============================================================================
	// Event Handlers
	// ============================================================================

	function handleRegionCreate(displayShape: RectShape) {
		if (!editor.image) return;

		// Convert display coordinates to image coordinates
		const imageShape = displayToRegion(displayShape, editor.displayContext);

		// Constrain to image bounds
		const constrained = constrainShapeToImageBounds(
			imageShape,
			editor.image.width,
			editor.image.height
		);

		const region: Region = {
			id: generateRegionId(),
			shape: constrained,
			label: ''
		};

		const command = new CreateRegionCommand(editor, region);
		history.execute(command);

		// Switch to select tool after drawing
		editor.setActiveTool('select');
	}

	function handleContainerResize(width: number, height: number) {
		editor.setContainerSize(width, height);
	}

	function handleSelectRegion(id: string | null) {
		editor.setSelectedRegionId(id);
		// Clear focus signal when selection changes via normal click
		focusLabelRegionId = null;
	}

	function handleDblClickRegion(id: string) {
		editor.setSelectedRegionId(id);
		focusLabelRegionId = id;
	}

	function handleUpdateRegion(
		id: string,
		updates: Partial<Pick<Region, 'label' | 'hint' | 'backContent'>>
	) {
		const region = editor.regions.find((r) => r.id === id);
		if (!region) return;

		const originalValues: Partial<Pick<Region, 'label' | 'hint' | 'backContent'>> = {};
		if ('label' in updates) originalValues.label = region.label;
		if ('hint' in updates) originalValues.hint = region.hint;
		if ('backContent' in updates) originalValues.backContent = region.backContent;

		const command = new UpdateRegionMetadataCommand(editor, id, originalValues, updates);
		history.execute(command);
	}

	function handleDeleteRegion(id: string) {
		deletingRegionId = id;
		deleteConfirmOpen = true;
	}

	function confirmDelete() {
		if (!deletingRegionId) return;
		const region = editor.regions.find((r) => r.id === deletingRegionId);
		if (region) {
			const command = new DeleteRegionCommand(editor, region, editor.selectedRegionId);
			history.execute(command);
		}
		deleteConfirmOpen = false;
		deletingRegionId = null;
	}

	function cancelDelete() {
		deleteConfirmOpen = false;
		deletingRegionId = null;
	}

	function handleRegionMove(id: string, newShape: RectShape) {
		// Live visual update during drag
		editor._updateRegion(id, { shape: newShape });
	}

	function handleRegionMoveEnd(id: string, originalShape: RectShape, newPosition: Point) {
		// Undo live mutation
		editor._updateRegion(id, { shape: originalShape });
		// Execute through command for undo support
		const command = new MoveRegionCommand(editor, id, originalShape, newPosition);
		history.execute(command);
	}

	function handleRegionResize(id: string, newShape: RectShape) {
		// Live visual update during drag
		editor._updateRegion(id, { shape: newShape });
	}

	function handleRegionResizeEnd(id: string, originalShape: RectShape, newShape: RectShape) {
		// Undo live mutation
		editor._updateRegion(id, { shape: originalShape });
		// Execute through command for undo support
		const command = new ResizeRegionCommand(editor, id, originalShape, newShape);
		history.execute(command);
	}

	function handleToolChange(tool: typeof editor.activeTool) {
		editor.setActiveTool(tool);
	}

	function handleUndo() {
		history.undo();
	}

	function handleRedo() {
		history.redo();
	}

	function handleRotate() {
		if (!editor.image) return;

		const currentRotation = editor.image.rotation;
		const newRotation = ((currentRotation + 90) % 360) as 0 | 90 | 180 | 270;

		const command = new RotateImageCommand(editor, currentRotation, newRotation);
		history.execute(command);
	}

	function handlePanChange(offset: { x: number; y: number }) {
		editor.setPanOffset(offset);
	}

	function handleZoomChange(zoom: number) {
		editor.setZoom(zoom);
	}

	function handleZoomFit() {
		editor.setZoom(1);
		editor.setPanOffset({ x: 0, y: 0 });
	}

	function handleImageLoad(url: string, width: number, height: number, assetId?: string) {
		editor.initialize({ url, width, height, rotation: 0, assetId }, [], [], 'hide_all_guess_one');
		history.clear();
	}

	// Resizable divider handlers
	function handleDividerMouseDown(e: MouseEvent) {
		e.preventDefault();
		isDragging = true;
	}

	function handleMouseMove(e: MouseEvent) {
		if (!isDragging || !containerRef) return;

		const containerRect = containerRef.getBoundingClientRect();
		const newWidth = containerRect.right - e.clientX;

		// Clamp between 200px and 500px
		panelWidth = Math.max(200, Math.min(500, newWidth));
	}

	function handleMouseUp() {
		isDragging = false;
	}

	// Clipboard paste handler
	function handlePaste(e: ClipboardEvent) {
		if (editor.hasImage) return;
		const items = e.clipboardData?.items;
		if (!items) return;
		for (const item of items) {
			if (item.type.startsWith('image/')) {
				e.preventDefault();
				const file = item.getAsFile();
				if (file) handlePastedFile(file);
				return;
			}
		}
	}

	async function handlePastedFile(file: File) {
		try {
			const asset = await uploadAsset(file);
			const url = assetUrl(asset.id);
			const width = asset.metadata?.width ?? 0;
			const height = asset.metadata?.height ?? 0;
			handleImageLoad(url, width, height, asset.id);
		} catch (err) {
			console.error('Paste upload failed:', err);
		}
	}

	// Keyboard shortcuts
	function handleKeydown(e: KeyboardEvent) {
		// Ignore if typing in an input
		if (
			(e.target as HTMLElement).tagName === 'INPUT' ||
			(e.target as HTMLElement).tagName === 'TEXTAREA'
		) {
			return;
		}

		// Undo: Cmd/Ctrl + Z
		if ((e.metaKey || e.ctrlKey) && e.key === 'z' && !e.shiftKey) {
			e.preventDefault();
			handleUndo();
			return;
		}

		// Redo: Cmd/Ctrl + Shift + Z
		if ((e.metaKey || e.ctrlKey) && e.shiftKey && e.key === 'z') {
			e.preventDefault();
			handleRedo();
			return;
		}

		// Tool shortcuts
		if (!e.metaKey && !e.ctrlKey && !e.altKey) {
			switch (e.key.toLowerCase()) {
				case 'v':
					editor.setActiveTool('select');
					break;
				case 'r':
					editor.setActiveTool('draw_region');
					break;
				case 'escape':
					editor.setSelectedRegionId(null);
					break;
				case 'delete':
				case 'backspace':
					if (editor.selectedRegionId) {
						handleDeleteRegion(editor.selectedRegionId);
					}
					break;
			}

			// Arrow key nudge for selected region
			if (editor.selectedRegionId && ['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight'].includes(e.key)) {
				e.preventDefault();
				const region = editor.regions.find(r => r.id === editor.selectedRegionId);
				if (!region || !editor.image) return;

				const step = e.shiftKey ? 10 : 1;
				let dx = 0, dy = 0;
				switch (e.key) {
					case 'ArrowUp': dy = -step; break;
					case 'ArrowDown': dy = step; break;
					case 'ArrowLeft': dx = -step; break;
					case 'ArrowRight': dx = step; break;
				}

				const newX = Math.max(0, Math.min(editor.image.width - region.shape.width, region.shape.x + dx));
				const newY = Math.max(0, Math.min(editor.image.height - region.shape.height, region.shape.y + dy));

				if (newX !== region.shape.x || newY !== region.shape.y) {
					const command = new MoveRegionCommand(
						editor,
						region.id,
						region.shape,
						{ x: newX, y: newY }
					);
					history.execute(command);
				}
			}
		}
	}
</script>

<svelte:window
	onkeydown={handleKeydown}
	onpaste={handlePaste}
	onmousemove={isDragging ? handleMouseMove : undefined}
	onmouseup={isDragging ? handleMouseUp : undefined}
/>

<div
	bind:this={containerRef}
	class="border-border bg-card flex h-full w-full overflow-hidden rounded-lg border"
	class:select-none={isDragging}
>
	<!-- Canvas area with floating toolbar -->
	<div class="relative flex-1 overflow-hidden">
		{#if editor.hasImage}
			<!-- Floating toolbar overlay -->
			<div
				data-debug="image occlusion editor toolbar wrapper"
				class="pointer-events-none absolute inset-x-0 top-0 z-10 flex {toolbarJustify[
					toolbarPosition
				]} p-2"
			>
				<div class="pointer-events-auto">
					<EditorToolbar
						activeTool={editor.activeTool}
						canUndo={history.canUndo}
						canRedo={history.canRedo}
						zoom={editor.zoom}
						position={toolbarPosition}
						onToolChange={handleToolChange}
						onUndo={handleUndo}
						onRedo={handleRedo}
						onRotate={handleRotate}
						onZoomChange={handleZoomChange}
						onZoomFit={handleZoomFit}
						onPositionChange={(pos) => (toolbarPosition = pos)}
					/>
				</div>
			</div>
			<div data-debug="image occlusion editor canvas" class="h-full w-full">
				<EditorCanvas
					image={editor.image}
					regions={canvasRegions}
					selectedRegionId={editor.selectedRegionId}
					displayContext={editor.displayContext}
					activeTool={editor.activeTool}
					onSelectRegion={handleSelectRegion}
					onDblClickRegion={handleDblClickRegion}
					onPanChange={handlePanChange}
					onZoomChange={handleZoomChange}
					onRegionCreate={handleRegionCreate}
					onRegionMove={handleRegionMove}
					onRegionMoveEnd={handleRegionMoveEnd}
					onRegionResize={handleRegionResize}
					onRegionResizeEnd={handleRegionResizeEnd}
					onContainerResize={handleContainerResize}
				/>
			</div>
			<!-- Floating status bar -->
			<div class="pointer-events-none absolute inset-x-0 bottom-0 z-10 flex justify-center p-2">
				<div
					class="border-border bg-card/90 pointer-events-auto inline-flex rounded-md border shadow-sm backdrop-blur-sm"
				>
					<StatusBar regionCount={editor.regions.length} mode={editor.mode} zoom={editor.zoom} />
				</div>
			</div>
		{:else}
			<ImageUploader onImageLoad={handleImageLoad} />
		{/if}
	</div>

	<!-- Resizable divider -->
	<div
		class="group bg-border hover:bg-primary/50 flex w-1 cursor-col-resize items-center justify-center transition-colors"
		class:bg-primary={isDragging}
		onmousedown={handleDividerMouseDown}
		role="separator"
		aria-orientation="vertical"
		tabindex="0"
	>
		<div
			class="bg-muted-foreground/30 group-hover:bg-primary h-8 w-0.5 rounded-full transition-colors"
			class:bg-primary={isDragging}
		></div>
	</div>

	<!-- Label Panel with dynamic width -->
	<div class="overflow-hidden" style="width: {panelWidth}px; min-width: {panelWidth}px;">
		<LabelPanel
			regions={editor.regions}
			selectedRegionId={editor.selectedRegionId}
			{focusLabelRegionId}
			title={cardTitle}
			titleError={errors?.title ?? false}
			regionLabelErrors={errors?.regionLabels}
			onSelectRegion={handleSelectRegion}
			onUpdateRegion={handleUpdateRegion}
			onDeleteRegion={handleDeleteRegion}
			onFilterChange={(ids) => (visibleRegionIds = ids)}
			onTitleChange={(v) => (cardTitle = v)}
		/>
	</div>
</div>

<ConfirmDialog
	bind:open={deleteConfirmOpen}
	title="Delete Region"
	description={deletingRegionLabel
		? `Delete region "${deletingRegionLabel}"? You can undo this with Ctrl+Z.`
		: 'Delete this region? You can undo this with Ctrl+Z.'}
	confirmLabel="Delete"
	onconfirm={confirmDelete}
	oncancel={cancelDelete}
/>
