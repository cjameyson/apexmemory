<script lang="ts">
	import type { ImageOcclusionCardDisplay } from '$lib/types/review';
	import type { Region } from '$lib/components/image-occlusion/types';
	import { EyeIcon, LightbulbIcon, CircleHelpIcon } from '@lucide/svelte';
	import { occlusionDebug, MASK_COLORS, type MaskColor } from '$lib/stores/debug-occlusion.svelte';

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

	let debug = $derived(occlusionDebug.state);
	let useDebug = $derived(debug.enabled);

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
	 * Get the CSS classes for a region overlay based on current state and debug settings.
	 */
	function getMaskClasses(isTarget: boolean): string {
		const masked = shouldShowMask(isTarget);
		const color: MaskColor = useDebug ? debug.maskColor : 'primary';
		const colors = MASK_COLORS[color];

		if (masked) {
			if (isTarget) {
				return `${colors.bg} border-2 ${colors.border}`;
			}
			return 'bg-slate-600 border border-slate-500';
		}

		if (isTarget && isRevealed) {
			return `${colors.bgRevealed} border-2 ${colors.borderRevealed}`;
		}

		return '';
	}

	/**
	 * Get animation classes for the target mask.
	 * Defaults: marching ants slow. Debug overrides all.
	 */
	function getAnimationClasses(isTarget: boolean): string {
		if (!isTarget || isRevealed) return '';

		const classes: string[] = [];

		const effectivePulse = useDebug ? debug.pulse : 'off';
		const effectiveAnts = useDebug ? debug.marchingAnts : 'slow';

		if (effectivePulse === 'subtle') classes.push('occlusion-pulse-subtle');
		else if (effectivePulse === 'pronounced') classes.push('occlusion-pulse-pronounced');

		if (effectiveAnts === 'slow') classes.push('occlusion-ants-slow');
		else if (effectiveAnts === 'medium') classes.push('occlusion-ants-medium');

		return classes.join(' ');
	}

	/**
	 * Get reveal transition classes for the persistent div.
	 */
	function getRevealTransitionClasses(isTarget: boolean): string {
		if (!isTarget || !useDebug) return 'transition-all duration-300';

		if (debug.reveal === 'fade-cross') return 'occlusion-reveal-fade-cross';
		if (debug.reveal === 'dissolve') return 'occlusion-reveal-dissolve';
		if (debug.reveal === 'slide-away') {
			const base = 'occlusion-reveal-slide-away';
			return isRevealed ? `${base} occlusion-slide-away-active` : base;
		}
		return 'transition-all duration-300';
	}

	/**
	 * Get entrance animation class.
	 */
	function getEntranceClass(): string {
		if (!useDebug) return '';
		if (debug.entrance === 'fade-in') return 'occlusion-enter-fade-in';
		if (debug.entrance === 'scale-up') return 'occlusion-enter-scale-up';
		return '';
	}

	/**
	 * Whether to show the indicator icon on the target mask.
	 */
	function showIndicator(isTarget: boolean): boolean {
		if (!isTarget || isRevealed) return false;
		if (useDebug) return debug.showIndicator;
		return true;
	}

	/**
	 * Get rotation CSS transform.
	 */
	let rotationStyle = $derived(
		display.imageRotation ? `transform:rotate(${display.imageRotation}deg)` : ''
	);
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
				class="absolute rounded-sm {getEntranceClass()}"
				style={regionStyle(region)}
			>
				<div
					class="w-full h-full rounded-sm {isTarget && isRevealed && display.revealStyle === 'show_label' ? 'rounded-tl-none' : ''} {getMaskClasses(isTarget)} {getAnimationClasses(isTarget)} {getRevealTransitionClasses(isTarget)}"
				>
					{#if showIndicator(isTarget)}
						<span class="absolute inset-0 flex items-center justify-center text-white/70">
							<CircleHelpIcon class="size-5" />
						</span>
					{/if}
				</div>

				<!-- PAWLS-style annotation label: above region, left-aligned, solid bg matches border -->
				{#if isTarget && isRevealed && display.revealStyle === 'show_label'}
					<span class="absolute bottom-full left-0 bg-emerald-400 rounded-t-[2px] px-1.5 py-px text-[10px] font-bold text-emerald-950 whitespace-nowrap leading-tight z-10">
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
