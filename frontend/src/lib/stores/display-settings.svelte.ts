/**
 * Display settings store for card rendering customization.
 * Currently uses hardcoded defaults; will read from user/notebook settings
 * once the backend integration is built.
 */

import { CircleQuestionMarkIcon, EyeIcon } from '@lucide/svelte';
import type { Component } from 'svelte';

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

export interface ImageOcclusionDisplaySettings {
	active_mask_color: string;
	inactive_mask_color: string;
	marching_ants: boolean;
	icon: string;
	reveal_color: string;
}

export interface DisplaySettings {
	image_occlusion: ImageOcclusionDisplaySettings;
}

// ---------------------------------------------------------------------------
// Defaults
// ---------------------------------------------------------------------------

export const DEFAULTS: DisplaySettings = {
	image_occlusion: {
		active_mask_color: 'primary',
		inactive_mask_color: 'slate',
		marching_ants: true,
		icon: 'circle-help',
		reveal_color: 'success'
	}
};

// ---------------------------------------------------------------------------
// Preset Lookup Maps
// ---------------------------------------------------------------------------

export const ACTIVE_MASK_PRESETS: Record<string, { bg: string; border: string }> = {
	primary: { bg: 'bg-primary', border: 'border-2 border-primary/50' },
	destructive: { bg: 'bg-destructive', border: 'border-2 border-destructive/50' },
	cloze: { bg: 'bg-cloze', border: 'border-2 border-cloze/50' },
	warning: { bg: 'bg-warning', border: 'border-2 border-warning/50' }
};

export const INACTIVE_MASK_PRESETS: Record<string, { bg: string; border: string }> = {
	slate: { bg: 'bg-slate-600', border: 'border border-slate-500' },
	muted: { bg: 'bg-muted', border: 'border border-muted-foreground/30' }
};

export const REVEAL_PRESETS: Record<string, { bg: string; border: string; label: string; labelText: string }> = {
	success: { bg: 'bg-success/20', border: 'border-2 border-success', label: 'bg-success', labelText: 'text-success-foreground' },
	good: { bg: 'bg-good/20', border: 'border-2 border-good', label: 'bg-good', labelText: 'text-good-foreground' }
};

export const ICON_PRESETS: Record<string, Component | null> = {
	'circle-help': CircleQuestionMarkIcon,
	'eye': EyeIcon,
	'none': null
};

// ---------------------------------------------------------------------------
// Store
// ---------------------------------------------------------------------------

function createDisplaySettingsStore() {
	let settings = $state<DisplaySettings>(structuredClone(DEFAULTS));

	/**
	 * Initialize settings by deep-merging server data over defaults.
	 * No-op until the backend provides user display_settings.
	 */
	function init(serverSettings?: Partial<DisplaySettings>) {
		if (!serverSettings) return;
		settings = {
			...DEFAULTS,
			...serverSettings,
			image_occlusion: {
				...DEFAULTS.image_occlusion,
				...serverSettings.image_occlusion
			}
		};
	}

	return {
		get current() {
			return settings;
		},
		get imageOcclusion() {
			return settings.image_occlusion;
		},
		init
	};
}

export const displaySettings = createDisplaySettingsStore();
