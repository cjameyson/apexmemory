const MAX_WIDTH = 2000;
const MAX_HEIGHT = 2000;

/**
 * Resize an image file if it exceeds MAX_WIDTH or MAX_HEIGHT.
 * Non-image files and GIFs (animated) are returned as-is.
 * Returns a Promise that resolves to either the resized file or original if resizing wasn't needed.
 */
export async function resizeImageIfNeeded(file: File): Promise<File> {
	// Don't resize GIFs (animated) or non-images
	if (!file.type.startsWith('image/') || file.type === 'image/gif') {
		return file;
	}

	return new Promise((resolve) => {
		const img = new window.Image();
		const url = URL.createObjectURL(file);

		img.onload = () => {
			URL.revokeObjectURL(url);

			if (img.width <= MAX_WIDTH && img.height <= MAX_HEIGHT) {
				resolve(file);
				return;
			}

			const ratio = Math.min(MAX_WIDTH / img.width, MAX_HEIGHT / img.height);
			const width = Math.round(img.width * ratio);
			const height = Math.round(img.height * ratio);

			const canvas = document.createElement('canvas');
			canvas.width = width;
			canvas.height = height;

			const ctx = canvas.getContext('2d');
			if (!ctx) {
				resolve(file);
				return;
			}

			ctx.drawImage(img, 0, 0, width, height);

			canvas.toBlob(
				(blob) => {
					if (!blob) {
						resolve(file);
						return;
					}
					resolve(new File([blob], file.name, { type: file.type }));
				},
				file.type,
				0.9
			);
		};

		img.onerror = () => {
			URL.revokeObjectURL(url);
			resolve(file);
		};

		img.src = url;
	});
}
