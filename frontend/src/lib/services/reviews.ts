import type { ApiStudyCard } from '$lib/api/types';
import type { StudyCard, CardDisplay, BasicCardDisplay, ClozeCardDisplay, ImageOcclusionCardDisplay } from '$lib/types/review';
import type { ImageOcclusionField } from '$lib/components/image-occlusion/types';
import type { JSONContent } from '@tiptap/core';
import { assetUrl } from '$lib/api/client';

const RATING_MAP: Record<number, 'again' | 'hard' | 'good' | 'easy'> = {
	1: 'again',
	2: 'hard',
	3: 'good',
	4: 'easy'
};

export function ratingToString(rating: 1 | 2 | 3 | 4): 'again' | 'hard' | 'good' | 'easy' {
	return RATING_MAP[rating];
}

export function toStudyCard(api: ApiStudyCard): StudyCard {
	return {
		id: api.id,
		factId: api.fact_id,
		notebookId: api.notebook_id,
		elementId: api.element_id,
		state: api.state,
		stability: api.stability,
		difficulty: api.difficulty,
		due: api.due,
		reps: api.reps,
		lapses: api.lapses,
		factType: api.fact_type,
		factContent: api.fact_content,
		intervals: api.intervals
	};
}

export function toStudyCards(apiList: ApiStudyCard[]): StudyCard[] {
	return apiList.map(toStudyCard);
}

async function readErrorMessage(response: Response, fallback: string): Promise<string> {
	try {
		const payload = await response.json();
		if (payload && typeof payload === 'object') {
			const message = (payload as { message?: unknown }).message;
			if (typeof message === 'string' && message.length > 0) {
				return message;
			}
			const error = (payload as { error?: unknown }).error;
			if (typeof error === 'string' && error.length > 0) {
				return error;
			}
		}
	} catch {
		// Ignore JSON parsing failures and use fallback.
	}
	return fallback;
}

function extractCards(payload: unknown): ApiStudyCard[] {
	if (Array.isArray(payload)) {
		return payload as ApiStudyCard[];
	}
	if (payload && typeof payload === 'object' && Array.isArray((payload as { data?: unknown }).data)) {
		return (payload as { data: ApiStudyCard[] }).data;
	}
	throw new Error('Invalid card response from server');
}

async function fetchCards(endpoint: string, failureMessage: string): Promise<StudyCard[]> {
	const response = await fetch(endpoint);
	if (!response.ok) {
		throw new Error(await readErrorMessage(response, failureMessage));
	}

	const payload = await response.json();
	return toStudyCards(extractCards(payload));
}

export async function fetchStudyCards(notebookId?: string): Promise<StudyCard[]> {
	const params = new URLSearchParams({ limit: '50' });
	if (notebookId) params.set('notebook_id', notebookId);

	return fetchCards(`/api/reviews/study?${params}`, 'Failed to load study cards. Please try again.');
}

export async function fetchPracticeCards(notebookId?: string): Promise<StudyCard[]> {
	const params = new URLSearchParams({ limit: '50' });
	if (notebookId) params.set('notebook_id', notebookId);

	return fetchCards(`/api/reviews/practice?${params}`, 'Failed to load practice cards. Please try again.');
}

/**
 * Extract display data from a study card's fact content.
 * Returns a discriminated union routed by card type.
 */
export function extractCardDisplay(card: StudyCard): CardDisplay {
	const content = card.factContent as {
		version?: number;
		fields?: Array<{ name: string; type: string; value: unknown }>;
	};
	const fields = content?.fields ?? [];

	if (card.factType === 'basic') {
		const frontField = fields.find((f) => f.name === 'front');
		const backField = fields.find((f) => f.name === 'back');
		const backExtraField = fields.find((f) => f.name === 'back_extra');
		const backExtra = backExtraField?.value ? String(backExtraField.value) : undefined;
		return {
			type: 'basic',
			front: frontField?.type === 'rich_text' && typeof frontField.value === 'object'
				? (frontField.value as JSONContent)
				: String(frontField?.value ?? ''),
			back: backField?.type === 'rich_text' && typeof backField.value === 'object'
				? (backField.value as JSONContent)
				: String(backField?.value ?? ''),
			backExtra
		} satisfies BasicCardDisplay;
	}

	if (card.factType === 'cloze') {
		const textField = fields.find((f) => f.type === 'cloze_text') ?? fields[0];
		const text = String(textField?.value ?? '');
		const elementId = card.elementId;

		let clozeAnswer = '';
		text.replace(/\{\{(c\d+)::([^}]+)\}\}/g, (_, id, answer) => {
			if (id === elementId) clozeAnswer = answer;
			return '';
		});

		const front = text.replace(/\{\{(c\d+)::([^}]+)\}\}/g, (_, id, answer) => {
			return id === elementId ? '[...]' : answer;
		});

		const backExtraField = fields.find((f) => f.name === 'back_extra');
		const backExtra = backExtraField?.value ? String(backExtraField.value) : undefined;
		return { type: 'cloze', front, clozeAnswer, backExtra } satisfies ClozeCardDisplay;
	}

	if (card.factType === 'image_occlusion') {
		const ioField = fields.find((f) => f.type === 'image_occlusion');
		const data = ioField?.value as ImageOcclusionField | undefined;

		if (!data?.image || !data?.regions?.length) {
			return { type: 'basic', front: data?.title ?? 'Image Occlusion', back: '' };
		}

		const resolvedUrl = data.image.assetId
			? assetUrl(data.image.assetId)
			: data.image.url;

		const targetRegion = data.regions.find((r) => r.id === card.elementId);
		const backExtra = targetRegion?.backExtra || undefined;

		return {
			type: 'image_occlusion',
			title: data.title || 'Image Occlusion',
			imageUrl: resolvedUrl,
			imageWidth: data.image.width,
			imageHeight: data.image.height,
			imageRotation: data.image.rotation ?? 0,
			regions: data.regions,
			targetRegionId: card.elementId,
			mode: data.mode ?? 'hide_all_guess_one',
			revealStyle: data.revealStyle ?? 'show_label',
			backExtra
		} satisfies ImageOcclusionCardDisplay;
	}

	return { type: 'basic', front: 'Unknown card type', back: '' };
}
