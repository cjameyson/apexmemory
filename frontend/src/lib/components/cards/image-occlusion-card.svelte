<script lang="ts">
	import type { ImageOcclusionCardDisplay } from '$lib/types/review';
	import type { Region } from '$lib/components/image-occlusion/types';
	import { EyeIcon, LightbulbIcon } from '@lucide/svelte';

	interface Props {
		display: ImageOcclusionCardDisplay;
		isRevealed: boolean;
		onReveal: () => void;
	}

	let { display, isRevealed, onReveal }: Props = $props();

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
	<div class="relative inline-block max-w-full" style="max-height:60vh;">
		<img
			src={display.imageUrl}
			alt={display.title}
			class="block max-w-full max-h-[60vh] object-contain rounded-lg"
			style={rotationStyle}
			draggable="false"
		/>

		<!-- Region overlays -->
		{#each visibleRegions as { region, isTarget } (region.id)}
			<div
				class="absolute transition-all duration-300 rounded-sm"
				style={regionStyle(region)}
			>
				{#if shouldShowMask(isTarget)}
					<!-- Masked state -->
					<div
						class="w-full h-full rounded-sm {isTarget
							? 'bg-sky-500 border-2 border-sky-300'
							: 'bg-slate-600 border border-slate-500'}"
					>
						{#if isTarget && display.mode === 'hide_all_guess_one'}
							<span class="absolute inset-0 flex items-center justify-center text-white font-bold text-lg">
								?
							</span>
						{/if}
					</div>
				{:else if isTarget && isRevealed}
					<!-- Revealed target -->
					<div class="w-full h-full border-2 border-emerald-400 rounded-sm bg-emerald-500/20"></div>
					{#if display.revealStyle === 'show_label'}
						<!-- Label floats below region, not clipped by it -->
						<div class="absolute left-1/2 -translate-x-1/2 top-full mt-1 whitespace-nowrap z-10">
							<span class="px-2 py-0.5 text-sm font-semibold text-emerald-300 bg-black/70 rounded text-center leading-tight">
								{targetRegion?.label ?? ''}
							</span>
						</div>
					{/if}
				{/if}
			</div>
		{/each}
	</div>

	<!-- Hint button (pre-reveal only) -->
	{#if !isRevealed && hasHint}
		<button
			type="button"
			onclick={() => (showHint = !showHint)}
			class="flex items-center gap-1.5 text-sm text-white/50 hover:text-white/70 transition-colors"
		>
			<LightbulbIcon class="size-4" />
			<span>{showHint ? 'Hide hint' : 'Show hint'}</span>
		</button>
		{#if showHint}
			<p class="text-sm text-amber-300/80 italic">{targetRegion?.hint}</p>
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
