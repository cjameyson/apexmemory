/**
 * Utility functions for the image occlusion editor.
 */

import { nanoid } from 'nanoid';

/**
 * Generate a unique region ID.
 * Format: m_<nanoid(10)> where 'm' stands for 'mask'
 */
export function generateRegionId(): string {
	return `m_${nanoid(10)}`;
}

/**
 * Generate a unique annotation ID.
 * Format: a_<nanoid(10)> where 'a' stands for 'annotation'
 */
export function generateAnnotationId(): string {
	return `a_${nanoid(10)}`;
}

/**
 * Create a debounced version of a function.
 * Useful for continuous operations like dragging.
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function debounce<T extends (...args: any[]) => void>(
	fn: T,
	delay: number
): (...args: Parameters<T>) => void {
	let timeoutId: ReturnType<typeof setTimeout> | null = null;

	return (...args: Parameters<T>) => {
		if (timeoutId !== null) {
			clearTimeout(timeoutId);
		}
		timeoutId = setTimeout(() => {
			fn(...args);
			timeoutId = null;
		}, delay);
	};
}

/**
 * Create a throttled version of a function.
 * Useful for high-frequency events like mousemove.
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function throttle<T extends (...args: any[]) => void>(
	fn: T,
	limit: number
): (...args: Parameters<T>) => void {
	let lastRun = 0;
	let pendingArgs: Parameters<T> | null = null;
	let timeoutId: ReturnType<typeof setTimeout> | null = null;

	return (...args: Parameters<T>) => {
		const now = Date.now();

		if (now - lastRun >= limit) {
			fn(...args);
			lastRun = now;
		} else {
			pendingArgs = args;
			if (timeoutId === null) {
				timeoutId = setTimeout(() => {
					if (pendingArgs !== null) {
						fn(...(pendingArgs as Parameters<T>));
						lastRun = Date.now();
						pendingArgs = null;
					}
					timeoutId = null;
				}, limit - (now - lastRun));
			}
		}
	};
}

/**
 * Clamp a number between min and max values.
 */
export function clamp(value: number, min: number, max: number): number {
	return Math.max(min, Math.min(max, value));
}

/**
 * Format zoom level as percentage string.
 */
export function formatZoom(zoom: number): string {
	return `${Math.round(zoom * 100)}%`;
}

/**
 * Format occlusion mode as human-readable string.
 */
export function formatOcclusionMode(mode: 'hide_all_guess_one' | 'hide_one_guess_one'): string {
	switch (mode) {
		case 'hide_all_guess_one':
			return 'Hide All, Guess One';
		case 'hide_one_guess_one':
			return 'Hide One, Guess One';
	}
}
