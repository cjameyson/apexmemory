<script lang="ts">
	import type { ImageOcclusionCardDisplay } from '$lib/types/review';
	import type { Region } from '$lib/components/image-occlusion/types';
	import { EyeIcon, LightbulbIcon } from '@lucide/svelte';
	import {
		displaySettings,
		ACTIVE_MASK_PRESETS,
		INACTIVE_MASK_PRESETS,
		REVEAL_PRESETS,
		ICON_PRESETS,
		DEFAULTS
	} from '$lib/stores/display-settings.svelte';

	interface Props {
		display: ImageOcclusionCardDisplay;
		isRevealed: boolean;
		onReveal: () => void;
		onHintUsed?: () => void;
	}

	let { display, isRevealed, onReveal, onHintUsed }: Props = $props();

	let showHint = $state(false);

	// Reset hint when card changes
	$effect(() => {
		display.targetRegionId;
		showHint = false;
	});

	let targetRegion = $derived(
		display.regions.find((r) => r.id === display.targetRegionId)
	);

	let hasHint = $derived(!!targetRegion?.hint);

	let ioSettings = $derived(displaySettings.imageOcclusion);

	/**
	 * Determine which regions to show as masked overlays.
	 * - hide_all_guess_one: all regions masked
	 * - hide_one_guess_one: only target region masked
	 */
	let visibleRegions = $derived.by((): Array<{ region: Region; isTarget: boolean }> => {
		if (display.mode === 'hide_all_guess_one') {
			return display.regions.map((r) => ({
				region: r,
				isTarget: r.id === display.targetRegionId
			}));
		}
		// hide_one_guess_one: only the target
		const target = display.regions.find((r) => r.id === display.targetRegionId);
		return target ? [{ region: target, isTarget: true }] : [];
	});

	/**
	 * Convert region shape (absolute image coords) to percentage-based CSS.
	 */
	function regionStyle(region: Region): string {
		const { x, y, width, height } = region.shape;
		const pctX = (x / display.imageWidth) * 100;
		const pctY = (y / display.imageHeight) * 100;
		const pctW = (width / display.imageWidth) * 100;
		const pctH = (height / display.imageHeight) * 100;
		return `left:${pctX}%;top:${pctY}%;width:${pctW}%;height:${pctH}%`;
	}

	/**
	 * Determine if a region's mask should be shown.
	 * After reveal, only non-target regions stay masked.
	 */
	function shouldShowMask(isTarget: boolean): boolean {
		if (!isRevealed) return true;
		return !isTarget;
	}

	/**
	 * Get the CSS classes for a region overlay based on current state.
	 */
	function getMaskClasses(isTarget: boolean): string {
		const masked = shouldShowMask(isTarget);

		if (masked) {
			if (isTarget) {
				const preset = ACTIVE_MASK_PRESETS[ioSettings.active_mask_color]
					?? ACTIVE_MASK_PRESETS[DEFAULTS.image_occlusion.active_mask_color];
				return `${preset.bg} ${preset.border}`;
			}
			const preset = INACTIVE_MASK_PRESETS[ioSettings.inactive_mask_color]
				?? INACTIVE_MASK_PRESETS[DEFAULTS.image_occlusion.inactive_mask_color];
			return `${preset.bg} ${preset.border}`;
		}

		if (isTarget && isRevealed) {
			const preset = REVEAL_PRESETS[ioSettings.reveal_color]
				?? REVEAL_PRESETS[DEFAULTS.image_occlusion.reveal_color];
			return `${preset.bg} ${preset.border}`;
		}

		return '';
	}

	/**
	 * Get animation classes for the target mask.
	 */
	function getAnimationClasses(isTarget: boolean): string {
		if (!isTarget || isRevealed) return '';
		return ioSettings.marching_ants ? 'occlusion-ants-slow' : '';
	}

	/**
	 * Whether to show the indicator icon on the target mask.
	 */
	function showIndicator(isTarget: boolean): boolean {
		if (!isTarget || isRevealed) return false;
		return ioSettings.icon !== 'none';
	}

	/**
	 * Get rotation CSS transform.
	 */
	let rotationStyle = $derived(
		display.imageRotation ? `transform:rotate(${display.imageRotation}deg)` : ''
	);

	let revealPreset = $derived(
		REVEAL_PRESETS[ioSettings.reveal_color]
			?? REVEAL_PRESETS[DEFAULTS.image_occlusion.reveal_color]
	);

	let IconComponent = $derived(ICON_PRESETS[ioSettings.icon] ?? null);
</script>

<div class="flex flex-col items-center gap-4 w-full">
	<!-- Title -->
	<h3 class="text-lg font-medium text-white/80">{display.title}</h3>

	<!-- Image container with overlays -->
	<div class="relative inline-block max-w-full" style="max-height:70vh;">
		<img
			src={display.imageUrl}
			alt={display.title}
			class="block max-w-full max-h-[70vh] object-contain rounded-lg"
			style={rotationStyle}
			draggable="false"
		/>

		<!-- Region overlays — single persistent div per region for smooth transitions -->
		{#each visibleRegions as { region, isTarget } (region.id)}
			<div
				class="absolute rounded-sm"
				style={regionStyle(region)}
			>
				<div
					class="w-full h-full rounded-sm transition-all duration-300 {isTarget && isRevealed && display.revealStyle === 'show_label' ? 'rounded-tl-none' : ''} {getMaskClasses(isTarget)} {getAnimationClasses(isTarget)}"
				>
					{#if showIndicator(isTarget) && IconComponent}
						<span class="absolute inset-0 flex items-center justify-center text-white/70">
							<IconComponent class="size-5 occlusion-icon-pulse" />
						</span>
					{/if}
				</div>

				<!-- PAWLS-style annotation label: above region, left-aligned, solid bg matches border -->
				{#if isTarget && isRevealed && display.revealStyle === 'show_label'}
					<span class="absolute bottom-full left-0 {revealPreset.label} rounded-t-[2px] px-1.5 py-px text-[10px] font-bold {revealPreset.labelText} whitespace-nowrap leading-tight z-10">
						{targetRegion?.label ?? ''}
					</span>
				{/if}
			</div>
		{/each}
	</div>

	<!-- Hint (pre-reveal only) -->
	{#if !isRevealed && hasHint}
		{#if showHint}
			<p class="text-sm text-amber-300/80 italic flex items-center gap-1.5">
				<LightbulbIcon class="size-4 shrink-0" />
				<span class="font-medium not-italic text-amber-300">Hint:</span> {targetRegion?.hint}
			</p>
		{:else}
			<button
				type="button"
				onclick={() => { showHint = true; onHintUsed?.(); }}
				class="flex items-center gap-1.5 text-sm text-white/50 hover:text-white/70 transition-colors"
			>
				<LightbulbIcon class="size-4" />
				<span>Show hint</span>
			</button>
		{/if}
	{/if}

	<!-- Reveal prompt (pre-reveal only) -->
	{#if !isRevealed}
		<button
			type="button"
			onclick={onReveal}
			class="flex items-center gap-2 text-white/50 hover:text-white/70 transition-colors"
		>
			<EyeIcon class="size-5" />
			<span>Tap or press space to reveal</span>
		</button>
	{/if}
</div>
