// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces

import type { User } from '$lib/api/types';

declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: User | null;
			sessionToken: string | null;
		}
		interface PageData {
			user: User | null;
		}
		interface PageState {
			focusMode?: {
				type: 'all' | 'notebook' | 'source';
				mode?: 'scheduled' | 'practice';
				notebookId?: string;
				notebookName?: string;
				notebookEmoji?: string;
				sourceId?: string;
				sourceName?: string;
				currentIndex?: number;
			};
			commandPalette?: boolean;
		}
		// interface Platform {}
	}
}

export {};
