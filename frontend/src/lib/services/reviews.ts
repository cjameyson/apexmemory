import type { ApiStudyCard } from '$lib/api/types';
import type { StudyCard, CardDisplay, BasicCardDisplay, ClozeCardDisplay, ImageOcclusionCardDisplay } from '$lib/types/review';
import type { ImageOcclusionField } from '$lib/components/image-occlusion/types';
import type { JSONContent } from '@tiptap/core';
import { assetUrl } from '$lib/api/client';
import { toast } from 'svelte-sonner';

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

export async function fetchStudyCards(notebookId?: string): Promise<StudyCard[]> {
	const params = new URLSearchParams({ limit: '50' });
	if (notebookId) params.set('notebook_id', notebookId);

	const res = await fetch(`/api/reviews/study?${params}`);
	if (!res.ok) {
		toast.error('Failed to load study cards. Please try again.');
		return [];
	}
	const data: ApiStudyCard[] = await res.json();
	return toStudyCards(data);
}

export async function fetchPracticeCards(notebookId?: string): Promise<StudyCard[]> {
	const params = new URLSearchParams({ limit: '50' });
	if (notebookId) params.set('notebook_id', notebookId);

	const res = await fetch(`/api/reviews/practice?${params}`);
	if (!res.ok) {
		toast.error('Failed to load practice cards. Please try again.');
		return [];
	}
	const response = await res.json();
	// Practice endpoint returns paginated response
	const data: ApiStudyCard[] = response.data ?? response;
	return toStudyCards(data);
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
		return {
			type: 'basic',
			front: frontField?.type === 'rich_text' && typeof frontField.value === 'object'
				? (frontField.value as JSONContent)
				: String(frontField?.value ?? ''),
			back: backField?.type === 'rich_text' && typeof backField.value === 'object'
				? (backField.value as JSONContent)
				: String(backField?.value ?? '')
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

		return { type: 'cloze', front, clozeAnswer } satisfies ClozeCardDisplay;
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
			revealStyle: data.revealStyle ?? 'show_label'
		} satisfies ImageOcclusionCardDisplay;
	}

	return { type: 'basic', front: 'Unknown card type', back: '' };
}
