/**
 * Dev-only debug store for experimenting with image occlusion mask effects.
 * Toggled via a floating panel in focus mode (D key).
 * Delete this file once the winning combination is chosen.
 */

export type MaskColor = 'primary' | 'blue' | 'amber' | 'violet' | 'rose';
export type PulseMode = 'off' | 'subtle' | 'pronounced';
export type MarchingAnts = 'off' | 'slow' | 'medium';
export type EntranceMode = 'none' | 'fade-in' | 'scale-up';
export type RevealMode = 'none' | 'fade-cross' | 'dissolve' | 'slide-away';

export interface OcclusionDebugState {
	enabled: boolean;
	maskColor: MaskColor;
	showIndicator: boolean;
	pulse: PulseMode;
	marchingAnts: MarchingAnts;
	entrance: EntranceMode;
	reveal: RevealMode;
}

const DEFAULTS: OcclusionDebugState = {
	enabled: false,
	maskColor: 'primary',
	showIndicator: true,
	pulse: 'off',
	marchingAnts: 'slow',
	entrance: 'none',
	reveal: 'none'
};

const STORAGE_KEY = 'debug-occlusion';

function createOcclusionDebugStore() {
	let state = $state<OcclusionDebugState>({ ...DEFAULTS });

	function init() {
		if (typeof window === 'undefined') return;
		try {
			const stored = localStorage.getItem(STORAGE_KEY);
			if (stored) {
				const parsed = JSON.parse(stored);
				state = { ...DEFAULTS, ...parsed };
			}
		} catch {
			// ignore corrupt storage
		}
	}

	function persist() {
		if (typeof window === 'undefined') return;
		localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
	}

	function set<K extends keyof OcclusionDebugState>(key: K, value: OcclusionDebugState[K]) {
		state[key] = value;
		persist();
	}

	function reset() {
		state = { ...DEFAULTS, enabled: state.enabled };
		persist();
	}

	function toggle() {
		state.enabled = !state.enabled;
		persist();
	}

	return {
		get state() {
			return state;
		},
		init,
		set,
		reset,
		toggle
	};
}

export const occlusionDebug = createOcclusionDebugStore();

// Color maps for mask presets
export const MASK_COLORS: Record<MaskColor, { bg: string; border: string; bgRevealed: string; borderRevealed: string }> = {
	primary: { bg: 'bg-primary', border: 'border-primary/50', bgRevealed: 'bg-emerald-500/20', borderRevealed: 'border-emerald-400' },
	blue: { bg: 'bg-sky-500', border: 'border-sky-300', bgRevealed: 'bg-emerald-500/20', borderRevealed: 'border-emerald-400' },
	amber: { bg: 'bg-amber-500', border: 'border-amber-300', bgRevealed: 'bg-emerald-500/20', borderRevealed: 'border-emerald-400' },
	violet: { bg: 'bg-violet-500', border: 'border-violet-300', bgRevealed: 'bg-emerald-500/20', borderRevealed: 'border-emerald-400' },
	rose: { bg: 'bg-rose-500', border: 'border-rose-300', bgRevealed: 'bg-emerald-500/20', borderRevealed: 'border-emerald-400' }
};
