<script lang="ts">
	import { SettingsIcon, XIcon, RotateCcwIcon } from '@lucide/svelte';
	import {
		occlusionDebug,
		type MaskColor,
		type PulseMode,
		type MarchingAnts,
		type EntranceMode,
		type RevealMode
	} from '$lib/stores/debug-occlusion.svelte';

	interface Props {
		visible: boolean;
		onToggle: () => void;
	}

	let { visible, onToggle }: Props = $props();

	let debug = $derived(occlusionDebug.state);

	const COLOR_OPTIONS: MaskColor[] = ['primary', 'blue', 'amber', 'violet', 'rose'];
	const PULSE_OPTIONS: PulseMode[] = ['off', 'subtle', 'pronounced'];
	const ANTS_OPTIONS: MarchingAnts[] = ['off', 'slow', 'medium'];
	const ENTRANCE_OPTIONS: EntranceMode[] = ['none', 'fade-in', 'scale-up'];
	const REVEAL_OPTIONS: RevealMode[] = ['none', 'fade-cross', 'dissolve', 'slide-away'];

	const COLOR_DOTS: Record<MaskColor, string> = {
		primary: 'bg-primary',
		blue: 'bg-sky-500',
		amber: 'bg-amber-500',
		violet: 'bg-violet-500',
		rose: 'bg-rose-500'
	};

	function cycleOption<T>(options: T[], current: T): T {
		const idx = options.indexOf(current);
		return options[(idx + 1) % options.length];
	}
</script>

{#if !visible}
	<!-- Collapsed: gear icon button -->
	<button
		type="button"
		onclick={onToggle}
		class="fixed top-20 right-4 z-50 p-2 rounded-lg bg-black/60 hover:bg-black/80 text-white/60 hover:text-white transition-colors"
		aria-label="Open occlusion debug panel"
	>
		<SettingsIcon class="size-5" />
	</button>
{:else}
	<!-- Expanded panel -->
	<div class="fixed top-20 right-4 z-50 w-60 bg-black/80 backdrop-blur-sm rounded-xl border border-white/10 shadow-2xl text-sm">
		<!-- Header -->
		<div class="flex items-center justify-between px-3 py-2 border-b border-white/10">
			<div class="flex items-center gap-2">
				<SettingsIcon class="size-3.5 text-white/50" />
				<span class="text-white/80 font-medium text-xs uppercase tracking-wider">Occlusion Debug</span>
			</div>
			<button
				type="button"
				onclick={onToggle}
				class="p-1 rounded text-white/40 hover:text-white/70 transition-colors"
				aria-label="Close debug panel"
			>
				<XIcon class="size-3.5" />
			</button>
		</div>

		<!-- Master toggle -->
		<div class="flex items-center justify-between px-3 py-2 border-b border-white/10">
			<span class="text-white/70">Enabled</span>
			<button
				type="button"
				onclick={() => occlusionDebug.toggle()}
				class="px-2.5 py-0.5 rounded text-xs font-medium transition-colors {debug.enabled
					? 'bg-emerald-500/30 text-emerald-300'
					: 'bg-white/10 text-white/40'}"
			>
				{debug.enabled ? 'on' : 'off'}
			</button>
		</div>

		{#if debug.enabled}
			<div class="px-3 py-2 space-y-2">
				<!-- Color -->
				<div class="flex items-center justify-between">
					<span class="text-white/60">Color</span>
					<button
						type="button"
						onclick={() => occlusionDebug.set('maskColor', cycleOption(COLOR_OPTIONS, debug.maskColor))}
						class="flex items-center gap-1.5 px-2.5 py-0.5 rounded bg-white/10 hover:bg-white/15 text-white/80 text-xs font-medium transition-colors"
					>
						<span class="size-2.5 rounded-full {COLOR_DOTS[debug.maskColor]}"></span>
						{debug.maskColor}
					</button>
				</div>

				<!-- Indicator Icon -->
				<div class="flex items-center justify-between">
					<span class="text-white/60">Icon</span>
					<button
						type="button"
						onclick={() => occlusionDebug.set('showIndicator', !debug.showIndicator)}
						class="px-2.5 py-0.5 rounded bg-white/10 hover:bg-white/15 text-white/80 text-xs font-medium transition-colors"
					>
						{debug.showIndicator ? 'on' : 'off'}
					</button>
				</div>

				<!-- Pulse -->
				<div class="flex items-center justify-between">
					<span class="text-white/60">Pulse</span>
					<button
						type="button"
						onclick={() => occlusionDebug.set('pulse', cycleOption(PULSE_OPTIONS, debug.pulse))}
						class="px-2.5 py-0.5 rounded bg-white/10 hover:bg-white/15 text-white/80 text-xs font-medium transition-colors"
					>
						{debug.pulse}
					</button>
				</div>

				<!-- Marching Ants -->
				<div class="flex items-center justify-between">
					<span class="text-white/60">Ants</span>
					<button
						type="button"
						onclick={() => occlusionDebug.set('marchingAnts', cycleOption(ANTS_OPTIONS, debug.marchingAnts))}
						class="px-2.5 py-0.5 rounded bg-white/10 hover:bg-white/15 text-white/80 text-xs font-medium transition-colors"
					>
						{debug.marchingAnts}
					</button>
				</div>

				<!-- Entrance -->
				<div class="flex items-center justify-between">
					<span class="text-white/60">Entrance</span>
					<button
						type="button"
						onclick={() => occlusionDebug.set('entrance', cycleOption(ENTRANCE_OPTIONS, debug.entrance))}
						class="px-2.5 py-0.5 rounded bg-white/10 hover:bg-white/15 text-white/80 text-xs font-medium transition-colors"
					>
						{debug.entrance}
					</button>
				</div>

				<!-- Reveal -->
				<div class="flex items-center justify-between">
					<span class="text-white/60">Reveal</span>
					<button
						type="button"
						onclick={() => occlusionDebug.set('reveal', cycleOption(REVEAL_OPTIONS, debug.reveal))}
						class="px-2.5 py-0.5 rounded bg-white/10 hover:bg-white/15 text-white/80 text-xs font-medium transition-colors"
					>
						{debug.reveal}
					</button>
				</div>
			</div>

			<!-- Reset -->
			<div class="px-3 py-2 border-t border-white/10">
				<button
					type="button"
					onclick={() => occlusionDebug.reset()}
					class="flex items-center gap-1.5 text-xs text-white/40 hover:text-white/60 transition-colors"
				>
					<RotateCcwIcon class="size-3" />
					Reset defaults
				</button>
			</div>
		{/if}
	</div>
{/if}
