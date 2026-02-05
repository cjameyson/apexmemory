import type { ApiAsset } from './types';
import { resizeImageIfNeeded } from '$lib/utils/image-resize';

/**
 * Upload a file as an asset via the SvelteKit proxy route.
 * Use this from client-side components (e.g., image upload in TipTap editor).
 * Images larger than 2000px are automatically resized before upload.
 */
export async function uploadAsset(file: File): Promise<ApiAsset> {
	const resized = await resizeImageIfNeeded(file);

	const formData = new FormData();
	formData.append('file', resized);

	const res = await fetch('/api/assets', {
		method: 'POST',
		body: formData
	});

	if (!res.ok) {
		const err = await res.json().catch(() => ({ message: 'Upload failed' }));
		throw new Error(err.message ?? 'Upload failed');
	}

	return res.json();
}

/**
 * Get the URL to serve an asset file via the SvelteKit proxy.
 * Use this to construct <img src> or other resource URLs.
 */
export function assetUrl(assetId: string): string {
	return `/api/assets/${assetId}/file`;
}
